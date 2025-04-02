package main

import (
	"fmt"
	"net/http"
	"time"
)

func makeRequest(url string, ch chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		ch <- "Rate-limited: 429 Too Many Requests"
	} else {
		ch <- fmt.Sprintf("Success: %d", resp.StatusCode)
	}
}

func main() {
	serverURL := "http://localhost:8080"
	ch := make(chan string)

	// Simulate multiple requests concurrently
	for range 20 {
		go makeRequest(serverURL, ch)
	}

	// Collect responses
	for range 20 {
		fmt.Println(<-ch)
	}

	// Wait for rate limit interval
	time.Sleep(5 * time.Second)

	for range 10 {
		go makeRequest(serverURL, ch)
	}

	for range 10 {
		fmt.Println(<-ch)
	}
}
