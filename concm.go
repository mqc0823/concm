package concm

import (
	"sync/atomic"
	"time"
)

type Pool struct {
	tasks   chan func()
	limiter chan struct{}
	// wg      sync.WaitGroup
	counter atomic.Uint64
}

func New() *Pool {
	return &Pool{}
}
func (p *Pool) WithMaxGoroutines(n int) *Pool {
	p.limiter = make(chan struct{}, n)
	return p
}
func (p *Pool) Start(n int) *Pool {
	p.tasks = make(chan func(), n)
	go func() {
		for {
			select {
			case p.limiter <- struct{}{}:
				go func(f func()) {
					// defer p.wg.Done()
					defer p.counter.Add(^uint64(0))
					f()
					<-p.limiter
				}(<-p.tasks)
			default:
				continue
			}
		}
	}()
	return p
}
func (p *Pool) Go(f func()) {
	// p.wg.Add(1)
	p.counter.Add(1)
	p.tasks <- f
}
func (p *Pool) Wait() {
	// p.wg.Wait()
	for {
		if p.counter.CompareAndSwap(0, 0) {
			return
		}
		time.Sleep(1 * time.Second)
	}

}
