package system

import (
	"context"
	"errors"
	"testing"

	"codedock/internal/engine/hetzner"
	"codedock/internal/models"
)

func (f *fakeProvider) Verify(ctx context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.verifyErr
}

func (f *fakeProvider) Offers(ctx context.Context) (models.ManagedOffers, error) {
	tier, _ := hetzner.TierByName("cx23")
	return models.ManagedOffers{
		Provider:  models.ManagedProviderHetzner,
		Locations: []models.ManagedLocation{{Name: "fsn1", City: "Falkenstein", Country: "DE"}},
		Offers:    []models.ManagedOffer{{Tier: tier, Prices: map[string]float64{"fsn1": 3.49}}},
		Images:    hetzner.Images(),
	}, nil
}

func TestManagedCredentialRejectsInvalidToken(t *testing.T) {
	env := setupManagedTest(t)
	env.fake.mu.Lock()
	env.fake.verifyErr = errors.New("unauthorized")
	env.fake.mu.Unlock()
	_, err := env.service.CreateCredential(context.Background(), env.adminID, env.orgID, models.CreateManagedCredentialRequest{Provider: models.ManagedProviderHetzner, Label: "bad", Token: "nope"})
	if err == nil {
		t.Fatal("expected invalid token to be rejected")
	}
}

func TestManagedOffersRequireMembership(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	if _, err := env.service.Offers(ctx, env.outsiderID, env.orgID, env.credID); err == nil {
		t.Fatal("expected outsider to be rejected")
	}
	offers, err := env.service.Offers(ctx, env.memberID, env.orgID, env.credID)
	if err != nil {
		t.Fatalf("offers: %v", err)
	}
	if len(offers.Offers) != 1 || offers.Offers[0].Prices["fsn1"] != 3.49 {
		t.Fatalf("unexpected offers %#v", offers)
	}
}

func TestManagedOperationProgressIsOrgScoped(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	review, err := env.service.ReviewProvision(ctx, env.adminID, env.orgID, models.ReviewManagedProvisionRequest{CredentialID: env.credID, Name: "edge-1", Region: "fsn1", ServerType: "cx23", Image: "ubuntu-24.04"})
	if err != nil {
		t.Fatalf("review: %v", err)
	}
	if _, err := env.service.GetOperation(ctx, env.outsiderID, review.Review.Operation.ID); err == nil {
		t.Fatal("expected outsider to be denied operation progress")
	}
	op, err := env.service.GetOperation(ctx, env.adminID, review.Review.Operation.ID)
	if err != nil {
		t.Fatalf("get operation: %v", err)
	}
	if op.Status != "REVIEWED" {
		t.Fatalf("expected reviewed operation, got %s", op.Status)
	}
}
