package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

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
	config, err := NewConfiguration()
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
		ch, err := LaunchPalWorldServer(config, logger)
		if err != nil {
			return err
		}
		// TODO
		<-ch
	}
}

func LaunchPalWorldServer(config *Configuration, logger *log.Logger) (<-chan error, error) {
	cmd := exec.Command(config.PalWorldCommandPath, config.PalWorldCommandArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Printf("Executing %s %s", config.PalWorldCommandPath, strings.Join(config.PalWorldCommandArgs, " "))
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	ch := make(chan error)
	go func() {
		err := cmd.Wait()
		ch <- err
	}()

	return ch, nil
}
