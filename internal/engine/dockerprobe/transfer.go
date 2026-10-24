package dockerprobe

import (
	"context"
	"fmt"
	"io"
	"strings"

	"golang.org/x/sync/errgroup"
)

func TransferImage(ctx context.Context, source, target Runner, image string, progress func(string)) error {
	if image == "" || strings.ContainsAny(image, " \t\n;|&$`'\"\\") {
		return fmt.Errorf("invalid image reference")
	}
	if progress != nil {
		progress("Saving image " + image)
	}
	reader, writer := io.Pipe()
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer writer.Close()
		return source.RunPipe(ctx, "docker save "+image+" | gzip -c", nil, writer)
	})
	group.Go(func() error {
		defer reader.Close()
		return target.RunPipe(ctx, "gunzip -c | docker load", reader, io.Discard)
	})
	if err := group.Wait(); err != nil {
		return fmt.Errorf("transfer image %s: %w", image, err)
	}
	return nil
}

func EnsureVolume(ctx context.Context, runner Runner, name string) error {
	if name == "" || strings.ContainsAny(name, " \t\n;|&$`'\"\\/") {
		return fmt.Errorf("invalid volume name")
	}
	existing, err := runner.Run(ctx, "docker volume ls -q --filter name="+name)
	if err != nil {
		return fmt.Errorf("inspect target volumes: %w", err)
	}
	for _, line := range strings.Split(existing, "\n") {
		if strings.TrimSpace(line) == name {
			return nil
		}
	}
	if _, err := runner.Run(ctx, "docker volume create "+name); err != nil {
		return fmt.Errorf("create target volume %s: %w", name, err)
	}
	return nil
}

func VolumeHasData(ctx context.Context, runner Runner, name string) (bool, error) {
	out, err := runner.Run(ctx, "docker run --rm -v "+name+":/codedock-data:ro alpine ls -A /codedock-data")
	if err != nil {
		return false, nil
	}
	return strings.TrimSpace(out) != "", nil
}

func TransferVolume(ctx context.Context, source, target Runner, sourceName, targetName string, progress func(string)) error {
	if err := EnsureVolume(ctx, target, targetName); err != nil {
		return err
	}
	if progress != nil {
		progress("Copying volume " + sourceName + " to " + targetName)
	}
	reader, writer := io.Pipe()
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer writer.Close()
		return source.RunPipe(ctx, "docker run --rm -v "+sourceName+":/codedock-data:ro alpine tar czf - -C /codedock-data .", nil, writer)
	})
	group.Go(func() error {
		defer reader.Close()
		return target.RunPipe(ctx, "docker run --rm -i -v "+targetName+":/codedock-data alpine tar xzf - -C /codedock-data", reader, io.Discard)
	})
	if err := group.Wait(); err != nil {
		return fmt.Errorf("transfer volume %s: %w", sourceName, err)
	}
	return nil
}

func RemoveVolume(ctx context.Context, runner Runner, name string) error {
	if _, err := runner.Run(ctx, "docker volume rm -f "+name); err != nil {
		return fmt.Errorf("remove volume %s: %w", name, err)
	}
	return nil
}

func ContainerState(ctx context.Context, runner Runner, id string) (string, error) {
	out, err := runner.Run(ctx, "docker inspect -f '{{.State.Status}}' "+id)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
