package runtimes

import (
	"codedock/internal/engine/bare"
	"codedock/internal/engine/kubernetes"
	"codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

func (s *Service) snapshot(ctx context.Context, id string) (string, error) {
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	runtime, err := s.store.Get(ctx, id)
	if err != nil {
		return "", err
	}
	raw, err := json.Marshal([]any{app, runtime})
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(raw)
	return hex.EncodeToString(hash[:]), nil
}
func (s *Service) validateTarget(ctx context.Context, app *models.AppService, request models.RuntimeReviewRequest) error {
	existing, err := s.store.Get(ctx, app.ID)
	if err != nil {
		return err
	}
	if existing.Revision != request.Revision || existing.Journal != "" {
		return fmt.Errorf("runtime configuration changed or recovery is pending")
	}
	if app.Status == models.AppServiceStatusRunning || app.Status == models.AppServiceStatusBuilding {
		return fmt.Errorf("stop the application before changing its destination")
	}
	if existing.Target.Kind == "docker" && app.ContainerID != "" {
		active, err := s.builder.DockerWorkloadActive(ctx, app.ID)
		if err != nil {
			return err
		}
		if active {
			return fmt.Errorf("existing Docker containers must be stopped before changing destination")
		}
	}
	if (existing.Target.Kind == "kubernetes" || (existing.Target.Kind == "bare" && strings.HasPrefix(app.ContainerID, "bare:"))) && existing.Revision > 0 && app.ContainerID != "" {
		observation, err := s.Observe(ctx, app.ID)
		if err != nil {
			return err
		}
		if observation.Desired > 0 || observation.Available > 0 {
			return fmt.Errorf("stop the existing cluster workload before changing its configuration")
		}
	}
	switch request.Target.Kind {
	case "docker":
		if request.Target.ClusterID != "" || len(request.Target.Volumes) > 0 || len(request.Target.NodeIDs) > 0 {
			return fmt.Errorf("Docker destinations cannot contain cluster settings")
		}
	case "bare":
		if err := bare.Validate(app, request.Target); err != nil {
			return err
		}
		if err := s.validateBare(ctx, app, request.Target); err != nil {
			return err
		}
		if err := s.native.Preflight(ctx, app, request.Target); err != nil {
			return err
		}
	case "kubernetes":
		if _, err := s.cluster(ctx, app, request.Target); err != nil {
			return err
		}
		copy := *app
		if copy.Replicas <= 0 {
			copy.Replicas = 1
		}
		if copy.InternalPort <= 0 {
			copy.InternalPort = 3000
		}
		image := copy.ImageRef
		if image == "" {
			if (copy.RegistryID == nil && request.Target.RegistryID == "") || request.Target.ImageRepository == "" {
				return fmt.Errorf("Git deployments require a registry and image repository")
			}
			image = request.Target.ImageRepository + ":review"
		}
		if err := kubernetes.ValidateWorkload(&models.KubernetesWorkload{App: copy, Target: request.Target, Image: image}); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported runtime destination")
	}
	return nil
}
func (s *Service) Review(ctx context.Context, user, id string, request models.RuntimeReviewRequest) (*models.OperationReview, error) {
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.validateTarget(ctx, app, request); err != nil {
		return nil, err
	}
	clusterRevision := 0
	if request.Target.Kind == "kubernetes" {
		cluster, err := s.cluster(ctx, app, request.Target)
		if err != nil {
			return nil, err
		}
		clusterRevision = cluster.Revision
	}
	payload, err := json.Marshal(models.RuntimeReviewPlan{ServiceID: id, Target: request.Target, Revision: request.Revision, ClusterRevision: clusterRevision})
	if err != nil {
		return nil, err
	}
	snapshot, err := s.snapshot(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.operations.Review(ctx, user, app.ProjectID, "runtime", id, string(payload), snapshot, "Save this application's deployment destination and placement. Future deployments use this runtime. Existing persistent data is retained; this operation does not copy it.")
}
func (s *Service) Apply(ctx context.Context, user, id, confirmation string) error {
	op, err := s.operations.Get(ctx, id)
	if err != nil {
		return err
	}
	if op.UserID != user || op.Kind != "runtime" {
		return fmt.Errorf("runtime operation not found")
	}
	var plan models.RuntimeReviewPlan
	if err := json.Unmarshal([]byte(op.Payload), &plan); err != nil {
		return err
	}
	snapshot, err := s.snapshot(ctx, plan.ServiceID)
	if err != nil {
		return err
	}
	return s.operations.Apply(ctx, id, user, confirmation, snapshot, func(ctx context.Context, op *models.Operation, progress func(string, string) error) error {
		operation, release, err := s.builder.BeginServiceOperation(ctx, plan.ServiceID)
		if err != nil {
			return err
		}
		defer release()
		ctx = operation
		latest, err := s.snapshot(ctx, plan.ServiceID)
		if err != nil {
			return err
		}
		if latest != op.Snapshot {
			return fmt.Errorf("application changed before applying runtime configuration")
		}
		app, err := s.apps.GetByID(ctx, plan.ServiceID)
		if err != nil {
			return err
		}
		if err := s.validateTarget(ctx, app, models.RuntimeReviewRequest{Target: plan.Target, Revision: plan.Revision}); err != nil {
			return err
		}
		if plan.Target.Kind == "kubernetes" {
			cluster, err := s.cluster(ctx, app, plan.Target)
			if err != nil {
				return err
			}
			if cluster.Revision != plan.ClusterRevision {
				return fmt.Errorf("cluster changed since review")
			}
		}
		if err := progress("CONFIGURING", "Saving reviewed runtime destination"); err != nil {
			return err
		}
		return s.store.Save(ctx, &models.ServiceRuntime{ServiceID: app.ID, ProjectID: app.ProjectID, Target: plan.Target}, plan.Revision)
	})
}
func (s *Service) ApplyForService(ctx context.Context, user, service, id, confirmation string) error {
	op, err := s.operations.Get(ctx, id)
	if err != nil {
		return err
	}
	if op.Target != service {
		return fmt.Errorf("runtime operation belongs to another service")
	}
	return s.Apply(ctx, user, id, confirmation)
}
func (s *Service) Operation(ctx context.Context, user, service, id string) (*models.Operation, error) {
	op, err := s.operations.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if op.UserID != user || op.Kind != "runtime" || op.Target != service {
		return nil, fmt.Errorf("runtime operation not found")
	}
	return op, nil
}
