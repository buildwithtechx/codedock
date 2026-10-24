package deploy

import (
	"path"
	"regexp"
	"strings"
)

type CutoverMode string

const (
	CutoverModeOverlap CutoverMode = "overlap"
	CutoverModeReplace CutoverMode = "replace"
)

type CutoverDecision struct {
	Mode    CutoverMode `json:"mode"`
	Reason  string      `json:"reason,omitempty"`
	Clashes []string    `json:"clashes,omitempty"`
}

type CutoverInput struct {
	IncumbentRunning bool
	Ports            []string
	NetworkMode      string
	PIDMode          string
	HasDevices       bool
	Volumes          []string
	IncumbentMounts  []string
	SharedOptIn      map[string]bool
	RequiredCPUCores float64
	RequiredMemoryMB int64
	AvailableCPUs    float64
	AvailableMemMB   int64
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)

func HasHostPortBindings(ports []string) bool {
	for _, p := range ports {
		mapping := p
		if idx := strings.LastIndex(mapping, "/"); idx != -1 {
			mapping = mapping[:idx]
		}
		if strings.Contains(mapping, ":") {
			return true
		}
	}
	return false
}

func NewcomerName(deploymentID, serviceName string) string {
	cleanID := nonAlphanumericRegex.ReplaceAllString(deploymentID, "")
	if len(cleanID) > 16 {
		cleanID = cleanID[:16]
	}
	return "codedock-cutover-" + cleanID + "-" + serviceName
}

func StableName(slug, serviceName string) string {
	return "codedock-" + slug + "-" + serviceName
}

func normalizeMountTarget(target string) string {
	cleaned := path.Clean(target)
	if !strings.HasPrefix(cleaned, "/") {
		cleaned = "/" + cleaned
	}
	return cleaned
}

func parseVolumeMountTarget(spec string) (string, bool) {
	parts := strings.Split(spec, ":")
	if len(parts) == 1 {
		return normalizeMountTarget(parts[0]), false
	}
	target := parts[1]
	readOnly := false
	if len(parts) > 2 {
		readOnly = strings.Contains(parts[2], "ro")
	}
	return normalizeMountTarget(target), readOnly
}

func ClashingRWDestinations(incumbentMounts, candidateVolumes []string, sharedOptIn map[string]bool) []string {
	liveRwMap := make(map[string]bool)
	for _, m := range incumbentMounts {
		target, ro := parseVolumeMountTarget(m)
		if !ro {
			liveRwMap[target] = true
		}
	}

	clashes := make([]string, 0)
	seen := make(map[string]bool)
	for _, spec := range candidateVolumes {
		target, ro := parseVolumeMountTarget(spec)
		if ro {
			continue
		}
		if liveRwMap[target] {
			if sharedOptIn != nil && sharedOptIn[target] {
				continue
			}
			if !seen[target] {
				seen[target] = true
				clashes = append(clashes, target)
			}
		}
	}
	return clashes
}

func DetermineCutoverMode(opts CutoverInput) CutoverDecision {
	if !opts.IncumbentRunning {
		return CutoverDecision{Mode: CutoverModeReplace, Reason: "no-incumbent"}
	}
	if HasHostPortBindings(opts.Ports) {
		return CutoverDecision{Mode: CutoverModeReplace, Reason: "host-ports"}
	}
	if opts.NetworkMode != "" || opts.PIDMode != "" {
		return CutoverDecision{Mode: CutoverModeReplace, Reason: "shared-namespaces"}
	}
	if opts.HasDevices {
		return CutoverDecision{Mode: CutoverModeReplace, Reason: "devices"}
	}

	if len(opts.Volumes) > 0 && len(opts.IncumbentMounts) > 0 {
		clashes := ClashingRWDestinations(opts.IncumbentMounts, opts.Volumes, opts.SharedOptIn)
		if len(clashes) > 0 {
			return CutoverDecision{
				Mode:    CutoverModeReplace,
				Reason:  "shared-rw-volume",
				Clashes: clashes,
			}
		}
	}

	if opts.RequiredCPUCores > 0 && opts.AvailableCPUs > 0 && opts.RequiredCPUCores > opts.AvailableCPUs {
		return CutoverDecision{Mode: CutoverModeReplace, Reason: "insufficient-capacity"}
	}
	if opts.RequiredMemoryMB > 0 && opts.AvailableMemMB > 0 && opts.RequiredMemoryMB > opts.AvailableMemMB {
		return CutoverDecision{Mode: CutoverModeReplace, Reason: "insufficient-capacity"}
	}

	return CutoverDecision{Mode: CutoverModeOverlap}
}
