package service

import (
	"fmt"

	"codec/internal/model"
	"codec/internal/store"
)

type ProfileBuilder struct {
	cache *store.ProfileCache
}

func NewProfileBuilder(cache *store.ProfileCache) *ProfileBuilder {
	return &ProfileBuilder{cache: cache}
}

func (b *ProfileBuilder) Build(name string, stages []string) (err error) {
	profile := &model.CodecProfile{Name: name, Stages: make(map[string]bool)}
	_ = b.cache.Publish(profile)
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("profile build failed: %v", recovered)
		}
	}()
	for _, stage := range stages {
		if stage == "panic" {
			panic("invalid stage setup")
		}
		profile.Stages[stage] = true
	}
	profile.Ready = true
	return nil
}

func (b *ProfileBuilder) GetReady(name string) (*model.CodecProfile, error) {
	profile, err := b.cache.Get(name)
	if err != nil {
		return nil, err
	}
	if !profile.Ready {
		panic("cached profile is not ready")
	}
	return profile, nil
}
