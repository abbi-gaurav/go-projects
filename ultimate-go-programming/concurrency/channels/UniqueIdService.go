package main

import "fmt"

func main() {
	id := make(chan string)

	go func() {
		var counter int64 = 0
		for {
			id <- fmt.Sprintf("%x", counter)
			counter += 1
		}
	}()

	x := <-id
	println(x)

	x = <-id
	println(x)
}
