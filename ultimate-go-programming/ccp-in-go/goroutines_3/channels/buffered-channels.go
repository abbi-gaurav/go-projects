package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	var stdOutBuffer bytes.Buffer
	defer stdOutBuffer.WriteTo(os.Stdout)

	intStream := make(chan int, 4)

	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdOutBuffer, "producer done")

		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdOutBuffer, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Fprintf(&stdOutBuffer, "received %v.\n", integer)
	}
}
