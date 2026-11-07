package compose

import (
	"testing"

	"github.com/docker/docker/api/types/container"
)

func TestParseComposeShmSize(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expected  int64
		expectErr bool
	}{
		{"nil", nil, 0, false},
		{"empty string", "", 0, false},
		{"bytes integer", 2147483648, 2147483648, false},
		{"bytes int64", int64(536870912), 536870912, false},
		{"size string 2g", "2g", 2147483648, false},
		{"size string 2gb", "2gb", 2147483648, false},
		{"size string 512m", "512m", 536870912, false},
		{"size string 512mb", "512mb", 536870912, false},
		{"size string with space", "  1024 k  ", 1024 * 1024, false},
		{"size float 1.5g", "1.5g", int64(1.5 * 1024 * 1024 * 1024), false},
		{"size string 1t", "1t", 1024 * 1024 * 1024 * 1024, false},
		{"invalid format", "invalid", 0, true},
		{"invalid unit", "2x", 0, true},
		{"zero bytes", 0, 0, true},
		{"negative bytes", -1024, 0, true},
		{"over 1 TiB", "2t", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseComposeShmSize(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error for input %v, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error for input %v: %v", tt.input, err)
				}
				if got != tt.expected {
					t.Fatalf("expected %d, got %d", tt.expected, got)
				}
			}
		})
	}
}

func TestApplyHostConfigShmSize(t *testing.T) {
	hc := &container.HostConfig{}
	ApplyHostConfigShmSize(hc, 2147483648)
	if hc.ShmSize != 2147483648 {
		t.Fatalf("expected ShmSize 2147483648, got %d", hc.ShmSize)
	}

	ApplyHostConfigShmSize(nil, 2147483648) // must not panic
}

func TestFormatShmSize(t *testing.T) {
	if s := FormatShmSize(2147483648); s != "2gb" {
		t.Fatalf("expected 2gb, got %s", s)
	}
	if s := FormatShmSize(536870912); s != "512mb" {
		t.Fatalf("expected 512mb, got %s", s)
	}
	if s := FormatShmSize(0); s != "" {
		t.Fatalf("expected empty string, got %s", s)
	}
}
