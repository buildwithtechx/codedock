package system

import (
	"context"
	"testing"

	"codedock/internal/models"
)

func TestManagedStreamRequiresMembership(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	review, err := env.service.ReviewProvision(ctx, env.adminID, env.orgID, models.ReviewManagedProvisionRequest{CredentialID: env.credID, Name: "edge-1", Region: "fsn1", ServerType: "cx23", Image: "ubuntu-24.04"})
	if err != nil {
		t.Fatalf("review: %v", err)
	}
	calls := 0
	err = env.service.StreamOperation(ctx, env.outsiderID, review.Review.Operation.ID, func(*models.Operation) error {
		calls++
		return nil
	})
	if err == nil {
		t.Fatal("expected outsider to be denied operation stream")
	}
	if calls != 0 {
		t.Fatalf("expected no emissions for outsider, got %d", calls)
	}
}

func TestManagedStreamEmitsSnapshotThenStops(t *testing.T) {
	env := setupManagedTest(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	review, err := env.service.ReviewProvision(ctx, env.adminID, env.orgID, models.ReviewManagedProvisionRequest{CredentialID: env.credID, Name: "edge-1", Region: "fsn1", ServerType: "cx23", Image: "ubuntu-24.04"})
	if err != nil {
		t.Fatalf("review: %v", err)
	}
	emitted := 0
	err = env.service.StreamOperation(ctx, env.memberID, review.Review.Operation.ID, func(op *models.Operation) error {
		emitted++
		if op.ID != review.Review.Operation.ID || op.Status != "REVIEWED" {
			t.Errorf("unexpected snapshot %#v", op)
		}
		cancel()
		return nil
	})
	if err != nil {
		t.Fatalf("stream: %v", err)
	}
	if emitted == 0 {
		t.Fatal("expected at least one snapshot")
	}
}
