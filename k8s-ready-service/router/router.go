package router

import "github.com/gorilla/mux"
import (
	"github.com/abbi-gaurav/go-learning-projects/k8s-ready-service/handlers"
	"log"
	"sync/atomic"
	"time"
)

func Router(buildTime, commit, release string) *mux.Router {
	isReady := atomic.Value{}
	isReady.Store(false)
	go func() {
		log.Println("Readyz probe is negative by default")
		time.Sleep(10 * time.Second)
		isReady.Store(true)
		log.Println("Readyz probe is positive")
	}()
	r := mux.NewRouter()
	r.HandleFunc("/home", handlers.Home(buildTime, commit, release)).Methods("GET")
	r.HandleFunc("/healthz", handlers.Healthz)
	r.HandleFunc("/readyz", handlers.Readyz(&isReady))
	return r
}
