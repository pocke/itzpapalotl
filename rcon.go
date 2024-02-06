package main

import (
	"fmt"
	"regexp"
	"time"

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

func RconShutdown(config *Configuration, wait time.Duration, msg string) error {
	conn, err := NewRconClient(config)
	if err != nil {
		return err
	}
	defer conn.Close()

	msg = escapeMessage(msg)
	_, err = conn.Execute(fmt.Sprintf("Shutdown %d %s", int(wait.Seconds()), msg))
	return err
}

func RconBroadcast(config *Configuration, msg string) error {
	conn, err := NewRconClient(config)
	if err != nil {
		return err
	}
	defer conn.Close()

	msg = escapeMessage(msg)
	_, err = conn.Execute(fmt.Sprintf("Broadcast %s", msg))
	return err
}

func escapeMessage(msg string) string {
	return regexp.MustCompile(`\s`).ReplaceAllString(msg, "_")
}
