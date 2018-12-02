package main

import (
	"sync"
	"math/rand"
	"time"
	"fmt"
)

const (
	numberOfGoroutines = 4
	taskLoad           = 10
)

var waitGroup sync.WaitGroup

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	tasks := make(chan string, taskLoad)

	waitGroup.Add(numberOfGoroutines)

	for gr := 1; gr <= numberOfGoroutines; gr++ {
		go worker(tasks, gr)
	}

	for post := 1; post <= taskLoad; post++ {
		tasks <- fmt.Sprintf("Task : %d", post)
	}

	close(tasks)

	waitGroup.Wait()
}

func worker(tasks chan string, worker int) {
	defer waitGroup.Done()

	for {
		task, ok := <-tasks

		if !ok {
			fmt.Printf("Worker %d : shutting down\n", worker)
			return
		}

		fmt.Printf("worker %d started task %s\n", worker, task)

		sleep := rand.Int63n(100)
		time.Sleep(time.Duration(sleep) * time.Millisecond)

		fmt.Printf("worker %d completed task %s\n", worker, task)
	}
}
