package deploy

import (
	"codedock.run/codedock/internal/engine/observability"
	"codedock.run/codedock/internal/models"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type scalingTestApps struct{ app *models.AppService }

func (s scalingTestApps) ListAll(context.Context) ([]*models.AppService, error) {
	return []*models.AppService{s.app}, nil
}

type scalingTestMetrics struct {
	cpu   float64
	calls int
}

func (s *scalingTestMetrics) GetHealth(context.Context, string) (*observability.ContainerHealth, error) {
	s.calls++
	return &observability.ContainerHealth{CPUUsagePercentage: s.cpu}, nil
}

type scalingTestStore struct {
	policies   []*models.AutoscalingPolicy
	reserved   bool
	from, to   int
	rolledBack bool
	decision   string
}

func (s *scalingTestStore) ListEnabled(context.Context) ([]*models.AutoscalingPolicy, error) {
	return s.policies, nil
}
func (s *scalingTestStore) Observe(_ context.Context, _ *models.AutoscalingPolicy, _ float64, reason string, _ bool) error {
	s.decision = reason
	return nil
}
func (s *scalingTestStore) SetReplicas(_ context.Context, _ *models.AutoscalingPolicy, from, to int) (bool, error) {
	s.from = from
	s.to = to
	return s.reserved, nil
}
func (s *scalingTestStore) RollbackReplicas(context.Context, string, int, int) error {
	s.rolledBack = true
	return nil
}

type scalingTestDeployments struct {
	failure  bool
	executed bool
}

func (s *scalingTestDeployments) CreateDeployment(_ context.Context, d *models.Deployment) (*models.Deployment, error) {
	if s.failure {
		return nil, errors.New("queue unavailable")
	}
	d.ID = "deployment"
	return d, nil
}
func (s *scalingTestDeployments) ExecuteDeploymentAsync(*models.Deployment) { s.executed = true }

func TestAutoscalingOptInLimitsAndCooldown(t *testing.T) {
	now := time.Now()
	for _, test := range []struct {
		name    string
		enabled bool
		current int
		cpu     float64
		last    string
		want    int
	}{
		{"disabled", false, 2, 99, "", 2}, {"scale up", true, 2, 99, "", 3}, {"scale down", true, 2, 0, "", 1},
		{"max", true, 5, 99, "", 5}, {"min", true, 1, 0, "", 1}, {"minimum reconciliation", true, 0, 50, "", 1},
		{"cooldown", true, 2, 99, now.Format(time.RFC3339Nano), 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			p := models.DefaultAutoscalingPolicy("app")
			p.Enabled = test.enabled
			p.LastScaledAt = test.last
			target, _ := scalingTarget(p, test.current, test.cpu, now)
			if target != test.want {
				t.Fatalf("target %d, want %d", target, test.want)
			}
		})
	}
	metrics := &scalingTestMetrics{cpu: 99}
	store := &scalingTestStore{}
	worker := NewAutoscalerWorker(scalingTestApps{}, metrics, &scalingTestDeployments{}, store)
	worker.checkAndScale(context.Background())
	if metrics.calls != 0 {
		t.Fatal("disabled services were sampled")
	}
}
func TestAutoscalingReservationAndFailedQueue(t *testing.T) {
	for _, failure := range []bool{false, true} {
		policy := models.DefaultAutoscalingPolicy("app")
		policy.Enabled = true
		store := &scalingTestStore{policies: []*models.AutoscalingPolicy{policy}, reserved: true}
		deployments := &scalingTestDeployments{failure: failure}
		app := &models.AppService{ID: "app", Status: models.AppServiceStatusRunning, ContainerID: "container", Replicas: 2}
		worker := NewAutoscalerWorker(scalingTestApps{app}, &scalingTestMetrics{cpu: 95}, deployments, store)
		worker.checkAndScale(context.Background())
		if store.from != 2 || store.to != 3 {
			t.Fatal("unexpected reservation")
		}
		if failure && (!store.rolledBack || deployments.executed || !strings.Contains(store.decision, "failed")) {
			t.Fatal("failed queue changed replicas or hid its failure")
		}
		if !failure && (!deployments.executed || !strings.Contains(store.decision, "deployment")) {
			t.Fatal("successful scaling was not observable")
		}
	}
}
