package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

func IsPortOpen(port int) bool {
	address := fmt.Sprintf(":%d", port)

	listener, err := net.Dial("tcp", address)
	if err != nil {
		return false
	}

	defer listener.Close()

	return true
}

func getOpenPortsSynchronously(startPort int, endPort int) []int {
	var openPorts []int
	for i := startPort; i <= endPort; i++ {
		if IsPortOpen(i) {
			openPorts = append(openPorts, i)
		}
	}
	return openPorts
}

func getOpenPortsConcurrently(startPort int, endPort int) []int {
	var wg sync.WaitGroup

	openPortsCh := make(chan int)

	for i := startPort; i <= endPort; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			if IsPortOpen(port) {
				openPortsCh <- port
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(openPortsCh)
	}()

	var openPorts []int
	for port := range openPortsCh {
		openPorts = append(openPorts, port)
	}

	sort.Ints(openPorts)

	return openPorts
}

func slicesEqual(slice1, slice2 []int) bool {
	// Check if the lengths are different
	if len(slice1) != len(slice2) {
		return false
	}

	// Check if elements are the same in both slices
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}

func main() {
	// Maximum valid port number (65535) for TCP/UDP
	// In networking, ports are 16-bit unsigned integers,
	// which means they can range from 0 to 65535
	startPort, endPort := 0, 65535

	syncOpenPorts := getOpenPortsSynchronously(startPort, endPort)
	asyncOpenPorts := getOpenPortsConcurrently(startPort, endPort)

	// prints true
	fmt.Println(slicesEqual(syncOpenPorts, asyncOpenPorts))
}
