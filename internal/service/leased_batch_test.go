package service

import (
	"strings"
	"testing"

	"codec/internal/store"
)

func TestBatchReleasesEachConversionLeaseBeforeNextInput(t *testing.T) {
	pool := store.NewLeasePool(2)
	outputs, err := ConvertBatchWithLeases(pool, []string{"one", "two", "three"}, func(input string) (string, error) {
		return strings.ToUpper(input), nil
	})
	if err != nil {
		t.Fatalf("three-item batch failed: %v", err)
	}
	if got := strings.Join(outputs, ","); got != "ONE,TWO,THREE" {
		t.Fatalf("batch outputs = %q", got)
	}
	active, committed := pool.Counts()
	if active != 0 || committed != 3 {
		t.Fatalf("lease counts active=%d committed=%d, want 0/3", active, committed)
	}
}
