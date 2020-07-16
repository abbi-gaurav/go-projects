package main

import (
	"code.cloudfoundry.org/lager"
	"context"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/broker"
	"github.com/pivotal-cf/brokerapi"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	brokerLogger := lager.NewLogger("sample-k8s-service-broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	serviceBroker := broker.NewBroker(brokerLogger)

	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerapi.BrokerCredentials{
		Username: "admin",
		Password: "nimda123",
	})

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	srv := &http.Server{
		Addr:    ":" + "8080",
		Handler: brokerAPI,
	}

	go func() {
		brokerLogger.Fatal("http-listen", srv.ListenAndServe())
	}()

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		brokerLogger.Info("os-interrupt")
	case syscall.SIGTERM:
		brokerLogger.Info("sigterm")
	}
	log.Println("The service is shutting down....")

	_ = srv.Shutdown(context.Background())

	brokerLogger.Info("Done")

}
