package compose

import "testing"

const nsPerSecond = int64(1_000_000_000)

func TestParseComposeDurationReadsComposeAndGoDurations(t *testing.T) {
	cases := []struct {
		value string
		want  int64
	}{
		{"90s", 90 * nsPerSecond},
		{"1m30s", 90 * nsPerSecond},
		{"1h2m3s", 3723 * nsPerSecond},
		{"500ms", nsPerSecond / 2},
		{"1.5s", 1_500_000_000},
		{"30", 30 * nsPerSecond},
		{" 90s ", 90 * nsPerSecond},
		{"0s", 0},
	}
	for _, tc := range cases {
		if got, ok := ParseComposeDuration(tc.value); !ok || got != tc.want {
			t.Errorf("ParseComposeDuration(%q) = %d, %v; want %d, true", tc.value, got, ok, tc.want)
		}
	}
}

func TestParseComposeDurationPreservesComposeSignsDecimalsAndUnits(t *testing.T) {
	cases := []struct {
		value string
		want  int64
	}{
		{"+90s", 90 * nsPerSecond},
		{"+1m30s", 90 * nsPerSecond},
		{".5s", nsPerSecond / 2},
		{"+.5s", nsPerSecond / 2},
		{"1.s", nsPerSecond},
		{"1m.5s", 60*nsPerSecond + nsPerSecond/2},
		{"1.5h", 5400 * nsPerSecond},
		{"1d2h", 93600 * nsPerSecond},
		{"1w", 604800 * nsPerSecond},
		{"1μs", 1000},
		{"+30", 30 * nsPerSecond},
		{".5", nsPerSecond / 2},
	}
	for _, tc := range cases {
		if got, ok := ParseComposeDuration(tc.value); !ok || got != tc.want {
			t.Errorf("ParseComposeDuration(%q) = %d, %v; want %d, true", tc.value, got, ok, tc.want)
		}
	}
}

func TestParseComposeDurationReadsPartialMatchAsNoDuration(t *testing.T) {
	for _, value := range []string{
		"1m30",
		"abc5s",
		"5s junk",
		"10 m",
		"1m 30s",
		"30S",
		"-5s",
		"+",
		".s",
		"1s+2s",
		"1..5s",
		"",
		"   ",
	} {
		if got, ok := ParseComposeDuration(value); ok {
			t.Errorf("ParseComposeDuration(%q) = %d, true; want not ok", value, got)
		}
	}
}

func TestParseComposeDurationRejectsOverflow(t *testing.T) {
	long := ""
	for i := 0; i < 400; i++ {
		long += "9"
	}
	for _, value := range []string{long + "s", "9223372036854775808ns", "99999999999999999999h"} {
		if _, ok := ParseComposeDuration(value); ok {
			t.Errorf("ParseComposeDuration(%q) ok; want rejected", value)
		}
	}
}

func TestValidateCatalogHealthcheckUsesComposeDurations(t *testing.T) {
	check := func(interval string) *CatalogHealthcheck {
		return &CatalogHealthcheck{Test: []string{"CMD", "true"}, Interval: interval}
	}
	for _, value := range []string{"30", "1m30s", "500ms"} {
		if failures := validateCatalogHealthcheck("web", check(value)); len(failures) != 0 {
			t.Errorf("interval %q should be accepted, got %v", value, failures)
		}
	}
	for _, value := range []string{"1m30", "soon", "0"} {
		if failures := validateCatalogHealthcheck("web", check(value)); len(failures) == 0 {
			t.Errorf("interval %q should be rejected", value)
		}
	}
}
