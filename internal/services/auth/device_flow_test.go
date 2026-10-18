package auth

import (
	"errors"
	"testing"
	"time"
)

func TestDeviceFlowService_Lifecycle(t *testing.T) {
	svc := NewDeviceFlowService()

	session, err := svc.CreateSession("Codedock Mobile")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if session.DeviceCode == "" || session.UserCode == "" {
		t.Fatalf("expected non-empty codes")
	}
	if session.Status != DeviceStatusPending {
		t.Fatalf("expected status pending, got %v", session.Status)
	}
	if session.ClientName != "Codedock Mobile" {
		t.Fatalf("expected client name Codedock Mobile, got %s", session.ClientName)
	}

	byDev, err := svc.GetByDeviceCode(session.DeviceCode)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if byDev.UserCode != session.UserCode {
		t.Fatalf("user code mismatch")
	}

	byUser, err := svc.GetByUserCode(session.UserCode)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if byUser.DeviceCode != session.DeviceCode {
		t.Fatalf("device code mismatch")
	}

	// Authorize approve
	authorized, err := svc.Authorize(session.UserCode, "user-123", "alice@example.com", "Alice", "pat_secret_123", true)
	if err != nil {
		t.Fatalf("expected authorize to succeed, got %v", err)
	}
	if authorized.Status != DeviceStatusApproved {
		t.Fatalf("expected status approved, got %v", authorized.Status)
	}
	if authorized.Token != "pat_secret_123" {
		t.Fatalf("token mismatch")
	}
	if authorized.UserEmail != "alice@example.com" {
		t.Fatalf("email mismatch")
	}

	// Second authorize fails
	_, err = svc.Authorize(session.UserCode, "user-123", "alice@example.com", "Alice", "pat_secret_123", true)
	if !errors.Is(err, ErrDeviceAlreadyDecided) {
		t.Fatalf("expected ErrDeviceAlreadyDecided, got %v", err)
	}
}

func TestDeviceFlowService_Deny(t *testing.T) {
	svc := NewDeviceFlowService()

	session, err := svc.CreateSession("VS Code Extension")
	if err != nil {
		t.Fatalf("create session error: %v", err)
	}

	denied, err := svc.Authorize(session.UserCode, "", "", "", "", false)
	if err != nil {
		t.Fatalf("authorize error: %v", err)
	}
	if denied.Status != DeviceStatusDenied {
		t.Fatalf("expected status denied, got %v", denied.Status)
	}
	if denied.Token != "" {
		t.Fatalf("expected empty token on deny")
	}
}

func TestDeviceFlowService_Expiration(t *testing.T) {
	svc := NewDeviceFlowService()

	session, err := svc.CreateSession("CLI")
	if err != nil {
		t.Fatalf("create session error: %v", err)
	}

	session.ExpiresAt = time.Now().Add(-1 * time.Minute)

	_, err = svc.GetByDeviceCode(session.DeviceCode)
	if !errors.Is(err, ErrDeviceSessionExpired) {
		t.Fatalf("expected ErrDeviceSessionExpired, got %v", err)
	}

	_, err = svc.GetByUserCode(session.UserCode)
	if !errors.Is(err, ErrDeviceSessionExpired) {
		t.Fatalf("expected ErrDeviceSessionExpired, got %v", err)
	}
}
