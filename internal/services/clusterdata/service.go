package clusterdata

import (
	"codedock.run/codedock/internal/engine/kubernetes"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services/operations"
	"context"
	"fmt"
	"net/http"
	"time"
)

type Store interface {
	Create(context.Context, *models.ClusterDataPlan) error
	Get(context.Context, string) (*models.ClusterDataPlan, error)
	List(context.Context, string) ([]models.ClusterData, error)
	Observe(context.Context, string, string, string) error
	Recover(context.Context) error
}
type Clusters interface {
	Get(context.Context, string) (*models.Cluster, error)
}
type Environments interface {
	Get(context.Context, string) (*models.EnvironmentConfig, error)
}
type Destinations interface {
	GetS3Destination(context.Context, string) (*models.S3Destination, error)
}
type Gate interface{ AcquireVolume(string) (func(), error) }
type Service struct {
	store        Store
	clusters     Clusters
	environments Environments
	destinations Destinations
	engine       *kubernetes.WorkloadRuntime
	commands     kubernetes.Commands
	operations   *operations.Service
	gate         Gate
	http         *http.Client
}

func NewService(store Store, clusters Clusters, environments Environments, destinations Destinations, commands kubernetes.Commands, operations *operations.Service, gate Gate, http *http.Client) *Service {
	return &Service{store, clusters, environments, destinations, kubernetes.NewWorkloadRuntime(commands), commands, operations, gate, http}
}
func (s *Service) Recover(ctx context.Context) error { return s.store.Recover(ctx) }
func (s *Service) cluster(ctx context.Context, project, id string) (*models.Cluster, error) {
	cluster, err := s.clusters.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if cluster.ProjectID != project || cluster.Status != "READY" || len(cluster.Nodes) == 0 {
		return nil, fmt.Errorf("select a ready cluster in this project")
	}
	return cluster, nil
}
func (s *Service) List(ctx context.Context, project, cluster string) ([]models.ClusterData, error) {
	target, err := s.cluster(ctx, project, cluster)
	if err != nil {
		return nil, err
	}
	records, err := s.store.List(ctx, cluster)
	if err != nil {
		return nil, err
	}
	observation, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	for index := range records {
		if records[index].Status != "READY" {
			continue
		}
		if err := s.engine.DataObserved(observation, target.Nodes[0], &records[index]); err != nil {
			records[index].Status = "DEGRADED"
			records[index].Error = fmt.Sprintf("Observed database state: %v", err)
			continue
		}
		roles, volumes, err := s.engine.DataInspect(observation, target.Nodes[0], &records[index])
		if err != nil {
			records[index].Status = "DEGRADED"
			records[index].Error = fmt.Sprintf("Observed database roles: %v", err)
			continue
		}
		records[index].ObservedRoles = roles
		records[index].ObservedVolumes = volumes
	}
	return records, nil
}
func (s *Service) owned(ctx context.Context, project, cluster, id string) (*models.ClusterDataPlan, error) {
	plan, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if plan.Record.ProjectID != project || plan.Record.ClusterID != cluster {
		return nil, fmt.Errorf("database target belongs to another project or cluster")
	}
	return plan, nil
}
func (s *Service) Credentials(ctx context.Context, project, cluster, id string) (*models.ClusterDataCredentials, error) {
	plan, err := s.owned(ctx, project, cluster, id)
	if err != nil {
		return nil, err
	}
	if plan.Record.Status != "READY" {
		return nil, fmt.Errorf("database is not ready")
	}
	namespace, name := kubernetes.DataIdentity(plan.Record.Spec)
	port := 6379
	username := "default"
	if plan.Record.Spec.Engine == "postgres" {
		port = 5432
		name += "-rw"
		username = "app"
	}
	return &models.ClusterDataCredentials{Username: username, Password: plan.Password, Host: name + "." + namespace + ".svc.cluster.local", Port: port}, nil
}
func (s *Service) Backups(ctx context.Context, project, clusterID, id string) ([]map[string]any, error) {
	plan, err := s.owned(ctx, project, clusterID, id)
	if err != nil {
		return nil, err
	}
	cluster, err := s.cluster(ctx, project, clusterID)
	if err != nil {
		return nil, err
	}
	if plan.Record.Spec.Engine != "postgres" {
		return nil, fmt.Errorf("Redis uses persistent AOF recovery")
	}
	return s.engine.DataBackups(ctx, cluster.Nodes[0], &plan.Record)
}

func (s *Service) VerifyConnection(ctx context.Context, project, clusterID, id string) error {
	plan, err := s.owned(ctx, project, clusterID, id)
	if err != nil {
		return err
	}
	if plan.Record.Status != "READY" {
		return fmt.Errorf("database is not ready")
	}
	cluster, err := s.cluster(ctx, project, clusterID)
	if err != nil {
		return err
	}
	check, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.engine.DataConnectionCheck(check, cluster.Nodes[0], &plan.Record, plan.Password)
}
