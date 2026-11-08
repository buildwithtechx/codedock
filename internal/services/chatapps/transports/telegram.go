package transports

import (
	"crypto/subtle"
	"errors"
)

// VerifyTelegramSecret verifies the secret token configured in setWebhook.
func VerifyTelegramSecret(configuredSecret, incomingHeader string) error {
	if configuredSecret == "" {
		return nil // No secret configured, pass-through
	}
	if subtle.ConstantTimeCompare([]byte(configuredSecret), []byte(incomingHeader)) != 1 {
		return errors.New("invalid telegram secret token header")
	}
	return nil
}
