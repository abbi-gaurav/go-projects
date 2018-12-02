package workerpool

import "sync"

type Work interface {
	Task()
}

type Pool struct {
	tasksChannel chan Work
	wg           sync.WaitGroup
}

func New(maxGoroutines int) *Pool {
	p := Pool{
		tasksChannel: make(chan Work),
	}
	p.wg.Add(maxGoroutines)

	for i := 0; i < maxGoroutines; i++ {
		go func() {
			//continuos loop
			//wait until work is available
			for w := range p.tasksChannel {
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

func (p *Pool) Run(w Work) {
	p.tasksChannel <- w
}

func (p *Pool) Shutdown() {
	close(p.tasksChannel)
	p.wg.Wait()
}
