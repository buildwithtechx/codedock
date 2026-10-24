package deployments

import (
	"testing"
)

func TestPlanBuildLimit(t *testing.T) {
	tests := []struct {
		tier     string
		expected int
	}{
		{"free", 1},
		{"hobby", 1},
		{"starter", 4},
		{"pro", 8},
		{"scale", 20},
		{"enterprise", 20},
	}

	for _, tt := range tests {
		t.Run(tt.tier, func(t *testing.T) {
			got := PlanBuildLimit(tt.tier)
			if got != tt.expected {
				t.Fatalf("expected limit %d for tier %s, got %d", tt.expected, tt.tier, got)
			}
		})
	}
}

func TestAssertBuildAdmission(t *testing.T) {
	t.Run("under limit succeeds", func(t *testing.T) {
		if err := AssertBuildAdmission(0, "free"); err != nil {
			t.Fatalf("expected admission to succeed, got %v", err)
		}
		if err := AssertBuildAdmission(3, "starter"); err != nil {
			t.Fatalf("expected admission to succeed, got %v", err)
		}
	})

	t.Run("at or over limit returns error", func(t *testing.T) {
		if err := AssertBuildAdmission(1, "free"); err != ErrBuildAdmissionLimitExceeded {
			t.Fatalf("expected ErrBuildAdmissionLimitExceeded, got %v", err)
		}
		if err := AssertBuildAdmission(8, "pro"); err != ErrBuildAdmissionLimitExceeded {
			t.Fatalf("expected ErrBuildAdmissionLimitExceeded, got %v", err)
		}
	})
}

func TestBuildAdmissionGuard(t *testing.T) {
	guard := NewBuildAdmissionGuard()
	scopeID := "project-123"

	if err := guard.Acquire(scopeID, "hobby"); err != nil {
		t.Fatalf("unexpected error on first acquire: %v", err)
	}
	if count := guard.ActiveCount(scopeID); count != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}

	if err := guard.Acquire(scopeID, "hobby"); err != ErrBuildAdmissionLimitExceeded {
		t.Fatalf("expected ErrBuildAdmissionLimitExceeded on contention, got %v", err)
	}

	guard.Release(scopeID)
	if count := guard.ActiveCount(scopeID); count != 0 {
		t.Fatalf("expected count 0 after release, got %d", count)
	}

	if err := guard.Acquire(scopeID, "hobby"); err != nil {
		t.Fatalf("unexpected error re-acquiring after release: %v", err)
	}
}
