package system

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codedock/internal/models"

	"github.com/labstack/echo/v4"
)

func TestSupportHandler_Lifecycle(t *testing.T) {
	e := echo.New()
	handler := NewSupportHandler()
	claims := &models.UserClaims{
		UserID: "user-123",
		Email:  "testuser@codedock.local",
		Role:   "owner",
	}

	// 1. GetSession
	req := httptest.NewRequest(http.MethodGet, "/api/cloud/support/session", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", claims)

	if err := handler.GetSession(c); err != nil {
		t.Fatalf("unexpected GetSession error: %v", err)
	}
	var session models.CloudSupportSession
	if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
		t.Fatalf("failed to decode session: %v", err)
	}
	if session.Account == nil || session.Account.Email != "testuser@codedock.local" {
		t.Fatalf("expected account with email testuser@codedock.local, got %+v", session.Account)
	}

	// 2. CreateTicket
	inputBody, _ := json.Marshal(models.CloudSupportCustomerInput{
		RequestID: "req-1",
		Subject:   "Deployment failing to route",
		Message:   "Our ingress gateway is giving a 502 error after restart.",
		Category:  models.CloudSupportCategoryDeployment,
	})
	req = httptest.NewRequest(http.MethodPost, "/api/cloud/support/mine", bytes.NewReader(inputBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set("user", claims)

	if err := handler.CreateTicket(c); err != nil {
		t.Fatalf("unexpected CreateTicket error: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rec.Code)
	}
	var receipt models.CloudSupportReceipt
	if err := json.Unmarshal(rec.Body.Bytes(), &receipt); err != nil {
		t.Fatalf("failed to decode receipt: %v", err)
	}
	if receipt.ID == "" {
		t.Fatalf("expected ticket ID in receipt")
	}

	// 3. ListTickets
	req = httptest.NewRequest(http.MethodGet, "/api/cloud/support/mine", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set("user", claims)

	if err := handler.ListTickets(c); err != nil {
		t.Fatalf("unexpected ListTickets error: %v", err)
	}
	var list models.CloudSupportCustomerList
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("failed to decode ticket list: %v", err)
	}
	if len(list.Tickets) != 1 || list.Tickets[0].ID != receipt.ID {
		t.Fatalf("expected 1 ticket with id %s, got %+v", receipt.ID, list.Tickets)
	}

	// 4. GetTicket
	req = httptest.NewRequest(http.MethodGet, "/api/cloud/support/mine/"+receipt.ID, nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set("user", claims)
	c.SetParamNames("id")
	c.SetParamValues(receipt.ID)

	if err := handler.GetTicket(c); err != nil {
		t.Fatalf("unexpected GetTicket error: %v", err)
	}
	var detail models.CloudSupportCustomerDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &detail); err != nil {
		t.Fatalf("failed to decode ticket detail: %v", err)
	}
	if detail.Ticket.Subject != "Deployment failing to route" {
		t.Fatalf("expected subject 'Deployment failing to route', got %s", detail.Ticket.Subject)
	}

	// 5. ReplyTicket
	replyBody, _ := json.Marshal(models.CloudSupportCustomerReply{
		RequestID: "req-2",
		Message:   "Here is additional diagnostic log output.",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/cloud/support/mine/"+receipt.ID+"/replies", bytes.NewReader(replyBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set("user", claims)
	c.SetParamNames("id")
	c.SetParamValues(receipt.ID)

	if err := handler.ReplyTicket(c); err != nil {
		t.Fatalf("unexpected ReplyTicket error: %v", err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rec.Code)
	}
	var replyDetail models.CloudSupportCustomerDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &replyDetail); err != nil {
		t.Fatalf("failed to decode detail after reply: %v", err)
	}
	if len(replyDetail.Messages) != 1 || replyDetail.Messages[0].Body != "Here is additional diagnostic log output." {
		t.Fatalf("expected 1 message in detail, got %+v", replyDetail.Messages)
	}

	// 6. UpdateStatus
	statusBody, _ := json.Marshal(models.CloudSupportCustomerStatus{
		Status: models.CloudSupportStatusResolved,
	})
	req = httptest.NewRequest(http.MethodPatch, "/api/cloud/support/mine/"+receipt.ID, bytes.NewReader(statusBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set("user", claims)
	c.SetParamNames("id")
	c.SetParamValues(receipt.ID)

	if err := handler.UpdateStatus(c); err != nil {
		t.Fatalf("unexpected UpdateStatus error: %v", err)
	}
	var resolvedDetail models.CloudSupportCustomerDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &resolvedDetail); err != nil {
		t.Fatalf("failed to decode detail after status update: %v", err)
	}
	if resolvedDetail.Ticket.Status != models.CloudSupportStatusResolved {
		t.Fatalf("expected status resolved, got %s", resolvedDetail.Ticket.Status)
	}
}
