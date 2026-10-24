package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authservices "codedock/internal/services/auth"

	"github.com/labstack/echo/v4"
)

func TestDeviceFlowHandler_CreateAndPoll(t *testing.T) {
	e := echo.New()
	svc := authservices.NewDeviceFlowService()
	h := NewDeviceFlowHandler(svc, nil, nil)

	// 1. Create Device Code
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

	// 2. Poll pending
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

	// 3. Verify user code
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
