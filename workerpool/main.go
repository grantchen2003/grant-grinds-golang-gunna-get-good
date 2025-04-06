package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a new worker pool with 3 workers
	wp := NewWorkerPool(3)

	// Add 8 jobs to the pool
	for i := 1; i <= 8; i++ {
		job := Job{
			Data:          fmt.Sprintf("Task %d", i),
			SleepDuration: 2 * time.Second,
		}
		wp.AddJob(job)
	}

	// Shutdown the worker pool and wait for all workers to finish
	wp.Shutdown()
}
