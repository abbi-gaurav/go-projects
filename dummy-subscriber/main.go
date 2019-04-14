package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {

	var (
		port = flag.Int("port", 8080, "tcp port on which to listen for http requests")
	)
	flag.Parse()

	http.Handle("/v1/events", mockHandler())

	log.Printf("HTTP server starting on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func mockHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch requestMethod := &r.Method; *requestMethod {
		case http.MethodPost:
			if r.Body == nil {
				http.Error(w, "Please send a CloudEvent in the HTTP request body", http.StatusBadRequest)
				return
			}

			wireRep := formatRequest(r)
			fmt.Println("String rep")
			fmt.Println(wireRep)
			defer r.Body.Close()
			w.WriteHeader(http.StatusCreated)
		default:
			http.Error(w, fmt.Sprintf("HTTP method '%v' is not supported", *requestMethod), http.StatusMethodNotAllowed)
		}
	})
}

// formatRequest generates ascii representation of a request
func formatRequest(req *http.Request) string {
	// Save a copy of this request for debugging.
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	return string(requestDump)
}
