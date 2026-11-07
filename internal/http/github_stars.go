package http

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type githubStarsCache struct {
	mu        sync.RWMutex
	stars     int
	fetchedAt time.Time
}

var starsCache = &githubStarsCache{}

func (s *Server) getGithubStars(c echo.Context) error {
	starsCache.mu.RLock()
	if time.Since(starsCache.fetchedAt) < 15*time.Minute && starsCache.stars > 0 {
		stars := starsCache.stars
		starsCache.mu.RUnlock()
		c.Response().Header().Set("Cache-Control", "public, max-age=900")
		return c.JSON(http.StatusOK, map[string]any{"stars": stars})
	}
	starsCache.mu.RUnlock()

	req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodGet, "https://api.github.com/repos/codedock/codedock", nil)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]any{"stars": 0})
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "Codedock")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		starsCache.mu.RLock()
		fallback := starsCache.stars
		starsCache.mu.RUnlock()
		return c.JSON(http.StatusOK, map[string]any{"stars": fallback})
	}
	defer resp.Body.Close()

	var payload struct {
		StargazersCount int `json:"stargazers_count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return c.JSON(http.StatusOK, map[string]any{"stars": 0})
	}

	starsCache.mu.Lock()
	starsCache.stars = payload.StargazersCount
	starsCache.fetchedAt = time.Now()
	starsCache.mu.Unlock()

	c.Response().Header().Set("Cache-Control", "public, max-age=900")
	return c.JSON(http.StatusOK, map[string]any{"stars": payload.StargazersCount})
}
