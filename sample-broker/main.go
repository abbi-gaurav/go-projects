package main

import (
	"code.cloudfoundry.org/lager"
	"context"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/broker"
	"github.com/abbi-gaurav/go-projects/sample-broker/internal/config"
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
	if err := config.InitConfig(); err != nil {
		log.Fatal("environment not loaded", err)
	}

	brokerLogger := setUpLogger()

	availableServiceTemplates, err := initializeCatalog(brokerLogger)

	serviceBroker := setUpBroker(err, brokerLogger, availableServiceTemplates)

	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerapi.BrokerCredentials{
		Username: config.AppConfig().UserName,
		Password: config.AppConfig().Password,
	})

	runServer(brokerAPI, brokerLogger)

}

func setUpBroker(err error, brokerLogger lager.Logger, availableServiceTemplates model.Services) *broker.K8SServiceBroker {
	service, err := middleware.New(brokerLogger)
	if err != nil {
		brokerLogger.Fatal("creating-middleware", err)
	}

	serviceBroker := broker.NewBroker(brokerLogger, availableServiceTemplates, service)
	return serviceBroker
}

func initializeCatalog(brokerLogger lager.Logger) (model.Services, error) {
	availableServiceTemplates, err := model.Parse(config.AppConfig().CatalogFilePath)
	if err != nil {
		brokerLogger.Fatal("parse", err)
	}
	brokerLogger.Info("parsed", lager.Data{"availableServiceTemplates": availableServiceTemplates})
	return availableServiceTemplates, err
}

func setUpLogger() lager.Logger {
	brokerLogger := lager.NewLogger("sample-k8s-service-broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))
	return brokerLogger
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
