package model

import (
	"strings"
	"time"
)

type Operation struct {
	ID            string    `json:"id"`
	Algorithm     string    `json:"algorithm"`
	Mode          string    `json:"mode"`
	InputHash     string    `json:"input_hash"`
	OutputPreview string    `json:"output_preview"`
	CreatedAt     time.Time `json:"created_at"`
}

func (o *Operation) Validate() error {
	o.Algorithm = strings.TrimSpace(o.Algorithm)
	o.Mode = strings.TrimSpace(o.Mode)
	if o.Algorithm == "" {
		return NewValidationError("algorithm", "算法名称不能为空")
	}
	if o.Mode != ModeEncode && o.Mode != ModeDecode {
		return NewValidationError("mode", "模式必须为 encode 或 decode")
	}
	return nil
}

type OperationFilter struct {
	Algorithm string
	Mode      string
	StartTime *time.Time
	EndTime   *time.Time
}

func (f OperationFilter) Match(o *Operation) bool {
	if f.Algorithm != "" && o.Algorithm != f.Algorithm {
		return false
	}
	if f.Mode != "" && o.Mode != f.Mode {
		return false
	}
	if f.StartTime != nil && o.CreatedAt.Before(*f.StartTime) {
		return false
	}
	if f.EndTime != nil && o.CreatedAt.After(*f.EndTime) {
		return false
	}
	return true
}
