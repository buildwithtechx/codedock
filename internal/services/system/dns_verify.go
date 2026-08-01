package system

import (
	"context"
	"errors"
	"fmt"
	"net"
)

func evaluateResolvedIPs(ips []net.IP, expectedIP string) (bool, string) {
	if len(ips) == 0 {
		return false, ""
	}

	if expectedIP == "" {
		return false, ips[0].String()
	}

	expected := net.ParseIP(expectedIP)
	if expected == nil {
		return false, ips[0].String()
	}
	for _, ip := range ips {
		if ip.Equal(expected) {
			return true, ip.String()
		}
	}

	return false, ips[0].String()
}

func VerifyDomain(ctx context.Context, domainName, expectedIP string) (bool, string, error) {
	if domainName == "" {
		return false, "", fmt.Errorf("domainName is required")
	}

	resolver := &net.Resolver{}
	network := "ip"
	if expected := net.ParseIP(expectedIP); expected != nil {
		if expected.To4() != nil {
			network = "ip4"
		} else {
			network = "ip6"
		}
	}
	ips, err := resolver.LookupIP(ctx, network, domainName)
	if err != nil {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) && dnsErr.IsNotFound {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to lookup IP for domain %s: %w", domainName, err)
	}
	verified, resolvedIP := evaluateResolvedIPs(ips, expectedIP)
	return verified, resolvedIP, nil
}
