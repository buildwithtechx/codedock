package stacks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	gemfileGem = regexp.MustCompile(`^\s*gem\s+['"]([^'"]+)['"]`)
	pythonPkg  = regexp.MustCompile(`^\s*([A-Za-z0-9_.\-]+)`)
	goModReq   = regexp.MustCompile(`^\s*([^\s/]+/[^\s]+)\s+v[0-9]`)
)

func DetectStack(sourceDir string) string {
	deps := DetectDeps(sourceDir)
	depSet := make(map[string]bool, len(deps))
	for _, dep := range deps {
		depSet[strings.ToLower(dep)] = true
	}
	for _, id := range DetectOrder {
		def := Registry[id]
		if stackMatches(def, sourceDir, depSet) {
			return id
		}
	}
	return "unknown"
}

func stackMatches(def StackDefinition, sourceDir string, deps map[string]bool) bool {
	for _, marker := range def.RootMarkers {
		if _, err := os.Stat(filepath.Join(sourceDir, marker)); err != nil {
			return false
		}
	}
	if len(def.Deps) > 0 {
		found := false
		for _, dep := range def.Deps {
			if deps[strings.ToLower(dep)] {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return len(def.RootMarkers) > 0 || len(def.Deps) > 0
}

func DetectDeps(sourceDir string) []string {
	set := map[string]bool{}
	if data, err := os.ReadFile(filepath.Join(sourceDir, "package.json")); err == nil {
		var pkg struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if json.Unmarshal(data, &pkg) == nil {
			for name := range pkg.Dependencies {
				set[name] = true
			}
			for name := range pkg.DevDependencies {
				set[name] = true
			}
		}
	}
	if data, err := os.ReadFile(filepath.Join(sourceDir, "composer.json")); err == nil {
		var composer struct {
			Require    map[string]string `json:"require"`
			RequireDev map[string]string `json:"require-dev"`
		}
		if json.Unmarshal(data, &composer) == nil {
			for name := range composer.Require {
				set[name] = true
			}
			for name := range composer.RequireDev {
				set[name] = true
			}
		}
	}
	if data, err := os.ReadFile(filepath.Join(sourceDir, "Gemfile")); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if m := gemfileGem.FindStringSubmatch(line); m != nil {
				set[m[1]] = true
			}
		}
	}
	if data, err := os.ReadFile(filepath.Join(sourceDir, "requirements.txt")); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if m := pythonPkg.FindStringSubmatch(line); m != nil {
				set[strings.ToLower(m[1])] = true
			}
		}
	}
	if data, err := os.ReadFile(filepath.Join(sourceDir, "go.mod")); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if m := goModReq.FindStringSubmatch(line); m != nil {
				set[m[1]] = true
			}
		}
	}
	return sortedKeys(set)
}

const detectFileMaxDepth = 3

var detectSkipDirs = map[string]bool{
	".git": true, "node_modules": true, "vendor": true,
	"dist": true, "build": true, "target": true, ".next": true,
}

func DetectFiles(sourceDir string) []string {
	set := map[string]bool{}
	root := filepath.Clean(sourceDir)
	var walk func(dir string, depth int)
	walk = func(dir string, depth int) {
		if depth > detectFileMaxDepth {
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() {
				if detectSkipDirs[name] {
					continue
				}
				walk(filepath.Join(dir, name), depth+1)
				continue
			}
			rel, err := filepath.Rel(root, filepath.Join(dir, name))
			if err != nil {
				continue
			}
			set[filepath.ToSlash(rel)] = true
		}
	}
	walk(root, 1)
	return sortedKeys(set)
}

func sortedKeys(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for key := range set {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}
