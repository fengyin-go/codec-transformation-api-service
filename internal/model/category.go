package model

import (
	"strings"
	"time"
)

type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *Category) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	c.Description = strings.TrimSpace(c.Description)
	if c.Name == "" {
		return NewValidationError("name", "分类名称不能为空")
	}
	return nil
}
