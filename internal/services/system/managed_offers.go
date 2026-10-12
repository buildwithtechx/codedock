package system

import (
	"context"
	"fmt"

	"codedock/internal/models"
)

func (s *managedService) Offers(ctx context.Context, userID, orgID, credentialID string) (*models.ManagedOffers, error) {
	if err := s.requireOrgRole(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	credential, err := s.credentialForOrg(ctx, orgID, credentialID)
	if err != nil {
		return nil, err
	}
	offers, err := s.providers(credential.Token).Offers(ctx)
	if err != nil {
		return nil, fmt.Errorf("load hetzner catalogue: %w", err)
	}
	return &offers, nil
}

func (s *managedService) GetOperation(ctx context.Context, userID, operationID string) (*models.Operation, error) {
	return s.authorizeManagedOperation(ctx, userID, operationID)
}
