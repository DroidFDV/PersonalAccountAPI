package workers

import (
	"sync"
)

type Manager struct {
	//TODO: заменить interface{} на что-то более разумное и определенное
	jobQueue   chan interface{}
	workerPool chan struct{}
	wg         sync.WaitGroup
}

func New(workerCount, jobQueueSize uint) *Manager {
	return &Manager{
		workerPool: make(chan struct{}, workerCount),
		jobQueue:   make(chan interface{}, jobQueueSize),
	}
}

func (m *Manager) StartPool() {

	for i := 0; i < cap(m.workerPool); i++ {
		m.workerPool <- struct{}{}
		m.wg.Add(1)
		go m.worker()
	}
}

func (m *Manager) Stop() {
	close(m.jobQueue)
	m.wg.Wait()
	close(m.workerPool)
}

func (m *Manager) AddJob(job interface{}) {
	m.jobQueue <- job
}

func (m *Manager) worker() {
	defer m.wg.Done()
	for job := range m.jobQueue {
		job.Work()
	}
}
