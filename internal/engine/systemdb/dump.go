package systemdb

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func DumpContainer(ctx context.Context, container, user, dbname, password string) ([]byte, error) {
	out, err := exec.CommandContext(ctx, "docker", "exec", "-e", "PGPASSWORD="+password, container,
		"pg_dump", "-U", user, "-d", dbname).Output()
	if err != nil {
		return nil, fmt.Errorf("pg_dump of container %s failed: %w", container, err)
	}
	return out, nil
}

func DumpURL(ctx context.Context, databaseURL string) ([]byte, error) {
	out, err := exec.CommandContext(ctx, "pg_dump", databaseURL).Output()
	if err != nil {
		return nil, fmt.Errorf("pg_dump failed (postgresql-client is required for external databases): %w", err)
	}
	return out, nil
}

func RestoreContainer(ctx context.Context, container, user, dbname, password string, sqlData []byte) error {
	cmd := exec.CommandContext(ctx, "docker", "exec", "-i", "-e", "PGPASSWORD="+password, container,
		"psql", "-U", user, "-d", dbname)
	cmd.Stdin = bytes.NewReader(sqlData)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("psql restore into container %s failed: %w", container, err)
	}
	return nil
}

func RestoreURL(ctx context.Context, databaseURL string, sqlData []byte) error {
	cmd := exec.CommandContext(ctx, "psql", databaseURL)
	cmd.Stdin = bytes.NewReader(sqlData)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("psql restore failed: %w", err)
	}
	return nil
}
