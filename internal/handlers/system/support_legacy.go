package system

import (
	"net/http"
	"sort"
	"time"

	"codedock/internal/models"

	"github.com/labstack/echo/v4"
)

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
			"id":       record.ID,
			"subject":  record.Subject,
			"category": record.Category,
			"status":   record.Status,
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
