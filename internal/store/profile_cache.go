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

// Publish 写入一条配置。缓存层本身不做 Ready 校验，是否在就绪后才发布
// 由调用方（ProfileBuilder）负责。Publish 直接覆盖同名旧值，因此即便上次
// 构建失败留下了残留（正常路径不会发生），下一次成功构建也能干净覆盖。
func (c *ProfileCache) Publish(profile *model.CodecProfile) error {
	c.mu.Lock()
	c.profiles[profile.Name] = profile
	c.mu.Unlock()
	return nil
}

// Get 读取一条配置。防御性约束：未就绪的 profile 一律视作不存在，
// 绝不把半成品泄漏给调用方。这样即使有脏状态残留在缓存里，读取侧也只会
// 拿到“记录不存在”，而不是一个没构建完成的配置。
func (c *ProfileCache) Get(name string) (*model.CodecProfile, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	profile := c.profiles[name]
	if profile == nil || !profile.Ready {
		return nil, ErrNotFound
	}
	return profile, nil
}
