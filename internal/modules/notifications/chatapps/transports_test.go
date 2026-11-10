package chatapps_test

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"

	"codedock/internal/modules/notifications/chatapps/transports"
)

func TestSlackSignatureVerification(t *testing.T) {
	secret := "test_slack_signing_secret_123"
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	body := []byte("command=%2Fcodedock&text=help&user_id=U12345")

	sigBase := fmt.Sprintf("v0:%s:%s", timestamp, string(body))
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(sigBase))
	validSig := "v0=" + hex.EncodeToString(mac.Sum(nil))

	// Valid signature
	err := transports.VerifySlackSignature(secret, validSig, timestamp, body)
	if err != nil {
		t.Fatalf("expected valid signature, got error: %v", err)
	}

	// Tampered body
	err = transports.VerifySlackSignature(secret, validSig, timestamp, []byte("tampered"))
	if err == nil {
		t.Fatalf("expected error for tampered body, got nil")
	}

	// Replay attack (timestamp 10 minutes ago)
	oldTimestamp := strconv.FormatInt(time.Now().Unix()-600, 10)
	err = transports.VerifySlackSignature(secret, validSig, oldTimestamp, body)
	if err == nil {
		t.Fatalf("expected replay attack error, got nil")
	}
}

func TestDiscordSignatureVerification(t *testing.T) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate ed25519 key: %v", err)
	}
	pubHex := hex.EncodeToString(pubKey)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	body := []byte(`{"type":1}`)

	msg := append([]byte(timestamp), body...)
	sig := ed25519.Sign(privKey, msg)
	sigHex := hex.EncodeToString(sig)

	// Valid signature
	err = transports.VerifyDiscordSignature(pubHex, sigHex, timestamp, body)
	if err != nil {
		t.Fatalf("expected valid discord signature, got: %v", err)
	}

	// Invalid signature
	err = transports.VerifyDiscordSignature(pubHex, "badbadbad", timestamp, body)
	if err == nil {
		t.Fatalf("expected error for bad signature, got nil")
	}
}

func TestTelegramSecretVerification(t *testing.T) {
	secret := "my_telegram_webhook_secret_999"

	if err := transports.VerifyTelegramSecret(secret, secret); err != nil {
		t.Fatalf("expected matching secret to pass, got: %v", err)
	}

	if err := transports.VerifyTelegramSecret(secret, "wrong_secret"); err == nil {
		t.Fatalf("expected mismatching secret to fail, got nil")
	}
}
