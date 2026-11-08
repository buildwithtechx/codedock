package deployments

import (
	"fmt"
	"sync"
)

type BuildAdmissionGuard struct {
	mu           sync.Mutex
	activeBuilds map[string]int // serviceID -> count
	globalActive int
	maxGlobal    int
	maxService   int
}

func NewBuildAdmissionGuard() *BuildAdmissionGuard {
	return &BuildAdmissionGuard{
		activeBuilds: make(map[string]int),
		maxGlobal:    DefaultMaxConcurrentBuilds,
		maxService:   DefaultServiceBuildLimit,
	}
}

func (g *BuildAdmissionGuard) Acquire(serviceID string, _ string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.globalActive >= g.maxGlobal {
		return fmt.Errorf("%w: global limit of %d concurrent builds reached", ErrAdmissionLimitExceeded, g.maxGlobal)
	}

	if g.activeBuilds[serviceID] >= g.maxService {
		return fmt.Errorf("%w: service build limit of %d reached", ErrAdmissionLimitExceeded, g.maxService)
	}

	g.activeBuilds[serviceID]++
	g.globalActive++
	return nil
}

func (g *BuildAdmissionGuard) Release(serviceID string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.activeBuilds[serviceID] > 0 {
		g.activeBuilds[serviceID]--
		if g.activeBuilds[serviceID] == 0 {
			delete(g.activeBuilds, serviceID)
		}
	}
	if g.globalActive > 0 {
		g.globalActive--
	}
}

func (g *BuildAdmissionGuard) Status() BuildAdmissionStatus {
	g.mu.Lock()
	defer g.mu.Unlock()

	return BuildAdmissionStatus{
		ActiveBuilds:        g.globalActive,
		MaxConcurrentBuilds: g.maxGlobal,
		QueueLength:         0,
	}
}
