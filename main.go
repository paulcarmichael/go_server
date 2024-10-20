package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const echo = "echo "

func main() {
	// define the port to bind to
	port := ":4173"

	// create a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", port)

	if err != nil {
		fmt.Printf("%s Failed to resolve UDP address: %v", getTimestamp(), err)
		os.Exit(1)
	}

	// create a UDP connection
	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		fmt.Printf("%s Failed to listen on UDP port: %v", getTimestamp(), err)
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Printf("%s Listening on UDP port %s\n", getTimestamp(), port)

	// create a data buffer
	// TODO what happens if the packet exceeds the buffer size?
	buffer := make([]byte, 1024)

	// handle data
	for {
		// read UDP packet
		n, addr, err := conn.ReadFromUDP(buffer)

		if err != nil {
			fmt.Printf("%s Error reading UDP packet: %v", getTimestamp(), err)
			continue
		}

		datagram := string(buffer[:n])

		// log the packet
		fmt.Printf("%s Read %d bytes from %s: %s", getTimestamp(), n, addr.String(), datagram)

		// respond
		response := echo + datagram

		_, err = conn.WriteToUDP([]byte(response), addr)
		if err != nil {
			fmt.Printf("%s Error sending response: %v", getTimestamp(), err)
		}

		fmt.Printf("%s Sent %d bytes to %s: %s", getTimestamp(), len(response), addr.String(), string(response))
	}
}

func getTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}
