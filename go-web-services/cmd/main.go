package main

import (
	"context"
	"github.com/abbi-gaurav/go-projects/go-web-services/cmd/crud/handlers"
	"github.com/abbi-gaurav/go-projects/go-web-services/internal/platform/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	log.Println("main: Started")

	readTimeout := 5 * time.Second
	writeTimeout := 10 * time.Second
	shutdownTimeout := 5 * time.Second
	dbConnectTimeout := 25 * time.Second

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "got:got2015@ds039441.mongolab.com:39441/gotraining"
	}

	log.Println("main : Started : Capturing masterDB")

	masterDB, err := db.New(dbHost, dbConnectTimeout)
	if err != nil {
		log.Fatalf("startup : register db failed %v", err)
	}
	defer masterDB.Close()

	host := os.Getenv("HOST")
	if host == "" {
		host = ":3000"
	}

	server := http.Server{
		Addr:           host,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
		Handler:        handlers.API(masterDB),
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		log.Printf("startup, listening on %s", host)
		log.Printf("shutdown: Listener closed %v", server.ListenAndServe())
		wg.Done()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)
	<-osSignals

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown did not complete %v : %v", shutdownTimeout, err)

		if err := server.Close(); err != nil {
			log.Printf("Shutdown: Error killing server: %v", err)
		}
	}

	wg.Wait()
	log.Printf("main completed")
}
