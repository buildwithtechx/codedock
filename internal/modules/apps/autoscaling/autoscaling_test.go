package autoscaling

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockAutoRepo struct {
	policies map[string]*AutoscalingPolicy
}

func (m *mockAutoRepo) GetByServiceID(ctx context.Context, serviceID string) (*AutoscalingPolicy, error) {
	return m.policies[serviceID], nil
}

func (m *mockAutoRepo) Upsert(ctx context.Context, policy *AutoscalingPolicy) error {
	m.policies[policy.ServiceID] = policy
	return nil
}

func TestAutoscalingPolicyDefaults(t *testing.T) {
	policy := DefaultAutoscalingPolicy("svc-abc")
	if policy.MinReplicas != 1 {
		t.Fatalf("expected min replicas 1, got %d", policy.MinReplicas)
	}
	if policy.MaxReplicas != 5 {
		t.Fatalf("expected max replicas 5, got %d", policy.MaxReplicas)
	}
	if policy.ScaleUpCPU != 80 {
		t.Fatalf("expected scale up cpu 80, got %f", policy.ScaleUpCPU)
	}
	if policy.ScaleDownCPU != 20 {
		t.Fatalf("expected scale down cpu 20, got %f", policy.ScaleDownCPU)
	}
	if policy.CooldownSeconds != 300 {
		t.Fatalf("expected cooldown seconds 300, got %d", policy.CooldownSeconds)
	}
}

func TestAutoscalingHandler(t *testing.T) {
	e := echo.New()
	repo := &mockAutoRepo{policies: make(map[string]*AutoscalingPolicy)}
	svc := NewAutoscalingService(repo, nil, nil)
	h := NewAutoscalingHandler(svc)

	t.Run("Get missing serviceId", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/services//autoscaling", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if err := h.Get(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("Get default policy", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/services/svc-1/autoscaling", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("serviceId")
		c.SetParamValues("svc-1")
		if err := h.Get(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("Update policy", func(t *testing.T) {
		body := `{"enabled": true, "minReplicas": 2, "maxReplicas": 10}`
		req := httptest.NewRequest(http.MethodPut, "/services/svc-1/autoscaling", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("serviceId")
		c.SetParamValues("svc-1")
		if err := h.Update(c); err != nil {
			t.Fatal(err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}
