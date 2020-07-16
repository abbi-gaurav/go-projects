package main

import (
	"fmt"
	"github.com/abbi-gaurav/go-projects/ultimate-go-programming/ccp-in-go/patterns/utils"
	sync2 "sync"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	out1, out2 := utils.Tee(done, utils.Take(done, utils.Repeat(done, 1, 2), 4))

	var wg sync2.WaitGroup
	wg.Add(2)

	go func() {
		for v1 := range out1 {
			fmt.Printf("v1 %d\n", v1)
		}
		wg.Done()
	}()

	go func() {
		for v2 := range out2 {
			fmt.Printf("v2 %d\n", v2)
		}
		wg.Done()
	}()

	wg.Wait()
}
