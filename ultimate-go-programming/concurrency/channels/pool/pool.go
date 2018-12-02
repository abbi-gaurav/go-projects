package pool

import (
	"sync"
	"io"
	"errors"
	"log"
)

type Pool struct {
	m         sync.Mutex
	resources chan io.Closer
	factory   func() (closer io.Closer, err error)
	closed    bool
}

var ErrPoolClosed = errors.New("pool has been closed")

func New(fn func() (closer io.Closer, err error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("size too small")
	}

	return &Pool{
		factory:   fn,
		resources: make(chan io.Closer, size),
	}, nil
}

func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.resources:
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil
	default:
		log.Println("Acquiring new resource...")
		return p.factory()
	}
}

func (p *Pool) Release(r io.Closer) {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		r.Close()
		return
	}

	select {
	case p.resources <- r:
		log.Println("Release ..", "In queue")
	default:
		log.Println("Release..", "Closing")
		r.Close()
	}
}

func (p *Pool) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	p.closed = true

	close(p.resources)

	for r := range p.resources {
		r.Close()
	}
}
