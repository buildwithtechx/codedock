package chatapps

import (
	chatappservices "codedock/internal/services/chatapps"
)

// Handler handles HTTP requests for Chat App connections, identities, pairing, and webhooks
type Handler struct {
	service    *chatappservices.Service
	dispatcher *chatappservices.Dispatcher
}

// NewHandler creates a new chat apps handler
func NewHandler(service *chatappservices.Service, dispatcher *chatappservices.Dispatcher) *Handler {
	return &Handler{
		service:    service,
		dispatcher: dispatcher,
	}
}
