package security

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Semver struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
}

func (s Semver) String() string {
	if s.Prerelease != "" {
		return fmt.Sprintf("%d.%d.%d-%s", s.Major, s.Minor, s.Patch, s.Prerelease)
	}
	return fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)
}

var semverPattern = regexp.MustCompile(`^v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:-([0-9A-Za-z.-]+))?`)

func ParseSemver(v string) (Semver, error) {
	v = strings.TrimSpace(v)
	matches := semverPattern.FindStringSubmatch(v)
	if matches == nil {
		return Semver{}, fmt.Errorf("invalid semver: %q", v)
	}

	major, _ := strconv.Atoi(matches[1])
	minor := 0
	if matches[2] != "" {
		minor, _ = strconv.Atoi(matches[2])
	}
	patch := 0
	if matches[3] != "" {
		patch, _ = strconv.Atoi(matches[3])
	}
	prerelease := ""
	if len(matches) > 4 {
		prerelease = matches[4]
	}

	return Semver{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: prerelease,
	}, nil
}

func CompareSemver(a, b Semver) int {
	if a.Major != b.Major {
		if a.Major < b.Major {
			return -1
		}
		return 1
	}
	if a.Minor != b.Minor {
		if a.Minor < b.Minor {
			return -1
		}
		return 1
	}
	if a.Patch != b.Patch {
		if a.Patch < b.Patch {
			return -1
		}
		return 1
	}
	// Normal version has higher precedence than pre-release version
	if a.Prerelease == "" && b.Prerelease != "" {
		return 1
	}
	if a.Prerelease != "" && b.Prerelease == "" {
		return -1
	}
	if a.Prerelease < b.Prerelease {
		return -1
	}
	if a.Prerelease > b.Prerelease {
		return 1
	}
	return 0
}

func MatchesConstraint(v Semver, op string, target Semver) bool {
	cmp := CompareSemver(v, target)
	switch op {
	case "<":
		return cmp < 0
	case "<=":
		return cmp <= 0
	case ">":
		return cmp > 0
	case ">=":
		return cmp >= 0
	case "=", "==":
		return cmp == 0
	case "!=":
		return cmp != 0
	default:
		return cmp == 0
	}
}

var constraintPattern = regexp.MustCompile(`^(<=|>=|<|>|=|!=)?\s*(v?\d+(?:\.\d+)*(?:-[0-9A-Za-z.-]+)?)`)

func MatchesSingleClause(v Semver, clause string) bool {
	clause = strings.TrimSpace(clause)
	if clause == "" || clause == "*" {
		return true
	}

	parts := strings.Fields(clause)
	for _, part := range parts {
		matches := constraintPattern.FindStringSubmatch(part)
		if matches == nil {
			continue
		}
		op := matches[1]
		targetStr := matches[2]
		target, err := ParseSemver(targetStr)
		if err != nil {
			return false
		}
		if !MatchesConstraint(v, op, target) {
			return false
		}
	}
	return true
}

func MatchesRange(vStr string, rangeStr string) bool {
	v, err := ParseSemver(vStr)
	if err != nil {
		return false
	}
	rangeStr = strings.TrimSpace(rangeStr)
	if rangeStr == "" || rangeStr == "*" {
		return true
	}

	// Clauses combined with || (OR)
	clauses := strings.Split(rangeStr, "||")
	for _, clause := range clauses {
		if MatchesSingleClause(v, clause) {
			return true
		}
	}
	return false
}
