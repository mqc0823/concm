package concm

import "sync"

type Pool struct {
	tasks   chan func()
	limiter chan struct{}
	wg      sync.WaitGroup
	// counter atomic.Uint64
}

func New() *Pool {
	return &Pool{}
}
func (p *Pool) WithMaxGoroutines(n int) *Pool {
	p.limiter = make(chan struct{}, n)
	return p
}
func (p *Pool) Start(n int) {
	p.tasks = make(chan func(), n)
	go func() {
		for {
			select {
			case p.limiter <- struct{}{}:
				go func(f func()) {
					defer p.wg.Done()
					// defer p.counter.Add(uint64(0) << 32)
					f()
					<-p.limiter
				}(<-p.tasks)
			default:
				continue
			}
		}
	}()
}
func (p *Pool) Go(f func()) {
	p.wg.Add(1)
	// p.counter.Add(uint64(1) << 32)
	p.tasks <- f
}
func (p *Pool) Wait() {
	p.wg.Wait()

}
