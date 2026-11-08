package security

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type ManagedImageInfo struct {
	Runtime    string `json:"runtime"`    // "node" or "bun"
	Repository string `json:"repository"` // "node" or "oven/bun"
	Version    string `json:"version"`    // e.g. "22.22.0"
	Major      int    `json:"major"`
	Variant    string `json:"variant"` // e.g. "-alpine", "-bookworm-slim"
}

type PatchedRuntimeImages struct {
	TargetVersion string `json:"targetVersion"`
	BuildImage     string `json:"buildImage"`
	RuntimeImage   string `json:"runtimeImage"`
}

type SecurityRuntimePatch struct {
	ReviewID       string `json:"reviewId"`
	Runtime        string `json:"runtime"`
	CurrentVersion string `json:"currentVersion"`
	TargetVersion  string `json:"targetVersion"`
	BuildImage     string `json:"buildImage"`
	RuntimeImage   string `json:"runtimeImage"`
}

var managedImageRefRegex = regexp.MustCompile(`^(?:(?:docker\.io|index\.docker\.io|registry-1\.docker\.io)/)?((?:library/)?node|oven/bun)(?::([^@]+))?(?:@sha256:[a-f0-9]{64})?$`)
var managedTagRegex = regexp.MustCompile(`^(\d+(?:\.\d+){0,2}|latest|lts|current)((?:-(?:alpine(?:\d+(?:\.\d+)*)?|slim|bookworm|bullseye|trixie|debian|distroless))*)$`)

// ManagedRuntimeImage parses an image reference and returns metadata if it is an official Node or Bun image.
func ManagedRuntimeImage(ref string) *ManagedImageInfo {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil
	}

	match := managedImageRefRegex.FindStringSubmatch(ref)
	if match == nil {
		return nil
	}

	repoToken := match[1]
	runtime := "node"
	repository := "node"
	if repoToken == "oven/bun" {
		runtime = "bun"
		repository = "oven/bun"
	}

	tag := match[2]
	if tag == "" {
		tag = "latest"
	}

	tagMatch := managedTagRegex.FindStringSubmatch(tag)
	if tagMatch == nil {
		return nil
	}

	baseVersionStr := tagMatch[1]
	variant := tagMatch[2]

	if !regexp.MustCompile(`^\d`).MatchString(baseVersionStr) {
		// Non-version tag like "latest", "lts", "current"
		return nil
	}

	parsed, err := ParseSemver(baseVersionStr)
	if err != nil {
		return nil
	}

	return &ManagedImageInfo{
		Runtime:    runtime,
		Repository: repository,
		Version:    parsed.String(),
		Major:      parsed.Major,
		Variant:    variant,
	}
}

