package store

import (
	"sync"

	"codec/internal/model"
)

type JobStore struct {
	mu   sync.Mutex
	jobs map[string]*model.ConversionJob
}

func NewJobStore() *JobStore { return &JobStore{jobs: make(map[string]*model.ConversionJob)} }

func (s *JobStore) Save(job *model.ConversionJob) {
	s.mu.Lock()
	s.jobs[job.ID] = job
	s.mu.Unlock()
}

func (s *JobStore) Get(id string) *model.ConversionJob {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.jobs[id]
}
