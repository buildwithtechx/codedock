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

func TestVerifyDomain_Localhost(t *testing.T) {
	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip", "localhost")
	if err != nil {
		t.Fatalf("unexpected lookup error: %v", err)
	}
	if len(ips) == 0 {
		t.Fatal("expected localhost to resolve to at least one IP")
	}

	verified, resolvedIP, err := VerifyDomain(context.Background(), "localhost", ips[0].String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolvedIP == "" {
		t.Errorf("expected non-empty resolvedIP for localhost")
	}
	if !verified {
		t.Errorf("expected verified to be true for localhost mapping")
	}

	ip := net.ParseIP(resolvedIP)
	if ip == nil || !ip.IsLoopback() {
		t.Errorf("expected localhost to resolve to a loopback IP, got %q", resolvedIP)
	}
}
