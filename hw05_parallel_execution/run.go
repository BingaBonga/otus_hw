package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	wg.Add(n)

	errorsRemain := int64(m)
	tasksLen := int64(len(tasks))
	currentIndexTask := int64(-1)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for {
				j := atomic.AddInt64(&currentIndexTask, 1)
				if j >= tasksLen || atomic.LoadInt64(&errorsRemain) <= 0 {
					return
				}

				if err := tasks[j](); err != nil {
					atomic.AddInt64(&errorsRemain, -1)
				}
			}
		}()
	}

	wg.Wait()
	if errorsRemain <= 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
