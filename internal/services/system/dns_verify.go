package system

import (
	"context"
	"fmt"
	"net"
)

func VerifyDomain(ctx context.Context, domainName, expectedIP string) (bool, string, error) {
	if domainName == "" {
		return false, "", fmt.Errorf("domainName is required")
	}

	resolver := &net.Resolver{}
	ips, err := resolver.LookupIP(ctx, "ip", domainName)
	if err != nil {
		return false, "", fmt.Errorf("failed to lookup IP for domain %s: %w", domainName, err)
	}

	if len(ips) == 0 {
		return false, "", nil
	}

	resolvedIP := ips[0].String()
	if expectedIP == "" {
		return true, resolvedIP, nil
	}

	for _, ip := range ips {
		if ip.String() == expectedIP {
			return true, resolvedIP, nil
		}
	}

	return false, resolvedIP, nil
}
