package main

import (
	"fmt"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/app"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/informer"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/internal/opts"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var onlyOneSignalHandler = make(chan struct{})

func main() {
	options := opts.ParseFlags()

	stopCh := setupSignalHandler()
	cakeInformer := informer.CreateInformer(options.ResyncDuration)
	application := app.New(options, cakeInformer)

	err := application.Run(stopCh)
	println(err)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", options.Port), application.ServerMux))

	<-stopCh

	application.ShutDown()

}

func setupSignalHandler() <-chan struct{} {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
