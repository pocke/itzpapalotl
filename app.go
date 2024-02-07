package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type App struct {
	config                 *Configuration
	logger                 *log.Logger
	startServerImmediately bool
}

func NewApp(args []string) (*App, error) {
	config, err := NewConfiguration(args)
	if err != nil {
		return nil, err
	}

	logger := log.Default()

	return &App{
		config:                 config,
		logger:                 logger,
		startServerImmediately: false,
	}, nil
}

// Wait for a UDP request on the PalWorld server port.
// If a request is received, it returns.
func (app *App) WaitUdpRequest() error {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(app.config.PalWorldServerPort))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	app.logger.Println("Start listening UDP request")
	var buf [1]byte
	_, _, err = conn.ReadFromUDP(buf[0:])
	if err != nil {
		return err
	}

	return nil
}

func (app *App) LaunchPalWorldServer(cancel context.CancelFunc) error {
	if app.startServerImmediately {
		app.logger.Println("Starting PalServer immediately without waiting for an UDP request")
		app.startServerImmediately = false
		return nil
	}

	cmd := exec.Command(app.config.PalServerCommand[0], app.config.PalServerCommand[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	app.logger.Printf("Executing %s", strings.Join(app.config.PalServerCommand, " "))
	err := cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			app.logger.Printf("An error occurred from the PalServer command: %s\n", err)
		}
		cancel()
	}()

	return nil
}

var userExistsRe = regexp.MustCompile(`\d+,\d+`)

func (app *App) UserExistenceCheck(ctx context.Context) {
	go func() {
		threshold := 30
		userEmptyCount := 0

		for {
			time.Sleep(1 * time.Minute)

			select {
			case <-ctx.Done():
				app.logger.Println("UserExistenceCheck is shutting down")
				return
			default:
				// do nothing
			}

			resp, err := RconShowPlayers(app.config)
			if err != nil {
				app.logger.Printf("An error occurred while executing ShowPlayers: %s\n", err)
				continue
			}

			if userExistsRe.MatchString(resp) {
				if userEmptyCount > 0 {
					app.logger.Printf("User exists, resetting userEmptyCount from %d to 0\n", userEmptyCount)
				}
				userEmptyCount = 0
			} else {
				userEmptyCount++
				app.logger.Printf("User does not exist, userEmptyCount is increased and now %d\n", userEmptyCount)

				if userEmptyCount >= threshold {
					app.logger.Println("User does not exist for a while, shutting down PalWorld server")
					if err := app.ShutdownPalWorldServer(1*time.Second, ""); err != nil {
						app.logger.Printf("An error occurred while shutting down PalWorld server: %s\n", err)
					}
					break
				}
			}
		}
	}()
}

func (app *App) MemoryUsageCheck(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		for {
			time.Sleep(1 * time.Minute)

			select {
			case <-ctx.Done():
				app.logger.Println("MemoryUsageCheck is shutting down")
				return
			default:
				// do nothing
			}

			out, err := exec.Command("ps", "-s", strconv.Itoa(os.Getpid()), "-o", "rss").Output()
			if err != nil {
				app.logger.Printf("An error occurred while executing ps: %s\n", err)
				continue
			}

			mem := 0
			for _, b := range regexp.MustCompile(`\d+`).FindAll(out, -1) {
				i, err := strconv.Atoi(string(b))
				if err != nil {
					app.logger.Printf("An error occurred while parsing memory usage: %s\n", err)
					continue
				}
				mem += i
			}

			if mem > app.config.MemoryThreshold {
				app.logger.Printf("Memory usage exceeds the threshold: %d > %d. Restarting the PalServer in 5 minutes\n", mem, app.config.MemoryThreshold)
				err := app.ShutdownPalWorldServer(5*time.Minute, "This server will reboot in 5 minutes due to high memory usage. Please save your work.")
				if err != nil {
					app.logger.Printf("An error occurred while shutting down PalWorld server: %s\n", err)
				}
				break
			}
		}

		time.Sleep(4 * time.Minute)
		app.logger.Println("PalServer will be restarted in 1 minute due to high memory usage")
		err := RconBroadcast(app.config, "Re-announcement. This server will reboot in 1 minute due to high memory usage. Please save your work.")
		if err != nil {
			app.logger.Printf("An error occurred while broadcasting: %s\n", err)
		}

		app.startServerImmediately = true
	}()
}

func (app *App) ShutdownPalWorldServer(wait time.Duration, msg string) error {
	err := RconShutdown(app.config, wait, msg)
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)
	return nil
}
