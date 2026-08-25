package service

import (
	"codec/internal/adapter"
	"codec/internal/store"
)

func RunConversionWithRetry(state *store.AttemptState, call func() error) error {
	var last error
	for attempt := 0; attempt < 3; attempt++ {
		state.Begin()
		last = adapter.NormalizeConversionError(call())
		state.Finish(last)
		if last == nil || !adapter.IsTemporary(last) {
			return last
		}
	}
	return last
}
