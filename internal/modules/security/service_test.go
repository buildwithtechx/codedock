package security

import (
	"context"
	"testing"
)

type mockAppRepo struct {
	services []*SecurityAppService
}

func (m *mockAppRepo) ListAppServices(ctx context.Context, projectID string) ([]*SecurityAppService, error) {
	var res []*SecurityAppService
	for _, s := range m.services {
		res = append(res, s)
	}
	return res, nil
}

func (m *mockAppRepo) UpdateAppService(ctx context.Context, svc *SecurityAppService) error {
	for i, s := range m.services {
		if s.ID == svc.ID {
			m.services[i] = svc
			return nil
		}
	}
	return nil
}

func (m *mockAppRepo) GetByID(ctx context.Context, id string) (*SecurityAppService, error) {
	for _, s := range m.services {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, nil
}

type mockDeployTrigger struct {
	created []string
}

func (m *mockDeployTrigger) TriggerRemediation(ctx context.Context, serviceID string) error {
	m.created = append(m.created, serviceID)
	return nil
}

func TestSecurityService(t *testing.T) {
	verifier, err := NewVerifier("")
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}
	feedManager := NewFeedManager(nil, verifier, "")

	appRepo := &mockAppRepo{
		services: []*SecurityAppService{
			{
				ID:             "svc-node",
				Name:           "web",
				BuildImage:     "node:22-alpine",
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
