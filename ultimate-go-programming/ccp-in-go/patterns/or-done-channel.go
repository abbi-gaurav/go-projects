package main

import (
	"fmt"
	"github.com/abbi-gaurav/go-projects/ultimate-go-programming/ccp-in-go/patterns/utils"
)

func main() {
	done := make(chan interface{})
	defer close(done)
	intChannel := utils.Generator(done, 1, 2, 3, 4, 5, 6, 7)

	for val := range utils.OrDone(done, intChannel) {
		fmt.Println(val)
	}
}
