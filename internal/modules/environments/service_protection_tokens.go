package environments

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func CreatePolicy(password string, sessionSeconds int) (*RouteProtectionPolicy, error) {
	if len(password) < MinPasswordLen || len(password) > MaxPasswordLen {
		return nil, fmt.Errorf("route password must be between %d and %d characters", MinPasswordLen, MaxPasswordLen)
	}
	if sessionSeconds < 300 || sessionSeconds > 604800 {
		sessionSeconds = DefaultSessionSec
	}

	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}

	verifier := pbkdf2.Key([]byte(password), salt, PBKDF2Iterations, 32, sha256.New)

	signingKey := make([]byte, 32)
	if _, err := rand.Read(signingKey); err != nil {
		return nil, fmt.Errorf("generate signing key: %w", err)
	}

	version := make([]byte, 16)
	if _, err := rand.Read(version); err != nil {
		return nil, fmt.Errorf("generate version: %w", err)
	}

	return &RouteProtectionPolicy{
		Version:        hex.EncodeToString(version),
		Salt:           hex.EncodeToString(salt),
		Verifier:       hex.EncodeToString(verifier),
		SigningKey:     hex.EncodeToString(signingKey),
		Iterations:     PBKDF2Iterations,
		SessionSeconds: sessionSeconds,
	}, nil
}

func RevokeRouteSessions(policy *RouteProtectionPolicy) *RouteProtectionPolicy {
	signingKey := make([]byte, 32)
	_, _ = rand.Read(signingKey)
	version := make([]byte, 16)
	_, _ = rand.Read(version)

	return &RouteProtectionPolicy{
		Version:        hex.EncodeToString(version),
		Salt:           policy.Salt,
		Verifier:       policy.Verifier,
		SigningKey:     hex.EncodeToString(signingKey),
		Iterations:     policy.Iterations,
		SessionSeconds: policy.SessionSeconds,
	}
}

func VerifyPassword(password string, policy *RouteProtectionPolicy) bool {
	salt, err := hex.DecodeString(policy.Salt)
	if err != nil {
		return false
	}
	expectedVerifier, err := hex.DecodeString(policy.Verifier)
	if err != nil {
		return false
	}
	computed := pbkdf2.Key([]byte(password), salt, policy.Iterations, 32, sha256.New)
	return subtle.ConstantTimeCompare(computed, expectedVerifier) == 1
}

func sign(signingKeyHex, message string) string {
	key, err := hex.DecodeString(signingKeyHex)
	if err != nil {
		return ""
	}
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func SignSessionToken(policy *RouteProtectionPolicy, hostname string, expires int64) string {
	msg := fmt.Sprintf("codedock-route-v1\nsession\n%s\n%s\n%d", hostname, policy.Version, expires)
	return fmt.Sprintf("%d.%s", expires, sign(policy.SigningKey, msg))
}

func VerifySessionToken(policy *RouteProtectionPolicy, hostname, token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}
	expires, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return false
	}
	now := time.Now().Unix()
	if expires <= now || expires > now+int64(policy.SessionSeconds) {
		return false
	}
	expected := SignSessionToken(policy, hostname, expires)
	return subtle.ConstantTimeCompare([]byte(token), []byte(expected)) == 1
}

func SignCSRFToken(policy *RouteProtectionPolicy, hostname, path string, expires int64) string {
	msg := fmt.Sprintf("codedock-route-v1\nform\n%s\n%s\n%d\n%s", hostname, policy.Version, expires, path)
	return fmt.Sprintf("%d.%s", expires, sign(policy.SigningKey, msg))
}

func VerifyCSRFToken(policy *RouteProtectionPolicy, hostname, path, token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}
	expires, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return false
	}
	now := time.Now().Unix()
	if expires <= now || expires > now+600 {
		return false
	}
	expected := SignCSRFToken(policy, hostname, path, expires)
	return subtle.ConstantTimeCompare([]byte(token), []byte(expected)) == 1
}

func NormalizeHostname(hostname string) (string, error) {
	h := strings.TrimSpace(strings.ToLower(hostname))
	h = strings.TrimPrefix(h, "http://")
	h = strings.TrimPrefix(h, "https://")
	if colon := strings.Index(h, ":"); colon != -1 {
		h = h[:colon]
	}
	if slash := strings.Index(h, "/"); slash != -1 {
		h = h[:slash]
	}
	if h == "" || strings.HasPrefix(h, "*.") || strings.Contains(h, " ") {
		return "", errors.New("choose an exact hostname, without a wildcard, port, or URL scheme")
	}
	return h, nil
}
