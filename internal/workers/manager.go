package workers

import (
	"context"
)

var jobQueueLen int = 1000

type Manager struct {
	jobQueue       chan func(context.Context) error
	logChan        chan error
	workerPoolSize int
}

func New(workerCount, queueLen int) *Manager {
	return &Manager{
		workerPoolSize: workerCount,
		jobQueue:       make(chan func(context.Context) error, queueLen),
		logChan:        make(chan error),
	}
}

func (m *Manager) StartPool() {
	for job := range m.jobQueue {
		go m.worker(job)
	}
}

func (m *Manager) Stop() {
	close(m.jobQueue)
	close(m.logChan)
}

func (m *Manager) SetJob(job func(ctx context.Context) error) {
	m.jobQueue <- job
}

func (m *Manager) worker(job func(ctx context.Context) error) {
	m.logChan <- job(context.TODO())
}

func (m *Manager) GetLog() error {
	return <-m.logChan
}

func Run(workerCount int) *Manager {
	manager := New(workerCount, jobQueueLen)
	go manager.StartPool()
	return manager
}
