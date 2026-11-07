package security

import (
	"testing"
)

func TestParseAndCompareSemver(t *testing.T) {
	v1, err := ParseSemver("22.23.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v1.Major != 22 || v1.Minor != 23 || v1.Patch != 3 {
		t.Fatalf("unexpected parse: %+v", v1)
	}

	v2, err := ParseSemver("v22.23.4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if CompareSemver(v1, v2) >= 0 {
		t.Fatalf("expected %s < %s", v1, v2)
	}
	if CompareSemver(v2, v1) <= 0 {
		t.Fatalf("expected %s > %s", v2, v1)
	}
	if CompareSemver(v1, v1) != 0 {
		t.Fatalf("expected %s == %s", v1, v1)
	}
}

func TestMatchesRange(t *testing.T) {
	// Node advisory range: "<22.23.4 || >=23.0.0 <24.22.1 || >=25.0.0 <26.12.1"
	nodeRange := "<22.23.4 || >=23.0.0 <24.22.1 || >=25.0.0 <26.12.1"

	t.Run("matches vulnerable lower version", func(t *testing.T) {
		if !MatchesRange("22.23.3", nodeRange) {
			t.Fatal("expected 22.23.3 to match vulnerable range")
		}
		if !MatchesRange("20.18.0", nodeRange) {
			t.Fatal("expected 20.18.0 to match vulnerable range")
		}
	})

	t.Run("does not match patched 22 release", func(t *testing.T) {
		if MatchesRange("22.23.4", nodeRange) {
			t.Fatal("expected 22.23.4 NOT to match vulnerable range")
		}
		if MatchesRange("22.24.0", nodeRange) {
			t.Fatal("expected 22.24.0 NOT to match vulnerable range")
		}
	})

	t.Run("matches vulnerable 23 release", func(t *testing.T) {
		if !MatchesRange("23.0.0", nodeRange) {
			t.Fatal("expected 23.0.0 to match vulnerable range")
		}
		if !MatchesRange("24.20.0", nodeRange) {
			t.Fatal("expected 24.20.0 to match vulnerable range")
		}
	})

	t.Run("does not match patched 24 release", func(t *testing.T) {
		if MatchesRange("24.22.1", nodeRange) {
			t.Fatal("expected 24.22.1 NOT to match vulnerable range")
		}
	})

	// Bun advisory range: "<1.4.4"
	bunRange := "<1.4.4"
	t.Run("bun version check", func(t *testing.T) {
		if !MatchesRange("1.4.3", bunRange) {
			t.Fatal("expected 1.4.3 to match")
		}
		if MatchesRange("1.4.4", bunRange) {
			t.Fatal("expected 1.4.4 NOT to match")
		}
		if MatchesRange("1.5.0", bunRange) {
			t.Fatal("expected 1.5.0 NOT to match")
		}
	})
}
