package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var void struct{} = struct{}{}

func main() {
	err := Main()
	if err != nil {
		panic(err)
	}
}

func Main() error {
	logger := log.Default()

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
		serverStatusCh, err := LaunchPalWorldServer(config, logger)
		if err != nil {
			return err
		}

		time.Sleep(10 * time.Second)

		logger.Println("Waiting for user existence check")
		shutdownUserCheckerCh, err := UserExistenceCheck(config, logger)
		if err != nil {
			return err
		}

		<-serverStatusCh
		logger.Println("PalWorld server is shutted down by some reason. Restarting...")
		shutdownUserCheckerCh <- void
	}
}

func LaunchPalWorldServer(config *Configuration, logger *log.Logger) (<-chan struct{}, error) {
	cmd := exec.Command(config.PalServerCommand[0], config.PalServerCommand[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Printf("Executing %s", strings.Join(config.PalServerCommand, " "))
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	ch := make(chan struct{})
	go func() {
		err := cmd.Wait()
		if err != nil {
			logger.Printf("An error occurred from the PalServer command: %s\n", err)
		}
		ch <- void
	}()

	return ch, nil
}

var userExistsRe = regexp.MustCompile(`\d+,\d+`)

func UserExistenceCheck(config *Configuration, logger *log.Logger) (chan<- struct{}, error) {
	shutdownCh := make(chan struct{})

	go func() {

		threshold := 30
		userEmptyCount := 0

		for {
			select {
			case <-shutdownCh:
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
					if err := ShutdownPalWorldServer(config, logger); err != nil {
						logger.Printf("An error occurred while shutting down PalWorld server: %s\n", err)
					}
					break
				}
			}
		}
	}()

	return shutdownCh, nil
}

func ShutdownPalWorldServer(config *Configuration, logger *log.Logger) error {
	err := RconShutdown(config)
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)
	return nil
}
