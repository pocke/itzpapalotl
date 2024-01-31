package main

import (
	"flag"
	"fmt"
)

type Configuration struct {
	PalWorldServerPort int
	PalServerCommand   []string
	RconPort           int
	AdminPassword      string
}

func NewConfiguration(args []string) (*Configuration, error) {
	fs := flag.NewFlagSet("itzpapalotl", flag.ExitOnError)
	fs.Usage = func() {
		o := fs.Output()
		fmt.Fprintln(o, "Usage: itzpapalotl [options] -- [palworld server command]")
		fs.PrintDefaults()
	}
	serverPort := fs.Int("server-port", 8211, "PalWorld server port")
	rocnPort := fs.Int("rcon-port", 25575, "RCON port")
	adminPassword := fs.String("admin-password", "", "Admin password")
	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	return &Configuration{
		PalWorldServerPort: *serverPort,
		PalServerCommand:   fs.Args(),
		RconPort:           *rocnPort,
		AdminPassword:      *adminPassword,
	}, nil

}
