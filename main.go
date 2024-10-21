package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func getTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

func main() {
	// check the required command line parameters
	if len(os.Args) < 3 {
		fmt.Printf("%s Usage: main <protocol> <mode> <port>\n", getTimestamp())
		fmt.Printf("%s protocol: udp | tcp\n", getTimestamp())
		fmt.Printf("%s mode: server\n", getTimestamp())
		fmt.Printf("%s port: numeric value\n", getTimestamp())
		return
	}

	protocol := os.Args[1]

	if protocol != "udp" && protocol != "tcp" {
		fmt.Printf("%s protocol %s is not supported\n", getTimestamp(), protocol)
		return
	}

	mode := os.Args[2]

	if mode != "server" {
		fmt.Printf("%s mode %s is not supported\n", getTimestamp(), mode)
		return
	}

	port, err := strconv.Atoi(os.Args[3])

	if err != nil {
		fmt.Printf("%s port %s is not a number\n", getTimestamp(), os.Args[3])
		return
	}

	if protocol == "udp" {
		// create a UDP address
		udpAddr, err := net.ResolveUDPAddr("udp", ":"+os.Args[3])

		if err != nil {
			fmt.Printf("%s failed to resolve address: %v", getTimestamp(), err)
			return
		}

		// create a UDP connection
		conn, err := net.ListenUDP("udp", udpAddr)

		if err != nil {
			fmt.Printf("%s failed to listen on UDP port: %v", getTimestamp(), err)
			os.Exit(1)
		}

		defer conn.Close()

		fmt.Printf("%s listening on UDP port %d\n", getTimestamp(), port)

		// create a data buffer
		// TODO what happens if the packet exceeds the buffer size?
		buffer := make([]byte, 1024)

		// handle data
		// TODO should each packet be handled in distinct goroutines?
		for {
			// read UDP packet
			n, addr, err := conn.ReadFromUDP(buffer)

			if err != nil {
				fmt.Printf("%s error reading packet: %v", getTimestamp(), err)
				continue
			}

			go handleUDPConnection(conn, addr, buffer[:n])
		}
	} else if protocol == "tcp" {
		listener, err := net.Listen("tcp", "0.0.0.0:"+os.Args[3])

		if err != nil {
			fmt.Printf("%s failed to listen on TCP port: %v", getTimestamp(), err)
			os.Exit(1)
		}
		defer listener.Close()

		fmt.Printf("%s listening on TCP port %d\n", getTimestamp(), port)

		// handle connections
		for {
			conn, err := listener.Accept()

			if err != nil {
				fmt.Printf("%s failed to accept connection: %v\n", getTimestamp(), err)
				continue
			}

			fmt.Printf("%s new connection from %s\n", getTimestamp(), conn.RemoteAddr())

			go handleTCPConnection(conn)
		}
	}
}

func handleUDPConnection(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	fmt.Printf("%s read %d bytes from %s: %s", getTimestamp(), len(data), addr.String(), data)

	// respond (totally optional with UDP - just echoing for now)
	_, err := conn.WriteToUDP([]byte(data), addr)
	if err != nil {
		fmt.Printf("%s error sending response: %v", getTimestamp(), err)
		return
	}

	fmt.Printf("%s sent %d bytes to %s: %s", getTimestamp(), len(data), addr.String(), string(data))
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	// read incoming data from client
	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadString('\n')

		if err != nil {
			fmt.Printf("%s failed to read data: %v\n", getTimestamp(), err)
			return
		}

		fmt.Printf("%s read %d bytes from %s: %s", getTimestamp(), len(data), conn.RemoteAddr(), data)

		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Printf("%s error sending response: %v", getTimestamp(), err)
			return
		}

		fmt.Printf("%s sent %d bytes to %s: %s", getTimestamp(), len(data), conn.RemoteAddr(), data)
	}
}
