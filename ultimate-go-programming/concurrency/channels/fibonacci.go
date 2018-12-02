package main

import "fmt"

func fibonacci(series chan int, quit chan int) {
	x, y := 0, 1

	for {
		select {
		case series <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("Quit")
			return
		}
	}
}

func main() {
	series := make(chan int)
	quit := make(chan int)

	go func() {
		for i := 0; i < 20; i++ {
			fmt.Println(<-series)
		}
		quit <- 0
	}()

	fibonacci(series, quit)
}
