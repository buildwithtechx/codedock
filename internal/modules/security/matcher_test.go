package security

import (
	"testing"
)

func TestNormalizeRepo(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"node", "docker.io/library/node"},
		{"node:22-alpine", "docker.io/library/node"},
		{"myorg/myapp:v1", "docker.io/myorg/myapp"},
		{"index.docker.io/myorg/myapp", "docker.io/myorg/myapp"},
		{"ghcr.io/org/repo:tag", "ghcr.io/org/repo"},
		{"registry.gitlab.com/group/proj@sha256:abcd", "registry.gitlab.com/group/proj"},
	}

	for _, c := range cases {
		got := NormalizeRepo(c.input)
		if got != c.expected {
			t.Errorf("NormalizeRepo(%q) = %q; want %q", c.input, got, c.expected)
		}
	}
}

func TestMatchAdvisories(t *testing.T) {
	nodeFixed := ">=22.23.4 <23 || >=24.22.1 <25"
	advisories := []SecurityAdvisory{
		{
			ID:       "NODE-BROTLI",
			Title:    "Node Brotli vulnerability",
			Severity: SeverityHigh,
			Action:   ActionRefreshImages,
			Targets: []Target{
				{
					Kind:     TargetKindRuntime,
					Runtime:  "node",
					Affected: "<22.23.4 || >=23.0.0 <24.22.1",
					Fixed:    &nodeFixed,
				},
			},
		},
		{
			ID:       "IMAGE-CVE",
			Title:    "Vulnerable Base Image",
			Severity: SeverityCritical,
			Action:   ActionRefreshImages,
			Targets: []Target{
				{
					Kind:            TargetKindImage,
					Repository:      "node",
					AffectedDigests: []string{"sha256:vulnerable123"},
					FixedDigests:    []string{"sha256:fixed456"},
					AffectedTags:    []string{"22-alpine"},
				},
			},
		},
	}

	t.Run("matches vulnerable runtime with same major fix available", func(t *testing.T) {
		services := []ServiceEvidence{
			{
				ServiceName:     "web",
				Image:           "node:22-alpine",
				RuntimeVersions: map[string]string{"node": "22.23.3"},
			},
		}

		findings := MatchAdvisories(advisories, services)
		if len(findings) == 0 {
			t.Fatal("expected finding for vulnerable runtime")
		}

		f := findings[0]
		if f.AdvisoryID != "NODE-BROTLI" || f.Confidence != ConfidenceConfirmed {
			t.Fatalf("unexpected finding: %+v", f)
		}
		if !f.CanRefresh || f.RuntimeUpgradeRequired {
			t.Fatalf("expected canRefresh true and runtimeUpgradeRequired false, got: %+v", f)
		}
	})

	t.Run("flags runtimeUpgradeRequired when major upgrade is needed", func(t *testing.T) {
		services := []ServiceEvidence{
			{
				ServiceName:     "api",
				Image:           "node:20-alpine",
				RuntimeVersions: map[string]string{"node": "20.18.0"},
			},
		}

		findings := MatchAdvisories(advisories, services)
		if len(findings) == 0 {
			t.Fatal("expected finding for node 20")
		}

		f := findings[0]
		if !f.RuntimeUpgradeRequired || f.CanRefresh {
			t.Fatalf("expected runtimeUpgradeRequired true and canRefresh false, got: %+v", f)
		}
	})

	t.Run("matches image digest with confirmed confidence", func(t *testing.T) {
		services := []ServiceEvidence{
			{
				ServiceName: "worker",
				Image:       "node:22-alpine",
				ImageDigest: "sha256:vulnerable123",
			},
		}

		findings := MatchAdvisories(advisories, services)
		matched := false
		for _, f := range findings {
			if f.AdvisoryID == "IMAGE-CVE" {
				matched = true
				if f.Confidence != ConfidenceConfirmed {
					t.Fatalf("expected confirmed confidence for matching digest, got %s", f.Confidence)
				}
			}
		}
		if !matched {
			t.Fatal("expected image advisory to match")
		}
	})

	t.Run("does not match when fixed digest is present", func(t *testing.T) {
		services := []ServiceEvidence{
			{
				ServiceName: "worker",
				Image:       "node:22-alpine",
				ImageDigest: "sha256:fixed456",
			},
		}

		findings := MatchAdvisories(advisories, services)
		for _, f := range findings {
			if f.AdvisoryID == "IMAGE-CVE" {
				t.Fatalf("expected fixed digest not to match, but got finding: %+v", f)
			}
		}
	})
}
