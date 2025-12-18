package workers

import (
	"sync"
	"testing"
	"time"
)

func TestNewWorkerPool(t *testing.T) {
	pool := NewWorkerPool(3)
	if pool == nil {
		t.Error("expected worker pool to be created")
	}
	if pool.jobs == nil {
		t.Error("expected jobs channel to be initialized")
	}
	pool.Shutdown()
}

func TestWorkerPool_Submit(t *testing.T) {
	pool := NewWorkerPool(2)
	defer pool.Shutdown()

	var executed bool
	var mu sync.Mutex

	job := func() {
		mu.Lock()
		executed = true
		mu.Unlock()
	}

	pool.Submit(job)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if !executed {
		t.Error("expected job to be executed")
	}
	mu.Unlock()
}

func TestWorkerPool_MultipleJobs(t *testing.T) {
	pool := NewWorkerPool(3)
	defer pool.Shutdown()

	var counter int
	var mu sync.Mutex

	jobCount := 10

	for i := 0; i < jobCount; i++ {
		pool.Submit(func() {
			mu.Lock()
			counter++
			mu.Unlock()
		})
	}

	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	if counter != jobCount {
		t.Errorf("expected %d jobs executed, got %d", jobCount, counter)
	}
	mu.Unlock()
}

func TestWorkerPool_Shutdown(t *testing.T) {
	pool := NewWorkerPool(2)

	var counter int
	var mu sync.Mutex

	for i := 0; i < 5; i++ {
		pool.Submit(func() {
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			counter++
			mu.Unlock()
		})
	}

	time.Sleep(50 * time.Millisecond)

	pool.Shutdown()

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	finalCount := counter
	mu.Unlock()

	if finalCount == 0 {
		t.Error("expected some jobs to be executed before shutdown")
	}
}

func TestWorkerPool_ConcurrentSubmission(t *testing.T) {
	pool := NewWorkerPool(5)
	defer pool.Shutdown()

	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	goroutines := 10
	jobsPerGoroutine := 5

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < jobsPerGoroutine; j++ {
				pool.Submit(func() {
					mu.Lock()
					counter++
					mu.Unlock()
				})
			}
		}()
	}

	wg.Wait()

	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	expectedCount := goroutines * jobsPerGoroutine
	if counter != expectedCount {
		t.Errorf("expected %d jobs, got %d", expectedCount, counter)
	}
	mu.Unlock()
}

func TestWorkerPool_JobOrder(t *testing.T) {
	pool := NewWorkerPool(1)
	defer pool.Shutdown()

	var results []int
	var mu sync.Mutex

	for i := 0; i < 5; i++ {
		value := i
		pool.Submit(func() {
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			results = append(results, value)
			mu.Unlock()
		})
	}

	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	if len(results) != 5 {
		t.Errorf("expected 5 results, got %d", len(results))
	}
	mu.Unlock()
}
