package main

import (
	"fmt"
	"math/rand"
)

func doWork(done <-chan interface{}) (<-chan interface{}, <-chan int) {
	hbStream := make(chan interface{}, 1)
	workStream := make(chan int)

	go func() {
		defer close(hbStream)
		defer close(workStream)

		for i := 0; i < 10; i++ {
			select {
			case hbStream <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case workStream <- rand.Intn(10):
			}
		}
	}()
	return hbStream, workStream
}

func main() {
	done := make(chan interface{})
	defer close(done)
	hb, results := doWork(done)

	for {
		select {
		case _, ok := <-hb:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("result %v\n", r)
			} else {
				return
			}
		}
	}
}
