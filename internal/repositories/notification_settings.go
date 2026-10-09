package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock/internal/models"
)

type NotificationSettingsRepository interface {
	GetNotificationSettings(ctx context.Context) (*models.NotificationSettings, error)
	UpdateNotificationSettings(ctx context.Context, cfg *models.NotificationSettings) error
	ListSubscriptions(ctx context.Context, userID, orgID string) ([]models.NotificationSubscription, error)
	UpsertSubscription(ctx context.Context, sub *models.NotificationSubscription) error
	ListDefaults(ctx context.Context, orgID string) ([]models.NotificationDefault, error)
	UpsertDefault(ctx context.Context, def *models.NotificationDefault) error
}

type NotificationSettingsRepo struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewNotificationSettingsRepo(db *sql.DB) *NotificationSettingsRepo {
	return &NotificationSettingsRepo{db: sqlx.NewDb(db, "pgx")}
}

const notificationSettingsColumns = `id, discord_webhook_url, discord_ping_enabled, discord_enabled, slack_webhook_url, slack_enabled, telegram_bot_token, telegram_chat_id, telegram_enabled, smtp_host, smtp_port, smtp_user, smtp_password, smtp_from_name, smtp_from_address, smtp_enabled, resend_api_key, resend_enabled, pushover_user_key, pushover_api_token, pushover_enabled, generic_webhook_url, generic_webhook_enabled, notification_alerts, created_at, updated_at`

func notificationSettingsPlaceholders() string {
	columns := strings.Split(notificationSettingsColumns, ",")
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "$" + strconv.Itoa(i+1)
	}
	return strings.Join(placeholders, ", ")
}

func scanNotificationSettings(scanner interface{ Scan(dest ...any) error }, cfg *models.NotificationSettings) error {
	return scanner.Scan(
		&cfg.ID, &cfg.DiscordWebhookURL, &cfg.DiscordPingEnabled, &cfg.DiscordEnabled, &cfg.SlackWebhookURL, &cfg.SlackEnabled, &cfg.TelegramBotToken, &cfg.TelegramChatID, &cfg.TelegramEnabled,
		&cfg.SMTPHost, &cfg.SMTPPort, &cfg.SMTPUser, &cfg.SMTPPassword, &cfg.SMTPFromName, &cfg.SMTPFromAddress, &cfg.SMTPEnabled,
		&cfg.ResendAPIKey, &cfg.ResendEnabled, &cfg.PushoverUserKey, &cfg.PushoverAPIToken, &cfg.PushoverEnabled, &cfg.GenericWebhookURL, &cfg.GenericWebhookEnabled,
		&cfg.NotificationAlerts,
		&cfg.CreatedAt, &cfg.UpdatedAt,
	)
}

func notificationSettingsArgs(cfg *models.NotificationSettings) []any {
	return []any{
		cfg.ID, cfg.DiscordWebhookURL, cfg.DiscordPingEnabled, cfg.DiscordEnabled, cfg.SlackWebhookURL, cfg.SlackEnabled, cfg.TelegramBotToken, cfg.TelegramChatID, cfg.TelegramEnabled,
		cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFromName, cfg.SMTPFromAddress, cfg.SMTPEnabled,
		cfg.ResendAPIKey, cfg.ResendEnabled, cfg.PushoverUserKey, cfg.PushoverAPIToken, cfg.PushoverEnabled, cfg.GenericWebhookURL, cfg.GenericWebhookEnabled,
		cfg.NotificationAlerts,
		cfg.CreatedAt, cfg.UpdatedAt,
	}
}

func (r *NotificationSettingsRepo) GetNotificationSettings(ctx context.Context) (*models.NotificationSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg models.NotificationSettings
	err := scanNotificationSettings(r.db.QueryRowContext(ctx, `SELECT `+notificationSettingsColumns+` FROM notification_settings WHERE id = 'global' LIMIT 1`), &cfg)
	if errors.Is(err, sql.ErrNoRows) {
		defaultSettings := &models.NotificationSettings{
			ID:                 "global",
			NotificationAlerts: true,
			SMTPPort:           587,
			CreatedAt:          time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:          time.Now().UTC().Format(time.RFC3339),
		}
		query := fmt.Sprintf(`INSERT INTO notification_settings (%s) VALUES (%s)`, notificationSettingsColumns, notificationSettingsPlaceholders())
		_, _ = r.db.ExecContext(ctx, query, notificationSettingsArgs(defaultSettings)...)
		return defaultSettings, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notification settings: %w", err)
	}
	return &cfg, nil
}

