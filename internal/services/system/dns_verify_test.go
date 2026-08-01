package system

import (
	"context"
	"testing"
)

func TestVerifyDomain_EmptyDomain(t *testing.T) {
	_, _, err := VerifyDomain(context.Background(), "", "127.0.0.1")
	if err == nil {
		t.Errorf("expected error for empty domainName, got nil")
	}
}

func TestVerifyDomain_Localhost(t *testing.T) {
	verified, resolvedIP, err := VerifyDomain(context.Background(), "localhost", "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolvedIP == "" {
		t.Errorf("expected non-empty resolvedIP for localhost")
	}
	if !verified && (resolvedIP == "127.0.0.1" || resolvedIP == "::1") {
		t.Errorf("expected verified to be true for localhost mapping")
	}
}
