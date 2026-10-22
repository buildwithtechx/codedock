package deploy

import (
	"testing"
)

func TestHasHostPortBindings(t *testing.T) {
	tests := []struct {
		name     string
		ports    []string
		expected bool
	}{
		{"empty", []string{}, false},
		{"bare container port", []string{"3000", "80"}, false},
		{"container port with protocol", []string{"3000/tcp", "80/udp"}, false},
		{"host port mapping", []string{"8080:80"}, true},
		{"host port with protocol", []string{"8080:80/tcp"}, true},
		{"ip and host port", []string{"127.0.0.1:8080:80"}, true},
		{"ip and host port with protocol", []string{"0.0.0.0:8080:80/tcp"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasHostPortBindings(tt.ports)
			if got != tt.expected {
				t.Errorf("HasHostPortBindings(%v) = %v, want %v", tt.ports, got, tt.expected)
			}
		})
	}
}

func TestClashingRWDestinations(t *testing.T) {
	incumbent := []string{"/host/data:/app/data", "/host/logs:/app/logs:ro"}
	candidate := []string{"/host/new-data:/app/data", "/host/logs:/app/logs"}

	// Without opt-in, /app/data clashes because incumbent has it RW and candidate has it RW
	clashes := ClashingRWDestinations(incumbent, candidate, nil)
	if len(clashes) != 1 || clashes[0] != "/app/data" {
		t.Fatalf("expected ['/app/data'], got %v", clashes)
	}

	// With opt-in for /app/data, no clashes
	optIn := map[string]bool{"/app/data": true}
	clashesWithOptIn := ClashingRWDestinations(incumbent, candidate, optIn)
	if len(clashesWithOptIn) != 0 {
		t.Fatalf("expected empty clashes with opt-in, got %v", clashesWithOptIn)
	}
}

func TestDetermineCutoverMode(t *testing.T) {
	t.Run("no incumbent replaces", func(t *testing.T) {
		dec := DetermineCutoverMode(CutoverInput{
			IncumbentRunning: false,
		})
		if dec.Mode != CutoverModeReplace || dec.Reason != "no-incumbent" {
			t.Errorf("unexpected decision: %+v", dec)
		}
	})

	t.Run("host ports replaces", func(t *testing.T) {
		dec := DetermineCutoverMode(CutoverInput{
			IncumbentRunning: true,
			Ports:            []string{"8080:80"},
		})
		if dec.Mode != CutoverModeReplace || dec.Reason != "host-ports" {
			t.Errorf("unexpected decision: %+v", dec)
		}
	})

	t.Run("shared namespaces replaces", func(t *testing.T) {
		dec := DetermineCutoverMode(CutoverInput{
			IncumbentRunning: true,
			NetworkMode:      "host",
		})
		if dec.Mode != CutoverModeReplace || dec.Reason != "shared-namespaces" {
			t.Errorf("unexpected decision: %+v", dec)
		}
	})

	t.Run("shared rw volume without opt-in replaces", func(t *testing.T) {
		dec := DetermineCutoverMode(CutoverInput{
			IncumbentRunning: true,
			Volumes:          []string{"/host/v:/app/data"},
			IncumbentMounts:  []string{"/host/v:/app/data"},
		})
		if dec.Mode != CutoverModeReplace || dec.Reason != "shared-rw-volume" {
			t.Errorf("unexpected decision: %+v", dec)
		}
		if len(dec.Clashes) != 1 || dec.Clashes[0] != "/app/data" {
			t.Errorf("expected clash on /app/data, got %v", dec.Clashes)
		}
	})

	t.Run("standard container overlaps", func(t *testing.T) {
		dec := DetermineCutoverMode(CutoverInput{
			IncumbentRunning: true,
			Ports:            []string{"3000"},
		})
		if dec.Mode != CutoverModeOverlap {
			t.Errorf("expected overlap mode, got %+v", dec)
		}
	})
}
