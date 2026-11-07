package deployments

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectMonorepoCandidateApps(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "codedock-monorepo-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	webDir := filepath.Join(tempDir, "apps", "web")
	apiDir := filepath.Join(tempDir, "apps", "api")
	docsDir := filepath.Join(tempDir, "packages", "docs")

	for _, d := range []string{webDir, apiDir, docsDir} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatalf("failed to create sub dir: %v", err)
		}
	}

	if err := os.WriteFile(filepath.Join(webDir, "package.json"), []byte(`{"name": "@myorg/web", "scripts": {"build": "next build"}}`), 0644); err != nil {
		t.Fatalf("failed to write web package.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(apiDir, "go.mod"), []byte("module myorg/api\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatalf("failed to write api go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "Dockerfile"), []byte("FROM alpine\n"), 0644); err != nil {
		t.Fatalf("failed to write docs Dockerfile: %v", err)
	}

	candidates := DetectMonorepoCandidateApps(tempDir)
	if len(candidates) != 3 {
		t.Fatalf("expected 3 candidate apps, got %d", len(candidates))
	}

	foundWeb := false
	foundAPI := false
	foundDocs := false

	for _, c := range candidates {
		if c.RootDirectory == "apps/web" {
			foundWeb = true
			if c.Name != "myorg-web" {
				t.Errorf("expected app name myorg-web, got %s", c.Name)
			}
		}
		if c.RootDirectory == "apps/api" {
			foundAPI = true
			if c.Framework != "Go" {
				t.Errorf("expected Go framework, got %s", c.Framework)
			}
		}
		if c.RootDirectory == "packages/docs" {
			foundDocs = true
			if c.BuildEngine != "dockerfile" {
				t.Errorf("expected dockerfile build engine, got %s", c.BuildEngine)
			}
		}
	}

	if !foundWeb || !foundAPI || !foundDocs {
		t.Fatalf("did not find all expected apps: web=%t, api=%t, docs=%t", foundWeb, foundAPI, foundDocs)
	}
}
