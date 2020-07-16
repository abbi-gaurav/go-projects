package main

import (
	"fmt"
	"github.com/abbi-gaurav/go-projects/ultimate-go-programming/ccp-in-go/patterns/utils"
	"time"
)

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	start := time.Now()

	<-utils.Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(7*time.Second),
		sig(1*time.Second),
	)

	fmt.Printf("time since start %v", time.Since(start))
}
