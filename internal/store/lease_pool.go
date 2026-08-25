package store

import (
	"errors"
	"sync"
)

var ErrLeaseExhausted = errors.New("conversion lease exhausted")

type LeasePool struct {
	mu        sync.Mutex
	capacity  int
	active    int
	committed int
}

type Lease struct {
	pool      *LeasePool
	committed bool
	closed    bool
}

func NewLeasePool(capacity int) *LeasePool { return &LeasePool{capacity: capacity} }

func (p *LeasePool) Acquire() (*Lease, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.active >= p.capacity {
		return nil, ErrLeaseExhausted
	}
	p.active++
	return &Lease{pool: p}, nil
}

func (l *Lease) Commit() {
	if l.closed || l.committed {
		return
	}
	l.pool.mu.Lock()
	l.pool.committed++
	l.pool.mu.Unlock()
	l.committed = true
}

func (l *Lease) Close() {
	if l.closed {
		return
	}
	l.pool.mu.Lock()
	l.pool.active--
	l.pool.mu.Unlock()
	l.closed = true
}

func (p *LeasePool) Counts() (int, int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.active, p.committed
}
