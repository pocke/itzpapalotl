package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	PalWorldServerPort  int
	PalWorldCommandPath string
	PalWorldCommandArgs []string
	RconPort            int
}

func NewConfiguration() (*Configuration, error) {
	palWorldServerPort, err := strconv.Atoi(os.Getenv("ITZPAPALOTL_PALWORLD_SERVER_PORT"))
	if err != nil {
		palWorldServerPort = 8211
	}

	palWorldCommandPath := os.Getenv("ITZPAPALOTL_PALWORLD_SERVER_PATH")
	if palWorldCommandPath == "" {
		return nil, errors.New("ITZPAPALOTL_PALWORLD_SERVER_PATH is not set")
	}

	palWorldCommandArgs := strings.Split(os.Getenv("ITZPAPALOTL_PALWORLD_SERVER_ARGS"), " ")

	rconPort, err := strconv.Atoi(os.Getenv("ITZPAPALOTL_RCON_PORT"))
	if err != nil {
		rconPort = 25575
	}

	return &Configuration{
		PalWorldServerPort:  palWorldServerPort,
		PalWorldCommandPath: palWorldCommandPath,
		PalWorldCommandArgs: palWorldCommandArgs,
		RconPort:            rconPort,
	}, nil
}
