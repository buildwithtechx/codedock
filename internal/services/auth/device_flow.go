package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
)

var (
	ErrDeviceSessionNotFound = errors.New("device session not found")
	ErrDeviceSessionExpired  = errors.New("device session has expired")
	ErrDeviceAlreadyDecided  = errors.New("device session has already been decided")
)

type DeviceSessionStatus string

const (
	DeviceStatusPending  DeviceSessionStatus = "pending"
	DeviceStatusApproved DeviceSessionStatus = "approved"
	DeviceStatusDenied   DeviceSessionStatus = "denied"
	DeviceStatusExpired  DeviceSessionStatus = "expired"
)

type DeviceSession struct {
	DeviceCode string              `json:"deviceCode"`
	UserCode   string              `json:"userCode"`
	ClientName string              `json:"clientName"`
	Status     DeviceSessionStatus `json:"status"`
	UserID     string              `json:"userId,omitempty"`
	Token      string              `json:"token,omitempty"`
	UserEmail  string              `json:"userEmail,omitempty"`
	UserName   string              `json:"userName,omitempty"`
	ExpiresAt  time.Time           `json:"expiresAt"`
	CreatedAt  time.Time           `json:"createdAt"`
}

type DeviceFlowService struct {
	mu       sync.RWMutex
	byDevice map[string]*DeviceSession
	byUser   map[string]*DeviceSession
}

func NewDeviceFlowService() *DeviceFlowService {
	s := &DeviceFlowService{
		byDevice: make(map[string]*DeviceSession),
		byUser:   make(map[string]*DeviceSession),
	}
	return s
}

func (s *DeviceFlowService) CreateSession(clientName string) (*DeviceSession, error) {
	if strings.TrimSpace(clientName) == "" {
		clientName = "Codedock Client"
	}

	deviceCodeBytes := make([]byte, 16)
	if _, err := rand.Read(deviceCodeBytes); err != nil {
		return nil, fmt.Errorf("generate device code: %w", err)
	}
	deviceCode := hex.EncodeToString(deviceCodeBytes)

	userCode, err := generateUserCode()
	if err != nil {
		return nil, fmt.Errorf("generate user code: %w", err)
	}

	session := &DeviceSession{
		DeviceCode: deviceCode,
		UserCode:   userCode,
		ClientName: clientName,
		Status:     DeviceStatusPending,
		ExpiresAt:  time.Now().Add(15 * time.Minute),
		CreatedAt:  time.Now(),
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupExpiredLocked()
	s.byDevice[deviceCode] = session
	s.byUser[userCode] = session

	return session, nil
}

func (s *DeviceFlowService) GetByDeviceCode(deviceCode string) (*DeviceSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.byDevice[deviceCode]
	if !exists {
		return nil, ErrDeviceSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		session.Status = DeviceStatusExpired
		return session, ErrDeviceSessionExpired
	}

	return session, nil
}

func (s *DeviceFlowService) GetByUserCode(userCode string) (*DeviceSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	norm := strings.ToUpper(strings.TrimSpace(userCode))
	session, exists := s.byUser[norm]
	if !exists {
		return nil, ErrDeviceSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		session.Status = DeviceStatusExpired
		return session, ErrDeviceSessionExpired
	}

	return session, nil
}

func (s *DeviceFlowService) Authorize(userCode string, userID, userEmail, userName, rawToken string, approve bool) (*DeviceSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	norm := strings.ToUpper(strings.TrimSpace(userCode))
	session, exists := s.byUser[norm]
	if !exists {
		return nil, ErrDeviceSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		session.Status = DeviceStatusExpired
		return nil, ErrDeviceSessionExpired
	}

	if session.Status != DeviceStatusPending {
		return nil, ErrDeviceAlreadyDecided
	}

	if approve {
		session.Status = DeviceStatusApproved
		session.UserID = userID
		session.UserEmail = userEmail
		session.UserName = userName
		session.Token = rawToken
	} else {
		session.Status = DeviceStatusDenied
	}

	return session, nil
}

func (s *DeviceFlowService) cleanupExpiredLocked() {
	now := time.Now()
	for dev, sess := range s.byDevice {
		if now.After(sess.ExpiresAt.Add(5 * time.Minute)) {
			delete(s.byDevice, dev)
			delete(s.byUser, sess.UserCode)
		}
	}
}

func generateUserCode() (string, error) {
	const charset = "BCDFGHJKLMNPQRSTVWXYZ23456789"
	part1 := make([]byte, 4)
	part2 := make([]byte, 4)

	for i := 0; i < 4; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		part1[i] = charset[idx.Int64()]
	}

	for i := 0; i < 4; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		part2[i] = charset[idx.Int64()]
	}

	return fmt.Sprintf("%s-%s", string(part1), string(part2)), nil
}
