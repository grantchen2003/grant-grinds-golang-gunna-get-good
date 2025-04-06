package main

type Worker struct {
	JobQueue    chan Job
	DoneChannel chan bool
}

func (w *Worker) Start() {
	for job := range w.JobQueue {
		job.Process()
	}

	// Signal that the worker is done
	w.DoneChannel <- true
}
