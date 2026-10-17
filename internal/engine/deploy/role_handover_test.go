package deploy

import (
	"strings"
	"testing"

	"codedock/internal/models"
)

func TestAssertRoleInventorySettled(t *testing.T) {
	active := []models.RoleContainerRef{
		{ContainerID: "cont-1234567890ab", Role: "worker", Replica: 1},
	}

	// Clean inventory
	if err := AssertRoleInventorySettled(active, []string{"cont-1234567890ab"}); err != nil {
		t.Fatalf("expected settled inventory, got error: %v", err)
	}

	// Unknown earlier containers
	err := AssertRoleInventorySettled(active, []string{"cont-1234567890ab", "unknown-9876543210"})
	if err == nil {
		t.Fatalf("expected 409 error when earlier role containers are unknown")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "earlier role containers still need review") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestPlanRoleReplicas(t *testing.T) {
	desired := []models.ServiceRoleSpec{
		{Name: "worker", Kind: "worker", Command: "python worker.py", Singleton: true},
	}
	existing := []models.RoleContainerRef{
		{ContainerID: "old-worker-1", Role: "worker", Replica: 1},
	}

	// Same image: reuses deep-equal replica
	plan, err := PlanRoleReplicas(desired, existing, "img:v1", "img:v1", false)
	if err != nil {
		t.Fatalf("unexpected plan error: %v", err)
	}
	if len(plan.Reused) != 1 || plan.Reused[0].ContainerID != "old-worker-1" {
		t.Fatalf("expected old-worker-1 to be reused, got %+v", plan.Reused)
	}
	if len(plan.Candidates) != 0 {
		t.Fatalf("expected 0 candidates on full reuse, got %d", len(plan.Candidates))
	}
	if !plan.PreserveWeb {
		t.Fatalf("expected PreserveWeb to be true when web unchanged")
	}

	// Different image: does not reuse, creates new candidate and stops old
	planNew, err := PlanRoleReplicas(desired, existing, "img:v1", "img:v2", false)
	if err != nil {
		t.Fatalf("unexpected plan error: %v", err)
	}
	if len(planNew.Reused) != 0 {
		t.Fatalf("expected 0 reused on image change")
	}
	if len(planNew.Candidates) != 1 || planNew.Candidates[0].Role != "worker" {
		t.Fatalf("expected 1 candidate for new worker")
	}
	if len(planNew.Stopped) != 1 || planNew.Stopped[0].ContainerID != "old-worker-1" {
		t.Fatalf("expected old-worker-1 to be stopped")
	}
}

func TestSettleHandover(t *testing.T) {
	journal := &models.RoleHandoverJournal{
		Version:  1,
		Revision: 1,
		Phase:    models.RoleHandoverPhaseCommitting,
	}

	// Committing phase settles forward instead of unwinding
	if err := SettleHandover(journal, models.RoleHandoverPhasePreparing); err != nil {
		t.Fatalf("unexpected settle error: %v", err)
	}
	if journal.Phase != models.RoleHandoverPhaseComplete {
		t.Fatalf("expected PhaseComplete after settling committing, got %v", journal.Phase)
	}

	// Kept containers refuse new rollouts
	journalKept := &models.RoleHandoverJournal{
		Version:  1,
		Revision: 2,
		Phase:    models.RoleHandoverPhaseKept,
	}
	if err := SettleHandover(journalKept, models.RoleHandoverPhasePreparing); err == nil {
		t.Fatalf("expected kept containers to refuse rollouts")
	}
}
