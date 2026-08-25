package service

import (
	"context"
	"testing"
	"time"
)

func TestParallelBatchReturnsValidOutputAndPerItemError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	results, err := ParallelBatch(ctx, []string{"alpha", ""})
	if err != nil {
		t.Fatalf("mixed batch did not finish: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("mixed batch result count = %d, want 2", len(results))
	}
	if results[0].Output != "encoded:alpha" || results[0].Error != "" {
		t.Fatalf("valid item result = %+v", results[0])
	}
	if results[1].Error != "input is empty" {
		t.Fatalf("invalid item result = %+v", results[1])
	}
}
