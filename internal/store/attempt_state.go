package store

import "sync"

type AttemptState struct {
	mu        sync.Mutex
	attempts  int
	committed int
}

func (s *AttemptState) Begin() { s.mu.Lock(); s.attempts++; s.mu.Unlock() }

func (s *AttemptState) Finish(err error) {
	s.mu.Lock()
	s.committed++
	s.mu.Unlock()
}

func (s *AttemptState) Counts() (int, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.attempts, s.committed
}
