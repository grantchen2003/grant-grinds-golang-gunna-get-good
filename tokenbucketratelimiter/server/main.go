package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %s\n", strings.Split(r.RemoteAddr, ":")[0])
	w.Write([]byte("Request served successfully"))
}

func main() {
	tokenBucketRateLimiter := NewTokenBucketRateLimiter()

	http.HandleFunc("/", tokenBucketRateLimiter.RateLimit(handler))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server", err)
	}
}
