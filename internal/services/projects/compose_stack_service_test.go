package projects

import (
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"
	"testing"
	"time"
)

type stackTestStore struct {
	ComposeStackStore
	stack    *models.ComposeStack
	saved    bool
	claimed  int
	finished chan string
}

func (s *stackTestStore) Get(context.Context, string, string) (*models.ComposeStack, error) {
	return s.stack, nil
}
func (s *stackTestStore) Save(context.Context, *models.ComposeStack, int) error {
	s.saved = true
	return nil
}
func (s *stackTestStore) Claim(_ context.Context, _ string, _ string, revision int) error {
	s.claimed = revision
	return nil
}
func (s *stackTestStore) Observe(_ context.Context, _ string, status, _ string, _ string) error {
	if status == "FAILED" || status == "INTERRUPTED" {
		s.finished <- status
	}
	return nil
}

type stackTestRuntime struct{ cancel bool }

func (r *stackTestRuntime) Canonical(context.Context, string, string, map[string]string) (string, error) {
	return `{"services":{"web":{"image":"nginx"}}}`, nil
}
func (r *stackTestRuntime) ValidatePorts(context.Context, string, string) error { return nil }
func (r *stackTestRuntime) Results(context.Context, string) ([]models.ComposeServiceResult, error) {
	return []models.ComposeServiceResult{{Name: "web", ContainerID: "previous", State: "running", Health: "healthy"}}, nil
}
func (r *stackTestRuntime) Apply(ctx context.Context, _ *models.ComposeStack, progress func(string) error) error {
	if err := progress("BUILDING"); err != nil {
		return err
	}
	if r.cancel {
		<-ctx.Done()
		return ctx.Err()
	}
	return fmt.Errorf("build failed; previous workload preserved")
}

func TestStackSaveRequiresExactReviewedConfiguration(t *testing.T) {
	store := &stackTestStore{}
	service := NewComposeStackService(store, &stackTestRuntime{})
	request := models.ComposeStackRequest{ID: "ca5a0f00-5733-4bc4-8e8c-a1b9a2f03a5a", Name: "Stack", EnvironmentID: "environment", Content: "services: {}"}
	review, err := service.Review(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	request.Digest = review.Digest
	request.EnvironmentID = "other"
	if _, err := service.Save(context.Background(), "project", request); err == nil || store.saved {
		t.Fatal("a review was reused for a different environment")
	}
}

func TestStackReportsFailedBuildAndCancelledOperations(t *testing.T) {
	for _, cancelled := range []bool{false, true} {
		store := &stackTestStore{stack: &models.ComposeStack{ID: "stack", Revision: 3}, finished: make(chan string, 1)}
		service := NewComposeStackService(store, &stackTestRuntime{cancel: cancelled})
		if err := service.Deploy(context.Background(), "project", "stack"); err != nil {
			t.Fatal(err)
		}
		if store.claimed != 3 {
			t.Fatal("deployment did not claim the reviewed revision")
		}
		expected := "FAILED"
		if cancelled {
			expected = "INTERRUPTED"
			if err := service.Cancel("project", "stack"); err != nil {
				t.Fatal(err)
			}
		}
		select {
		case status := <-store.finished:
			if status != expected {
				t.Fatalf("incorrect terminal state: %s", status)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("operation did not finish")
		}
	}
}

func TestStackBackupFailurePreventsRuntimeActivation(t *testing.T) {
	store := &stackTestStore{stack: &models.ComposeStack{ID: "stack", ProjectID: "project", Revision: 1}, finished: make(chan string, 1)}
	service := NewComposeStackService(store, &stackTestRuntime{cancel: true})
	service.BeforeDeployment = func(ctx context.Context, project, service string) error {
		if project != "project" || service != "" {
			t.Error("incorrect backup scope")
		}
		return fmt.Errorf("required backup failed")
	}
	if err := service.Deploy(context.Background(), "project", "stack"); err != nil {
		t.Fatal(err)
	}
	select {
	case status := <-store.finished:
		if status != "FAILED" {
			t.Fatal(status)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("runtime activated despite failed backup")
	}
}
