package system

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"codedock/internal/models"

	"github.com/labstack/echo/v4"
)

type SupportHandler struct {
	mu       sync.RWMutex
	tickets  map[string]*models.CloudSupportTicketRecord
	messages map[string][]*models.CloudSupportMessageRecord
}

func NewSupportHandler() *SupportHandler {
	h := &SupportHandler{
		tickets:  make(map[string]*models.CloudSupportTicketRecord),
		messages: make(map[string][]*models.CloudSupportMessageRecord),
	}
	return h
}

func generateSupportID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return fmt.Sprintf("SUP-%s", strings.ToUpper(hex.EncodeToString(b)))
}

func generateMessageID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("msg-%s", hex.EncodeToString(b))
}

func (h *SupportHandler) GetSession(c echo.Context) error {
	claims, ok := c.Get("user").(*models.UserClaims)
	if !ok || claims == nil || claims.UserID == "" {
		return c.JSON(http.StatusOK, models.CloudSupportSession{Account: nil})
	}

	name := claims.Email
	if parts := strings.Split(claims.Email, "@"); len(parts) > 0 {
		name = parts[0]
	}

	account := &models.CloudSupportAccountInfo{
		Key:   claims.UserID,
		ID:    claims.UserID,
		Name:  &name,
		Email: claims.Email,
	}

	return c.JSON(http.StatusOK, models.CloudSupportSession{Account: account})
}

func (h *SupportHandler) ListTickets(c echo.Context) error {
	claims, ok := c.Get("user").(*models.UserClaims)
	if !ok || claims == nil || claims.UserID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	statusFilter := models.CloudSupportStatus(strings.TrimSpace(c.QueryParam("status")))
	search := strings.ToLower(strings.TrimSpace(c.QueryParam("search")))
	limit := 25
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}
	before := strings.TrimSpace(c.QueryParam("before"))

	h.mu.RLock()
	var matching []*models.CloudSupportTicketRecord
	for _, t := range h.tickets {
		if t.AccountID != claims.UserID {
			continue
		}
		if statusFilter != "" && t.Status != statusFilter {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(t.Subject), search) && !strings.Contains(strings.ToLower(t.Message), search) {
			continue
		}
		matching = append(matching, t)
	}
	h.mu.RUnlock()

	sort.Slice(matching, func(i, j int) bool {
		return matching[i].CreatedAt.After(matching[j].CreatedAt)
	})

	startIndex := 0
	if before != "" {
		for i, t := range matching {
			if t.ID == before {
				startIndex = i + 1
				break
			}
		}
	}

	var results []models.CloudSupportCustomerTicket
	var nextCursor *string

	endIndex := startIndex + limit
	if endIndex > len(matching) {
		endIndex = len(matching)
	}

	for i := startIndex; i < endIndex; i++ {
		t := matching[i]
		results = append(results, models.CloudSupportCustomerTicket{
			ID:        t.ID,
			Subject:   t.Subject,
			Category:  t.Category,
			Status:    t.Status,
			CreatedAt: t.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: t.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	if endIndex < len(matching) {
		cursor := matching[endIndex-1].ID
		nextCursor = &cursor
	}

	if results == nil {
		results = []models.CloudSupportCustomerTicket{}
	}

	return c.JSON(http.StatusOK, models.CloudSupportCustomerList{
		Tickets:    results,
		NextCursor: nextCursor,
	})
}

func (h *SupportHandler) UpdateStatus(c echo.Context) error {
	claims, ok := c.Get("user").(*models.UserClaims)
	if !ok || claims == nil || claims.UserID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	id := c.Param("id")

	var input models.CloudSupportCustomerStatus
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if input.Status != models.CloudSupportStatusOpen && input.Status != models.CloudSupportStatusResolved {
		return echo.NewHTTPError(http.StatusBadRequest, "Status must be 'open' or 'resolved'")
	}

	h.mu.Lock()
	record, exists := h.tickets[id]
	if !exists || record.AccountID != claims.UserID {
		h.mu.Unlock()
		return echo.NewHTTPError(http.StatusNotFound, "Ticket not found")
	}

	now := time.Now().UTC()
	record.Status = input.Status
	record.UpdatedAt = now
	msgs := h.messages[id]
	h.mu.Unlock()

	var messageItems []models.CloudSupportCustomerMessageItem
	for _, m := range msgs {
		messageItems = append(messageItems, models.CloudSupportCustomerMessageItem{
			ID:        m.ID,
			Author:    m.Author,
			Body:      m.Body,
			CreatedAt: m.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	if messageItems == nil {
		messageItems = []models.CloudSupportCustomerMessageItem{}
	}

	return c.JSON(http.StatusOK, models.CloudSupportCustomerDetail{
		Ticket: models.CloudSupportCustomerTicketDetail{
			CloudSupportCustomerTicket: models.CloudSupportCustomerTicket{
				ID:        record.ID,
				Subject:   record.Subject,
				Category:  record.Category,
				Status:    record.Status,
				CreatedAt: record.CreatedAt.UTC().Format(time.RFC3339),
				UpdatedAt: record.UpdatedAt.UTC().Format(time.RFC3339),
			},
			Message: record.Message,
		},
		Messages: messageItems,
	})
}
