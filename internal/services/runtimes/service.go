package runtimes

import (
	"codedock.run/codedock/internal/engine/deploy"
	"codedock.run/codedock/internal/engine/kubernetes"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services/operations"
	"context"
	"fmt"
	"io"
	"sort"
)

type Store interface {
	Get(context.Context, string) (*models.ServiceRuntime, error)
	Save(context.Context, *models.ServiceRuntime, int) error
	Begin(context.Context, string, int, string) error
	Observe(context.Context, string, string, string, bool) error
	Pending(context.Context) ([]string, error)
}
type Apps interface {
	GetByID(context.Context, string) (*models.AppService, error)
	Update(context.Context, *models.AppService) error
}
type Clusters interface {
	Get(context.Context, string) (*models.Cluster, error)
}
type Servers interface {
	GetByID(context.Context, string) (*models.Server, error)
}
type Gate interface{ AcquireVolume(string) (func(), error) }
type Service struct {
	store      Store
	apps       Apps
	clusters   Clusters
	servers    Servers
	engine     *kubernetes.WorkloadRuntime
	builder    *deploy.Deployer
	operations *operations.Service
	gate       Gate
}

func NewService(store Store, apps Apps, clusters Clusters, servers Servers, engine *kubernetes.WorkloadRuntime, builder *deploy.Deployer, operations *operations.Service, gate Gate) *Service {
	return &Service{store, apps, clusters, servers, engine, builder, operations, gate}
}
func (s *Service) Get(ctx context.Context, id string) (*models.ServiceRuntime, error) {
	return s.store.Get(ctx, id)
}
func (s *Service) Handles(ctx context.Context, id string) (bool, error) {
	target, err := s.store.Get(ctx, id)
	if err != nil {
		return false, err
	}
	return target.Target.Kind == "kubernetes", nil
}
func (s *Service) cluster(ctx context.Context, app *models.AppService, target models.RuntimeTarget) (*models.Cluster, error) {
	cluster, err := s.clusters.Get(ctx, target.ClusterID)
	if err != nil {
		return nil, err
	}
	if cluster.ProjectID != app.ProjectID || cluster.Status != "READY" || len(cluster.Nodes) == 0 {
		return nil, fmt.Errorf("select a ready cluster belonging to this project")
	}
	nodes := map[string]bool{}
	for _, node := range cluster.Nodes {
		if err := kubernetes.ValidateNode(node); err != nil {
			return nil, err
		}
		server, err := s.servers.GetByID(ctx, node.ServerID)
		if err != nil {
			return nil, err
		}
		if server.OrganizationID != cluster.OrganizationID {
			return nil, fmt.Errorf("cluster server organization changed")
		}
		nodes[node.ServerID] = true
	}
	selected := map[string]bool{}
	for _, id := range target.NodeIDs {
		if !nodes[id] || selected[id] {
			return nil, fmt.Errorf("placement must select distinct nodes belonging to this cluster")
		}
		selected[id] = true
	}
	return cluster, nil
}
func (s *Service) lockCluster(cluster *models.Cluster) (func(), error) {
	keys := []string{"cluster:" + cluster.ID}
	for _, node := range cluster.Nodes {
		keys = append(keys, "server:"+node.ServerID)
	}
	sort.Strings(keys)
	releases := []func(){}
	releaseAll := func() {
		for i := len(releases) - 1; i >= 0; i-- {
			releases[i]()
		}
	}
	for _, key := range keys {
		release, err := s.gate.AcquireVolume(key)
		if err != nil {
			releaseAll()
			return nil, err
		}
		releases = append(releases, release)
	}
	return releaseAll, nil
}
func (s *Service) Deploy(ctx context.Context, app *models.AppService, source string, logs io.Writer) (string, error) {
	target, err := s.store.Get(ctx, app.ID)
	if err != nil {
		return "", err
	}
	cluster, err := s.cluster(ctx, app, target.Target)
	if err != nil {
		return "", err
	}
	release, err := s.lockCluster(cluster)
	if err != nil {
		return "", err
	}
	defer release()
	cluster, err = s.cluster(ctx, app, target.Target)
	if err != nil {
		return "", err
	}
	if err := s.engine.VerifyCluster(ctx, cluster, true); err != nil {
		return "", err
	}
	workload, err := s.builder.PrepareClusterWorkload(ctx, app, target.Target, source, logs)
	if err != nil {
		return "", err
	}
	if err := kubernetes.ValidateWorkload(workload); err != nil {
		return "", err
	}
	node := cluster.Nodes[0]
	if err := s.engine.Namespace(ctx, node, app); err != nil {
		return "", err
	}
	journal, err := s.engine.Snapshot(ctx, node, workload)
	if err != nil {
		return "", err
	}
	manifest, err := kubernetes.WorkloadManifest(workload)
	if err != nil {
		return "", err
	}
	if err := s.engine.Validate(ctx, node, manifest); err != nil {
		return "", err
	}
	if err := s.store.Begin(ctx, app.ID, target.Revision, journal); err != nil {
		return "", err
	}
	err = s.engine.Apply(ctx, node, manifest)
	if err == nil {
		err = s.engine.Ready(ctx, node, &workload.App)
	}
	if err != nil {
		return "", s.recoverFailure(app, node, journal, err)
	}
	app.Status = models.AppServiceStatusRunning
	app.ContainerID = "kubernetes:" + app.ID
	if err := s.apps.Update(ctx, app); err != nil {
		return "", s.recoverFailure(app, node, journal, err)
	}
	desiredStore, ok := s.store.(DesiredStore)
	if !ok {
		return "", s.recoverFailure(app, node, journal, fmt.Errorf("desired workload storage unavailable"))
	}
	if err := desiredStore.CommitDesired(ctx, app.ID, target.Revision, &models.DesiredRuntime{Revision: target.Revision, Workload: *workload, Manifest: manifest}); err != nil {
		return "", s.recoverFailure(app, node, journal, err)
	}
	return app.ContainerID, nil
}
