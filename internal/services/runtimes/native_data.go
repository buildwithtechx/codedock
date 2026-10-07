package runtimes

import (
	"context"
	"fmt"
	"io"
)

func (s *Service) BackupNativeData(ctx context.Context, id string, output io.Writer) error {
	if s.native == nil {
		return fmt.Errorf("native runtime unavailable")
	}
	app, runtime, err := s.nativeTarget(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return fmt.Errorf("service does not use the native runtime")
	}
	unlock, err := s.gate.AcquireVolume("server:" + runtime.Target.BareNode.ServerID)
	if err != nil {
		return err
	}
	defer unlock()
	return s.native.BackupData(ctx, app, runtime.Target, output)
}

func (s *Service) RestoreNativeData(ctx context.Context, id string, archive io.Reader) error {
	if s.native == nil {
		return fmt.Errorf("native runtime unavailable")
	}
	app, runtime, err := s.nativeTarget(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return fmt.Errorf("service does not use the native runtime")
	}
	if runtime.Journal != "" {
		return fmt.Errorf("recover the interrupted native release before restoring data")
	}
	unlock, err := s.gate.AcquireVolume("server:" + runtime.Target.BareNode.ServerID)
	if err != nil {
		return err
	}
	defer unlock()
	return s.native.RestoreData(ctx, app, runtime.Target, archive)
}