// SecurityRuntimeImages computes the minimum patched version satisfying fixed constraints on the same major.
func SecurityRuntimeImages(runtime, currentVersion, fixed, buildImage, runtimeImage string) (*PatchedRuntimeImages, error) {
	runtime = strings.ToLower(strings.TrimSpace(runtime))
	if runtime != "node" && runtime != "bun" {
		return nil, fmt.Errorf("unsupported runtime %q for security runtime patching", runtime)
	}

	currentParsed, err := ParseSemver(currentVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid current runtime version %q: %w", currentVersion, err)
	}
	currentMajor := currentParsed.Major

	builderInfo := ManagedRuntimeImage(buildImage)
	if builderInfo == nil {
		return nil, errors.New("build image is not an official managed image or cannot be automatically patched")
	}

	runtimeInfo := ManagedRuntimeImage(runtimeImage)
	if runtimeInfo == nil {
		return nil, errors.New("runtime image is not an official managed image or cannot be automatically patched")
	}

	if builderInfo.Runtime != runtime || runtimeInfo.Runtime != runtime {
		return nil, fmt.Errorf("image runtime (%s/%s) does not match advisory runtime %s", builderInfo.Runtime, runtimeInfo.Runtime, runtime)
	}

	if builderInfo.Major != currentMajor || runtimeInfo.Major != currentMajor {
		return nil, fmt.Errorf("image major version does not match current major %d", currentMajor)
	}

	// Calculate floor: never downgrade below max(currentVersion, builder.version, runtime.version)
	floor := currentParsed
	if bVer, err := ParseSemver(builderInfo.Version); err == nil && CompareSemver(bVer, floor) > 0 {
		floor = bVer
	}
	if rVer, err := ParseSemver(runtimeInfo.Version); err == nil && CompareSemver(rVer, floor) > 0 {
		floor = rVer
	}

	// If floor already satisfies fixed constraint, no version bump needed
	if MatchesRange(floor.String(), fixed) {
		return &PatchedRuntimeImages{
			TargetVersion: floor.String(),
			BuildImage:     fmt.Sprintf("%s:%s%s", builderInfo.Repository, floor.String(), builderInfo.Variant),
			RuntimeImage:   fmt.Sprintf("%s:%s%s", runtimeInfo.Repository, floor.String(), runtimeInfo.Variant),
		}, nil
	}

	// Find candidates on the same major version satisfying the fixed constraint
	var candidates []Semver
	clauses := strings.Split(fixed, "||")
	for _, clause := range clauses {
		clause = strings.TrimSpace(clause)
		if clause == "" {
			continue
		}

		parts := strings.Fields(clause)
		for _, part := range parts {
			matches := constraintPattern.FindStringSubmatch(part)
			if matches == nil {
				continue
			}
			targetStr := matches[2]
			cand, err := ParseSemver(targetStr)
			if err != nil {
				continue
			}
			// Only consider candidates on currentMajor, non-prerelease, >= floor, satisfying fixed constraint
			if cand.Major == currentMajor && cand.Prerelease == "" && CompareSemver(cand, floor) >= 0 {
				if MatchesRange(cand.String(), fixed) {
					candidates = append(candidates, cand)
				}
			}
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no patch available on major version %d satisfying %q; major runtime upgrade required", currentMajor, fixed)
	}

	// Sort candidates in ascending order and pick the minimum
	sort.Slice(candidates, func(i, j int) bool {
		return CompareSemver(candidates[i], candidates[j]) < 0
	})
	selected := candidates[0]

	targetVersion := selected.String()
	newBuildImage := fmt.Sprintf("%s:%s%s", builderInfo.Repository, targetVersion, builderInfo.Variant)
	newRuntimeImage := fmt.Sprintf("%s:%s%s", runtimeInfo.Repository, targetVersion, runtimeInfo.Variant)

	return &PatchedRuntimeImages{
		TargetVersion: targetVersion,
		BuildImage:     newBuildImage,
		RuntimeImage:   newRuntimeImage,
	}, nil
}

type reviewDigestPayload struct {
	ProjectID      string `json:"projectId"`
	DeploymentID   string `json:"deploymentId"`
	AdvisoryID     string `json:"advisoryId"`
	Revision       int64  `json:"revision"`
	CurrentVersion string `json:"currentVersion"`
	TargetVersion  string `json:"targetVersion"`
	BuildImage     string `json:"buildImage"`
	RuntimeImage   string `json:"runtimeImage"`
}

func calculateReviewID(payload reviewDigestPayload) string {
	bytes, _ := json.Marshal(payload)
	sum := sha256.Sum256(bytes)
	return hex.EncodeToString(sum[:])
}

// PlanSecurityRuntimePatch produces a deterministic SecurityRuntimePatch with a verifiable ReviewID.
func PlanSecurityRuntimePatch(
	projectID, deploymentID, advisoryID string,
	revision int64,
	runtime, currentVersion, fixed, buildImage, runtimeImage string,
) (*SecurityRuntimePatch, error) {
	patched, err := SecurityRuntimeImages(runtime, currentVersion, fixed, buildImage, runtimeImage)
	if err != nil {
		return nil, err
	}

	reviewID := calculateReviewID(reviewDigestPayload{
		ProjectID:      projectID,
		DeploymentID:   deploymentID,
		AdvisoryID:     advisoryID,
		Revision:       revision,
		CurrentVersion: currentVersion,
		TargetVersion:  patched.TargetVersion,
		BuildImage:     patched.BuildImage,
		RuntimeImage:   patched.RuntimeImage,
	})

	return &SecurityRuntimePatch{
		ReviewID:       reviewID,
		Runtime:        runtime,
		CurrentVersion: currentVersion,
		TargetVersion:  patched.TargetVersion,
		BuildImage:     patched.BuildImage,
		RuntimeImage:   patched.RuntimeImage,
	}, nil
}

// ValidateReviewID verifies that the patch ReviewID matches the current state.
func ValidateReviewID(
	patch *SecurityRuntimePatch,
	projectID, deploymentID, advisoryID string,
	revision int64,
) bool {
	if patch == nil || patch.ReviewID == "" {
		return false
	}

	expected := calculateReviewID(reviewDigestPayload{
		ProjectID:      projectID,
		DeploymentID:   deploymentID,
		AdvisoryID:     advisoryID,
		Revision:       revision,
		CurrentVersion: patch.CurrentVersion,
		TargetVersion:  patch.TargetVersion,
		BuildImage:     patch.BuildImage,
		RuntimeImage:   patch.RuntimeImage,
	})

	return patch.ReviewID == expected
}
