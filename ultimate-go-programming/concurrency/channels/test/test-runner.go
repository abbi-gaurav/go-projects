package main

import (
	"time"
	"log"
	"github.com/abbi-gaurav/go-projects/ultimate-go-programming/concurrency/channels/runner"
	"os"
)

const timeout = 4 * time.Second

func main() {
	log.Println("Starting work...")

	r := runner.New(timeout)

	r.Add(createTask(), createTask(), createTask())

	if err := r.Start(); err != nil {
		switch err {
		case runner.ErrTimeout:
			log.Println("Timeout...")
			os.Exit(1)
		case runner.ErrInterrupt:
			log.Println("Interrupt")
			os.Exit(2)
		}
	}
	log.Println("Process ended...")
}
func createTask() func(int) {
	return func(id int) {
		log.Printf("Task.. %d", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}
