package tickerpool

import (
	"fmt"
	"sync/atomic"
	"time"

	"golang.org/x/sync/syncmap"
)

// TickerPool represents a pool of workers set to perform a task periodically. A TickerPool has a
// set of workers which can grow and shrink dynamically.
type TickerPool struct {
	Interval       time.Duration
	workers        syncmap.Map
	workerTotal    int64
	workerQueue    chan func()
	workerInterval time.Duration
	generate       chan bool
}

// NewTickerPool creates a TickerPool with a set interval.
func NewTickerPool(interval time.Duration) (*TickerPool, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("interval cannot be zero or negative")
	}

	tp := &TickerPool{
		Interval:    interval,
		workerQueue: make(chan func()),
	}

	return tp, nil
}

// Add schedules a new worker to spin up on the TickerPool's interval. The task will not fire right
// away but as soon as the new worker interval is calculated and the task is assigned a position.
func (tp *TickerPool) Add(name string, task func()) (exists bool) {
	_, exists = tp.workers.LoadOrStore(name, task)
	if !exists {
		atomic.AddInt64(&tp.workerTotal, 1)

		if atomic.LoadInt64(&tp.workerTotal) == 1 {
			// if this is the first worker in the pool, the queue isn't running to start it
			go tp.queue()
		}
	}
	return
}

// Remove simply removes a worker from the pool, no other action is required as if there is a queue
// running, it will stop itself if there are no more workers.
func (tp *TickerPool) Remove(name string) {
	tp.workers.Delete(name)
	atomic.AddInt64(&tp.workerTotal, -1)
}

// Exists checks if a worker exists in the pool
func (tp *TickerPool) Exists(name string) bool {
	_, exists := tp.workers.Load(name)
	return exists
}

// calculateWorkerInterval is a simple function to divide an interval by the amount of workers
func calculateWorkerInterval(interval time.Duration, workers int64) time.Duration {
	return time.Duration(interval.Nanoseconds() / int64(workers))
}

// queue's job is to collect a list of workers to run in this period. It will first consume
// the newWorkers list and add them to the workers list. It will then fire off a loop through the
// workers list pushing each one down
func (tp *TickerPool) queue() {
	workers := atomic.LoadInt64(&tp.workerTotal)
	if workers == 0 {
		return
	}

	// calculate the new interval, this is the pool's interval dividied by the amount of workers
	tp.workerInterval = calculateWorkerInterval(tp.Interval, workers)

	// this isn't a time.Ticker because we want to iterate through the workers
	tp.workers.Range(func(key, value interface{}) bool {
		worker := value.(func())

		go worker()
		time.Sleep(tp.workerInterval)
		return true
	})

	tp.queue()
}
