package security

import (
	"testing"
)

func TestManagedRuntimeImage(t *testing.T) {
	tests := []struct {
		name        string
		ref         string
		wantRuntime string
		wantMajor   int
		wantVariant string
		wantNil     bool
	}{
		{
			name:        "bun alpine",
			ref:         "oven/bun:1.3.1-alpine",
			wantRuntime: "bun",
			wantMajor:   1,
			wantVariant: "-alpine",
		},
		{
			name:        "node bookworm-slim",
			ref:         "node:22.22.0-bookworm-slim",
			wantRuntime: "node",
			wantMajor:   22,
			wantVariant: "-bookworm-slim",
		},
		{
			name:        "docker.io library node",
			ref:         "docker.io/library/node:20.10.0-alpine",
			wantRuntime: "node",
			wantMajor:   20,
			wantVariant: "-alpine",
		},
		{
			name:    "custom third party registry",
			ref:     "ghcr.io/myorg/bun:1.3.1",
			wantNil: true,
		},
		{
			name:    "canary tag without semver",
			ref:     "oven/bun:canary",
			wantNil: true,
		},
		{
			name:    "latest tag",
			ref:     "node:latest",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ManagedRuntimeImage(tt.ref)
			if tt.wantNil {
				if got != nil {
					t.Fatalf("expected nil for %s, got %+v", tt.ref, got)
				}
				return
			}
			if got == nil {
				t.Fatalf("expected non-nil for %s", tt.ref)
			}
			if got.Runtime != tt.wantRuntime {
				t.Errorf("runtime = %v, want %v", got.Runtime, tt.wantRuntime)
			}
			if got.Major != tt.wantMajor {
				t.Errorf("major = %v, want %v", got.Major, tt.wantMajor)
			}
			if got.Variant != tt.wantVariant {
				t.Errorf("variant = %v, want %v", got.Variant, tt.wantVariant)
			}
		})
	}
}

func TestSecurityRuntimeImages(t *testing.T) {
	t.Run("preserves alpine variant and updates version", func(t *testing.T) {
		patched, err := SecurityRuntimeImages(
			"bun",
			"1.3.1",
			">=1.4.4",
			"oven/bun:1.3.1-alpine",
			"oven/bun:1.3.1-alpine",
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if patched.TargetVersion != "1.4.4" {
			t.Errorf("targetVersion = %s, want 1.4.4", patched.TargetVersion)
		}
		if patched.BuildImage != "oven/bun:1.4.4-alpine" {
			t.Errorf("buildImage = %s, want oven/bun:1.4.4-alpine", patched.BuildImage)
		}
		if patched.RuntimeImage != "oven/bun:1.4.4-alpine" {
			t.Errorf("runtimeImage = %s, want oven/bun:1.4.4-alpine", patched.RuntimeImage)
		}
	})

	t.Run("selects candidate in same major on disjoint ranges", func(t *testing.T) {
		patched, err := SecurityRuntimeImages(
			"node",
			"22.22.0",
			">=22.23.4 <23 || >=24.22.1 <25",
			"node:22.22.0-bookworm-slim",
			"node:22.22.0-bookworm-slim",
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if patched.TargetVersion != "22.23.4" {
			t.Errorf("targetVersion = %s, want 22.23.4", patched.TargetVersion)
		}
		if patched.BuildImage != "node:22.23.4-bookworm-slim" {
			t.Errorf("buildImage = %s, want node:22.23.4-bookworm-slim", patched.BuildImage)
		}
	})

	t.Run("does not downgrade if floor already satisfies advisory", func(t *testing.T) {
		patched, err := SecurityRuntimeImages(
			"bun",
			"1.5.0",
			">=1.4.4",
			"oven/bun:1.5.0-alpine",
			"oven/bun:1.5.0-alpine",
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if patched.TargetVersion != "1.5.0" {
			t.Errorf("targetVersion = %s, want 1.5.0", patched.TargetVersion)
		}
		if patched.BuildImage != "oven/bun:1.5.0-alpine" {
			t.Errorf("buildImage = %s, want oven/bun:1.5.0-alpine", patched.BuildImage)
		}
	})

	t.Run("refuses automatic patch across major version boundary", func(t *testing.T) {
		_, err := SecurityRuntimeImages(
			"bun",
			"1.3.1",
			">=2.0.0",
			"oven/bun:1.3.1-alpine",
			"oven/bun:1.3.1-alpine",
		)
		if err == nil {
			t.Fatal("expected error for major version boundary jump")
		}
	})

	t.Run("refuses non-managed custom images", func(t *testing.T) {
		_, err := SecurityRuntimeImages(
			"node",
			"22.22.0",
			">=22.23.4 <23",
			"ghcr.io/custom/node:22.22.0",
			"ghcr.io/custom/node:22.22.0",
		)
		if err == nil {
			t.Fatal("expected error for custom docker image")
		}
	})
}

func TestPlanSecurityRuntimePatchAndValidateReviewID(t *testing.T) {
	patch, err := PlanSecurityRuntimePatch(
		"proj-123",
		"dep-456",
		"ADV-001",
		1,
		"bun",
		"1.3.1",
		">=1.4.4",
		"oven/bun:1.3.1-alpine",
		"oven/bun:1.3.1-alpine",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(patch.ReviewID) != 64 {
		t.Fatalf("reviewId length = %d, want 64 hex chars", len(patch.ReviewID))
	}

	if !ValidateReviewID(patch, "proj-123", "dep-456", "ADV-001", 1) {
		t.Fatal("ValidateReviewID failed for valid patch")
	}

	if ValidateReviewID(patch, "proj-diff", "dep-456", "ADV-001", 1) {
		t.Fatal("ValidateReviewID should fail for changed projectID")
	}

	if ValidateReviewID(patch, "proj-123", "dep-456", "ADV-001", 2) {
		t.Fatal("ValidateReviewID should fail for changed revision")
	}
}
