package clusterdata

import (
	"codedock/internal/engine/kubernetes"
	"codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type BindingApps interface {
	GetByID(context.Context, string) (*models.AppService, error)
}

type BindingVars interface {
	ListByService(context.Context, string) ([]*models.Variable, error)
	Create(context.Context, *models.Variable) error
}

type ClusterBindingRequest struct {
	AppID      string `json:"appId"`
	DatabaseID string `json:"databaseId"`
}

type ClusterBindingPlan struct {
	AppID      string            `json:"appId"`
	ProjectID  string            `json:"projectId"`
	DatabaseID string            `json:"databaseId"`
	Keys       map[string]string `json:"keys"`
}

func clusterBindingVars(plan *models.ClusterDataPlan) map[string]string {
	vars := make(map[string]string)
	namespace, name := kubernetes.DataIdentity(plan.Record.Spec)
	prefix := plan.Record.Spec.Name
	if prefix == "" {
		prefix = plan.Record.ID
	}
	if plan.Record.Spec.Engine == "postgres" {
		host := name + "-rw." + namespace + ".svc.cluster.local"
		conn := fmt.Sprintf("postgresql://app:%s@%s:5432/app", plan.Password, host)
		vars[prefix+"_DATABASE_URL"] = conn
		vars["DATABASE_URL"] = conn
		vars[prefix+"_HOST"] = host
		vars[prefix+"_PASSWORD"] = plan.Password
	} else {
		host := name + "." + namespace + ".svc.cluster.local"
		conn := fmt.Sprintf("redis://:%s@%s:6379", plan.Password, host)
		vars[prefix+"_REDIS_URL"] = conn
		vars["REDIS_URL"] = conn
		vars[prefix+"_HOST"] = host
		vars[prefix+"_PASSWORD"] = plan.Password
	}
	return vars
}

func (s *Service) ReviewBinding(ctx context.Context, user, project, cluster, appID, databaseID string, apps BindingApps, vars BindingVars) (*models.OperationReview, error) {
	if appID == "" || databaseID == "" {
		return nil, fmt.Errorf("select an application and a cluster database")
	}
	app, err := apps.GetByID(ctx, appID)
	if err != nil {
		return nil, err
	}
	if app.ProjectID != project {
		return nil, fmt.Errorf("application belongs to another project")
	}
	plan, err := s.owned(ctx, project, cluster, databaseID)
	if err != nil {
		return nil, err
	}
	if plan.Record.Status != "READY" {
		return nil, fmt.Errorf("database is not ready")
	}
	keys := clusterBindingVars(plan)
	masked := make(map[string]string, len(keys))
	for key := range keys {
		masked[key] = "***"
	}
	payload, err := json.Marshal(ClusterBindingPlan{AppID: appID, ProjectID: project, DatabaseID: databaseID, Keys: keys})
	if err != nil {
		return nil, err
	}
	existing, err := vars.ListByService(ctx, appID)
	if err != nil {
		return nil, err
	}
	replaced := 0
	for _, variable := range existing {
		if _, ok := keys[variable.Key]; ok {
			replaced++
		}
	}
	snapshotData, err := json.Marshal([]any{app.ID, app.ProjectID, app.UpdatedAt, plan.Record.ID, plan.Record.Status, plan.Record.UpdatedAt, len(existing)})
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(snapshotData)
	effects := fmt.Sprintf("Bind cluster database %s to application %s by writing %d connection variables (%d replace existing values). Credentials are stored as service variables; redeploy to apply them.", plan.Record.Spec.Name, app.Name, len(keys), replaced)
	_ = masked
	return s.operations.Review(ctx, user, project, "cluster-binding", appID, string(payload), hex.EncodeToString(sum[:]), effects)
}

func (s *Service) ApplyBinding(ctx context.Context, user, id, confirmation string, apps BindingApps, vars BindingVars) error {
	op, err := s.operations.Get(ctx, id)
	if err != nil {
		return err
	}
	if op.Kind != "cluster-binding" || op.UserID != user {
		return fmt.Errorf("binding operation not found")
	}
	var plan ClusterBindingPlan
	if err := json.Unmarshal([]byte(op.Payload), &plan); err != nil {
		return err
	}
	app, err := apps.GetByID(ctx, plan.AppID)
	if err != nil {
		return err
	}
	if app.ProjectID != plan.ProjectID || app.ProjectID != op.ProjectID {
		return fmt.Errorf("application project changed")
	}
	stored, err := s.store.Get(ctx, plan.DatabaseID)
	if err != nil {
		return err
	}
	if stored.Record.ProjectID != plan.ProjectID || stored.Record.Status != "READY" {
		return fmt.Errorf("database is not ready")
	}
	fresh := clusterBindingVars(stored)
	for key, value := range fresh {
		plan.Keys[key] = value
	}
	existing, err := vars.ListByService(ctx, plan.AppID)
	if err != nil {
		return err
	}
	snapshotData, err := json.Marshal([]any{app.ID, app.ProjectID, app.UpdatedAt, stored.Record.ID, stored.Record.Status, stored.Record.UpdatedAt, len(existing)})
	if err != nil {
		return err
	}
	sum := sha256.Sum256(snapshotData)
	return s.operations.Apply(ctx, id, user, confirmation, hex.EncodeToString(sum[:]), func(ctx context.Context, op *models.Operation, progress func(string, string) error) error {
		if err := progress("CONFIGURING", "Writing reviewed cluster database bindings"); err != nil {
			return err
		}
		for key, value := range plan.Keys {
			if strings.TrimSpace(key) == "" {
				return fmt.Errorf("invalid binding variable")
			}
			if err := vars.Create(ctx, &models.Variable{ServiceID: plan.AppID, Key: key, Value: value, IsSecret: true}); err != nil {
				return err
			}
		}
		return progress("VERIFIED", fmt.Sprintf("Wrote %d reviewed binding variables; redeploy to apply them", len(plan.Keys)))
	})
}
