package projects

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"time"

	"codedock.run/codedock/internal/engine/compose"
	"codedock.run/codedock/internal/models"
	"github.com/google/uuid"
)

type ComposeStackStore interface {
	List(context.Context, string) ([]models.ComposeStack, error)
	Get(context.Context, string, string) (*models.ComposeStack, error)
	Save(context.Context, *models.ComposeStack, int) error
	Claim(context.Context, string, string, int) error
	Observe(context.Context, string, string, string, string) error
	Recover(context.Context) error
}

type ComposeStackRuntime interface {
	Canonical(context.Context, string, string, map[string]string) (string, error)
	ValidatePorts(context.Context, string, string) error
	Apply(context.Context, *models.ComposeStack, func(string) error) error
	Results(context.Context, string) ([]models.ComposeServiceResult, error)
}

type ComposeStackService struct {
	store   ComposeStackStore
	runtime ComposeStackRuntime
	mu      sync.Mutex
	running map[string]context.CancelFunc
}

func NewComposeStackService(store ComposeStackStore, runtime ComposeStackRuntime) *ComposeStackService {
	return &ComposeStackService{store: store, runtime: runtime, running: map[string]context.CancelFunc{}}
}

func (s *ComposeStackService) Recover(ctx context.Context) error { return s.store.Recover(ctx) }
func (s *ComposeStackService) List(ctx context.Context, project string) ([]models.ComposeStack, error) {
	return s.store.List(ctx, project)
}

func (s *ComposeStackService) Review(ctx context.Context, request models.ComposeStackRequest) (*models.ComposeReview, error) {
	if _, err := uuid.Parse(request.ID); err != nil {
		return nil, fmt.Errorf("stack ID must be a UUID")
	}
	if strings.TrimSpace(request.Name) == "" || len(request.Name) > 100 {
		return nil, fmt.Errorf("stack name is required and must not exceed 100 characters")
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	source, err := compose.ResolveRepositoryContexts(request)
	if err != nil {
		return nil, err
	}
	config, err := s.runtime.Canonical(ctx, request.ID, source, request.Variables)
	if err != nil {
		return nil, err
	}
	if err := s.runtime.ValidatePorts(ctx, request.ID, config); err != nil {
		return nil, err
	}
	var document struct {
		Services map[string]json.RawMessage `json:"services"`
	}
	if err := json.Unmarshal([]byte(config), &document); err != nil {
		return nil, err
	}
	services := []string{}
	for name := range document.Services {
		services = append(services, name)
	}
	sort.Strings(services)
	return &models.ComposeReview{Config: config, Digest: composeReviewDigest(request, config), Services: services, Effects: []string{"Save the complete resolved configuration, encrypted on disk.", "Build and pull all images before changing existing containers.", "Start services in dependency order and wait for health checks.", "Reuse owned named volumes and networks on retry. No volumes are deleted.", "Activation failures or cancellation may leave partially updated services; inspect each service before retrying.", "Services removed from the configuration remain running until explicitly removed."}}, nil
}

func composeReviewDigest(request models.ComposeStackRequest, config string) string {
	digest := sha256.Sum256([]byte(fmt.Sprintf("%q|%q|%q|%q|%d", request.ID, request.EnvironmentID, request.Name, config, request.Revision)))
	return hex.EncodeToString(digest[:])
}

func (s *ComposeStackService) Save(ctx context.Context, project string, request models.ComposeStackRequest) (*models.ComposeStack, error) {
	review, err := s.Review(ctx, request)
	if err != nil {
		return nil, err
	}
	if request.Digest != review.Digest {
		return nil, fmt.Errorf("configuration changed since review; review again before saving")
	}
	stack := &models.ComposeStack{ID: request.ID, ProjectID: project, EnvironmentID: request.EnvironmentID, Name: request.Name, Config: review.Config, Status: "saved", Results: "[]"}
	if err := s.store.Save(ctx, stack, request.Revision); err != nil {
		return nil, err
	}
	return stack, nil
}

func (s *ComposeStackService) Deploy(ctx context.Context, project, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running[id] != nil {
		return fmt.Errorf("stack already deploying")
	}
	stack, err := s.store.Get(ctx, project, id)
	if err != nil {
		return err
	}
	if err := s.store.Claim(ctx, project, id, stack.Revision); err != nil {
		return err
	}
	operation, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	s.running[id] = cancel
	go func() {
		defer cancel()
		defer func() { s.mu.Lock(); delete(s.running, id); s.mu.Unlock() }()
		err := s.runtime.Apply(operation, stack, func(phase string) error { return s.store.Observe(operation, id, phase, "", "[]") })
		status, message := "READY", ""
		if err != nil {
			status, message = "FAILED", err.Error()
			if operation.Err() != nil {
				status = "INTERRUPTED"
			}
		}
		finalCtx, finish := context.WithTimeout(context.Background(), 15*time.Second)
		defer finish()
		results, observedErr := s.runtime.Results(finalCtx, id)
		if observedErr == nil && status == "READY" {
			if validationErr := validateStackResults(stack.Config, results); validationErr != nil {
				status, message = "DEGRADED", validationErr.Error()
			}
		}
		if observedErr != nil {
			status, message = "OBSERVATION_FAILED", message+" Runtime observation unavailable: "+observedErr.Error()
		}
		data, marshalErr := json.Marshal(results)
		if marshalErr != nil {
			status, message = "FAILED", marshalErr.Error()
			data = []byte("[]")
		}
		if err := s.store.Observe(finalCtx, id, status, message, string(data)); err != nil {
			slog.Error("save stack result", "stack", id, "error", err)
		}
	}()
	return nil
}

func (s *ComposeStackService) Cancel(project, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cancel := s.running[id]
	if cancel == nil {
		return fmt.Errorf("stack is not deploying")
	}
	cancel()
	return nil
}
