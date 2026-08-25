package service

import (
	"testing"

	"codec/internal/store"
)

func TestFailedProfileBuildDoesNotPublishHalfInitializedProfile(t *testing.T) {
	builder := NewProfileBuilder(store.NewProfileCache())
	if err := builder.Build("broken", []string{"base64", "panic"}); err == nil {
		t.Fatal("broken profile build returned no error")
	}
	var brokenProfile interface{}
	var brokenErr error
	var panicValue interface{}
	func() {
		defer func() { panicValue = recover() }()
		brokenProfile, brokenErr = builder.GetReady("broken")
	}()
	if panicValue != nil || brokenErr == nil {
		t.Fatalf("half-initialized profile remained visible or panicked: profile=%+v panic=%v", brokenProfile, panicValue)
	}
	if err := builder.Build("working", []string{"base64", "hex"}); err != nil {
		t.Fatalf("next profile build failed: %v", err)
	}
	profile, err := builder.GetReady("working")
	if err != nil || !profile.Ready || !profile.Stages["base64"] || !profile.Stages["hex"] {
		t.Fatalf("working profile = %+v err=%v", profile, err)
	}
}
