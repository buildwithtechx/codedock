package mail

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	jailNameRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$`)

	numberJailsRegex = regexp.MustCompile(`(?i)Number of jail:\s*(\d+)`)
	jailListRegex    = regexp.MustCompile(`(?i)Jail list:\s*([^\r\n]*)`)

	jailHeaderRegex    = regexp.MustCompile(`(?i)Status for the jail:\s*([^\r\n]+)`)
	currentFailedRegex = regexp.MustCompile(`(?i)Currently failed:\s*(\d+)`)
	totalFailedRegex   = regexp.MustCompile(`(?i)Total failed:\s*(\d+)`)
	currentBannedRegex = regexp.MustCompile(`(?i)Currently banned:\s*(\d+)`)
	totalBannedRegex   = regexp.MustCompile(`(?i)Total banned:\s*(\d+)`)
	bannedIPListRegex  = regexp.MustCompile(`(?i)Banned IP list:\s*([^\r\n]*)`)
)

func ValidateJailName(jail string) error {
	jail = strings.TrimSpace(jail)
	if !jailNameRegex.MatchString(jail) {
		return fmt.Errorf("invalid jail name: %q", jail)
	}
	return nil
}

func ValidateIP(ipStr string) error {
	ipStr = strings.TrimSpace(ipStr)
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %q", ipStr)
	}
	return nil
}

func ParseFail2banStatus(output string) ([]string, error) {
	numMatch := numberJailsRegex.FindStringSubmatch(output)
	if numMatch == nil {
		return nil, errors.New("invalid fail2ban status: missing 'Number of jail'")
	}
	expectedCount, err := strconv.Atoi(numMatch[1])
	if err != nil || expectedCount > 64 {
		return nil, fmt.Errorf("invalid fail2ban status count: %v", numMatch[1])
	}

	listMatch := jailListRegex.FindStringSubmatch(output)
	if listMatch == nil {
		return nil, errors.New("invalid fail2ban status: missing 'Jail list'")
	}

	rawList := strings.TrimSpace(listMatch[1])
	if rawList == "" {
		if expectedCount != 0 {
			return nil, fmt.Errorf("expected %d jails but list is empty", expectedCount)
		}
		return []string{}, nil
	}

	parts := strings.Split(rawList, ",")
	jails := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		if err := ValidateJailName(trimmed); err != nil {
			return nil, err
		}
		jails = append(jails, trimmed)
	}

	if len(jails) != expectedCount {
		return nil, fmt.Errorf("jail count mismatch: header reported %d but found %d", expectedCount, len(jails))
	}

	return jails, nil
}

func ParseFail2banJailStatus(jail string, output string) (*Fail2banJail, error) {
	headerMatch := jailHeaderRegex.FindStringSubmatch(output)
	if headerMatch == nil || strings.TrimSpace(headerMatch[1]) != jail {
		return nil, fmt.Errorf("output does not match expected jail %q", jail)
	}

	cfMatch := currentFailedRegex.FindStringSubmatch(output)
	if cfMatch == nil {
		return nil, errors.New("missing 'Currently failed' count")
	}
	currentFailed, _ := strconv.Atoi(cfMatch[1])

	tfMatch := totalFailedRegex.FindStringSubmatch(output)
	if tfMatch == nil {
		return nil, errors.New("missing 'Total failed' count")
	}
	totalFailed, _ := strconv.Atoi(tfMatch[1])

	cbMatch := currentBannedRegex.FindStringSubmatch(output)
	if cbMatch == nil {
		return nil, errors.New("missing 'Currently banned' count")
	}
	currentBanned, _ := strconv.Atoi(cbMatch[1])

	tbMatch := totalBannedRegex.FindStringSubmatch(output)
	if tbMatch == nil {
		return nil, errors.New("missing 'Total banned' count")
	}
	totalBanned, _ := strconv.Atoi(tbMatch[1])

	bannedListMatch := bannedIPListRegex.FindStringSubmatch(output)
	if bannedListMatch == nil {
		return nil, errors.New("missing 'Banned IP list'")
	}

	rawIPs := strings.TrimSpace(bannedListMatch[1])
	var bannedIPs []string
	if rawIPs != "" {
		fields := strings.Fields(rawIPs)
		for _, f := range fields {
			trimmed := strings.TrimSpace(f)
			if trimmed == "" {
				continue
			}
			if err := ValidateIP(trimmed); err != nil {
				return nil, err
			}
			bannedIPs = append(bannedIPs, trimmed)
		}
	} else {
		bannedIPs = []string{}
	}

	if len(bannedIPs) != currentBanned {
		return nil, fmt.Errorf("banned IP count mismatch: header reported %d but found %d", currentBanned, len(bannedIPs))
	}

	return &Fail2banJail{
		Jail:            jail,
		CurrentlyFailed: currentFailed,
		TotalFailed:     totalFailed,
		CurrentlyBanned: currentBanned,
		TotalBanned:     totalBanned,
		BannedIPs:       bannedIPs,
	}, nil
}

func FormatFail2banSnapshot(jails []string, details map[string]Fail2banJail, canManage bool) *Fail2banStatusResponse {
	return &Fail2banStatusResponse{
		Jails:     jails,
		Details:   details,
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
		CanManage: canManage,
	}
}
