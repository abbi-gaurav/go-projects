package main

import (
	"fmt"
	heartbeat2 "github.com/abbi-gaurav/go-learning-projects/ultimate-go-programming/ccp-in-go/conc-at-scale/heartbeat"
	"time"
)

func main() {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second

	heartbeat, results := heartbeat2.DoWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results: %v\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}
}
