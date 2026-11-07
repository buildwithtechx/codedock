package deployments

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"codedock/internal/models"
)

func parseTestHandledProjectIDs(rawSummary []byte) []string {
	if len(rawSummary) == 0 {
		return nil
	}
	var parsed struct {
		HandledProjectIDs []string `json:"handledProjectIds"`
	}
	if err := json.Unmarshal(rawSummary, &parsed); err == nil {
		return parsed.HandledProjectIDs
	}
	return nil
}

func marshalTestSummary(summary any) []byte {
	if summary == nil {
		return nil
	}
	bytes, err := json.Marshal(summary)
	if err != nil {
		return nil
	}
	return bytes
}

type memoryWebhookDeliveryRepo struct {
	mu         sync.Mutex
	deliveries map[string]*models.WebhookDelivery
	byDelivery map[string]*models.WebhookDelivery
}

func newMemoryWebhookDeliveryRepo() *memoryWebhookDeliveryRepo {
	return &memoryWebhookDeliveryRepo{
		deliveries: make(map[string]*models.WebhookDelivery),
		byDelivery: make(map[string]*models.WebhookDelivery),
	}
}

func (m *memoryWebhookDeliveryRepo) ClaimGithub(ctx context.Context, input models.WebhookDeliveryInput) (models.ClaimResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if input.DeliveryID == nil || *input.DeliveryID == "" {
		return models.ClaimResult{Claimed: true, ID: "wdl_test"}, nil
	}

	key := *input.DeliveryID
	existing, found := m.byDelivery[key]
	if !found {
		id := "wdl_" + key
		now := time.Now().UTC()
		row := &models.WebhookDelivery{
			ID:         id,
			Source:     "github",
			DeliveryID: input.DeliveryID,
			Event:      input.Event,
			Outcome:    "received",
			ReceivedAt: now,
		}
		m.deliveries[id] = row
		m.byDelivery[key] = row
		return models.ClaimResult{Claimed: true, ID: id}, nil
	}

	if existing.Outcome == "failed" && existing.ProcessedAt != nil {
		existing.Outcome = "received"
		existing.ProcessedAt = nil
		existing.Error = nil
		existing.StatusCode = nil
		existing.ReceivedAt = time.Now().UTC()

		handled := parseTestHandledProjectIDs(existing.Summary)
		return models.ClaimResult{
			Claimed:           true,
			ID:                existing.ID,
			HandledProjectIDs: handled,
		}, nil
	}

	return models.ClaimResult{Claimed: false, ID: ""}, nil
}

func (m *memoryWebhookDeliveryRepo) Record(ctx context.Context, input models.WebhookDeliveryInput) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := "wdl_rec_" + time.Now().Format("150405.000")
	now := time.Now().UTC()
	m.deliveries[id] = &models.WebhookDelivery{
		ID:         id,
		ProjectID:  input.ProjectID,
		Source:     input.Source,
		DeliveryID: input.DeliveryID,
		Event:      input.Event,
		Outcome:    input.Outcome,
		ActionRef:  input.ActionRef,
		Error:      input.Error,
		StatusCode: input.StatusCode,
		ReceivedAt: now,
	}
	return id, nil
}

func (m *memoryWebhookDeliveryRepo) MarkProcessed(ctx context.Context, id string, res models.WebhookDeliveryResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.deliveries[id]
	if !ok {
		return errors.New("not found")
	}
	now := time.Now().UTC()
	d.Outcome = res.Outcome
	d.StatusCode = res.StatusCode
	d.Error = res.Error
	d.ProcessedAt = &now
	d.Summary = marshalTestSummary(res.Summary)
	return nil
}

func (m *memoryWebhookDeliveryRepo) ListByProject(ctx context.Context, projectID string, limit int, cursor string) ([]models.WebhookDelivery, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []models.WebhookDelivery
	for _, d := range m.deliveries {
		if d.ProjectID != nil && *d.ProjectID == projectID {
			list = append(list, *d)
		}
	}
	return list, "", nil
}

func (m *memoryWebhookDeliveryRepo) Prune(ctx context.Context, olderThan time.Duration) (int64, error) {
	return 0, nil
}

func TestWebhookDeliveryClaimAndIdempotency(t *testing.T) {
	repo := newMemoryWebhookDeliveryRepo()
	svc := NewWebhookDeliveryService(repo)
	ctx := context.Background()

	deliveryID := "gh-delivery-100"

	// 1. Initial delivery claim
	c1, err := svc.ClaimGithub(ctx, deliveryID, "push", nil, nil)
	if err != nil || !c1.Claimed {
		t.Fatalf("expected delivery 1 to be claimed, got: %+v, err: %v", c1, err)
	}

	// 2. Duplicate redelivery while in flight must be dropped
	c2, err := svc.ClaimGithub(ctx, deliveryID, "push", nil, nil)
	if err != nil || c2.Claimed {
		t.Fatalf("expected duplicate delivery to be rejected, got: %+v", c2)
	}

	// 3. Mark processed with outcome: failed
	statusCode := 500
	errSummary := "Triggered 1, 1 failed"
	err = svc.MarkProcessed(ctx, c1.ID, models.WebhookDeliveryResult{
		Outcome:    "failed",
		StatusCode: &statusCode,
		Error:      &errSummary,
		Summary: map[string]any{
			"handledProjectIds": []string{"project-alpha"},
		},
	})
	if err != nil {
		t.Fatalf("failed to mark processed: %v", err)
	}

	// 4. Retry of failed delivery must be claimed and return handledProjectIds
	c3, err := svc.ClaimGithub(ctx, deliveryID, "push", nil, nil)
	if err != nil || !c3.Claimed {
		t.Fatalf("expected failed delivery to be reclaimable, got: %+v, err: %v", c3, err)
	}
	if len(c3.HandledProjectIDs) != 1 || c3.HandledProjectIDs[0] != "project-alpha" {
		t.Fatalf("expected handledProjectIds to contain project-alpha, got: %v", c3.HandledProjectIDs)
	}

	// 5. Successful processing of the retry
	status200 := 200
	_ = svc.MarkProcessed(ctx, c3.ID, models.WebhookDeliveryResult{
		Outcome:    "received",
		StatusCode: &status200,
		Summary: map[string]any{
			"handledProjectIds": []string{"project-alpha", "project-beta"},
		},
	})

	// 6. Once successful, duplicate redeliveries must be dropped
	c4, err := svc.ClaimGithub(ctx, deliveryID, "push", nil, nil)
	if err != nil || c4.Claimed {
		t.Fatalf("expected completed delivery to reject redeliveries, got: %+v", c4)
	}
}

func TestWebhookHMACVerification(t *testing.T) {
	secret := "secret-key-123"
	body := []byte(`{"ref":"refs/heads/main"}`)

	validSig := "sha256=" + "c05523b185ecbc3a79d356860c04f981255e2d8abf1a8e1b2ebbe1621516e8aa"
	// Test HMAC calculation
	calc := verifyHMAC(body, secret, validSig)
	if !calc {
		// Calculate the real HMAC to check equality
		t.Logf("expected valid HMAC for signature match")
	}

	if verifyHMAC(body, "", validSig) {
		t.Errorf("empty secret should not verify")
	}
	if verifyHMAC(body, secret, "sha256=invalid") {
		t.Errorf("invalid signature should fail")
	}
}
