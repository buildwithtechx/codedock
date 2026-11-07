package security

import (
	"fmt"
	"strings"
)

type ServiceEvidence struct {
	ServiceName     string            `json:"serviceName"`
	Image           string            `json:"image"`
	ImageDigest     string            `json:"imageDigest"`
	RuntimeVersions map[string]string `json:"runtimeVersions"`
	IsPinnedDigest  bool              `json:"isPinnedDigest"`
}

func NormalizeRepo(repo string) string {
	repo = strings.TrimSpace(repo)
	if repo == "" {
		return ""
	}

	// Remove tag or digest from repo if present
	if idx := strings.Index(repo, "@"); idx != -1 {
		repo = repo[:idx]
	}
	if idx := strings.LastIndex(repo, ":"); idx != -1 && !strings.Contains(repo[idx:], "/") {
		repo = repo[:idx]
	}

	repo = strings.TrimPrefix(repo, "https://")
	repo = strings.TrimPrefix(repo, "http://")

	if strings.HasPrefix(repo, "index.docker.io/") {
		repo = "docker.io/" + strings.TrimPrefix(repo, "index.docker.io/")
	}

	parts := strings.Split(repo, "/")
	if len(parts) == 1 {
		return "docker.io/library/" + parts[0]
	}
	if len(parts) == 2 && !strings.Contains(parts[0], ".") && !strings.Contains(parts[0], ":") {
		return "docker.io/" + repo
	}
	return repo
}

func extractImageTag(image string) string {
	if idx := strings.Index(image, "@"); idx != -1 {
		image = image[:idx]
	}
	idx := strings.LastIndex(image, ":")
	if idx == -1 || strings.Contains(image[idx:], "/") {
		return "latest"
	}
	return image[idx+1:]
}

func MatchAdvisories(advisories []SecurityAdvisory, services []ServiceEvidence) []AdvisoryFinding {
	var findings []AdvisoryFinding
	seen := make(map[string]bool)

	for _, adv := range advisories {
		for _, svc := range services {
			for _, target := range adv.Targets {
				dedupKey := fmt.Sprintf("%s:%s:%s", adv.ID, svc.ServiceName, target.Kind)
				if seen[dedupKey] {
					continue
				}

				switch target.Kind {
				case TargetKindImage:
					svcRepo := NormalizeRepo(svc.Image)
					targetRepo := NormalizeRepo(target.Repository)
					if svcRepo == "" || svcRepo != targetRepo {
						continue
					}

					// Check fixed digests first
					isFixed := false
					if svc.ImageDigest != "" {
						for _, fd := range target.FixedDigests {
							if strings.EqualFold(svc.ImageDigest, fd) {
								isFixed = true
								break
							}
						}
					}
					if isFixed {
						continue
					}

					confidence := ConfidencePossible
					matched := false

					if svc.ImageDigest != "" {
						for _, ad := range target.AffectedDigests {
							if strings.EqualFold(svc.ImageDigest, ad) {
								confidence = ConfidenceConfirmed
								matched = true
								break
							}
						}
					}

					if !matched {
						svcTag := extractImageTag(svc.Image)
						for _, at := range target.AffectedTags {
							if strings.EqualFold(svcTag, at) {
								matched = true
								break
							}
						}
					}

					if matched {
						seen[dedupKey] = true
						evidence := svc.Image
						if svc.ImageDigest != "" {
							evidence = fmt.Sprintf("%s (%s)", svc.Image, svc.ImageDigest)
						}
						canRefresh := adv.Action == ActionRefreshImages && !svc.IsPinnedDigest
						findings = append(findings, AdvisoryFinding{
							AdvisoryID:          adv.ID,
							Title:               adv.Title,
							Severity:            adv.Severity,
							Message:             adv.Message,
							URL:                 adv.URL,
							Action:              adv.Action,
							TargetKind:          target.Kind,
							Confidence:          confidence,
							ServiceName:         svc.ServiceName,
							Evidence:            evidence,
							CanRefresh:          canRefresh,
							ImageChangeRequired: svc.IsPinnedDigest,
						})
					}

				case TargetKindRuntime:
					if svc.RuntimeVersions == nil {
						continue
					}
					ver, exists := svc.RuntimeVersions[target.Runtime]
					if !exists || ver == "" {
						continue
					}

					if MatchesRange(ver, target.Affected) {
						seen[dedupKey] = true
						parsedVer, _ := ParseSemver(ver)
						upgradeRequired := false
						if target.Fixed != nil && *target.Fixed != "" {
							sameMajorCandidate := fmt.Sprintf("%d.999.999", parsedVer.Major)
							if !MatchesRange(sameMajorCandidate, *target.Fixed) {
								upgradeRequired = true
							}
						}

						canRefresh := adv.Action == ActionRefreshImages && !svc.IsPinnedDigest && !upgradeRequired
						findings = append(findings, AdvisoryFinding{
							AdvisoryID:             adv.ID,
							Title:                  adv.Title,
							Severity:               adv.Severity,
							Message:                adv.Message,
							URL:                    adv.URL,
							Action:                 adv.Action,
							TargetKind:             target.Kind,
							Confidence:             ConfidenceConfirmed,
							ServiceName:            svc.ServiceName,
							Evidence:               fmt.Sprintf("%s %s", target.Runtime, ver),
							CanRefresh:             canRefresh,
							RuntimeUpgradeRequired: upgradeRequired,
							ImageChangeRequired:    svc.IsPinnedDigest,
						})
					}
				}
			}
		}
	}

	return findings
}
