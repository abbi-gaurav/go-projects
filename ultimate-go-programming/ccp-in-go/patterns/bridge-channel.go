package main

import (
	"fmt"
	"github.com/abbi-gaurav/go-learning-projects/ultimate-go-programming/ccp-in-go/patterns/utils"
)

func genVals() <-chan <-chan interface{} {
	chanStream := make(chan (<-chan interface{}))

	go func() {
		defer close(chanStream)
		for i := 0; i < 10; i++ {
			stream := make(chan interface{}, 1)
			stream <- i
			close(stream)
			chanStream <- stream
		}
	}()
	return chanStream
}

func main() {
	done := make(chan interface{})
	defer close(done)

	for v := range utils.Bridge(done, genVals()) {
		fmt.Printf("%v ", v)
	}

}
