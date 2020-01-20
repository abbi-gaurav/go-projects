package main

import (
	"context"
	"github.com/abbi-gaurav/go-projects/k8s-ready-service/router"
	"github.com/abbi-gaurav/go-projects/k8s-ready-service/version"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Printf("Starting the service\n, commit: %s, build time: %s, release: %s",
		version.Commit, version.BuildTime, version.Release)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port is not set")
	}
	rtr := router.Router(version.BuildTime, version.Commit, version.Release)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: rtr,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Println("Got os interrupt...")
	case syscall.SIGTERM:
		log.Println("Got SIGTERM")
	}
	log.Println("The service is shutting down....")

	srv.Shutdown(context.Background())

	log.Println("Done..")
}
