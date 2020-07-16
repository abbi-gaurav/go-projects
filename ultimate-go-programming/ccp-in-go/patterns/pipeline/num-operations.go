package main

import "fmt"
import "github.com/abbi-gaurav/go-projects/ultimate-go-programming/ccp-in-go/patterns/utils"

func multiply(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
	processedStream := make(chan int)

	go func() {
		defer close(processedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case processedStream <- i * multiplier:
			}
		}
	}()
	return processedStream
}

func add(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
	processedStream := make(chan int)

	go func() {
		defer close(processedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case processedStream <- i + additive:
			}
		}
	}()
	return processedStream
}

func main() {
	done := make(chan interface{})
	defer close(done)

	intStream := utils.ToInt(done, utils.Generator(done, 1, 2, 3, 4))

	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}
