package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	//nolint:depguard
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("keep calm and make math", func(t *testing.T) {
		tasksCount := 1000
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				answerQuestionOfLife := 0
				for i := range rand.Intn(10_000_000) {
					answerQuestionOfLife += i
				}

				atomic.AddInt32(&runTasksCount, 1)

				if answerQuestionOfLife%2 == 1 {
					return fmt.Errorf("wrong answer on life question %d in gorutine %d", answerQuestionOfLife, i)
				}

				return nil
			})
		}

		workersCount := 10
		maxErrorsCount := 100

		_ = Run(tasks, workersCount, maxErrorsCount)
		require.GreaterOrEqual(t, runTasksCount, int32(maxErrorsCount), "tasks started less then errors received")
	})

	t.Run("tasks without errors with Eventually", func(t *testing.T) {
		workersCount := 5
		maxErrorsCount := 1

		syncChan := make(chan struct{}, 5)
		defer close(syncChan)

		tasksCount := 45
		tasks := make([]Task, 0, tasksCount)

		once := sync.Once{}
		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				<-syncChan
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		require.Eventually(t, func() bool {
			once.Do(func() {
				go func() {
					_ = Run(tasks, workersCount, maxErrorsCount)
				}()
			})

			if int32(tasksCount) == atomic.LoadInt32(&runTasksCount) {
				return true
			}

			for range workersCount {
				syncChan <- struct{}{}
			}

			return false
		}, time.Second+time.Millisecond*50, time.Millisecond*100, "tasks were run sequentially?")
	})
}
