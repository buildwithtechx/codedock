package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"codedock/internal/models"
	"codedock/internal/utils"
)

func (r *UserRepo) GetAppearance(ctx context.Context, userID string) (*models.AppearancePreferences, error) {
	var raw sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT appearance FROM users WHERE id = $1`, userID).Scan(&raw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if !raw.Valid || strings.TrimSpace(raw.String) == "" {
		return &models.AppearancePreferences{
			ActiveTheme: "default",
			Mode:        "dark",
		}, nil
	}
	var pref models.AppearancePreferences
	if err := json.Unmarshal([]byte(raw.String), &pref); err != nil {
		return nil, fmt.Errorf("decode appearance preferences: %w", err)
	}
	return &pref, nil
}

func (r *UserRepo) SetAppearance(ctx context.Context, userID string, pref *models.AppearancePreferences) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var val any
	if pref != nil {
		data, err := json.Marshal(pref)
		if err != nil {
			return fmt.Errorf("encode appearance preferences: %w", err)
		}
		val = string(data)
	}

	res, err := r.db.ExecContext(ctx, `UPDATE users SET appearance = $1, updated_at = NOW() WHERE id = $2`, val, userID)
	if err != nil {
		return fmt.Errorf("update appearance: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return utils.NewNotFoundError("User", userID)
	}
	return nil
}
