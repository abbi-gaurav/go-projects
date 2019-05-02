package main

import (
	"fmt"
	"time"
)

func main() {
	doWork :=
		func(
			done <-chan interface{},
			strings <-chan string,
		) <-chan interface{} {
			terminated := make(chan interface{})
			go func() {
				defer fmt.Println("do work existed")
				defer close(terminated)

				for {
					select {
					case s := <-strings:
						fmt.Println(s)
					case <-done:
						return
					}
				}
			}()
			return terminated
		}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("cancelling doWork goroutine")
		close(done)
	}()

	<-terminated
	fmt.Println("done.")
}
