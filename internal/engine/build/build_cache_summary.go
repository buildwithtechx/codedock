package build

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

var (
	buildKitStepRegex     = regexp.MustCompile(`^#(\d+)\s`)
	buildKitActionRegex   = regexp.MustCompile(`\]\s+(RUN|COPY|ADD)\s`)
	buildKitCachedRegex   = regexp.MustCompile(`^#\d+\s+CACHED`)
	buildKitDoneRegex     = regexp.MustCompile(`^#\d+\s+DONE`)
)

type BuildKitCacheObserver struct {
	mu            sync.Mutex
	cacheLines    map[string]string
	cacheSteps    map[string]bool
	cacheHits     map[string]bool
	cacheFinished map[string]bool
}

func NewBuildKitCacheObserver() *BuildKitCacheObserver {
	return &BuildKitCacheObserver{
		cacheLines:    make(map[string]string),
		cacheSteps:    make(map[string]bool),
		cacheHits:     make(map[string]bool),
		cacheFinished: make(map[string]bool),
	}
}

func (o *BuildKitCacheObserver) Observe(data string, streamID string) {
	if streamID == "" {
		streamID = "default"
	}
	o.mu.Lock()
	defer o.mu.Unlock()

	combined := o.cacheLines[streamID] + data
	lines := strings.Split(combined, "\n")
	if len(lines) > 0 {
		remainder := lines[len(lines)-1]
		if len(remainder) > 8192 {
			remainder = remainder[len(remainder)-8192:]
		}
		o.cacheLines[streamID] = remainder
		lines = lines[:len(lines)-1]
	}

	for _, line := range lines {
		trimmed := strings.TrimRight(line, "\r")
		m := buildKitStepRegex.FindStringSubmatch(trimmed)
		if len(m) < 2 {
			continue
		}
		stepID := fmt.Sprintf("%s:%s", streamID, m[1])
		if buildKitActionRegex.MatchString(trimmed) {
			o.cacheSteps[stepID] = true
		}
		if buildKitCachedRegex.MatchString(trimmed) {
			o.cacheHits[stepID] = true
		}
		if buildKitDoneRegex.MatchString(trimmed) {
			o.cacheFinished[stepID] = true
		}
	}
}

func (o *BuildKitCacheObserver) Summary() string {
	o.mu.Lock()
	defer o.mu.Unlock()

	var total, hits int
	for stepID := range o.cacheSteps {
		if o.cacheHits[stepID] || o.cacheFinished[stepID] {
			total++
			if o.cacheHits[stepID] {
				hits++
			}
		}
	}
	if total == 0 {
		return ""
	}
	misses := total - hits
	return fmt.Sprintf("Docker layer cache: %d hit(s), %d miss(es).", hits, misses)
}
