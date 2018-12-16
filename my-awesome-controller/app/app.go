package app

import (
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/controller"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/internal/opts"
	"github.com/abbi-gaurav/go-learning-projects/my-awesome-controller/pkg/client/informers/externalversions/awesome.controller.io/v1"
	"net/http"
)

type Application struct {
	cakeController *controller.CakeController
	ServerMux      *http.ServeMux
}

func New(opts *opts.Options, informer v1.CakeInformer) *Application {
	cakeController := controller.New(opts, informer)

	serverMux := http.NewServeMux()

	return &Application{
		cakeController: cakeController,
		ServerMux:      serverMux,
	}
}

func (a *Application) Run(stopCh <-chan struct{}) error {
	err := a.cakeController.Run(2, stopCh)
	return err
}
