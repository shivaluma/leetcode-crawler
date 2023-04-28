package workerpool

import (
	"time"
)

const (
	timeoutInSeconds = 10
)

type WorkerPool struct {
	taskQueue  chan func()
	maxWorkers int
}

func New(maxWorkers int) *WorkerPool {
	return &WorkerPool{
		taskQueue:  make(chan func()),
		maxWorkers: maxWorkers,
	}
}

func (wp *WorkerPool) Enqueue(job func()) {
	wp.taskQueue <- job
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.maxWorkers; i++ {
		go wp.run()
	}
}

func (wp *WorkerPool) run() {
	for job := range wp.taskQueue {
		job()
		time.Sleep(time.Second * 2)
	}
}
