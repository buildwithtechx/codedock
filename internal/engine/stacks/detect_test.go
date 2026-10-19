package stacks

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDetectStackLaravel(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "composer.json", `{"require":{"laravel/framework":"^11"}}`)
	writeFile(t, dir, "artisan", "")
	if got := DetectStack(dir); got != "laravel" {
		t.Fatalf("expected laravel, got %q", got)
	}
}

func TestDetectStackNextAndGo(t *testing.T) {
	next := t.TempDir()
	writeFile(t, next, "package.json", `{"dependencies":{"next":"14"}}`)
	if got := DetectStack(next); got != "nextjs" {
		t.Fatalf("expected nextjs, got %q", got)
	}

	golang := t.TempDir()
	writeFile(t, golang, "go.mod", "module example.com/app\n\ngo 1.22\n\nrequire github.com/labstack/echo/v4 v4.12.0\n")
	if got := DetectStack(golang); got != "go" {
		t.Fatalf("expected go, got %q", got)
	}
}

func TestDetectDepsParsesManifests(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "Gemfile", "source \"https://rubygems.org\"\ngem \"solid_queue\"\ngem 'sidekiq', '~> 7'\n")
	writeFile(t, dir, "requirements.txt", "Flask==3.0.0\n# comment\n")
	deps := DetectDeps(dir)
	for _, want := range []string{"solid_queue", "sidekiq", "flask"} {
		if !contains(deps, want) {
			t.Errorf("expected dep %q in %v", want, deps)
		}
	}
}

func TestResolveRolesForSource(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "composer.json", `{"require":{"laravel/framework":"^11","laravel/horizon":"^5"}}`)
	writeFile(t, dir, "artisan", "")
	id, roles := ResolveRolesForSource(dir)
	if id != "laravel" {
		t.Fatalf("expected laravel, got %q", id)
	}
	names := roleNames(roles)
	if !contains(names, "horizon") || contains(names, "queue") {
		t.Fatalf("expected horizon without queue, got %v", names)
	}
}

func TestDetectFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "bin/jobs", "")
	writeFile(t, dir, "node_modules/ignored.js", "")
	files := DetectFiles(dir)
	if !contains(files, "bin/jobs") {
		t.Fatalf("expected bin/jobs in %v", files)
	}
	if contains(files, "node_modules/ignored.js") {
		t.Fatalf("node_modules should be skipped, got %v", files)
	}
}

func TestDetectStackComposeDisambiguation(t *testing.T) {
	for _, composeFile := range []string{"docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml"} {
		dir := t.TempDir()
		writeFile(t, dir, "Dockerfile", "FROM node:22-alpine\n")
		writeFile(t, dir, composeFile, "services:\n  web:\n    build: .\n")
		writeFile(t, dir, "package.json", `{"dependencies":{"next":"14"}}`)
		if got := DetectStack(dir); got != "docker-compose" {
			t.Fatalf("expected docker-compose for %s alongside Dockerfile, got %q", composeFile, got)
		}
	}

	bareCompose := t.TempDir()
	writeFile(t, bareCompose, "compose.yaml", "services:\n  redis:\n    image: redis:alpine\n")
	if got := DetectStack(bareCompose); got != "docker-compose" {
		t.Fatalf("expected docker-compose for bare compose, got %q", got)
	}

	frameworkDevCompanion := t.TempDir()
	writeFile(t, frameworkDevCompanion, "package.json", `{"dependencies":{"next":"14"}}`)
	writeFile(t, frameworkDevCompanion, "compose.yaml", "services:\n  db:\n    image: postgres:16\n")
	if got := DetectStack(frameworkDevCompanion); got != "nextjs" {
		t.Fatalf("expected nextjs when compose is dev companion without Dockerfile, got %q", got)
	}
}
