package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type TcpServer struct{}

func (ts *TcpServer) Start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer ln.Close()

	fmt.Println("TCP server started on :3000")

	for {
		// ln.Accept() only establishes the connection, and it is
		// executed once per connection even if multiple chunks are
		// sent over a connection, as it is responsible solely for
		// accepting the connection, not for handling the data transfer.
		conn, err := ln.Accept()

		if err != nil {
			// log.Fatal logs the given message and then calls
			// os.Exit(1), which immediately terminates the program
			log.Fatal(err)
		}

		go ts.read(conn)
	}
}

func (ts *TcpServer) read(conn net.Conn) {
	// Even though the client closes the connection,
	// adding defer conn.Close() in the server ensures
	// proper resource cleanup, and calling conn.Close()
	// on an already closed connection won't cause any errors in Go.
	defer conn.Close()

	// Create a buffer of 2048 bytes to hold
	// incoming data from the connection
	buf := make([]byte, 2048)

	// Start an infinite loop to continuously
	// read data from the connection
	for {
		// Read data from the connection into the buffer
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				// Client has closed the connection
				fmt.Println("Client closed the connection.")
				return
			}
			// If an error occurs (e.g., connection closed
			// or read issue), log the error and return
			log.Printf("Connection read error: %v", err)
			return
		}

		fmt.Printf("Received %d bytes over the network\n", n)

		// Slice the buffer to only include the valid
		// portion of the buffer that was filled with data
		// Note: This will never cause an index error as long as
		// `n` is within the bounds of the buffer (which it is).
		data := buf[:n]

		// process data...
		fmt.Println(string(data))
	}
}

func makeRequest() error {
	data := []byte(strings.Repeat("Hi", 5))

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return fmt.Errorf("Failed to connect to server: %w", err)
	}
	defer conn.Close()

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
		return fmt.Errorf("Failed to send data: %w", err)
	}

	fmt.Printf("Wrote %d bytes over the network\n", n)

	return nil
}

func main() {
	go func() {
		time.Sleep(2 * time.Second)

		err := makeRequest()
		if err != nil {
			log.Printf("Client error: %v", err)
		}
	}()

	ts := &TcpServer{}
	ts.Start()
}
