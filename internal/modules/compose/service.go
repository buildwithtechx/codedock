package compose

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

	dockerplatform "codedock/internal/platform/docker"
	"github.com/google/uuid"
)

type ComposeStackRuntime interface {
	Canonical(context.Context, string, string, map[string]string) (string, error)
	ValidatePorts(context.Context, string, string) error
	Apply(context.Context, *ComposeStack, func(string) error) error
	Results(context.Context, string) ([]ComposeServiceResult, error)
}

type ComposeStackService struct {
	BeforeDeployment func(context.Context, string, string) error
	store            ComposeStackStore
	runtime          ComposeStackRuntime
	mu               sync.Mutex
	running          map[string]context.CancelFunc
}

func NewComposeStackService(store ComposeStackStore, runtime ComposeStackRuntime) *ComposeStackService {
	return &ComposeStackService{store: store, runtime: runtime, running: map[string]context.CancelFunc{}}
}

func (s *ComposeStackService) Recover(ctx context.Context) error { return s.store.Recover(ctx) }
func (s *ComposeStackService) List(ctx context.Context, project string) ([]ComposeStack, error) {
	return s.store.List(ctx, project)
}

func (s *ComposeStackService) Review(ctx context.Context, request ComposeStackRequest) (*ComposeReview, error) {
	if _, err := uuid.Parse(request.ID); err != nil {
		return nil, fmt.Errorf("stack ID must be a UUID")
	}
	if strings.TrimSpace(request.Name) == "" || len(request.Name) > 100 {
		return nil, fmt.Errorf("stack name is required and must not exceed 100 characters")
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	source, err := dockerplatform.ResolveRepositoryContexts(dockerplatform.RepositoryContextRequest{
		RepositoryURL: request.RepositoryURL,
		Branch:        request.Branch,
		RootDirectory: request.RootDirectory,
		Content:       request.Content,
	})
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

	var changes *ComposeSourceChanges
	if request.ProjectID != "" && request.ID != "" {
		if existing, err := s.store.Get(ctx, request.ProjectID, request.ID); err == nil && existing != nil && existing.Config != "" {
			var oldDoc struct {
				Services map[string]json.RawMessage `json:"services"`
			}
			if err := json.Unmarshal([]byte(existing.Config), &oldDoc); err == nil {
				changes = &ComposeSourceChanges{
					Added:   []string{},
					Changed: []string{},
					Removed: []string{},
				}
				for name, current := range document.Services {
					old, exists := oldDoc.Services[name]
					if !exists {
						changes.Added = append(changes.Added, name)
					} else if string(current) != string(old) {
						changes.Changed = append(changes.Changed, name)
					}
				}
				for name := range oldDoc.Services {
					if _, exists := document.Services[name]; !exists {
						changes.Removed = append(changes.Removed, name)
					}
				}
				sort.Strings(changes.Added)
				sort.Strings(changes.Changed)
				sort.Strings(changes.Removed)
			}
		}
	}

	effects := []string{
		"Save the complete resolved configuration, encrypted on disk.",
		"Build and pull all images before changing existing containers.",
		"Start services in dependency order and wait for health checks.",
		"Reuse owned named volumes and networks on retry. No volumes are deleted.",
		"Activation failures or cancellation may leave partially updated services; inspect each service before retrying.",
	}
	if changes != nil && len(changes.Removed) > 0 {
		effects = append(effects, fmt.Sprintf("Removed services (%s) will be stopped, unrouted, and pruned safely.", strings.Join(changes.Removed, ", ")))
	} else {
		effects = append(effects, "Services removed from the configuration are safely retired and unrouted.")
	}
	return &ComposeReview{
		Config:   config,
		Digest:   composeReviewDigest(request, config),
		Services: services,
		Effects:  effects,
		Changes:  changes,
	}, nil
}

func composeReviewDigest(request ComposeStackRequest, config string) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%s:%s:%d:%s", request.ID, request.EnvironmentID, request.Name, request.RepositoryURL, request.Branch, request.Revision, config)))
	return hex.EncodeToString(sum[:])
}

func (s *ComposeStackService) Save(ctx context.Context, project string, request ComposeStackRequest) (*ComposeStack, error) {
	review, err := s.Review(ctx, request)
	if err != nil {
		return nil, err
	}
	if request.Digest == "" || request.Digest != review.Digest {
		return nil, fmt.Errorf("configuration changed during review; re-review before saving")
	}
	stack := &ComposeStack{
		ID:            request.ID,
		ProjectID:     project,
		EnvironmentID: request.EnvironmentID,
		Name:          request.Name,
		Config:        review.Config,
	}
	if err := s.store.Save(ctx, stack, request.Revision); err != nil {
		return nil, err
	}
	return s.store.Get(ctx, project, stack.ID)
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
		var err error
		if s.BeforeDeployment != nil {
			err = s.store.Observe(operation, id, "BACKING_UP", "", "[]")
			if err == nil {
				err = s.BeforeDeployment(operation, stack.ProjectID, "")
			}
		}
		if err == nil {
			err = s.runtime.Apply(operation, stack, func(phase string) error { return s.store.Observe(operation, id, phase, "", "[]") })
		}
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

func (s *ComposeStackService) Restart(ctx context.Context, project, id string) error {
	slog.Info("Restarting stack", "project", project, "id", id)
	return s.Deploy(ctx, project, id)
}

func (s *ComposeStackService) Stop(ctx context.Context, project, id string) error {
	slog.Info("Stopping stack", "project", project, "id", id)
	return s.store.Observe(ctx, id, "STOPPED", "Stack stopped manually", "")
}
