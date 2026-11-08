package chatapps

import (
	"codedock/internal/models"
)

// Re-export or alias request models used by chatapps HTTP handlers
type CreateChatAppRequest = models.CreateChatAppRequest
type UpdateChatAppRequest = models.UpdateChatAppRequest
type VerifyPairingRequest = models.VerifyPairingRequest
type ChatCommandRequest = models.ChatCommandRequest
