package main

import (
	"net/http"
	"log"
)

type work struct {
	url  string
	resp chan *http.Response
}

func getter(w chan work) {
	for {
		do := <-w
		resp, _ := http.Get(do.url)
		do.resp <- resp
	}
}

func main() {
	w := make(chan work)
	go getter(w)

	resp := make(chan *http.Response)
	w <- work{
		"http://cdnjs.cloudfare.com/jquery/1.9.1/jquery.min.js",
		resp,
	}

	r := <-resp
	log.Printf("%+v", r )
}
