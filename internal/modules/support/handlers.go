package support

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

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type SupportHandler struct {
	mu       sync.RWMutex
	tickets  map[string]*CloudSupportTicketRecord
	messages map[string][]*CloudSupportMessageRecord
}

func NewSupportHandler() *SupportHandler {
	return &SupportHandler{
		tickets:  make(map[string]*CloudSupportTicketRecord),
		messages: make(map[string][]*CloudSupportMessageRecord),
	}
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
	userID, email := utils.ExtractClaims(c)
	if userID == "" {
		return c.JSON(http.StatusOK, CloudSupportSession{Account: nil})
	}

	name := email
	if parts := strings.Split(email, "@"); len(parts) > 0 {
		name = parts[0]
	}

	account := &CloudSupportAccountInfo{
		Key:   userID,
		ID:    userID,
		Name:  &name,
		Email: email,
	}

	return c.JSON(http.StatusOK, CloudSupportSession{Account: account})
}

func (h *SupportHandler) ListTickets(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	statusFilter := CloudSupportStatus(strings.TrimSpace(c.QueryParam("status")))
	search := strings.ToLower(strings.TrimSpace(c.QueryParam("search")))
	limit := 25
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}
	before := strings.TrimSpace(c.QueryParam("before"))

	h.mu.RLock()
	var matching []*CloudSupportTicketRecord
	for _, t := range h.tickets {
		if t.AccountID != userID {
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

	var results []CloudSupportCustomerTicket
	var nextCursor *string

	endIndex := startIndex + limit
	if endIndex > len(matching) {
		endIndex = len(matching)
	}

	for i := startIndex; i < endIndex; i++ {
		t := matching[i]
		results = append(results, CloudSupportCustomerTicket{
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
		results = []CloudSupportCustomerTicket{}
	}

	return c.JSON(http.StatusOK, CloudSupportCustomerList{
		Tickets:    results,
		NextCursor: nextCursor,
	})
}

func (h *SupportHandler) UpdateStatus(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	id := c.Param("id")

	var input CloudSupportCustomerStatus
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if input.Status != CloudSupportStatusOpen && input.Status != CloudSupportStatusResolved {
		return echo.NewHTTPError(http.StatusBadRequest, "Status must be 'open' or 'resolved'")
	}

	h.mu.Lock()
	record, exists := h.tickets[id]
	if !exists || record.AccountID != userID {
		h.mu.Unlock()
		return echo.NewHTTPError(http.StatusNotFound, "Ticket not found")
	}

	now := time.Now().UTC()
	record.Status = input.Status
	record.UpdatedAt = now
	msgs := h.messages[id]
	h.mu.Unlock()

	var messageItems []CloudSupportCustomerMessageItem
	for _, m := range msgs {
		messageItems = append(messageItems, CloudSupportCustomerMessageItem{
			ID:        m.ID,
			Author:    m.Author,
			Body:      m.Body,
			CreatedAt: m.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	if messageItems == nil {
		messageItems = []CloudSupportCustomerMessageItem{}
	}

	return c.JSON(http.StatusOK, CloudSupportCustomerDetail{
		Ticket: CloudSupportCustomerTicketDetail{
			CloudSupportCustomerTicket: CloudSupportCustomerTicket{
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
