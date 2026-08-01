package system

import (
	"context"
	"net"
	"testing"
)

func TestVerifyDomain_EmptyDomain(t *testing.T) {
	_, _, err := VerifyDomain(context.Background(), "", "127.0.0.1")
	if err == nil {
		t.Errorf("expected error for empty domainName, got nil")
	}
}

func TestEvaluateResolvedIPs_AllMatch(t *testing.T) {
	ips := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.1")}

	verified, resolvedIP := evaluateResolvedIPs(ips, "127.0.0.1")
	if !verified {
		t.Fatalf("expected verified to be true")
	}
	if resolvedIP != "127.0.0.1" {
		t.Fatalf("expected resolvedIP to match expected IP, got %q", resolvedIP)
	}
}

func TestEvaluateResolvedIPs_MixedRecords(t *testing.T) {
	ips := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}

	verified, resolvedIP := evaluateResolvedIPs(ips, "127.0.0.1")
	if !verified {
		t.Fatalf("expected verification to succeed when the expected address is present")
	}
	if resolvedIP != "127.0.0.1" {
		t.Fatalf("expected resolvedIP to report the matching address, got %q", resolvedIP)
	}
}
