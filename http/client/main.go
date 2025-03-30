package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	data := []byte("Hello, server!")

	// http.Post is blocking
	// bytes.NewBuffer(data) wraps the []byte data into a *bytes.Buffer, which implements
	// the io.Reader interface required by http.Post() to send the request body.
	resp, err := http.Post("http://localhost:8080", "text/plain", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error making request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal(fmt.Sprintf("Error: Received status code %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	fmt.Println("Response from server:", string(body))
}
