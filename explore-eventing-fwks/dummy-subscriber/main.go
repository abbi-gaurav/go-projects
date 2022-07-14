package main

import (
	"context"
	"github.com/abbi-gaurav/go-projects/explore-eventing-fwks/dummy-subscriber/internal/config"
	"github.com/abbi-gaurav/go-projects/explore-eventing-fwks/dummy-subscriber/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("starting the service...")

	config.InitConfig()
	rtr := router.New()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	srv := &http.Server{
		Addr:    ":" + "8080",
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
