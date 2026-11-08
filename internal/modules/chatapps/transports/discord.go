package transports

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
)

// VerifyDiscordSignature verifies incoming Discord interactions via Ed25519.
func VerifyDiscordSignature(publicKeyHex, signatureHex, timestamp string, body []byte) error {
	if publicKeyHex == "" {
		return errors.New("discord public key is not configured")
	}
	if signatureHex == "" || timestamp == "" {
		return errors.New("missing discord signature or timestamp headers")
	}

	pubKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil || len(pubKeyBytes) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid ed25519 public key: %w", err)
	}

	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil || len(sigBytes) != ed25519.SignatureSize {
		return fmt.Errorf("invalid ed25519 signature format: %w", err)
	}

	msg := append([]byte(timestamp), body...)
	if !ed25519.Verify(pubKeyBytes, msg, sigBytes) {
		return errors.New("invalid discord interaction signature")
	}
	return nil
}
