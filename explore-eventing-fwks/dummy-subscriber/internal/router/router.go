package router

import (
	"github.com/abbi-gaurav/go-projects/explore-eventing-fwks/dummy-subscriber/internal/config"
	"github.com/abbi-gaurav/go-projects/explore-eventing-fwks/dummy-subscriber/internal/handlers"
	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func New() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/v1/events", handlers.EventHandler)
	if config.AppConfig().LogRequest {
		gh.LoggingHandler(os.Stdout, router)
	}
	return router
}
