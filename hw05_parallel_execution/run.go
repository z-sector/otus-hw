package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrErrorsCountGor      = errors.New("errors count of goroutines")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m < 1 {
		return ErrErrorsLimitExceeded
	}
	if n < 1 {
		return ErrErrorsCountGor
	}

	taskCh := make(chan Task)
	wg := &sync.WaitGroup{}
	errC := newErrCounter(m)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go process(taskCh, errC, wg)
	}

	for _, task := range tasks {
		if errC.isExceeded() {
			break
		}
		taskCh <- task
	}
	close(taskCh)

	wg.Wait()

	if errC.count.Load() > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

type errCounter struct {
	limit int32
	count atomic.Int32
}

func (e *errCounter) inc() {
	e.count.Add(1)
}

func (e *errCounter) isExceeded() bool {
	return e.count.Load() >= e.limit
}

func newErrCounter(limit int) *errCounter {
	return &errCounter{
		limit: int32(limit),
		count: atomic.Int32{},
	}
}

func process(ch <-chan Task, errC *errCounter, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range ch {
		if err := task(); err != nil {
			errC.inc()
		}
	}
}
