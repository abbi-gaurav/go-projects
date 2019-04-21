package main

import "fmt"

func main() {
	stringStream := make(chan string)
	go func() {
		stringStream <- "hi channels"
	}()

	salutation, ok := <-stringStream
	fmt.Println(salutation)
	fmt.Println(ok)
}
