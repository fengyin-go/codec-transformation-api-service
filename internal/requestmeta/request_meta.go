package requestmeta

import (
	"sync"

	"codec/internal/model"
)

type Recorder interface {
	RecordLater(*model.RequestMeta, <-chan struct{})
}

type RequestMetaPool struct {
	mu   sync.Mutex
	idle *model.RequestMeta
}

func (p *RequestMetaPool) Get() *model.RequestMeta {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.idle == nil {
		return &model.RequestMeta{}
	}
	m := p.idle
	p.idle = nil
	return m
}

func (p *RequestMetaPool) Put(m *model.RequestMeta) {
	p.mu.Lock()
	p.idle = m
	p.mu.Unlock()
}

func CaptureRequestTrace(p *RequestMetaPool, recorder Recorder, client string, payload []byte, release <-chan struct{}) <-chan struct{} {
	meta := p.Get()
	meta.ClientID = client
	meta.Payload = payload
	done := make(chan struct{})
	go func() {
		recorder.RecordLater(meta, release)
		close(done)
	}()
	p.Put(meta)
	return done
}
