package app

import (
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/controller"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/db"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/internal/opts"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/informers/externalversions/awesome.controller.io/v1"
	"net/http"
)

type Application struct {
	cakeController *controller.CakeController
	ServerMux      *http.ServeMux
	Database       db.DB
}

func New(opts *opts.Options, informer v1.CakeInformer) *Application {
	database := db.New(opts.DbType)
	cakeController := controller.New(informer, database)

	serverMux := http.NewServeMux()

	return &Application{
		cakeController: cakeController,
		ServerMux:      serverMux,
		Database:       database,
	}
}

func (a *Application) Run(stopCh <-chan struct{}) error {
	err := a.cakeController.Run(2, stopCh)
	return err
}

func (a *Application) ShutDown() {
	a.cakeController.ShutDown()
}
