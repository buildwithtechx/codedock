package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
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

func TestDeviceFlowHandler_CreateAndPoll(t *testing.T) {
	e := echo.New()
	svc := NewDeviceFlowService()
	h := NewDeviceFlowHandler(svc, nil, nil)

	reqBody := `{"clientName":"Test Client"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/device/code", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.CreateDeviceCode(c); err != nil {
		t.Fatalf("CreateDeviceCode failed: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var createResp struct {
		Data struct {
			DeviceCode string `json:"deviceCode"`
			UserCode   string `json:"userCode"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if createResp.Data.DeviceCode == "" || createResp.Data.UserCode == "" {
		t.Fatalf("expected non-empty device and user codes")
	}

	pollBody, _ := json.Marshal(map[string]string{"deviceCode": createResp.Data.DeviceCode})
	pollReq := httptest.NewRequest(http.MethodPost, "/api/auth/device/token", bytes.NewBuffer(pollBody))
	pollReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	pollRec := httptest.NewRecorder()
	pollCtx := e.NewContext(pollReq, pollRec)

	if err := h.PollDeviceToken(pollCtx); err != nil {
		t.Fatalf("PollDeviceToken failed: %v", err)
	}
	if pollRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", pollRec.Code)
	}

	var pollResp struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(pollRec.Body.Bytes(), &pollResp)
	if pollResp.Data.Status != "pending" {
		t.Fatalf("expected pending, got %s", pollResp.Data.Status)
	}

	verifyReq := httptest.NewRequest(http.MethodGet, "/api/auth/device/verify?code="+createResp.Data.UserCode, nil)
	verifyRec := httptest.NewRecorder()
	verifyCtx := e.NewContext(verifyReq, verifyRec)

	if err := h.VerifyUserCode(verifyCtx); err != nil {
		t.Fatalf("VerifyUserCode failed: %v", err)
	}
	if verifyRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", verifyRec.Code)
	}
}
