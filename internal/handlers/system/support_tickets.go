package system

import (
	"net/http"
	"strings"
	"time"

	"codedock/internal/models"

	"github.com/labstack/echo/v4"
)

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
