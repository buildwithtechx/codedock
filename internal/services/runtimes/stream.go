package runtimes

import (
	"codedock/internal/engine/kubernetes"
	"context"
	"fmt"
	"io"
)

func (s *Service) StreamLogs(ctx context.Context, id, pod string, output io.Writer) error {
	if s.streamer == nil {
		return fmt.Errorf("live log transport unavailable")
	}
	app, cluster, err := s.target(ctx, id)
	if err != nil {
		return err
	}
	return s.engine.StreamLogs(ctx, s.streamer, cluster.Nodes[0], app, pod, output)
}

func (s *Service) StreamExec(ctx context.Context, id, pod string, command []string, input io.Reader, output io.Writer) error {
	if s.streamer == nil {
		return fmt.Errorf("live terminal transport unavailable")
	}
	app, cluster, err := s.target(ctx, id)
	if err != nil {
		return err
	}
	return s.engine.StreamExec(ctx, s.streamer, cluster.Nodes[0], app, pod, command, input, output)
}

func (s *Service) InstanceGraph(ctx context.Context, id string) (*kubernetes.InstanceGraph, error) {
	app, cluster, err := s.target(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.engine.InstanceGraph(ctx, cluster.Nodes[0], app)
}

func (s *Service) ResourceSummary(ctx context.Context, id string) ([]kubernetes.ResourceControl, error) {
	app, cluster, err := s.target(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.engine.ResourceSummary(ctx, cluster.Nodes[0], app)
}
