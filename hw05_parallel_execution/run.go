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
	tasksLen := len(tasks)

	for i := 0; i < n; i++ {
		go func() {
			for j := i; j < tasksLen && atomic.LoadInt64(&errorsRemain) > 0; j += n {
				err := tasks[j]()
				if err != nil {
					atomic.AddInt64(&errorsRemain, -1)
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()
	if errorsRemain <= 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
