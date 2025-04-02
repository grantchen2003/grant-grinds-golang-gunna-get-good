package main

import (
	"sync"
	"testing"
)

func TestConcurrentSafeMap(t *testing.T) {
	m := NewConcurrentSafeMap[string, int]()
	var wg sync.WaitGroup
	n := 1000 // Number of concurrent operations

	// Concurrent writes
	wg.Add(n)
	for i := range n {
		go func(i int) {
			defer wg.Done()
			m.Set(string(rune('A'+(i%26))), i)
		}(i)
	}

	wg.Wait()

	// Concurrent reads
	wg.Add(n)
	for i := range n {
		go func(i int) {
			defer wg.Done()
			m.Get(string(rune('A' + (i % 26))))
		}(i)
	}

	wg.Wait()

	// Concurrent existence checks
	wg.Add(n)
	for i := range n {
		go func(i int) {
			defer wg.Done()
			m.Exists(string(rune('A' + (i % 26))))
		}(i)
	}

	wg.Wait()

	// Ensure map is still usable after concurrent access
	m.Set("Z", 999)
	if val, exists := m.Get("Z"); !exists || val != 999 {
		t.Errorf("Expected key 'Z' to be 999, got %d, exists: %v", val, exists)
	}
}
