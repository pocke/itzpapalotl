package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	err := Main()
	if err != nil {
		panic(err)
	}
}

func Main() error {
	logger := log.Default()
	// This ctx will be done when the PalServer is shutted down
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Println("Starting Itzpapalotl")
	logger.Println("Loading configuration")
	config, err := NewConfiguration(os.Args[1:])
	if err != nil {
		return err
	}

	for {
		logger.Println("Waiting for UDP request")
		err := WaitUdpRequest(config, logger)
		if err != nil {
			return err
		}

		logger.Println("Launching PalWorld server")
		err = LaunchPalWorldServer(cancel, config, logger)
		if err != nil {
			return err
		}

		time.Sleep(10 * time.Second)

		logger.Println("Waiting for user existence check")
		UserExistenceCheck(ctx, config, logger)
		MemoryUsageCheck(ctx, cancel, config, logger)

		<-ctx.Done()
		logger.Println("PalWorld server is shutted down by some reason. Restarting...")
	}
}

func LaunchPalWorldServer(cancel context.CancelFunc, config *Configuration, logger *log.Logger) error {
	cmd := exec.Command(config.PalServerCommand[0], config.PalServerCommand[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Printf("Executing %s", strings.Join(config.PalServerCommand, " "))
	err := cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			logger.Printf("An error occurred from the PalServer command: %s\n", err)
		}
		cancel()
	}()

	return nil
}

var userExistsRe = regexp.MustCompile(`\d+,\d+`)

func UserExistenceCheck(ctx context.Context, config *Configuration, logger *log.Logger) {
	go func() {
		threshold := 30
		userEmptyCount := 0

		for {
			select {
			case <-ctx.Done():
				logger.Println("UserExistenceCheck is shutting down")
				return
			default:
				// do nothing
			}

			time.Sleep(1 * time.Minute)

			resp, err := RconShowPlayers(config)
			if err != nil {
				logger.Printf("An error occurred while executing ShowPlayers: %s\n", err)
				continue
			}

			if userExistsRe.MatchString(resp) {
				if userEmptyCount > 0 {
					logger.Printf("User exists, resetting userEmptyCount from %d to 0\n", userEmptyCount)
				}
				userEmptyCount = 0
			} else {
				userEmptyCount++
				logger.Printf("User does not exist, userEmptyCount is increased and now %d\n", userEmptyCount)

				if userEmptyCount >= threshold {
					logger.Println("User does not exist for a while, shutting down PalWorld server")
					if err := ShutdownPalWorldServer(config, logger, 1*time.Second, ""); err != nil {
						logger.Printf("An error occurred while shutting down PalWorld server: %s\n", err)
					}
					break
				}
			}
		}
	}()
}

func MemoryUsageCheck(ctx context.Context, cancel context.CancelFunc, config *Configuration, logger *log.Logger) {
	go func() {
		for {
			time.Sleep(1 * time.Minute)

			select {
			case <-ctx.Done():
				logger.Println("MemoryUsageCheck is shutting down")
				return
			default:
				// do nothing
			}

			out, err := exec.Command("ps", "-s", strconv.Itoa(os.Getpid()), "-o", "rss").Output()
			if err != nil {
				logger.Printf("An error occurred while executing ps: %s\n", err)
				continue
			}

			mem := 0
			for _, b := range regexp.MustCompile(`\d+`).FindAll(out, -1) {
				i, err := strconv.Atoi(string(b))
				if err != nil {
					logger.Printf("An error occurred while parsing memory usage: %s\n", err)
					continue
				}
				mem += i
			}

			if mem > config.MemoryThreshold {
				logger.Printf("Memory usage exceeds the threshold: %d > %d. Restarting the PalServer in 5 minutes\n", mem, config.MemoryThreshold)
				err := ShutdownPalWorldServer(config, logger, 5*time.Minute, "This server will reboot in 5 minutes due to high memory usage. Please save your work.")
				if err != nil {
					logger.Printf("An error occurred while shutting down PalWorld server: %s\n", err)
				}
				break
			}
		}

		time.Sleep(4 * time.Minute)
		logger.Println("PalServer will be restarted in 1 minute due to high memory usage")
		err := RconBroadcast(config, "Re-announcement. This server will reboot in 1 minute due to high memory usage. Please save your work.")
		if err != nil {
			logger.Printf("An error occurred while broadcasting: %s\n", err)
		}
	}()
}

func ShutdownPalWorldServer(config *Configuration, logger *log.Logger, wait time.Duration, msg string) error {
	err := RconShutdown(config, wait, msg)
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)
	return nil
}
