package git

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type gitAppRepo struct {
	db    *sqlx.DB
	vault Encryptor
}

func NewGitAppRepository(db *sql.DB, vault Encryptor) GitAppRepository {
	return &gitAppRepo{db: sqlx.NewDb(db, "pgx"), vault: vault}
}

func (r *gitAppRepo) ListGithubApps(ctx context.Context) ([]GithubApp, error) {
	query := `SELECT id, name, app_id, installation_id, client_id, client_secret, webhook_secret, private_key, is_public, created_at, updated_at FROM github_apps`
	var apps []GithubApp
	if err := r.db.SelectContext(ctx, &apps, query); err != nil {
		return nil, err
	}
	if apps == nil {
		apps = make([]GithubApp, 0)
	}
	for i := range apps {
		if apps[i].ClientSecret != "" {
			apps[i].ClientSecret = "********"
		}
		if apps[i].WebhookSecret != "" {
			apps[i].WebhookSecret = "********"
		}
		if apps[i].PrivateKey != "" {
			apps[i].PrivateKey = "********"
		}
	}
	return apps, nil
}

func (r *gitAppRepo) GetGithubApp(ctx context.Context, id string) (*GithubApp, error) {
	query := `SELECT id, name, app_id, installation_id, client_id, client_secret, webhook_secret, private_key, is_public, created_at, updated_at FROM github_apps WHERE id = $1`
	var a GithubApp
	if err := r.db.GetContext(ctx, &a, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAppNotFound
		}
		return nil, err
	}
	if a.ClientSecret != "" {
		a.ClientSecret = "********"
	}
	if a.WebhookSecret != "" {
		a.WebhookSecret = "********"
	}
	if a.PrivateKey != "" {
		a.PrivateKey = "********"
	}
	return &a, nil
}

func (r *gitAppRepo) SaveGithubApp(ctx context.Context, app *GithubApp) error {
	if app.ID == "" {
		app.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	if app.CreatedAt.IsZero() {
		app.CreatedAt = now
	}
	app.UpdatedAt = now

	columns := []string{"id", "name", "app_id", "installation_id", "client_id", "client_secret", "webhook_secret", "private_key", "is_public", "created_at", "updated_at"}
	values := []any{app.ID, app.Name, app.AppID, app.InstallationID, app.ClientID, app.ClientSecret, app.WebhookSecret, app.PrivateKey, app.IsPublic, app.CreatedAt, app.UpdatedAt}

	placeholders := make([]string, len(columns))
	updates := make([]string, 0, len(columns))
	for i, col := range columns {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		if col != "id" && col != "created_at" && col != "updated_at" {
			updates = append(updates, fmt.Sprintf("%s=excluded.%s", col, col))
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO github_apps (%s)
		VALUES (%s)
		ON CONFLICT(id) DO UPDATE SET
			%s,
			updated_at=CURRENT_TIMESTAMP
	`, strings.Join(columns, ", "), strings.Join(placeholders, ", "), strings.Join(updates, ",\n\t\t\t"))

	_, err := r.db.ExecContext(ctx, query, values...)
	return err
}

func (r *gitAppRepo) DeleteGithubApp(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM github_apps WHERE id = $1", id)
	return err
}
