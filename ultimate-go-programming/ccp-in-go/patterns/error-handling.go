package main

import (
	"fmt"
	"net/http"
)

type result struct {
	err      error
	response *http.Response
}

func checkStatus(done <-chan interface{}, urls ...string) <-chan result {
	results := make(chan result)
	go func() {
		defer close(results)

		for _, url := range urls {
			resp, err := http.Get(url)
			res := result{err: err, response: resp}
			select {
			case <-done:
				return
			case results <- res:
			}
		}
	}()
	return results
}

func main() {
	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}

	for result := range checkStatus(done, urls...) {
		if result.err != nil {
			fmt.Printf("error : %v\n", result.err)
			continue
		}
		fmt.Printf("response %v\n", result.response)
	}
}
