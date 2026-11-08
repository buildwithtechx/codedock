package transports

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

// VerifySlackSignature verifies an incoming Slack request using the Slack signing secret.
// Formula: hex(HMAC-SHA256("v0:" + timestamp + ":" + body, signingSecret)) == signature
func VerifySlackSignature(signingSecret, signature, timestamp string, body []byte) error {
	if signingSecret == "" {
		return errors.New("slack signing secret is not configured")
	}
	if signature == "" || timestamp == "" {
		return errors.New("missing slack signature or timestamp headers")
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %w", err)
	}

	// Protect against replay attacks (reject timestamps older than 5 minutes)
	now := time.Now().Unix()
	if math.Abs(float64(now-ts)) > 300 {
		return errors.New("slack request timestamp out of bounds (replay protection)")
	}

	sigBase := fmt.Sprintf("v0:%s:%s", timestamp, string(body))
	mac := hmac.New(sha256.New, []byte(signingSecret))
	mac.Write([]byte(sigBase))
	expectedSig := "v0=" + hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return errors.New("invalid slack signature")
	}
	return nil
}
