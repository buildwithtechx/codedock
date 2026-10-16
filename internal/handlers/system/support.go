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

func (h *SupportHandler) CreateTicket(c echo.Context) error {
	claims, ok := c.Get("user").(*models.UserClaims)
	if !ok || claims == nil || claims.UserID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	var input models.CloudSupportCustomerInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	subject := strings.TrimSpace(input.Subject)
	message := strings.TrimSpace(input.Message)
	if subject == "" || message == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Subject and message are required")
	}

	category := input.Category
	if category == "" {
		category = models.CloudSupportCategoryGeneral
	}

	now := time.Now().UTC()
	ticketID := generateSupportID()

	record := &models.CloudSupportTicketRecord{
		ID:        ticketID,
		AccountID: claims.UserID,
		Subject:   subject,
		Category:  category,
		Status:    models.CloudSupportStatusOpen,
		Message:   message,
		CreatedAt: now,
		UpdatedAt: now,
	}

	h.mu.Lock()
	h.tickets[ticketID] = record
	h.messages[ticketID] = []*models.CloudSupportMessageRecord{}
	h.mu.Unlock()

	return c.JSON(http.StatusCreated, models.CloudSupportReceipt{
		ID:        ticketID,
		CreatedAt: now.Format(time.RFC3339),
	})
}

func (h *SupportHandler) GetTicket(c echo.Context) error {
	claims, ok := c.Get("user").(*models.UserClaims)
	if !ok || claims == nil || claims.UserID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	id := c.Param("id")

	h.mu.RLock()
	record, exists := h.tickets[id]
	msgs := h.messages[id]
	h.mu.RUnlock()

	if !exists || record.AccountID != claims.UserID {
		return echo.NewHTTPError(http.StatusNotFound, "Ticket not found")
	}

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

	detail := models.CloudSupportCustomerDetail{
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
	}

	return c.JSON(http.StatusOK, detail)
}

