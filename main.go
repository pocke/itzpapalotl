package main

import (
	"context"
	"os"
	"time"
)

func main() {
	err := Main()
	if err != nil {
		panic(err)
	}
}

func Main() error {

	app, err := NewApp(os.Args[1:])
	if err != nil {
		return err
	}

	app.logger.Println("Starting Itzpapalotl")

	for {
		err := inLoop(app)
		if err != nil {
			return err
		}
	}
}

func inLoop(app *App) error {
	// This ctx will be done when the PalServer is shutted down
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.logger.Println("Waiting for UDP request")
	err := app.WaitUdpRequest()
	if err != nil {
		return err
	}

	app.logger.Println("Launching PalWorld server")
	err = app.LaunchPalWorldServer(cancel)
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)

	app.logger.Println("Waiting for user existence check")
	app.UserExistenceCheck(ctx)
	app.MemoryUsageCheck(ctx, cancel)

	<-ctx.Done()
	app.logger.Println("PalWorld server is shutted down by some reason. Restarting...")

	return nil
}
