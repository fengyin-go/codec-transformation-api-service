package service

import (
	"testing"

	"codec/internal/store"
)

func TestSuccessfulRetryRejectsLateCallbackAndKeepsOneEffect(t *testing.T) {
	jobStore := store.NewJobStore()
	runner := NewJobRunner(jobStore)
	releaseOld := make(chan struct{})
	id, done := runner.StartRetry(releaseOld)
	close(releaseOld)
	<-done

	job := jobStore.Get(id)
	if job == nil || job.Version != 2 || job.Status != "succeeded" {
		t.Fatalf("job regressed after late callback: %+v", job)
	}
	if got := runner.EffectCount(); got != 1 {
		t.Fatalf("conversion side effect count = %d, want 1", got)
	}
}
