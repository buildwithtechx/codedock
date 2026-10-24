package clusterdata

import (
	"codedock/internal/engine/kubernetes"
	"codedock/internal/models"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (s *Service) Review(ctx context.Context, user, project, clusterID string, request models.ClusterDataRequest) (*models.OperationReview, error) {
	if request.Action == "restore" {
		value, err := time.Parse(time.RFC3339, request.RestoreTime)
		if err != nil || value.After(time.Now()) {
			return nil, fmt.Errorf("select a past UTC recovery time")
		}
		request.RestoreTime = value.UTC().Format(time.RFC3339)
	}
	cluster, err := s.cluster(ctx, project, clusterID)
	if err != nil {
		return nil, err
	}
	plan := &models.ClusterDataPlan{Action: request.Action, Record: models.ClusterData{ClusterID: clusterID, ProjectID: project}}
	if request.Action == "operators" {
		plan.Operators, err = s.operatorManifests(ctx, clusterID, request)
		if err != nil {
			return nil, err
		}
	} else {
		if request.Action != "create" && request.Action != "restore" && request.Action != "recover" && request.Action != "backup" {
			return nil, fmt.Errorf("unsupported cluster database action")
		}
		if request.Action == "recover" || request.Action == "backup" {
			plan, err = s.owned(ctx, project, clusterID, request.Spec.ID)
			if err != nil {
				return nil, err
			}
			plan.Action = request.Action
		} else {
			if err := kubernetes.ValidateDataSpec(request.Spec, len(cluster.Nodes)); err != nil {
				return nil, err
			}
			environment, err := s.environments.Get(ctx, request.Spec.EnvironmentID)
			if err != nil || environment.ProjectID != project {
				return nil, fmt.Errorf("select an environment in this project")
			}
			if existing, err := s.store.Get(ctx, request.Spec.ID); err == nil {
				if existing.Record.ProjectID != project || existing.Record.ClusterID != clusterID || existing.Record.Spec != request.Spec || existing.SourceID != request.SourceID || (request.Action == "restore" && existing.RestoreTime != request.RestoreTime) {
					return nil, fmt.Errorf("database ID already belongs to another saved configuration")
				}
				plan = existing
				plan.Action = request.Action
			} else if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			} else {
				secret := make([]byte, 32)
				if _, err := rand.Read(secret); err != nil {
					return nil, err
				}
				plan.Password = hex.EncodeToString(secret)
				plan.Record.ID = request.Spec.ID
				plan.Record.Spec = request.Spec
				plan.SourceID = request.SourceID
				destination, err := s.destination(ctx, request.Spec.S3DestinationID)
				if err != nil {
					return nil, err
				}
				var source *models.S3Destination
				sourceName := ""
				if request.Action == "restore" {
					original, err := s.owned(ctx, project, clusterID, request.SourceID)
					if err != nil {
						return nil, err
					}
					if request.Spec.Engine != "postgres" || original.Record.Spec.Engine != "postgres" || original.Record.ID == request.Spec.ID || original.Record.Spec.S3DestinationID == "" {
						return nil, fmt.Errorf("restore requires a separate PostgreSQL target and a source with object-store backups")
					}
					if postgresMajor(request.Spec.Image) != postgresMajor(original.Record.Spec.Image) {
						return nil, fmt.Errorf("PostgreSQL recovery requires the same major version")
					}
					restoreTime, err := time.Parse(time.RFC3339, request.RestoreTime)
					if err != nil || restoreTime.After(time.Now()) {
						return nil, fmt.Errorf("select a past UTC recovery time")
					}
					plan.RestoreTime = restoreTime.UTC().Format(time.RFC3339)
					backups, err := s.engine.DataBackups(ctx, cluster.Nodes[0], &original.Record)
					if err != nil {
						return nil, err
					}
					completed := false
					for _, backup := range backups {
						if backup["phase"] == "completed" {
							completed = true
						}
					}
					if !completed {
						return nil, fmt.Errorf("source requires a completed object-store backup before recovery")
					}
					source, err = s.destination(ctx, original.Record.Spec.S3DestinationID)
					if err != nil {
						return nil, err
					}
					_, sourceName = kubernetes.DataIdentity(original.Record.Spec)
				}
				plan.Manifest, err = kubernetes.DataManifest(plan, destination, source, sourceName)
				if err != nil {
					return nil, err
				}
				if err := s.store.Create(ctx, plan); err != nil {
					return nil, err
				}
				plan, err = s.store.Get(ctx, plan.Record.ID)
				if err != nil {
					return nil, err
				}
			}
		}
		if plan.Record.Spec.Engine == "postgres" {
			if err := s.operatorsReady(ctx, cluster); err != nil {
				return nil, err
			}
		}
		if request.Action == "backup" {
			if plan.Record.Spec.Engine != "postgres" || plan.Record.Spec.S3DestinationID == "" {
				return nil, fmt.Errorf("base backups require PostgreSQL and an object-store destination")
			}
			plan.Manifest, _, err = kubernetes.DataBackupManifest(&plan.Record)
			if err != nil {
				return nil, err
			}
		}
	}
	payload, err := json.Marshal(plan)
	if err != nil {
		return nil, err
	}
	snapshot, err := s.snapshot(ctx, cluster, plan.Record.ID)
	if err != nil {
		return nil, err
	}
	target := "cluster-data:" + plan.Record.ID
	effects := fmt.Sprintf("%s %s database %s with %d instance(s), %d GiB per instance and generated credentials. Preserve existing PVCs and the source database. Interruptions can leave partially reconciled resources; review recovery to resume. Redis replicas use a fixed persisted primary and do not provide Sentinel failover.", request.Action, plan.Record.Spec.Engine, plan.Record.Spec.Name, plan.Record.Spec.Instances, plan.Record.Spec.StorageGiB)
	if request.Action == "restore" {
		effects += " Recover PostgreSQL into a separate target at " + plan.RestoreTime + "; source data is preserved."
	}
	if request.Action == "operators" {
		target = "cluster:" + clusterID
		effects = "Install reviewed cert-manager, CloudNativePG and Barman Cloud manifests, including cluster-wide RBAC, CRDs and admission services. Refuse resources owned by another installation. Partial installations remain for an explicit reviewed retry."
		for _, manifest := range plan.Operators {
			effects += " " + manifest.Name + " SHA256 " + manifest.SHA256
		}
	}
	return s.operations.Review(ctx, user, project, "cluster-data", target, string(payload), snapshot, effects)
}
func (s *Service) destination(ctx context.Context, id string) (*models.S3Destination, error) {
	if id == "" {
		return nil, nil
	}
	return s.destinations.GetS3Destination(ctx, id)
}
func postgresMajor(image string) string {
	_, tag, _ := strings.Cut(image, ":")
	major, _, _ := strings.Cut(tag, ".")
	return major
}
func (s *Service) snapshot(ctx context.Context, cluster *models.Cluster, id string) (string, error) {
	config := ""
	if id != "" {
		plan, err := s.store.Get(ctx, id)
		if err != nil {
			return "", err
		}
		config = plan.Record.Config
	}
	data, err := json.Marshal([]any{cluster.ID, cluster.ProjectID, cluster.Revision, cluster.Version, cluster.Nodes, config})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
