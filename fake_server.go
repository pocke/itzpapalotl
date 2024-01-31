package main

import (
	"log"
	"net"
	"strconv"
)

// Wait for a UDP request on the PalWorld server port.
// If a request is received, it returns.
func WaitUdpRequest(config *Configuration, logger *log.Logger) error {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(config.PalWorldServerPort))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	logger.Println("Start listening UDP request")
	var buf [1]byte
	_, _, err = conn.ReadFromUDP(buf[0:])
	if err != nil {
		return err
	}

	return nil
}
