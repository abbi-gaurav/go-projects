package main

import (
	"math/rand"
	"fmt"
)

func main() {
	c := make(chan int)
	d := make(chan bool)

	go func(src chan int) {
		for {
			select {
			case src <- rand.Intn(100):
			case <-d:
				src = nil
			}
		}
	}(c)

	fmt.Printf("%d\n", <-c)
	fmt.Printf("%d\n", <-c)
	d <- true
	fmt.Printf("%d\n", <-c)
}
