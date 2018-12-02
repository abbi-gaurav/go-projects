package main

import (
	"net/http"
	"log"
)

type response struct {
	resp *http.Response
	url  string
}

func get(url string, r chan response) {
	if resp, err := http.Get(url); err == nil {
		r <- response{resp, url}
	}
}

func main() {
	first := make(chan response)

	for _, url := range [] string{
		"http://code.jquery.com/jquery-1.9.1.min.js",
		"https://ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js",} {
		go get(url, first)
	}

	responseData := <-first

	log.Printf("%+v", responseData)
}
