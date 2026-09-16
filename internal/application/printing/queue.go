package printing

import (
	"sync"

	domainErrors "sv-printer/internal/domain/errors"
	"sv-printer/internal/domain/job"
)

type PrintQueue interface {
	Enqueue(j *job.PrintJob) error
	Dequeue() *job.PrintJob
	GetJob(id string) (*job.PrintJob, error)
	UpdateStatus(id string, status job.PrintJobStatus) error
}

type InMemoryQueue struct {
	mu    sync.RWMutex
	jobs  map[string]*job.PrintJob
	queue []string
}

func NewInMemoryQueue() PrintQueue {
	return &InMemoryQueue{
		jobs:  make(map[string]*job.PrintJob),
		queue: make([]string, 0),
	}
}

func (q *InMemoryQueue) Enqueue(j *job.PrintJob) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, exists := q.jobs[j.ID]; exists {
		return domainErrors.ErrJobAlreadyExists
	}

	j.Status = job.StatusQueued
	q.jobs[j.ID] = j
	q.queue = append(q.queue, j.ID)

	return nil
}

func (q *InMemoryQueue) Dequeue() *job.PrintJob {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) == 0 {
		return nil
	}

	id := q.queue[0]
	q.queue = q.queue[1:]

	if j, ok := q.jobs[id]; ok {
		return j
	}

	return nil
}

func (q *InMemoryQueue) GetJob(id string) (*job.PrintJob, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	j, ok := q.jobs[id]
	if !ok {
		return nil, domainErrors.ErrJobNotFound
	}

	return j, nil
}

func (q *InMemoryQueue) UpdateStatus(id string, status job.PrintJobStatus) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	j, ok := q.jobs[id]
	if !ok {
		return domainErrors.ErrJobNotFound
	}

	j.Status = status
	return nil
}
