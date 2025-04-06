package main

import (
	"fmt"
	"time"
)

type Job struct {
	Data          string
	SleepDuration time.Duration
}

func (job *Job) Process() {
	time.Sleep(job.SleepDuration)
	fmt.Println("Processed job with data:", job.Data)
}
