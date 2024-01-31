package main

import (
	"fmt"

	"github.com/gorcon/rcon"
)

func NewRconClient(config *Configuration) (*rcon.Conn, error) {
	return rcon.Dial(fmt.Sprintf("127.0.0.1:%d", config.RconPort), config.AdminPassword)
}

func RconShowPlayers(config *Configuration) (string, error) {
	conn, err := NewRconClient(config)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return conn.Execute("ShowPlayers")
}

func RconShutdown(config *Configuration) error {
	conn, err := NewRconClient(config)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Execute("Shutdown 1")
	return err
}
