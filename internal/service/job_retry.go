package service

import (
	"sync"

	"codec/internal/model"
	"codec/internal/store"
)

type JobRunner struct {
	store   *store.JobStore
	mu      sync.Mutex
	effects map[string]struct{}
}

func NewJobRunner(jobStore *store.JobStore) *JobRunner {
	return &JobRunner{store: jobStore, effects: make(map[string]struct{})}
}

func (r *JobRunner) effect(key string) {
	r.mu.Lock()
	r.effects[key] = struct{}{}
	r.mu.Unlock()
}

func (r *JobRunner) StartRetry(releaseOld <-chan struct{}) (string, <-chan struct{}) {
	id := "conversion-1"
	r.store.Save(&model.ConversionJob{ID: id, Version: 1, Status: "running"})
	done := make(chan struct{})
	go func() {
		<-releaseOld
		r.effect(id + "-attempt-1")
		r.store.Save(&model.ConversionJob{ID: id, Version: 1, Status: "running"})
		close(done)
	}()
	r.effect(id + "-attempt-2")
	r.store.Save(&model.ConversionJob{ID: id, Version: 2, Status: "succeeded"})
	return id, done
}

func (r *JobRunner) EffectCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.effects)
}