func (r *NotificationSettingsRepo) UpdateNotificationSettings(ctx context.Context, cfg *models.NotificationSettings) error {
	if cfg.ID == "" {
		cfg.ID = "global"
	}
	cfg.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.mu.Lock()
	defer r.mu.Unlock()
	query := fmt.Sprintf(`INSERT INTO notification_settings (%s)
	          VALUES (%s)
	          ON CONFLICT(id) DO UPDATE SET
	          discord_webhook_url = excluded.discord_webhook_url,
	          discord_ping_enabled = excluded.discord_ping_enabled,
	          discord_enabled = excluded.discord_enabled,
	          slack_webhook_url = excluded.slack_webhook_url,
	          slack_enabled = excluded.slack_enabled,
	          telegram_bot_token = excluded.telegram_bot_token,
	          telegram_chat_id = excluded.telegram_chat_id,
	          telegram_enabled = excluded.telegram_enabled,
	          smtp_host = excluded.smtp_host,
	          smtp_port = excluded.smtp_port,
	          smtp_user = excluded.smtp_user,
	          smtp_password = excluded.smtp_password,
	          smtp_from_name = excluded.smtp_from_name,
	          smtp_from_address = excluded.smtp_from_address,
	          smtp_enabled = excluded.smtp_enabled,
	          resend_api_key = excluded.resend_api_key,
	          resend_enabled = excluded.resend_enabled,
	          pushover_user_key = excluded.pushover_user_key,
	          pushover_api_token = excluded.pushover_api_token,
	          pushover_enabled = excluded.pushover_enabled,
	          generic_webhook_url = excluded.generic_webhook_url,
	          generic_webhook_enabled = excluded.generic_webhook_enabled,
	          notification_alerts = excluded.notification_alerts,
	          updated_at = excluded.updated_at`, notificationSettingsColumns, notificationSettingsPlaceholders())
	_, err := r.db.ExecContext(ctx, query, notificationSettingsArgs(cfg)...)
	if err != nil {
		return fmt.Errorf("failed to update notification settings: %w", err)
	}
	return nil
}

func (r *NotificationSettingsRepo) ListSubscriptions(ctx context.Context, userID, orgID string) ([]models.NotificationSubscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	query := `SELECT id, user_id, organization_id, category, channels, enabled, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at, to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as updated_at
	          FROM notification_subscriptions
	          WHERE user_id = $1`
	args := []any{userID}
	if orgID != "" {
		query += ` AND organization_id = $2`
		args = append(args, orgID)
	}
	query += ` ORDER BY category ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []models.NotificationSubscription
	for rows.Next() {
		var sub models.NotificationSubscription
		var channelsRaw []byte
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.OrganizationID, &sub.Category, &channelsRaw, &sub.Enabled, &sub.CreatedAt, &sub.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		if len(channelsRaw) > 0 {
			_ = json.Unmarshal(channelsRaw, &sub.Channels)
		}
		if sub.Channels == nil {
			sub.Channels = []string{}
		}
		subs = append(subs, sub)
	}
	if subs == nil {
		subs = []models.NotificationSubscription{}
	}
	return subs, nil
}

func (r *NotificationSettingsRepo) UpsertSubscription(ctx context.Context, sub *models.NotificationSubscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if sub.ID == "" {
		sub.ID = uuid.NewString()
	}
	channelsJSON, err := json.Marshal(sub.Channels)
	if err != nil {
		channelsJSON = []byte("[]")
	}

	query := `INSERT INTO notification_subscriptions (id, user_id, organization_id, category, channels, enabled, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	          ON CONFLICT(user_id, organization_id, category) DO UPDATE SET
	          channels = excluded.channels,
	          enabled = excluded.enabled,
	          updated_at = CURRENT_TIMESTAMP`

	_, err = r.db.ExecContext(ctx, query, sub.ID, sub.UserID, sub.OrganizationID, sub.Category, channelsJSON, sub.Enabled)
	if err != nil {
		return fmt.Errorf("failed to upsert notification subscription: %w", err)
	}
	return nil
}

func (r *NotificationSettingsRepo) ListDefaults(ctx context.Context, orgID string) ([]models.NotificationDefault, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	query := `SELECT organization_id, category, channels, enabled, to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as updated_at
	          FROM notification_defaults
	          WHERE organization_id = $1
	          ORDER BY category ASC`

	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification defaults: %w", err)
	}
	defer rows.Close()

	var defs []models.NotificationDefault
	for rows.Next() {
		var def models.NotificationDefault
		var channelsRaw []byte
		if err := rows.Scan(&def.OrganizationID, &def.Category, &channelsRaw, &def.Enabled, &def.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan notification default: %w", err)
		}
		if len(channelsRaw) > 0 {
			_ = json.Unmarshal(channelsRaw, &def.Channels)
		}
		if def.Channels == nil {
			def.Channels = []string{}
		}
		defs = append(defs, def)
	}
	if defs == nil {
		defs = []models.NotificationDefault{}
	}
	return defs, nil
}

func (r *NotificationSettingsRepo) UpsertDefault(ctx context.Context, def *models.NotificationDefault) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	channelsJSON, err := json.Marshal(def.Channels)
	if err != nil {
		channelsJSON = []byte("[]")
	}

	query := `INSERT INTO notification_defaults (organization_id, category, channels, enabled, updated_at)
	          VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
	          ON CONFLICT(organization_id, category) DO UPDATE SET
	          channels = excluded.channels,
	          enabled = excluded.enabled,
	          updated_at = CURRENT_TIMESTAMP`

	_, err = r.db.ExecContext(ctx, query, def.OrganizationID, def.Category, channelsJSON, def.Enabled)
	if err != nil {
		return fmt.Errorf("failed to upsert notification default: %w", err)
	}
	return nil
}
