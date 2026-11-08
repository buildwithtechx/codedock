package chatapps

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"codedock/internal/models"

	"github.com/google/uuid"
)

func (s *Service) CreatePairingChallenge(
	ctx context.Context,
	connectionID string,
	provider models.ChatAppProvider,
	extUserID string,
	extChannelID string,
) (*models.ChatAppChallenge, error) {
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("generate challenge token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// 6-digit verification code
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return nil, fmt.Errorf("generate verification code: %w", err)
	}
	code := fmt.Sprintf("%06d", n.Int64()+100000)

	challenge := &models.ChatAppChallenge{
		ID:                "chal_" + uuid.NewString()[:12],
		ConnectionID:      connectionID,
		Provider:          provider,
		ExternalUserID:    extUserID,
		ExternalChannelID: extChannelID,
		Token:             token,
		Code:              code,
		ExpiresAt:         time.Now().UTC().Add(15 * time.Minute),
	}

	if err := s.repo.CreateChallenge(ctx, challenge); err != nil {
		return nil, fmt.Errorf("save pairing challenge: %w", err)
	}
	return challenge, nil
}

func (s *Service) GetChallenge(ctx context.Context, token string) (*models.ChatAppChallenge, error) {
	ch, err := s.repo.GetChallengeByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if ch == nil {
		return nil, errors.New("challenge not found or expired")
	}
	if time.Now().UTC().After(ch.ExpiresAt) {
		return nil, errors.New("challenge has expired; please request a new link in chat")
	}
	return ch, nil
}

func (s *Service) VerifyPairing(ctx context.Context, userID string, token string, code string) (*models.ChatAppIdentity, error) {
	ch, err := s.repo.VerifyChallenge(ctx, token, code)
	if err != nil {
		return nil, fmt.Errorf("verify pairing: %w", err)
	}

	identity := &models.ChatAppIdentity{
		ID:               "ident_" + uuid.NewString()[:12],
		ConnectionID:     ch.ConnectionID,
		UserID:           userID,
		Provider:         ch.Provider,
		ExternalUserID:   ch.ExternalUserID,
		ExternalUsername: ch.ExternalUserID,
		Role:             "operator",
		Permissions:      []string{"projects:read", "deployments:write", "logs:read"},
		Active:           true,
		PairedAt:         time.Now().UTC(),
	}

	if err := s.repo.CreateIdentity(ctx, identity); err != nil {
		return nil, fmt.Errorf("create paired identity: %w", err)
	}
	return identity, nil
}
