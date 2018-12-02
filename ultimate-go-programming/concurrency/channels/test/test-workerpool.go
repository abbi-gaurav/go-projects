package main

import (
	"log"
	"time"
	"github.com/go-learning-projects/ultimate-go-programming/concurrency/channels/workerpool"
	"sync"
)

var names = [] string{
	"ss",
	"ga",
	"sa",
	"ap",
	"ra",
}

type namePrinter struct {
	name string
}

func (n *namePrinter) Task() {
	log.Println(n.name)
	time.Sleep(time.Second)
}

func main() {
	p := workerpool.New(2)

	var wg sync.WaitGroup
	wg.Add(100 * len(names))

	for i := 99; i < 100; i++ {
		for _, name := range names {
			np := namePrinter{
				name: name,
			}

			go func() {
				p.Run(&np)
				wg.Done()
			}()
		}
	}

	wg.Wait()

	p.Shutdown()

}
