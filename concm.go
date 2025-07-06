package concm

import "sync"

type Pool struct {
	tasks   chan func()
	limiter chan struct{}
	wg      sync.WaitGroup
	over    chan struct{}
	// counter atomic.Uint64
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
	LOOP:
		for {
			select {
			case p.limiter <- struct{}{}:
				go func(f func()) {
					defer func() {
						p.wg.Done()
						<-p.limiter
					}()
					// defer p.counter.Add(^uint64(0))
					f()

				}(<-p.tasks)
			case <-p.over:
				break LOOP
				// default:
				// 	continue
			}
		}
	}()
	return p
}
func (p *Pool) Go(f func()) {
	p.wg.Add(1)
	// p.counter.Add(1)
	p.tasks <- f
}
func (p *Pool) Wait() {
	p.wg.Wait()
	p.over <- struct{}{}
	// for {
	// 	if p.counter.CompareAndSwap(0, 0) {
	// 		return
	// 	}
	// 	time.Sleep(1 * time.Second)
	// }

}
