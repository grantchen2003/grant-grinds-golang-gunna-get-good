package main

import "fmt"

type WorkerPool struct {
	JobQueue    chan Job
	WorkerCount int
	Workers     []*Worker
	DoneChannel chan bool
}

func NewWorkerPool(workerCount int) *WorkerPool {
	var workers = make([]*Worker, workerCount)

	wp := &WorkerPool{
		JobQueue:    make(chan Job),
		WorkerCount: workerCount,
		Workers:     []*Worker{},
		DoneChannel: make(chan bool),
	}

	// Create workers and start them
	for i := range workerCount {
		worker := &Worker{JobQueue: wp.JobQueue, DoneChannel: wp.DoneChannel}
		workers[i] = worker
		go worker.Start()
	}

	return wp
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.JobQueue <- job
}

func (wp *WorkerPool) Shutdown() {
	// We first close the JobQueue to signal to the workers
	// that no more jobs will be sent, allowing them to stop
	// processing and exit their loops. Once the workers detect
	// the closed channel, they finish their remaining tasks
	// and signal via DoneChannel that they are done. This
	// ensures that wp.Shutdown doesn't block indefinitely and
	// only proceeds once all workers have completed their
	// tasks, preventing premature termination or inconsistency
	close(wp.JobQueue)

	for range wp.WorkerCount {
		// Waiting on DoneChannel ensures that the
		// program doesn't proceed until all workers
		// have finished their tasks, preventing premature
		// termination or inconsistency.
		<-wp.DoneChannel
	}

	// No need to close wp.DoneChannel because it's only used
	// for receiving completion signals, and closing it is
	// unnecessary and could lead to confusion or errors if
	// other parts of the program interact with it.

	fmt.Println("All workers have stopped.")
}
