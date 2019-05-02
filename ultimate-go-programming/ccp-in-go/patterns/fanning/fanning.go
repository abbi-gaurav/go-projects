package fanning

import "sync"

type processor = func(done <-chan interface{}, stream <-chan interface{}, id int) <-chan interface{}

func FanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))

	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func FanOut(done <-chan interface{}, p processor, factor int, input <-chan interface{}) []<-chan interface{} {
	streams := make([]<-chan interface{}, factor)
	for i := 0; i < factor; i++ {
		streams[i] = p(done, input, i)
	}
	return streams
}
