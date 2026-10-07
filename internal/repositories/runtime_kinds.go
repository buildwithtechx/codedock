package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

func (r *RuntimeRepo) SyncKinds(ctx context.Context) error {
	var rows []models.ServiceRuntime
	if err := r.db.SelectContext(ctx, &rows, `SELECT * FROM service_runtimes`); err != nil {
		return err
	}
	for _, row := range rows {
		decoded, err := r.vault.Decrypt(row.Config)
		if err != nil {
			return err
		}
		var target models.RuntimeTarget
		if err := json.Unmarshal([]byte(decoded), &target); err != nil {
			return err
		}
		if target.Kind == "" {
			target.Kind = "docker"
		}
		if target.Kind != "docker" && target.Kind != "kubernetes" && target.Kind != "bare" {
			return fmt.Errorf("invalid saved runtime kind")
		}
		if row.RuntimeKind == target.Kind {
			continue
		}
		if _, err := r.db.ExecContext(ctx, `UPDATE service_runtimes SET runtime_kind=? WHERE service_id=? AND encrypted_config=?`, target.Kind, row.ServiceID, row.Config); err != nil {
			return err
		}
	}
	return nil
}
