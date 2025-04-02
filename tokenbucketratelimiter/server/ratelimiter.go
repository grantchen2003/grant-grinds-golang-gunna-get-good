package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type TokenBucketRateLimiter struct {
	mutex                   sync.Mutex
	capacity                int
	tokens                  int
	tokensPerRefillInterval int
	refillInterval          time.Duration
}

func NewTokenBucketRateLimiter() *TokenBucketRateLimiter {
	tbrl := &TokenBucketRateLimiter{
		mutex:                   sync.Mutex{},
		capacity:                10,
		tokens:                  10,
		tokensPerRefillInterval: 5,
		refillInterval:          5 * time.Second,
	}

	go tbrl.refill()

	return tbrl
}

func (tbrl *TokenBucketRateLimiter) refill() {
	for {
		time.Sleep(tbrl.refillInterval)

		tbrl.mutex.Lock()
		tbrl.tokens += tbrl.tokensPerRefillInterval
		if tbrl.tokens > tbrl.capacity {
			tbrl.tokens = tbrl.capacity
		}
		tbrl.mutex.Unlock()
	}
}

func (tbrl *TokenBucketRateLimiter) RateLimit(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tbrl.mutex.Lock()

		if tbrl.tokens == 0 {
			tbrl.mutex.Unlock()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			log.Println("Rate limited")
			return
		}

		tbrl.tokens -= 1
		tbrl.mutex.Unlock()
		handler(w, r)
	}
}
