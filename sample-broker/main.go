package main

import (
	"code.cloudfoundry.org/lager"
	"context"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/broker"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/middleware"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/model"
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

	availableServiceTemplates, err := model.Parse("/Users/d066419/go/src/github.com/abbi-gaurav/go-projects/sample-broker/examples/catalog.yaml")
	if err != nil {
		brokerLogger.Fatal("parse", err)
	}
	brokerLogger.Info("parsed", lager.Data{"availableServiceTemplates": availableServiceTemplates})

	service, err := middleware.New(brokerLogger)
	if err != nil {
		brokerLogger.Fatal("creating-middleware", err)
	}

	serviceBroker := broker.NewBroker(brokerLogger, availableServiceTemplates, service)

	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerapi.BrokerCredentials{
		Username: "admin",
		Password: "nimda123",
	})

	runServer(brokerAPI, brokerLogger)

}

func runServer(brokerAPI http.Handler, brokerLogger lager.Logger) {
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
