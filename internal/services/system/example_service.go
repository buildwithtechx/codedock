package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"codedock/internal/models"
)

type ExampleService struct {
	contentsURL string
	manifestURL string
	rawBase     string
	cache       []models.ExampleApp
	lastFetched time.Time
	mu          sync.RWMutex
}

type exampleManifestEntry struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
}

type exampleManifest struct {
	Templates []exampleManifestEntry `json:"templates"`
}

func NewExampleService() *ExampleService {
	return NewExampleServiceWithURLs(
		"https://api.github.com/repos/buildwithtechx/codedock-examples/contents",
		"https://raw.githubusercontent.com/buildwithtechx/codedock-examples/main/templates.json",
		"https://raw.githubusercontent.com/buildwithtechx/codedock-examples/main/",
	)
}

func NewExampleServiceWithURLs(contentsURL, manifestURL, rawBase string) *ExampleService {
	return &ExampleService{contentsURL: contentsURL, manifestURL: manifestURL, rawBase: rawBase}
}

func (s *ExampleService) ListExamples() ([]models.ExampleApp, error) {
	s.mu.RLock()
	cacheValid := time.Since(s.lastFetched) < 1*time.Hour && len(s.cache) > 0
	cached := s.cache
	s.mu.RUnlock()

	if cacheValid {
		return cached, nil
	}

	manifest := fetchExampleManifest(s.manifestURL)
	contents, err := fetchExampleContents(s.contentsURL)
	if err != nil {
		if len(cached) > 0 {
			return cached, nil
		}
		if len(manifest) > 0 {
			examples := []models.ExampleApp{}
			for id, entry := range manifest {
				logo := ""
				if entry.Logo != "" {
					logo = s.rawBase + strings.TrimPrefix(entry.Logo, "/")
				}
				name := entry.Name
				if name == "" {
					name = formatExampleName(id)
				}
				desc := entry.Description
				if desc == "" {
					desc = "Deploy " + name + " example app"
				}
				examples = append(examples, models.ExampleApp{
					ID:          id,
					Name:        name,
					Description: desc,
					Repo:        "https://github.com/buildwithtechx/codedock-examples/tree/main/" + id,
					Logo:        logo,
				})
			}
			sort.Slice(examples, func(i, j int) bool {
				return examples[i].Name < examples[j].Name
			})
			s.mu.Lock()
			s.cache = examples
			s.lastFetched = time.Now()
			s.mu.Unlock()
			return examples, nil
		}
		return nil, err
	}
	assets := exampleAssetDirs(manifest)

	examples := []models.ExampleApp{}
	for _, entry := range contents {
		if entry.Type != "dir" || strings.HasPrefix(entry.Name, ".") {
			continue
		}
		if assets[entry.Name] {
			continue
		}
		examples = append(examples, s.describeExample(entry, manifest))
	}

	s.mu.Lock()
	s.cache = examples
	s.lastFetched = time.Now()
	s.mu.Unlock()

	return examples, nil
}

func exampleAssetDirs(manifest map[string]exampleManifestEntry) map[string]bool {
	assets := map[string]bool{}
	for _, entry := range manifest {
		if dir, _, found := strings.Cut(entry.Logo, "/"); found && dir != "" {
			assets[dir] = true
		}
	}
	return assets
}

func (s *ExampleService) describeExample(entry exampleContentEntry, manifest map[string]exampleManifestEntry) models.ExampleApp {
	example := models.ExampleApp{
		ID:          entry.Name,
		Name:        formatExampleName(entry.Name),
		Description: "Deploy " + formatExampleName(entry.Name) + " example app",
		Repo:        entry.HtmlUrl,
	}
	if described, ok := manifest[entry.Name]; ok {
		if described.Name != "" {
			example.Name = described.Name
		}
		if described.Description != "" {
			example.Description = described.Description
		}
		if described.Logo != "" {
			example.Logo = s.rawBase + strings.TrimPrefix(described.Logo, "/")
		}
	}
	return example
}

type exampleContentEntry struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	HtmlUrl string `json:"html_url"`
}

func fetchExampleContents(url string) ([]exampleContentEntry, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Codedock")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var contents []exampleContentEntry
	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		return nil, err
	}
	return contents, nil
}

func fetchExampleManifest(url string) map[string]exampleManifestEntry {
	manifest := map[string]exampleManifestEntry{}
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return manifest
	}
	req.Header.Set("User-Agent", "Codedock")
	resp, err := client.Do(req)
	if err != nil {
		return manifest
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return manifest
	}

	var document exampleManifest
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return manifest
	}
	for _, entry := range document.Templates {
		if entry.ID != "" {
			manifest[entry.ID] = entry
		}
	}
	return manifest
}

func formatExampleName(name string) string {
	parts := strings.Split(name, "-")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, " ")
}
