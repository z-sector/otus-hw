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
	errC := &atomic.Int32{}
	mInt32 := int32(m)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go process(taskCh, errC, wg)
	}

	for _, task := range tasks {
		if errC.Load() >= mInt32 {
			break
		}
		taskCh <- task
	}
	close(taskCh)

	wg.Wait()

	if errC.Load() > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func process(ch <-chan Task, errC *atomic.Int32, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range ch {
		if err := task(); err != nil {
			errC.Add(1)
		}
	}
}
