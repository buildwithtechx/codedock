package deployments

import (
	"errors"
	"strings"
	"sync"

	"codedock/internal/config"
)

var (
	ErrBuildAdmissionLimitExceeded = errors.New("concurrent build admission ceiling reached; wait for running builds to complete or upgrade your plan")
)

type BuildAdmissionGuard struct {
	mu           sync.Mutex
	activeBuilds map[string]int
}

func NewBuildAdmissionGuard() *BuildAdmissionGuard {
	return &BuildAdmissionGuard{
		activeBuilds: make(map[string]int),
	}
}

func PlanBuildLimit(tier string) int {
	switch strings.ToLower(strings.TrimSpace(tier)) {
	case "hobby", "free":
		return 1
	case "starter":
		return 4
	case "pro":
		return 8
	case "scale", "enterprise":
		return 20
	default:
		if defaultLimit := config.Get().Limits.MaxConcurrentBuilds; defaultLimit > 0 {
			return defaultLimit
		}
		return 2
	}
}

func AssertBuildAdmission(activeBuildsCount int, planTier string) error {
	limit := PlanBuildLimit(planTier)
	if activeBuildsCount >= limit {
		return ErrBuildAdmissionLimitExceeded
	}
	return nil
}

func (g *BuildAdmissionGuard) Acquire(scopeID, planTier string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	active := g.activeBuilds[scopeID]
	if err := AssertBuildAdmission(active, planTier); err != nil {
		return err
	}
	g.activeBuilds[scopeID] = active + 1
	return nil
}

func (g *BuildAdmissionGuard) Release(scopeID string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	active := g.activeBuilds[scopeID]
	if active > 1 {
		g.activeBuilds[scopeID] = active - 1
	} else {
		delete(g.activeBuilds, scopeID)
	}
}

func (g *BuildAdmissionGuard) ActiveCount(scopeID string) int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.activeBuilds[scopeID]
}
