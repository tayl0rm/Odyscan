package scanner

import (
	"odyscan/config"
	"testing"
)

func TestPullImageFromArtifactRegistry(t *testing.T) {
	cfg := &config.Config{ImageName: "/ga-test-project-503ca/core/apline-edge@sha256:6062e4763b0aefcb2f5e29789efc188383cd1ee4c1ce3ff50d012ef260922a22"}
	err := PullImageFromArtifactRegistry(cfg)
	if err != nil {
		t.Errorf("PullImageFromArtifactRegistry failed: %v", err)
	}
}
