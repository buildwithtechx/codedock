package security

import (
	"context"
	"testing"

	"codedock/internal/models"
)

type mockAppRepo struct {
	services []*models.AppService
}

func (m *mockAppRepo) Create(ctx context.Context, svc *models.AppService) error { return nil }
func (m *mockAppRepo) GetByID(ctx context.Context, id string) (*models.AppService, error) {
	for _, s := range m.services {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, nil
}
func (m *mockAppRepo) ListByEnvironment(ctx context.Context, environmentID string) ([]*models.AppService, error) {
	return nil, nil
}
func (m *mockAppRepo) ListByProject(ctx context.Context, projectID string) ([]*models.AppService, error) {
	var res []*models.AppService
	for _, s := range m.services {
		if s.ProjectID == projectID {
			res = append(res, s)
		}
	}
	return res, nil
}
func (m *mockAppRepo) ListByOrganization(ctx context.Context, organizationID string, search string) ([]*models.AppService, error) {
	return nil, nil
}
func (m *mockAppRepo) ListAll(ctx context.Context, search string) ([]*models.AppService, error) {
	return nil, nil
}
func (m *mockAppRepo) SearchIDsByName(ctx context.Context, term string, limit int) ([]string, error) {
	return nil, nil
}
func (m *mockAppRepo) Update(ctx context.Context, svc *models.AppService) error { return nil }
func (m *mockAppRepo) Delete(ctx context.Context, id string) error              { return nil }
func (m *mockAppRepo) CreateWebhook(ctx context.Context, w *models.Webhook) error {
	return nil
}
func (m *mockAppRepo) ListWebhooksByService(ctx context.Context, serviceID string) ([]*models.Webhook, error) {
	return nil, nil
}
func (m *mockAppRepo) DeleteWebhook(ctx context.Context, id, serviceID string) error { return nil }
func (m *mockAppRepo) CreateLogDrain(ctx context.Context, d *models.LogDrain) error  { return nil }
func (m *mockAppRepo) ListLogDrainsByService(ctx context.Context, serviceID string) ([]*models.LogDrain, error) {
	return nil, nil
}
func (m *mockAppRepo) DeleteLogDrain(ctx context.Context, id, serviceID string) error { return nil }

type mockDeployTrigger struct {
	created []*models.Deployment
}

func (m *mockDeployTrigger) CreateDeployment(ctx context.Context, d *models.Deployment) (*models.Deployment, error) {
	m.created = append(m.created, d)
	return d, nil
}

func TestSecurityService(t *testing.T) {
	verifier, err := NewVerifier("")
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}
	feedManager := NewFeedManager(nil, verifier, "")

	appRepo := &mockAppRepo{
		services: []*models.AppService{
			{
				ID:             "svc-node",
				ProjectID:      "proj-1",
				Name:           "web",
				ImageRef:       "node:22-alpine",
				Stack:          "node",
				RuntimeVersion: "22.23.3", // matches CVE-2026-111497-node (< 22.23.4)
			},
		},
	}
	deployTrigger := &mockDeployTrigger{}

	svc := NewSecurityService(nil, feedManager, nil, appRepo, nil, deployTrigger)

	t.Run("GetAdvisories loads baseline", func(t *testing.T) {
		resp, err := svc.GetAdvisories(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Advisories) < 3 {
			t.Fatalf("expected at least 3 baseline advisories, got %d", len(resp.Advisories))
		}
		if resp.Stale {
			t.Fatal("expected baseline not to be stale")
		}
	})

	t.Run("CheckProject detects vulnerable node service", func(t *testing.T) {
		res, err := svc.CheckProject(context.Background(), "proj-1", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Findings) == 0 {
			t.Fatal("expected finding for vulnerable node service")
		}

		finding := res.Findings[0]
		if finding.AdvisoryID != "CVE-2026-111497-node" {
			t.Fatalf("expected CVE-2026-111497-node, got %s", finding.AdvisoryID)
		}
		if !finding.CanRefresh {
			t.Fatal("expected canRefresh to be true for same major node fix")
		}
	})

	t.Run("ApplyRemediation queues attempt and triggers deploy", func(t *testing.T) {
		attempt, err := svc.ApplyRemediation(context.Background(), "proj-1", "CVE-2026-111497-node", "dep-old", nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if attempt.Status != "queued" {
			t.Fatalf("expected queued status, got %s", attempt.Status)
		}
		if len(deployTrigger.created) != 1 {
			t.Fatalf("expected 1 deployment created, got %d", len(deployTrigger.created))
		}
	})

	t.Run("ApplyRemediation rejects non-existent advisory", func(t *testing.T) {
		_, err := svc.ApplyRemediation(context.Background(), "proj-1", "NON-EXISTENT", "dep-old", nil, nil)
		if err == nil {
			t.Fatal("expected error for non-existent advisory")
		}
	})

	t.Run("ApplyRemediation validates runtime patch reviewId and persists settings", func(t *testing.T) {
		appRepo.services[0].RuntimeVersion = "22.23.3"
		appRepo.services[0].BuildImage = ""
		appRepo.services[0].RuntimeImage = ""

		res, err := svc.CheckProject(context.Background(), "proj-1", true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Findings) == 0 || res.Findings[0].RuntimePatch == nil {
			t.Fatal("expected runtime patch on finding")
		}

		patch := res.Findings[0].RuntimePatch
		invalidPatchID := "0000000000000000000000000000000000000000000000000000000000000000"
		_, err = svc.ApplyRemediation(context.Background(), "proj-1", "CVE-2026-111497-node", "dep-old", nil, &invalidPatchID)
		if err == nil {
			t.Fatal("expected error for mismatched runtimePatchId")
		}

		// Apply with correct runtimePatchId
		attempt, err := svc.ApplyRemediation(context.Background(), "proj-1", "CVE-2026-111497-node", "dep-old", nil, &patch.ReviewID)
		if err != nil {
			t.Fatalf("unexpected error applying valid patch: %v", err)
		}
		if attempt.Status != "queued" {
			t.Fatalf("expected status queued, got %s", attempt.Status)
		}

		// Check service persistence
		updatedSvc, err := appRepo.GetByID(context.Background(), "svc-node")
		if err != nil {
			t.Fatalf("failed to get updated service: %v", err)
		}
		if updatedSvc.RuntimeVersion != patch.TargetVersion {
			t.Errorf("persisted runtimeVersion = %s, want %s", updatedSvc.RuntimeVersion, patch.TargetVersion)
		}
		if updatedSvc.BuildImage != patch.BuildImage {
			t.Errorf("persisted buildImage = %s, want %s", updatedSvc.BuildImage, patch.BuildImage)
		}
	})
}
