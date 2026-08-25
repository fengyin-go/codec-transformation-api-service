package store

import (
	"sync"

	"codec/internal/model"
)

type ProfileCache struct {
	mu       sync.Mutex
	profiles map[string]*model.CodecProfile
}

func NewProfileCache() *ProfileCache {
	return &ProfileCache{profiles: make(map[string]*model.CodecProfile)}
}

func (c *ProfileCache) Publish(profile *model.CodecProfile) error {
	c.mu.Lock()
	c.profiles[profile.Name] = profile
	c.mu.Unlock()
	return nil
}

func (c *ProfileCache) Get(name string) (*model.CodecProfile, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	profile := c.profiles[name]
	if profile == nil {
		return nil, ErrNotFound
	}
	return profile, nil
}
