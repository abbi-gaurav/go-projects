package main

import (
	"fmt"
	"github.com/abbi-gaurav/go-learning-projects/ultimate-go-programming/ccp-in-go/conc-at-scale/error-handling/intermediate"
	"log"
	"os"
)

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v] ", key))
	log.Printf("#%v", err)
	fmt.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := intermediate.RunJob("1")
	if err != nil {
		msg := "Unexpected issue, please report bug"
		if _, ok := err.(intermediate.InterimError); ok {
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}
