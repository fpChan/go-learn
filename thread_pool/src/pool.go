package main

import (
	"fmt"
	"log"
	"sync"
)

type Job interface {
	handler() error
}

type worker struct {
	ownerPool    *Pool
	jobChannel   chan Job
	stopWorkerCh chan struct{}
}

func newWorker(pool *Pool) *worker {
	return &worker{
		ownerPool:    pool,
		jobChannel:   make(chan Job),
		stopWorkerCh: make(chan struct{}),
	}
}

func (w *worker) start() {
	go func() {
		for {
			w.ownerPool.workerPool <- w
			select {
			case job := <-w.jobChannel:
				if job == nil {
					fmt.Println("Params Error: the job is nil.")
				}
				err := job.handler()
				w.ownerPool.JobDone()
				if err != nil {
					log.Fatal("Exec Error: an error occurred during execute handle func. ", err)
					return
				}
			case <-w.stopWorkerCh:
				w.stopWorkerCh <- struct{}{}
				return
			}
		}
	}()
}

type Pool struct {
	jobQueue   chan Job     // job queue as Producer
	workerPool chan *worker // work queue as Consumer
	stopPoolCh chan struct{}
	wg         sync.WaitGroup // wait the job done
}

func NewPool(numWorkers int, jobQueueLen int) *Pool {
	jobQueue := make(chan Job, jobQueueLen)
	workerPool := make(chan *worker, numWorkers)

	pool := &Pool{
		jobQueue:   jobQueue,
		workerPool: workerPool,
		stopPoolCh: make(chan struct{}),
	}

	for i := 0; i < numWorkers; i++ {
		worker := newWorker(pool)
		worker.start()
	}

	go pool.dispatch()
	return pool
}

func (p *Pool) dispatch() {
	for {
		select {
		case job := <-p.jobQueue:
			p.WaitCount(1)
			worker := <-p.workerPool
			worker.jobChannel <- job
		case <-p.stopPoolCh:
			for i := 0; i < cap(p.workerPool); i++ {
				worker := <-p.workerPool
				worker.stopWorkerCh <- struct{}{}
				<-worker.stopWorkerCh
			}
			p.stopPoolCh <- struct{}{}
			return
		}
	}
}

func (p *Pool) JobDone() {
	p.wg.Done()
}

func (p *Pool) WaitCount(count int) {
	p.wg.Add(count)
}

func (p *Pool) WaitAll() {
	p.wg.Wait()
}

func (p *Pool) Release() {
	p.WaitAll()
	p.stopPoolCh <- struct{}{}
	<-p.stopPoolCh
}
