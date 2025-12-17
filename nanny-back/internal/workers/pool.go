package workers

import (
	"log"
)

type Job func()

type WorkerPool struct {
	jobs chan Job
}

func NewWorkerPool(workerCount int) *WorkerPool {
	pool := &WorkerPool{
		jobs: make(chan Job),
	}

	for i := 0; i < workerCount; i++ {
		go pool.worker(i)
	}

	return pool
}

func (p *WorkerPool) worker(id int) {
	for job := range p.jobs {
		log.Printf("Worker %d started job", id)
		job()
		log.Printf("Worker %d finished job", id)
	}
}

func (p *WorkerPool) Submit(job Job) {
	p.jobs <- job
}

func (p *WorkerPool) Shutdown() {
	close(p.jobs)
}
