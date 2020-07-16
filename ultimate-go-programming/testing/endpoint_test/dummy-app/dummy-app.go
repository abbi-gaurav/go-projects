package main

import (
	"github.com/abbi-gaurav/go-projects/ultimate-go-programming/testing/endpoint_test/handlers"
	"net/http"
)

func main() {
	handlers.Routes()

	http.ListenAndServe(":4000", nil)
}
