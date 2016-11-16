package goraft

import (
	"fmt"
	"testing"
	"time"
)

func TestTick(t *testing.T) {
	i := 1
	job := newTickJob(func() {
		i++
		fmt.Println("hello", i)
	})

	job.start()
	time.Sleep(time.Second)

	job.pause()
	time.Sleep(5 * time.Second)

	fmt.Println("try resume")
	job.resume()
	time.Sleep(time.Second)

	job.stop()
	time.Sleep(time.Second)
}
