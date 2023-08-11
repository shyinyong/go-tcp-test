package chat

// 协程池结构
type Pool struct {
	size     int
	workerCh chan struct{}
}

func NewPool(size int) *Pool {
	return &Pool{
		size:     size,
		workerCh: make(chan struct{}, size),
	}
}

func (p *Pool) Submit(task func()) {
	p.workerCh <- struct{}{}
	go func() {
		task()
		<-p.workerCh
	}()
}

func (p *Pool) Close() {
	close(p.workerCh)
}
