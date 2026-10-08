package operations

import (
	"codedock.run/codedock/internal/models"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"sync"
	"time"
)

type Store interface {
	Create(context.Context, *models.Operation) error
	Get(context.Context, string) (*models.Operation, error)
	List(context.Context, string, string) ([]models.Operation, error)
	Claim(context.Context, string, string) error
	Observe(context.Context, string, string, string, string, string) error
	Recover(context.Context) error
}

type Service struct {
	store   Store
	mu      sync.Mutex
	running map[string]context.CancelFunc
}

func NewService(store Store) *Service {
	return &Service{store: store, running: map[string]context.CancelFunc{}}
}

func (s *Service) Recover(ctx context.Context) error { return s.store.Recover(ctx) }

func (s *Service) Get(ctx context.Context, id string) (*models.Operation, error) {
	return s.store.Get(ctx, id)
}

func (s *Service) List(ctx context.Context, project, user string) ([]models.Operation, error) {
	return s.store.List(ctx, project, user)
}

func (s *Service) Review(ctx context.Context, user, project, kind, target, payload, snapshot, effects string) (*models.OperationReview, error) {
	if user == "" || kind == "" || target == "" || snapshot == "" {
		return nil, fmt.Errorf("review requires an owner, operation, target and current snapshot")
	}
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("generate confirmation: %w", err)
	}
	token := hex.EncodeToString(secret)
	hash := sha256.Sum256([]byte(token))
	op := &models.Operation{ID: uuid.NewString(), UserID: user, ProjectID: project, Kind: kind, Target: target, Payload: payload, Snapshot: snapshot, Effects: effects, Status: "REVIEWED", ExpiresAt: time.Now().Add(10 * time.Minute).Unix(), TokenHash: hex.EncodeToString(hash[:])}
	if err := s.store.Create(ctx, op); err != nil {
		return nil, err
	}
	return &models.OperationReview{Operation: op, Confirmation: token}, nil
}

func (s *Service) Apply(ctx context.Context, id, user, token, snapshot string, run func(context.Context, *models.Operation, func(string, string) error) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	op, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	hash := sha256.Sum256([]byte(token))
	if op.UserID != user || subtle.ConstantTimeCompare([]byte(op.TokenHash), []byte(hex.EncodeToString(hash[:]))) != 1 {
		return fmt.Errorf("invalid confirmation or operation owner")
	}
	if op.Status == "COMPLETED" || s.running[id] != nil {
		return nil
	}
	if op.Snapshot != snapshot {
		return fmt.Errorf("target changed since review; prepare again")
	}
	if err := s.store.Claim(ctx, id, snapshot); err != nil {
		return err
	}
	execution, cancel := context.WithTimeout(context.Background(), 24*time.Hour)
	s.running[id] = cancel
	go func() {
		defer cancel()
		defer func() { s.mu.Lock(); delete(s.running, id); s.mu.Unlock() }()
		progress := func(phase, log string) error {
			if err := execution.Err(); err != nil {
				return err
			}
			current, err := s.store.Get(execution, id)
			if err == nil && current != nil && current.Status == "CANCELLING" {
				cancel()
				return errors.New("operation cancelled")
			}
			return s.store.Observe(execution, id, "RUNNING", phase, "", log+"\n")
		}
		err := runSafely(execution, op, progress, run)
		status, message := "COMPLETED", ""
		if err != nil {
			status, message = "FAILED", err.Error()
		}
		if execution.Err() != nil {
			status = "INTERRUPTED"
			message = "Execution was interrupted. Inspect target state before preparing a retry. " + message
		}
		final, finish := context.WithTimeout(context.Background(), 10*time.Second)
		defer finish()
		if err := s.store.Observe(final, id, status, "FINISHED", message, ""); err != nil {
			slog.Error("persist operation result", "operation", id, "error", err)
		}
	}()
	return nil
}

func (s *Service) Cancel(ctx context.Context, id, user string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	op, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	if op.UserID != user {
		return fmt.Errorf("operation belongs to another user")
	}
	switch op.Status {
	case "COMPLETED", "FAILED", "INTERRUPTED":
		return fmt.Errorf("operation is not running")
	}
	if err := s.store.Observe(ctx, id, "CANCELLING", op.Phase, "", "Cancellation requested.\n"); err != nil {
		return err
	}
	if cancel := s.running[id]; cancel != nil {
		cancel()
	}
	return nil
}

func runSafely(ctx context.Context, op *models.Operation, progress func(string, string) error, run func(context.Context, *models.Operation, func(string, string) error) error) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("operation stopped unexpectedly; inspect target before retrying")
		}
	}()
	return run(ctx, op, progress)
}
