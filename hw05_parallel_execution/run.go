package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Error struct {
	current int
	limit   int
}
type manager struct {
	wg    sync.WaitGroup
	mutex sync.RWMutex
	queue chan Task
	error Error
}

func (m *manager) errorIncrement() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.error.current++
}

func (m *manager) pushToQueue(task Task) {
	m.queue <- task
}

func (m *manager) closeQueue() {
	close(m.queue)
}

func newManager(errorLimit int) manager {
	return manager{
		error: Error{
			limit:   errorLimit,
			current: 0,
		},
		queue: make(chan Task),
	}
}

func (m *manager) isLimitExceeded() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.error.limit > 0 && m.error.current >= m.error.limit
}

func runTask(manager *manager) {
	defer manager.wg.Done()

	for task := range manager.queue {
		err := task()
		if err != nil {
			manager.errorIncrement()
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	manager := newManager(m)
	manager.wg.Add(n)
	defer manager.wg.Wait()
	defer manager.closeQueue()

	for i := 0; i < n; i++ {
		go runTask(&manager)
	}

	for _, task := range tasks {
		if manager.isLimitExceeded() {
			break
		}
		manager.pushToQueue(task)
	}

	if manager.isLimitExceeded() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
