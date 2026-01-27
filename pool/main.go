package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
Task:
Implement a worker pool that processes jobs concurrently and returns results.

Requirements:
- N workers
- jobs channel
- results channel
- graceful shutdown
*/

type Pool struct {
	count    int
	queue    chan func() int
	results  chan int
	jobsWg   sync.WaitGroup
	readerWg sync.WaitGroup
}

func NewPool(ctx context.Context, count int) *Pool {
	p := &Pool{
		count:   count,
		queue:   make(chan func() int),
		results: make(chan int, 50),
	}

	p.jobsWg.Add(count)
	for i := 0; i < count; i++ {
		go p.startWorker(ctx)
	}

	p.readerWg.Add(1)
	go p.startReader()

	return p
}

func (p *Pool) Add(job func() int) {
	p.queue <- job
}

func (p *Pool) Close() {
	close(p.queue)
	p.jobsWg.Wait()
	close(p.results)
	p.readerWg.Wait()
}

func (p *Pool) startReader() {
	defer p.readerWg.Done()

	for res := range p.results {
		fmt.Println("res: ", res)
	}
}

func (p *Pool) startWorker(ctx context.Context) {
	defer p.jobsWg.Done()

	for {
		select {
		case job, ok := <-p.queue:
			if !ok {
				return
			}
			p.results <- job()
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool := NewPool(ctx, 5)
	defer pool.Close()

	jobs := []func() int{
		func() int { return 1 },
		func() int { return 2 },
		func() int {
			time.Sleep(5 * time.Second)
			return 3
		},
		func() int {
			time.Sleep(10 * time.Second)
			return 4
		},
		func() int { return 5 },
		func() int { return 6 },
		func() int {
			time.Sleep(3 * time.Second)
			return 7
		},
		func() int {
			time.Sleep(7 * time.Second)
			return 8
		},
	}

	for _, job := range jobs {
		pool.Add(job)
	}
}
