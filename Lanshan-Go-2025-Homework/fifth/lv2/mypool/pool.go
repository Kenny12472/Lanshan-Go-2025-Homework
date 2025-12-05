package mypool

import "sync"

type Task func()

type Pool struct {
	tasks chan Task
	wg    sync.WaitGroup
}

func New(workerNum int) *Pool {
	p := &Pool{
		tasks: make(chan Task, 1024),
	}

	for i := 0; i < workerNum; i++ {
		go func() {
			for task := range p.tasks {
				task()
				p.wg.Done()
			}
		}()
	}

	return p
}

func (p *Pool) Submit(task Task) {
	p.wg.Add(1)
	p.tasks <- task
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Close() {
	close(p.tasks)
}
