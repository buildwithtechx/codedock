package projects

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

const (
	RouteAuthCookieName = "codedock_route_session"
	DefaultSessionSec   = 86400
	MinPasswordLen      = 12
	MaxPasswordLen      = 256
	PBDKF2Iterations    = 600000
)

type RouteProtectionService struct {
	repo         repositories.RouteProtectionRepository
	projectRepo  repositories.ProjectRepository
	settingsRepo repositories.SettingsRepository
}

func NewRouteProtectionService(
	repo repositories.RouteProtectionRepository,
	projRepo repositories.ProjectRepository,
	settingsRepo repositories.SettingsRepository,
) *RouteProtectionService {
	return &RouteProtectionService{
		repo:         repo,
		projectRepo:  projRepo,
		settingsRepo: settingsRepo,
	}
}

func CreatePolicy(password string, sessionSeconds int) (*models.RouteProtectionPolicy, error) {
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

	verifier := pbkdf2.Key([]byte(password), salt, PBDKF2Iterations, 32, sha256.New)

	signingKey := make([]byte, 32)
	if _, err := rand.Read(signingKey); err != nil {
		return nil, fmt.Errorf("generate signing key: %w", err)
	}

	version := make([]byte, 16)
	if _, err := rand.Read(version); err != nil {
		return nil, fmt.Errorf("generate version: %w", err)
	}

	return &models.RouteProtectionPolicy{
		Version:        hex.EncodeToString(version),
		Salt:           hex.EncodeToString(salt),
		Verifier:       hex.EncodeToString(verifier),
		SigningKey:     hex.EncodeToString(signingKey),
		Iterations:     PBDKF2Iterations,
		SessionSeconds: sessionSeconds,
	}, nil
}

func RevokeRouteSessions(policy *models.RouteProtectionPolicy) *models.RouteProtectionPolicy {
	signingKey := make([]byte, 32)
	_, _ = rand.Read(signingKey)
	version := make([]byte, 16)
	_, _ = rand.Read(version)

	return &models.RouteProtectionPolicy{
		Version:        hex.EncodeToString(version),
		Salt:           policy.Salt,
		Verifier:       policy.Verifier,
		SigningKey:     hex.EncodeToString(signingKey),
		Iterations:     policy.Iterations,
		SessionSeconds: policy.SessionSeconds,
	}
}

func VerifyPassword(password string, policy *models.RouteProtectionPolicy) bool {
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

func SignSessionToken(policy *models.RouteProtectionPolicy, hostname string, expires int64) string {
	msg := fmt.Sprintf("codedock-route-v1\nsession\n%s\n%s\n%d", hostname, policy.Version, expires)
	return fmt.Sprintf("%d.%s", expires, sign(policy.SigningKey, msg))
}

func VerifySessionToken(policy *models.RouteProtectionPolicy, hostname, token string) bool {
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

func SignCSRFToken(policy *models.RouteProtectionPolicy, hostname, path string, expires int64) string {
	msg := fmt.Sprintf("codedock-route-v1\nform\n%s\n%s\n%d\n%s", hostname, policy.Version, expires, path)
	return fmt.Sprintf("%d.%s", expires, sign(policy.SigningKey, msg))
}

func VerifyCSRFToken(policy *models.RouteProtectionPolicy, hostname, path, token string) bool {
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

func (s *RouteProtectionService) List(ctx context.Context, projectID string) ([]models.RouteProtectionStatus, error) {
	rows, err := s.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	result := make([]models.RouteProtectionStatus, 0, len(rows))
	for _, row := range rows {
		st := "pending"
		if (!row.Enabled || row.Secret != nil) && row.AppliedRevision != nil && *row.AppliedRevision == row.Revision {
			st = "applied"
		}
		result = append(result, models.RouteProtectionStatus{
			Hostname: row.Hostname,
			Enabled:  row.Enabled,
			Status:   st,
		})
	}
	return result, nil
}

func (s *RouteProtectionService) GetPolicyByHostname(ctx context.Context, hostname string) (*models.RouteProtection, *models.RouteProtectionPolicy, error) {
	clean, err := NormalizeHostname(hostname)
	if err != nil {
		return nil, nil, err
	}
	rp, err := s.repo.GetByHostname(ctx, clean)
	if err != nil || rp == nil || !rp.Enabled || rp.Secret == nil {
		return nil, nil, nil
	}
	var policy models.RouteProtectionPolicy
	if err := json.Unmarshal([]byte(*rp.Secret), &policy); err != nil {
		return nil, nil, err
	}
	return rp, &policy, nil
}

func (s *RouteProtectionService) Change(ctx context.Context, projectID string, input models.RouteProtectionInput) (*models.RouteProtectionStatus, error) {
	clean, err := NormalizeHostname(input.Hostname)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.Get(ctx, projectID, clean)
	if err != nil {
		return nil, err
	}

	switch input.Action {
	case "protect":
		policy, err := CreatePolicy(input.Password, DefaultSessionSec)
		if err != nil {
			return nil, err
		}
		raw, _ := json.Marshal(policy)
		secretStr := string(raw)
		rev := fmt.Sprintf("%x", sha256.Sum256([]byte(secretStr+clean)))
		rp := &models.RouteProtection{
			ProjectID:       projectID,
			Hostname:        clean,
			Enabled:         true,
			Secret:          &secretStr,
			Revision:        rev,
			AppliedRevision: &rev,
		}
		if err := s.repo.Save(ctx, rp); err != nil {
			return nil, err
		}
		return &models.RouteProtectionStatus{Hostname: clean, Enabled: true, Status: "applied"}, nil

	case "revoke":
		if existing == nil || existing.Secret == nil {
			return nil, errors.New("no active route protection to revoke")
		}
		var policy models.RouteProtectionPolicy
		if err := json.Unmarshal([]byte(*existing.Secret), &policy); err != nil {
			return nil, err
		}
		rotated := RevokeRouteSessions(&policy)
		raw, _ := json.Marshal(rotated)
		secretStr := string(raw)
		rev := fmt.Sprintf("%x", sha256.Sum256([]byte(secretStr+clean)))
		existing.Secret = &secretStr
		existing.Revision = rev
		existing.AppliedRevision = &rev
		if err := s.repo.Save(ctx, existing); err != nil {
			return nil, err
		}
		return &models.RouteProtectionStatus{Hostname: clean, Enabled: true, Status: "applied"}, nil

	case "retry":
		if existing == nil {
			return nil, errors.New("route protection record not found")
		}
		_ = s.repo.MarkApplied(ctx, projectID, clean, existing.Revision)
		return &models.RouteProtectionStatus{Hostname: clean, Enabled: existing.Enabled, Status: "applied"}, nil

	case "unlock":
		if existing != nil {
			existing.Enabled = false
			existing.Revision = fmt.Sprintf("unlocked-%d", time.Now().Unix())
			existing.AppliedRevision = &existing.Revision
			if err := s.repo.Save(ctx, existing); err != nil {
				return nil, err
			}
		}
		return &models.RouteProtectionStatus{Hostname: clean, Enabled: false, Status: "applied"}, nil

	default:
		return nil, fmt.Errorf("unsupported route protection action: %s", input.Action)
	}
}
