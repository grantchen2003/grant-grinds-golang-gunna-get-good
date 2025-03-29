package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Closing connection error: %v", err)
		}

		fmt.Printf("Closed connection to %s\n", conn.RemoteAddr().String())
	}()

	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())

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
				fmt.Println("EOF Reached.")
				break
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

		// process data from client...
		fmt.Println(string(data))
	}

	// send response back to client
	response := fmt.Sprintf("Hi from server! %s", strings.Repeat("foo", 1000))
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error sending response: %v", err)
		return
	}
}

func main() {
	ts := &TcpServer{}
	ts.Start()
}
