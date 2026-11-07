package deploy

import (
	"testing"
)

func TestExtractBuildExitCode(t *testing.T) {
	cases := []struct {
		input string
		want  int
	}{
		{"command returned a non-zero code: 137", 137},
		{"process exited with code 1", 1},
		{"The command '/bin/sh -c npm run build' returned a non-zero code: 2", 2},
		{"no exit code here", 0},
	}
	for _, c := range cases {
		if got := ExtractBuildExitCode(c.input); got != c.want {
			t.Errorf("ExtractBuildExitCode(%q) = %d, want %d", c.input, got, c.want)
		}
	}
}

func TestClassifyResourceIssue(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"FATAL ERROR: Reached heap limit Allocation failed - JavaScript heap out of memory", ResourceIssueProcessHeap},
		{"Ineffective mark-compacts near heap limit", ResourceIssueProcessHeap},
		{"fatal error: runtime: out of memory", ResourceIssueMemoryAllocation},
		{"fork/exec: cannot allocate memory", ResourceIssueMemoryAllocation},
		{"write /tmp/build: no space left on device", ResourceIssueStorageFull},
		{"build failed with ENOSPC error", ResourceIssueStorageFull},
		{"container exited with code 137", ResourceIssueProcessKilled},
		{"some regular syntax error", ""},
	}
	for _, c := range cases {
		if got := ClassifyResourceIssue(c.input); got != c.want {
			t.Errorf("ClassifyResourceIssue(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
