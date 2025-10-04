package workers

import (
	"context"
	"log/slog"
)

type Manager struct {
	jobQueue       chan func(context.Context) error
	workerPoolSize int
}

func New(workerCount, queueLen int) *Manager {
	return &Manager{
		workerPoolSize: workerCount,
		jobQueue:       make(chan func(context.Context) error, queueLen),
	}
}

func (m *Manager) StartPool() {
	for i := 0; i < m.workerPoolSize; i++ {
		go m.worker()
	}
}

func (m *Manager) worker() {
	for job := range m.jobQueue {
		if err := job(context.Background()); err != nil {
			slog.Error("Worker failed to execute job", slog.Any("error", err))
		}
	}
}

func (m *Manager) SetJob(job func(ctx context.Context) error) {
	select {
	case m.jobQueue <- job:
	default:
		slog.Error("Task rejected", slog.Any("reason", "Job queue full"))
	}
}

func (m *Manager) Stop() {
	close(m.jobQueue)
}

func Run(workerCount, queueLen int) *Manager {
	manager := New(workerCount, queueLen)
	go manager.StartPool()
	return manager
}
