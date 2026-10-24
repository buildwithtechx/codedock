package runtimes

import (
	"codedock/internal/engine/bare"
	"codedock/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"time"
)

type NativeProjects interface {
	Get(context.Context, string) (*models.ProjectConfig, error)
}

func (s *Service) SetBare(runtime *bare.Runtime, projects NativeProjects) {
	s.native = runtime
	s.projects = projects
}
func (s *Service) validateBare(ctx context.Context, app *models.AppService, target models.RuntimeTarget) error {
	if s.native == nil || s.projects == nil {
		return fmt.Errorf("native runtime unavailable")
	}
	if err := bare.ValidateIdentity(app, target); err != nil {
		return err
	}
	project, err := s.projects.Get(ctx, app.ProjectID)
	if err != nil {
		return err
	}
	server, err := s.servers.GetByID(ctx, target.BareNode.ServerID)
	if err != nil {
		return err
	}
	if server.IsLocal || server.OrganizationID != project.OrganizationID {
		return fmt.Errorf("native server must belong to the application organization")
	}
	return nil
}
func (s *Service) nativeTarget(ctx context.Context, id string) (*models.AppService, *models.ServiceRuntime, error) {
	runtime, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if runtime.Target.Kind != "bare" {
		return nil, runtime, nil
	}
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if err := s.validateBare(ctx, app, runtime.Target); err != nil {
		return nil, nil, err
	}
	return app, runtime, nil
}
func (s *Service) deployBare(ctx context.Context, app *models.AppService, runtime *models.ServiceRuntime, logs io.Writer) (string, error) {
	if err := bare.Validate(app, runtime.Target); err != nil {
		return "", err
	}
	if err := s.validateBare(ctx, app, runtime.Target); err != nil {
		return "", err
	}
	unlock, err := s.gate.AcquireVolume("server:" + runtime.Target.BareNode.ServerID)
	if err != nil {
		return "", err
	}
	defer unlock()
	if runtime.Journal != "" {
		return "", fmt.Errorf("recover the interrupted native release before deploying")
	}
	variables, err := s.builder.NativeEnvironment(ctx, app, logs)
	if err != nil {
		return "", err
	}
	journal, err := json.Marshal(models.NativeJournal{PreviousApp: *app, ReleaseID: uuid.NewString()})
	if err != nil {
		return "", err
	}
	runtime.Journal = string(journal)
	var activation models.NativeJournal
	if err := json.Unmarshal(journal, &activation); err != nil {
		return "", err
	}
	if err := s.store.Begin(ctx, app.ID, runtime.Revision, runtime.Journal); err != nil {
		return "", err
	}
	if err := s.native.Deploy(ctx, app, runtime.Target, variables, activation.ReleaseID); err != nil {
		return "", s.recoverBare(app, runtime, err)
	}
	app.Status = models.AppServiceStatusRunning
	app.ContainerID = "bare:" + app.ID + ":" + activation.ReleaseID
	if err := s.apps.Update(ctx, app); err != nil {
		return "", s.recoverBare(app, runtime, err)
	}
	if err := s.native.Finalize(ctx, app, runtime.Target); err != nil {
		return "", s.recoverBare(app, runtime, err)
	}
	if err := s.store.Observe(ctx, app.ID, "READY", "", true); err != nil {
		return "", err
	}
	return app.ContainerID, nil
}
func (s *Service) recoverBare(app *models.AppService, runtime *models.ServiceRuntime, cause error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()
	rollback := s.native.Recover(ctx, app, runtime.Target)
	if rollback == nil {
		rollback = s.restoreNativeApp(ctx, app, runtime)
	}
	status := "FAILED"
	if rollback != nil {
		status = "RECOVERY_REQUIRED"
	}
	persist := s.store.Observe(ctx, app.ID, status, errors.Join(cause, rollback).Error(), rollback == nil)
	return errors.Join(cause, rollback, persist)
}
func (s *Service) nativeLifecycle(ctx context.Context, app *models.AppService, runtime *models.ServiceRuntime, action string) error {
	if runtime.Journal != "" {
		return fmt.Errorf("recover the interrupted native release before lifecycle changes")
	}
	unlock, err := s.gate.AcquireVolume("server:" + runtime.Target.BareNode.ServerID)
	if err != nil {
		return err
	}
	defer unlock()
	if err := s.native.Lifecycle(ctx, app, runtime.Target, action); err != nil {
		return err
	}
	if action == "stop" || action == "remove" {
		app.Status = models.AppServiceStatusStopped
	} else {
		app.Status = models.AppServiceStatusRunning
	}
	return s.apps.Update(ctx, app)
}

func (s *Service) restoreNativeApp(ctx context.Context, app *models.AppService, runtime *models.ServiceRuntime) error {
	var journal models.NativeJournal
	if err := json.Unmarshal([]byte(runtime.Journal), &journal); err != nil {
		return err
	}
	if journal.PreviousApp.ID != app.ID || journal.PreviousApp.ProjectID != app.ProjectID {
		return fmt.Errorf("native recovery journal ownership changed")
	}
	release, err := s.native.CurrentRelease(ctx, app, runtime.Target)
	if err != nil {
		return err
	}
	if release == journal.ReleaseID {
		app.Status = models.AppServiceStatusRunning
		app.ContainerID = "bare:" + app.ID + ":" + release
		return s.apps.Update(ctx, app)
	}
	return s.apps.Update(ctx, &journal.PreviousApp)
}
