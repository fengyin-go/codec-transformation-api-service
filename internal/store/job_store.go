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

// Save commits job unless a newer version is already stored for the same
// job, in which case the stale write is dropped and false is returned. This
// keeps a late callback from regressing a state already committed by a newer
// attempt.
func (s *JobStore) Save(job *model.ConversionJob) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cur := s.jobs[job.ID]; cur != nil && job.Version < cur.Version {
		return false
	}
	s.jobs[job.ID] = job
	return true
}

func (s *JobStore) Get(id string) *model.ConversionJob {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.jobs[id]
}
