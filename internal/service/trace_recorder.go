package service

import (
	"sync"

	"codec/internal/model"
)

type RequestTrace struct {
	ClientID string
	Payload  string
}

type TraceRecorder struct {
	mu     sync.Mutex
	traces []RequestTrace
}

func NewTraceRecorder() *TraceRecorder { return &TraceRecorder{} }

func (r *TraceRecorder) RecordLater(meta *model.RequestMeta, release <-chan struct{}) {
	<-release
	r.mu.Lock()
	r.traces = append(r.traces, RequestTrace{ClientID: meta.ClientID, Payload: string(meta.Payload)})
	r.mu.Unlock()
}

func (r *TraceRecorder) Traces() []RequestTrace {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]RequestTrace(nil), r.traces...)
}
