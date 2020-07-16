package main

import (
	"log"
	"io"
	"sync/atomic"
	"sync"
	"github.com/abbi-gaurav/go-projects/ultimate-go-programming/concurrency/channels/pool"
	"time"
	"math/rand"
)

const (
	maxGoroutines   = 25
	pooledResources = 2
)

type dbConnection struct {
	ID int32
}

func (dbConnection *dbConnection) Close() error {
	log.Println("Close: Connection", dbConnection.ID)
	return nil
}

var idCounter int32

func CreateConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCounter, 1)
	log.Println("Create new connection", id)

	return &dbConnection{id}, nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(maxGoroutines)

	p, err := pool.New(CreateConnection, pooledResources)

	if err != nil {
		log.Panic(err)
	}

	for query := 0; query < maxGoroutines; query++ {
		go func(q int) {
			performQueries(q, p)
			wg.Done()
		}(query)
	}

	wg.Wait()
	log.Println("Shutdown")
	p.Close()
}
func performQueries(query int, p *pool.Pool) {
	conn, err := p.Acquire()

	if err != nil {
		log.Println(err)
		return
	}

	defer p.Release(conn)
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	log.Printf("QID[%d] CID[%d]\n", query, conn.(*dbConnection).ID)
}
