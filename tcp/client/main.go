package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	tcpServerAddress, err := net.ResolveTCPAddr("tcp", ":3000")
	if err != nil {
		log.Fatal("ResolveTCPAddr failed:", err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, tcpServerAddress)
	if err != nil {
		log.Fatal("Failed to connect to server: %w", err)
	}
	defer conn.Close()

	data := []byte(strings.Repeat("Hi", 2000))

	// When using conn.Write(data), the data is written to an
	// internal buffer, and if the buffer is full, the write
	// operation will block until the operating system has
	// transmitted enough data over the network to free up space.
	// Once space becomes available in the buffer, more data
	// can be written to the buffer, and this process continues until
	// all the data is sent. The data is sent in chunks over a single
	// connection, with the operating system managing the flow and
	// buffering of the data without creating multiple network requests.
	n, err := conn.Write(data)
	if err != nil {
		log.Fatal("Failed to send data: %w", err)
	}

	fmt.Printf("Wrote %d bytes over the network\n", n)

	// Signal that no more data will be written (close the writing half of the connection)
	if err := conn.CloseWrite(); err != nil {
		log.Fatal("Failed to close write side of the connection: %w", err)
	}

	fmt.Println("EOF sent")

	// Read the server's response
	for {
		buf := make([]byte, 2048)
		n, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				// We have received all the data
				fmt.Println("EOF Reached.")
				break
			}

			log.Fatal("Failed to read response: %w", err)
			return
		}

		// process data from server...
		data := buf[:n]
		fmt.Println(string(data))
	}
}
