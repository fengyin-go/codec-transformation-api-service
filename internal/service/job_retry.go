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

// effect records that the conversion side effect for key fired. It is
// idempotent: a late callback whose Save was dropped (because a newer attempt
// already committed) does not fire the effect, and an attempt whose Save
// landed the terminal status only fires once even if reached by two paths.
func (r *JobRunner) effect(key string) {
	r.mu.Lock()
	r.effects[key] = struct{}{}
	r.mu.Unlock()
}

// StartRetry models a conversion that is retried: the first attempt's result
// callback may arrive late, after the retry has already committed the final
// state. releaseOld signals that the first attempt's callback is now
// available. The retry must remain authoritative: a late, stale callback can
// neither regress the committed state nor re-run the conversion effect.
func (r *JobRunner) StartRetry(releaseOld <-chan struct{}) (string, <-chan struct{}) {
	id := "conversion-1"
	r.store.Save(&model.ConversionJob{ID: id, Version: 1, Status: "running"})

	done := make(chan struct{})
	go func() {
		defer close(done)
		<-releaseOld
		// Replay the first attempt's outcome. The retry has already moved the
		// job forward to Version 2, so this stale Save is dropped by the store
		// and the conversion effect must not fire again.
		if r.store.Save(&model.ConversionJob{ID: id, Version: 1, Status: "running"}) {
			r.effect(id)
		}
	}()

	// Retry attempt commits the terminal state first.
	if r.store.Save(&model.ConversionJob{ID: id, Version: 2, Status: "succeeded"}) {
		r.effect(id)
	}
	return id, done
}

func (r *JobRunner) EffectCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.effects)
}
