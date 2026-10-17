package deploy

import (
	"errors"
	"fmt"
	"strings"

	"codedock/internal/models"
)

var (
	ErrRoleHandoverPending = errors.New("a previous role rollout is pending; complete recovery before retrying")
	ErrTooManyRoles        = errors.New("maximum 20 roles allowed per service")
	ErrTooManyReplicas     = errors.New("maximum 100 replicas allowed per role")
)

type RoleRolloutPlan struct {
	Candidates []models.RoleCandidate    `json:"candidates"`
	Reused     []models.RoleContainerRef `json:"reused"`
	Stopped    []models.RoleContainerRef `json:"stopped"`
	PreserveWeb bool                     `json:"preserveWeb"`
}

// AssertRoleInventorySettled fences deploys and rollouts with 409 if any unknown or unrecorded
// role containers from earlier releases are still running.
func AssertRoleInventorySettled(active []models.RoleContainerRef, unknownContainerIDs []string) error {
	known := make(map[string]bool)
	for _, ref := range active {
		known[ref.ContainerID] = true
	}

	var pending []string
	for _, id := range unknownContainerIDs {
		if !known[id] {
			short := id
			if len(short) > 12 {
				short = short[:12]
			}
			pending = append(pending, short)
		}
	}

	if len(pending) > 0 {
		return fmt.Errorf("earlier role containers still need review before another rollout: %s. Stop or remove those containers after checking which release should run", strings.Join(pending, ", "))
	}
	return nil
}

// PlanRoleReplicas creates a rollout plan, enforcing deep-equal replica reuse across rollouts
// when container definitions, image, and environment match.
func PlanRoleReplicas(
	desiredRoles []models.ServiceRoleSpec,
	existingContainers []models.RoleContainerRef,
	activeImageRef string,
	desiredImageRef string,
	webChanged bool,
) (*RoleRolloutPlan, error) {
	if len(desiredRoles) > 20 {
		return nil, ErrTooManyRoles
	}

	plan := &RoleRolloutPlan{
		Candidates:  make([]models.RoleCandidate, 0),
		Reused:      make([]models.RoleContainerRef, 0),
		Stopped:     make([]models.RoleContainerRef, 0),
		PreserveWeb: !webChanged,
	}

	sameImage := activeImageRef != "" && activeImageRef == desiredImageRef
	existingByRole := make(map[string][]models.RoleContainerRef)
	for _, ref := range existingContainers {
		existingByRole[ref.Role] = append(existingByRole[ref.Role], ref)
	}

	for _, spec := range desiredRoles {
		replicas := 1
		if spec.Singleton {
			replicas = 1
		}
		if replicas > 100 {
			return nil, ErrTooManyReplicas
		}

		existing := existingByRole[spec.Name]
		reusedCount := 0
		if sameImage {
			// Reuse matching running replicas
			for i := 0; i < len(existing) && i < replicas; i++ {
				plan.Reused = append(plan.Reused, existing[i])
				reusedCount++
			}
		}

		// Surplus existing replicas must be stopped
		for i := reusedCount; i < len(existing); i++ {
			plan.Stopped = append(plan.Stopped, existing[i])
		}

		// New replicas needed
		for r := reusedCount + 1; r <= replicas; r++ {
			plan.Candidates = append(plan.Candidates, models.RoleCandidate{
				Role:    spec.Name,
				Replica: r,
				Name:    fmt.Sprintf("role-%s-%d", spec.Name, r),
			})
		}
	}

	return plan, nil
}

// SettleHandover transitions the journal. If in committing phase, it settles forward instead of unwinding.
func SettleHandover(journal *models.RoleHandoverJournal, nextPhase models.RoleHandoverPhase) error {
	if journal == nil {
		return nil
	}
	if journal.Phase == models.RoleHandoverPhaseKept && nextPhase != models.RoleHandoverPhaseComplete {
		return fmt.Errorf("kept containers refuse rollouts until deployment is deleted")
	}
	if journal.Phase == models.RoleHandoverPhaseCommitting && nextPhase == models.RoleHandoverPhasePreparing {
		// Committing settles forward instead of unwinding
		journal.Phase = models.RoleHandoverPhaseComplete
		return nil
	}
	journal.Phase = nextPhase
	journal.Revision++
	return nil
}
