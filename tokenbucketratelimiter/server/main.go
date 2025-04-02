package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type TokenBucket struct {
	mutex                   sync.Mutex
	capacity                int
	tokens                  int
	tokensPerRefillInterval int
	refillInterval          time.Duration
}

func NewTokenBucket() *TokenBucket {
	tb := &TokenBucket{
		mutex:                   sync.Mutex{},
		capacity:                10,
		tokens:                  10,
		tokensPerRefillInterval: 5,
		refillInterval:          5 * time.Second,
	}

	go tb.refill()

	return tb
}

func (tb *TokenBucket) refill() {
	for {
		time.Sleep(tb.refillInterval)

		tb.mutex.Lock()
		tb.tokens += tb.tokensPerRefillInterval
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.mutex.Unlock()
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %s\n", strings.Split(r.RemoteAddr, ":")[0])
	w.Write([]byte("Request served successfully"))
}

func rateLimit(tb *TokenBucket, handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tb.mutex.Lock()

		if tb.tokens == 0 {
			tb.mutex.Unlock()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			log.Println("Rate limited")
			return
		}

		tb.tokens -= 1
		tb.mutex.Unlock()
		handler(w, r)
	}
}

func main() {
	tokenBucket := NewTokenBucket()

	http.HandleFunc("/", rateLimit(tokenBucket, handler))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server", err)
	}
}
