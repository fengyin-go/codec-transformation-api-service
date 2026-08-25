// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"codec/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Algorithm
	CreateAlgorithm(a *model.Algorithm) error
	GetAlgorithm(id string) (*model.Algorithm, error)
	GetAlgorithmByName(name string) (*model.Algorithm, error)
	ListAlgorithms() []*model.Algorithm
	UpdateAlgorithm(a *model.Algorithm) error
	DeleteAlgorithm(id string) error

	// Operation
	CreateOperation(o *model.Operation) error
	GetOperation(id string) (*model.Operation, error)
	ListOperations() []*model.Operation
	DeleteOperation(id string) error

	// Category
	CreateCategory(c *model.Category) error
	GetCategory(id string) (*model.Category, error)
	GetCategoryByName(name string) (*model.Category, error)
	ListCategories() []*model.Category
	UpdateCategory(c *model.Category) error
	DeleteCategory(id string) error
}
