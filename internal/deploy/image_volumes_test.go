package deploy

import (
	"context"
	"testing"
)

func TestCleanVolumeDestination(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"/app/data", "/app/data", false},
		{"/var/lib/postgresql/data/", "/var/lib/postgresql/data", false},
		{"/data//sub", "/data/sub", false},
		{"data/relative", "", true},
		{"/app:invalid", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		got, err := CleanVolumeDestination(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("CleanVolumeDestination(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("CleanVolumeDestination(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestImageVolumeNameDeterministic(t *testing.T) {
	name1 := ImageVolumeName("svc-123", "/var/lib/postgresql/data")
	name2 := ImageVolumeName("svc-123", "/var/lib/postgresql/data")
	if name1 != name2 {
		t.Fatalf("expected deterministic volume name, got %s vs %s", name1, name2)
	}
	if len(name1) < 20 {
		t.Fatalf("volume name too short: %s", name1)
	}

	diff := ImageVolumeName("svc-456", "/var/lib/postgresql/data")
	if diff == name1 {
		t.Fatalf("expected different volume name for different workload, got %s", diff)
	}
}

func TestResolveImageVolumeBindsNilClient(t *testing.T) {
	binds := []string{"/host/data:/app/data"}
	res, err := ResolveImageVolumeBinds(context.Background(), nil, "svc-1", "proj-1", binds, "some-img:tag", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0] != binds[0] {
		t.Fatalf("expected original binds preserved, got %v", res)
	}
}