func (h *SupportHandler) ReplyTicket(c echo.Context) error {
	claims, ok := c.Get("user").(*models.UserClaims)
	if !ok || claims == nil || claims.UserID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	id := c.Param("id")

	var input models.CloudSupportCustomerReply
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	body := strings.TrimSpace(input.Message)
	if body == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Message body cannot be empty")
	}

	h.mu.Lock()
	record, exists := h.tickets[id]
	if !exists || record.AccountID != claims.UserID {
		h.mu.Unlock()
		return echo.NewHTTPError(http.StatusNotFound, "Ticket not found")
	}

	now := time.Now().UTC()
	record.UpdatedAt = now
	msgRecord := &models.CloudSupportMessageRecord{
		ID:        generateMessageID(),
		TicketID:  id,
		Author:    models.CloudSupportAuthorCustomer,
		Body:      body,
		CreatedAt: now,
	}
	h.messages[id] = append(h.messages[id], msgRecord)
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

	return c.JSON(http.StatusCreated, models.CloudSupportCustomerDetail{
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

func (h *SupportHandler) ListLegacyTickets(c echo.Context) error {
	claims, ok := c.Get("user").(*models.UserClaims)
	if !ok || claims == nil || claims.UserID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	h.mu.RLock()
	var matching []*models.CloudSupportTicketRecord
	for _, t := range h.tickets {
		if t.AccountID == claims.UserID {
			matching = append(matching, t)
		}
	}
	h.mu.RUnlock()

	sort.Slice(matching, func(i, j int) bool {
		return matching[i].CreatedAt.After(matching[j].CreatedAt)
	})

	type legacyMessage struct {
		ID         string `json:"id"`
		Sender     string `json:"sender"`
		SenderName string `json:"senderName"`
		Content    string `json:"content"`
		CreatedAt  string `json:"createdAt"`
	}
	type legacyTicket struct {
		ID        string          `json:"id"`
		Subject   string          `json:"subject"`
		Category  string          `json:"category"`
		Status    string          `json:"status"`
		Messages  []legacyMessage `json:"messages"`
		CreatedAt string          `json:"createdAt"`
		UpdatedAt string          `json:"updatedAt"`
	}

	var tickets []legacyTicket
	for _, t := range matching {
		h.mu.RLock()
		msgs := h.messages[t.ID]
		h.mu.RUnlock()

		legacyMsgs := []legacyMessage{
			{
				ID:         "msg-0",
				Sender:     "user",
				SenderName: claims.Email,
				Content:    t.Message,
				CreatedAt:  t.CreatedAt.UTC().Format(time.RFC3339),
			},
		}
		for _, m := range msgs {
			sender := "staff"
			if m.Author == models.CloudSupportAuthorCustomer {
				sender = "user"
			}
			legacyMsgs = append(legacyMsgs, legacyMessage{
				ID:         m.ID,
				Sender:     sender,
				SenderName: string(m.Author),
				Content:    m.Body,
				CreatedAt:  m.CreatedAt.UTC().Format(time.RFC3339),
			})
		}

		tickets = append(tickets, legacyTicket{
			ID:        t.ID,
			Subject:   t.Subject,
			Category:  string(t.Category),
			Status:    string(t.Status),
			Messages:  legacyMsgs,
			CreatedAt: t.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: t.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}
	if tickets == nil {
		tickets = []legacyTicket{}
	}

	return c.JSON(http.StatusOK, echo.Map{"tickets": tickets})
}

func (h *SupportHandler) CreateLegacyTicket(c echo.Context) error {
	claims, ok := c.Get("user").(*models.UserClaims)
	if !ok || claims == nil || claims.UserID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	var input struct {
		Subject           string `json:"subject"`
		Category          string `json:"category"`
		Message           string `json:"message"`
		AttachDiagnostics bool   `json:"attachDiagnostics"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	now := time.Now().UTC()
	ticketID := generateSupportID()
	record := &models.CloudSupportTicketRecord{
		ID:        ticketID,
		AccountID: claims.UserID,
		Subject:   input.Subject,
		Category:  models.CloudSupportCategory(input.Category),
		Status:    models.CloudSupportStatusOpen,
		Message:   input.Message,
		CreatedAt: now,
		UpdatedAt: now,
	}

	h.mu.Lock()
	h.tickets[ticketID] = record
	h.messages[ticketID] = []*models.CloudSupportMessageRecord{}
	h.mu.Unlock()

	return c.JSON(http.StatusCreated, echo.Map{
		"ticket": echo.Map{
			"id":        record.ID,
			"subject":   record.Subject,
			"category":  record.Category,
			"status":    record.Status,
			"messages": []echo.Map{
				{
					"id":         "msg-0",
					"sender":     "user",
					"senderName": claims.Email,
					"content":    record.Message,
					"createdAt":  now.Format(time.RFC3339),
				},
			},
			"createdAt": now.Format(time.RFC3339),
			"updatedAt": now.Format(time.RFC3339),
		},
	})
}

func (h *SupportHandler) GetLegacyTicket(c echo.Context) error {
	claims, ok := c.Get("user").(*models.UserClaims)
	if !ok || claims == nil || claims.UserID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	id := c.Param("id")
	h.mu.RLock()
	record, exists := h.tickets[id]
	msgs := h.messages[id]
	h.mu.RUnlock()

	if !exists || record.AccountID != claims.UserID {
		return echo.NewHTTPError(http.StatusNotFound, "Ticket not found")
	}

	legacyMsgs := []echo.Map{
		{
			"id":         "msg-0",
			"sender":     "user",
			"senderName": claims.Email,
			"content":    record.Message,
			"createdAt":  record.CreatedAt.UTC().Format(time.RFC3339),
		},
	}
	for _, m := range msgs {
		sender := "staff"
		if m.Author == models.CloudSupportAuthorCustomer {
			sender = "user"
		}
		legacyMsgs = append(legacyMsgs, echo.Map{
			"id":         m.ID,
			"sender":     sender,
			"senderName": string(m.Author),
			"content":    m.Body,
			"createdAt":  m.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"ticket": echo.Map{
			"id":        record.ID,
			"subject":   record.Subject,
			"category":  record.Category,
			"status":    record.Status,
			"messages":  legacyMsgs,
			"createdAt": record.CreatedAt.UTC().Format(time.RFC3339),
			"updatedAt": record.UpdatedAt.UTC().Format(time.RFC3339),
		},
	})
}

func (h *SupportHandler) ReplyLegacyTicket(c echo.Context) error {
	claims, ok := c.Get("user").(*models.UserClaims)
	if !ok || claims == nil || claims.UserID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	id := c.Param("id")
	var input struct {
		Message string `json:"message"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	h.mu.Lock()
	record, exists := h.tickets[id]
	if !exists || record.AccountID != claims.UserID {
		h.mu.Unlock()
		return echo.NewHTTPError(http.StatusNotFound, "Ticket not found")
	}

	now := time.Now().UTC()
	record.UpdatedAt = now
	msgRecord := &models.CloudSupportMessageRecord{
		ID:        generateMessageID(),
		TicketID:  id,
		Author:    models.CloudSupportAuthorCustomer,
		Body:      input.Message,
		CreatedAt: now,
	}
	h.messages[id] = append(h.messages[id], msgRecord)
	msgs := h.messages[id]
	h.mu.Unlock()

	legacyMsgs := []echo.Map{
		{
			"id":         "msg-0",
			"sender":     "user",
			"senderName": claims.Email,
			"content":    record.Message,
			"createdAt":  record.CreatedAt.UTC().Format(time.RFC3339),
		},
	}
	for _, m := range msgs {
		sender := "staff"
		if m.Author == models.CloudSupportAuthorCustomer {
			sender = "user"
		}
		legacyMsgs = append(legacyMsgs, echo.Map{
			"id":         m.ID,
			"sender":     sender,
			"senderName": string(m.Author),
			"content":    m.Body,
			"createdAt":  m.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"ticket": echo.Map{
			"id":        record.ID,
			"subject":   record.Subject,
			"category":  record.Category,
			"status":    record.Status,
			"messages":  legacyMsgs,
			"createdAt": record.CreatedAt.UTC().Format(time.RFC3339),
			"updatedAt": record.UpdatedAt.UTC().Format(time.RFC3339),
		},
	})
}
