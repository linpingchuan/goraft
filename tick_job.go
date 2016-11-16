package goraft

import (
	"fmt"
	"sync"
	"time"
)

const (
	running int = iota
	pause
	stop
)

type TickJob struct {
	duration time.Duration
	job      func()
	// for pause
	cond    *sync.Cond
	signal  chan int
	state   int
	started bool
}

func newTickJob(job func()) *TickJob {
	return &TickJob{
		duration: 500 * time.Millisecond,
		job:      job,
		cond:     sync.NewCond(&sync.Mutex{}),
		signal:   make(chan int, 1),
		state:    stop,
	}
}

func (job *TickJob) start() {
	job.started = true

	go func() {
		for {
			select {
			case state := <-job.signal:
				switch state {
				case stop:
					fmt.Println("stop...", state)
					job.started = false
					return
				case pause:
					fmt.Println("pause...", state)

					job.cond.L.Lock()
					for job.state != running {
						job.cond.Wait()
					}
					job.cond.L.Unlock()
				}
			default:
			}

			// do work here
			job.job()

			time.Sleep(job.duration)
		}
	}()
}

func (job *TickJob) stop() {
	job.state = stop
	job.signal <- stop
}

func (job *TickJob) pause() {
	job.state = pause
	job.signal <- pause
}

func (job *TickJob) resume() {
	job.state = running
	job.cond.Signal()
}
