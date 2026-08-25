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

// Build 构建编解码配置并写入缓存。
//
// 关键约束：只有在全部阶段构建成功后才把 profile 发布到缓存。构建过程中
// 出错（包括遇到非法阶段触发的 panic）会被 recover 转成普通 error 返回，
// 但绝不写入半成品。这样：
//   - 失败的配置不会残留在缓存里、不对外可见；
//   - 读取同名配置时不会被上次的半成品拖垮；
//   - 同名配置在失败后仍可被下一次正常构建覆盖使用。
func (b *ProfileBuilder) Build(name string, stages []string) (err error) {
	profile := &model.CodecProfile{Name: name, Stages: make(map[string]bool)}
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
	return b.cache.Publish(profile)
}

// GetReady 读取已就绪的配置。
//
// 未找到时返回 store.ErrNotFound；缓存里存在但未就绪（不应发生的脏状态）
// 时返回 store.ErrConflict。这里不再用 panic 表达“未就绪”，避免一次读取
// 把整个进程拖垮——把异常退化成 error，由上层决定如何响应。
func (b *ProfileBuilder) GetReady(name string) (*model.CodecProfile, error) {
	profile, err := b.cache.Get(name)
	if err != nil {
		return nil, err
	}
	if !profile.Ready {
		return nil, store.ErrConflict
	}
	return profile, nil
}
