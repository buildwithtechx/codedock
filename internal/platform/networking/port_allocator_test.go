package networking

import (
	"testing"
)

func TestParseEphemeralPortRange(t *testing.T) {
	tests := []struct {
		input    string
		expected PortRange
	}{
		{"32768 60999\n", PortRange{Start: 32768, End: 60999}},
		{"40000\t50000", PortRange{Start: 40000, End: 50000}},
		{"invalid", DefaultLinuxEphemeralRange},
		{"", DefaultLinuxEphemeralRange},
		{"-5 100", DefaultLinuxEphemeralRange},
		{"70000 80000", DefaultLinuxEphemeralRange},
	}

	for _, tc := range tests {
		got := ParseEphemeralPortRange(tc.input)
		if got != tc.expected {
			t.Errorf("ParseEphemeralPortRange(%q) = %+v, want %+v", tc.input, got, tc.expected)
		}
	}
}

func TestPickHostPort(t *testing.T) {
	ephemeral := &PortRange{Start: 32768, End: 60999}

	t.Run("reuses preferred port when free", func(t *testing.T) {
		port, err := PickHostPort(map[int]bool{}, AllocateHostPortOptions{
			Preferred:      20500,
			EphemeralRange: ephemeral,
		})
		if err != nil || port != 20500 {
			t.Fatalf("expected 20500, got %d (err: %v)", port, err)
		}
	})

	t.Run("skips preferred port when occupied and falls back to range", func(t *testing.T) {
		port, err := PickHostPort(map[int]bool{20500: true}, AllocateHostPortOptions{
			Preferred:      20500,
			RangeStart:     20000,
			RangeEnd:       20999,
			EphemeralRange: ephemeral,
		})
		if err != nil || port != 20000 {
			t.Fatalf("expected 20000, got %d (err: %v)", port, err)
		}
	})

	t.Run("reuses occupied preferred port when ownership proved", func(t *testing.T) {
		port, err := PickHostPort(map[int]bool{20500: true}, AllocateHostPortOptions{
			Preferred:              20500,
			ReuseOccupiedPreferred: true,
			EphemeralRange:         ephemeral,
		})
		if err != nil || port != 20500 {
			t.Fatalf("expected 20500, got %d (err: %v)", port, err)
		}
	})

	t.Run("does not let proven ownership override avoid claim", func(t *testing.T) {
		port, err := PickHostPort(map[int]bool{20500: true}, AllocateHostPortOptions{
			Preferred:              20500,
			ReuseOccupiedPreferred: true,
			Avoid:                  []int{20500},
			RangeStart:             20000,
			RangeEnd:               20999,
			EphemeralRange:         ephemeral,
		})
		if err != nil || port != 20000 {
			t.Fatalf("expected 20000, got %d (err: %v)", port, err)
		}
	})

	t.Run("skips occupied and avoided ports", func(t *testing.T) {
		occupied := map[int]bool{20000: true, 20001: true}
		avoid := []int{20002, 20003}
		port, err := PickHostPort(occupied, AllocateHostPortOptions{
			RangeStart:     20000,
			RangeEnd:       20999,
			Avoid:          avoid,
			EphemeralRange: ephemeral,
		})
		if err != nil || port != 20004 {
			t.Fatalf("expected 20004, got %d (err: %v)", port, err)
		}
	})

	t.Run("does not reuse preferred port inside ephemeral range", func(t *testing.T) {
		port, err := PickHostPort(map[int]bool{}, AllocateHostPortOptions{
			Preferred:      32769,
			RangeStart:     20000,
			RangeEnd:       20999,
			EphemeralRange: ephemeral,
		})
		if err != nil || port != 20000 {
			t.Fatalf("expected 20000, got %d (err: %v)", port, err)
		}
	})

	t.Run("refuses ephemeral range when range is fully exhausted", func(t *testing.T) {
		_, err := PickHostPort(map[int]bool{30000: true, 30001: true}, AllocateHostPortOptions{
			RangeStart:     30000,
			RangeEnd:       30001,
			EphemeralRange: ephemeral,
		})
		if err == nil {
			t.Fatalf("expected error when range exhausted, got nil")
		}
	})
}

func TestSuggestAlternativeHostPort(t *testing.T) {
	ephemeral := &PortRange{Start: 32768, End: 60999}
	occupied := map[int]bool{20000: true}
	excluded := []int{20001, 20002}

	port, err := SuggestAlternativeHostPort(occupied, excluded, AllocateHostPortOptions{
		RangeStart:     20000,
		RangeEnd:       20999,
		EphemeralRange: ephemeral,
	})
	if err != nil || port != 20003 {
		t.Fatalf("expected 20003, got %d (err: %v)", port, err)
	}
}
