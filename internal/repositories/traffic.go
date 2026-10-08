package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"codedock.run/codedock/internal/models"
)

type TrafficRepository interface {
	RecordBatch(ctx context.Context, samples []models.TrafficSample) error
	Summary(ctx context.Context, projectID, domain, from, to string) (*models.AnalyticsSummary, error)
	Overview(ctx context.Context, projectID, domain, from, to, stepMinutes string) (*models.AnalyticsOverview, error)
	Geo(ctx context.Context, projectID, from, to string) (*models.AnalyticsGeo, error)
	PathsEnabled(ctx context.Context, projectID string) (bool, error)
	SetPathsEnabled(ctx context.Context, projectID string, enabled bool) error
	RetentionCutoff(days int) string
	DeleteBefore(ctx context.Context, cutoff string) error
}

type postgresTrafficRepository struct {
	db *sql.DB
}

func NewTrafficRepository(db *sql.DB) TrafficRepository {
	return &postgresTrafficRepository{db: db}
}

func bucketKey(sample models.TrafficSample) string {
	parsed, err := time.Parse(time.RFC3339, sample.Time)
	if err != nil {
		parsed = time.Now().UTC()
	}
	return parsed.UTC().Truncate(time.Minute).Format("2006-01-02T15:04")
}

func dayKey(sample models.TrafficSample) string {
	parsed, err := time.Parse(time.RFC3339, sample.Time)
	if err != nil {
		parsed = time.Now().UTC()
	}
	return parsed.UTC().Format("2006-01-02")
}

func (r *postgresTrafficRepository) RecordBatch(ctx context.Context, samples []models.TrafficSample) error {
	if len(samples) == 0 {
		return nil
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, sample := range samples {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO traffic_buckets (bucket_minute, project_id, domain, path, status, requests, bytes, duration_ms)
			VALUES ($1, $2, $3, $4, $5, 1, $6, $7)
			ON CONFLICT(bucket_minute, project_id, domain, path, status)
			DO UPDATE SET requests = traffic_buckets.requests + 1, bytes = traffic_buckets.bytes + excluded.bytes, duration_ms = traffic_buckets.duration_ms + excluded.duration_ms
		`, bucketKey(sample), sample.ProjectID, sample.Domain, sample.Path, sample.Status, sample.Bytes, sample.DurationMs); err != nil {
			return fmt.Errorf("record traffic bucket: %w", err)
		}
		if sample.ClientIP != "" {
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO traffic_visitors (day, project_id, ip, country, requests, bytes)
				VALUES ($1, $2, $3, '', 1, $4)
				ON CONFLICT(day, project_id, ip)
				DO UPDATE SET requests = traffic_visitors.requests + 1, bytes = traffic_visitors.bytes + excluded.bytes
			`, dayKey(sample), sample.ProjectID, sample.ClientIP, sample.Bytes); err != nil {
				return fmt.Errorf("record traffic visitor: %w", err)
			}
		}
	}
	return tx.Commit()
}

func (r *postgresTrafficRepository) Summary(ctx context.Context, projectID, domain, from, to string) (*models.AnalyticsSummary, error) {
	clause, args := trafficRange(projectID, domain, from, to)
	var summary models.AnalyticsSummary
	summary.ProjectID, summary.From, summary.To = projectID, from, to
	var requests, errors int
	var bytes, duration int64
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) as n, COALESCE(SUM(requests),0), COALESCE(SUM(bytes),0), COALESCE(SUM(duration_ms),0),
			COALESCE(SUM(CASE WHEN status >= 500 THEN requests ELSE 0 END),0)
		FROM traffic_buckets `+clause, args...).Scan(new(int), &requests, &bytes, &duration, &errors)
	if err != nil {
		return nil, fmt.Errorf("traffic summary: %w", err)
	}
	summary.Requests = requests
	summary.Bytes = bytes
	if requests > 0 {
		summary.ErrorRate = float64(errors) / float64(requests)
		summary.AvgDurationMs = float64(duration) / float64(requests)
	}
	return &summary, nil
}

func (r *postgresTrafficRepository) Overview(ctx context.Context, projectID, domain, from, to, stepMinutes string) (*models.AnalyticsOverview, error) {
	clause, args := trafficRange(projectID, domain, from, to)
	overview := &models.AnalyticsOverview{ProjectID: projectID, From: from, To: to}
	seriesRows, err := r.db.QueryContext(ctx, `
		SELECT bucket_minute, SUM(requests), SUM(bytes), SUM(CASE WHEN status >= 500 THEN requests ELSE 0 END)
		FROM traffic_buckets `+clause+` GROUP BY bucket_minute ORDER BY bucket_minute ASC LIMIT 1500
	`, args...)
	if err != nil {
		return nil, fmt.Errorf("traffic series: %w", err)
	}
	defer seriesRows.Close()
	raw := map[string]*models.AnalyticsSeriesPoint{}
	var order []string
	for seriesRows.Next() {
		var minute string
		var point models.AnalyticsSeriesPoint
		if err := seriesRows.Scan(&minute, &point.Requests, &point.Bytes, &point.Errors); err != nil {
			return nil, fmt.Errorf("scan traffic series: %w", err)
		}
		point.Time = minute
		raw[minute] = &models.AnalyticsSeriesPoint{Time: minute, Requests: point.Requests, Bytes: point.Bytes, Errors: point.Errors}
		order = append(order, minute)
	}
	overview.Series = downsampleSeries(raw, order, stepMinutes)
	statusRows, err := r.db.QueryContext(ctx, `
		SELECT status, SUM(requests) FROM traffic_buckets `+clause+` GROUP BY status ORDER BY 2 DESC LIMIT 20
	`, args...)
	if err != nil {
		return nil, fmt.Errorf("traffic statuses: %w", err)
	}
	defer statusRows.Close()
	for statusRows.Next() {
		var entry models.AnalyticsStatusBreakdown
		if err := statusRows.Scan(&entry.Status, &entry.Requests); err != nil {
			return nil, fmt.Errorf("scan traffic statuses: %w", err)
		}
		overview.Statuses = append(overview.Statuses, entry)
	}
	pathRows, err := r.db.QueryContext(ctx, `
		SELECT path, SUM(requests), SUM(bytes) FROM traffic_buckets `+clause+` AND path <> '' GROUP BY path ORDER BY 2 DESC LIMIT 25
	`, args...)
	if err != nil {
		return nil, fmt.Errorf("traffic paths: %w", err)
	}
	defer pathRows.Close()
	for pathRows.Next() {
		var entry models.AnalyticsTopPath
		if err := pathRows.Scan(&entry.Path, &entry.Requests, &entry.Bytes); err != nil {
			return nil, fmt.Errorf("scan traffic paths: %w", err)
		}
		overview.TopPaths = append(overview.TopPaths, entry)
	}
	return overview, nil
}

func (r *postgresTrafficRepository) Geo(ctx context.Context, projectID, from, to string) (*models.AnalyticsGeo, error) {
	geo := &models.AnalyticsGeo{ProjectID: projectID, From: from, To: to}
	rows, err := r.db.QueryContext(ctx, `
		SELECT COALESCE(NULLIF(country,''),'unknown'), SUM(requests), COUNT(DISTINCT ip), SUM(bytes)
		FROM traffic_visitors WHERE project_id = $1 AND day >= $2 AND day <= $3
		GROUP BY 1 ORDER BY 2 DESC LIMIT 50
	`, projectID, dayBound(from), dayBound(to))
	if err != nil {
		return nil, fmt.Errorf("traffic geo: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var entry models.AnalyticsGeoEntry
		if err := rows.Scan(&entry.Country, &entry.Requests, &entry.Visitors, &entry.Bytes); err != nil {
			return nil, fmt.Errorf("scan traffic geo: %w", err)
		}
		geo.Countries = append(geo.Countries, entry)
	}
	return geo, nil
}

func (r *postgresTrafficRepository) PathsEnabled(ctx context.Context, projectID string) (bool, error) {
	var enabled int
	err := r.db.QueryRowContext(ctx, `SELECT enabled FROM traffic_paths WHERE project_id = $1`, projectID).Scan(&enabled)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("paths toggle: %w", err)
	}
	return enabled != 0, nil
}

func (r *postgresTrafficRepository) SetPathsEnabled(ctx context.Context, projectID string, enabled bool) error {
	value := 0
	if enabled {
		value = 1
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO traffic_paths (project_id, enabled, updated_at) VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT(project_id) DO UPDATE SET enabled = excluded.enabled, updated_at = CURRENT_TIMESTAMP
	`, projectID, value)
	if err != nil {
		return fmt.Errorf("save paths toggle: %w", err)
	}
	return nil
}

func (r *postgresTrafficRepository) RetentionCutoff(days int) string {
	return time.Now().UTC().AddDate(0, 0, -days).Format("2006-01-02T15:04")
}

func (r *postgresTrafficRepository) DeleteBefore(ctx context.Context, cutoff string) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM traffic_buckets WHERE bucket_minute < $1`, cutoff); err != nil {
		return fmt.Errorf("retire traffic buckets: %w", err)
	}
	day := cutoff
	if len(cutoff) >= 10 {
		day = cutoff[:10]
	}
	if _, err := r.db.ExecContext(ctx, `DELETE FROM traffic_visitors WHERE day < $1`, day); err != nil {
		return fmt.Errorf("retire traffic visitors: %w", err)
	}
	return nil
}

func trafficRange(projectID, domain, from, to string) (string, []any) {
	clause := `WHERE project_id = $1`
	args := []any{projectID}
	if domain != "" {
		clause += fmt.Sprintf(` AND domain = $%d`, len(args)+1)
		args = append(args, domain)
	}
	if from != "" {
		clause += fmt.Sprintf(` AND bucket_minute >= $%d`, len(args)+1)
		args = append(args, from)
	}
	if to != "" {
		clause += fmt.Sprintf(` AND bucket_minute <= $%d`, len(args)+1)
		args = append(args, to)
	}
	return clause, args
}

func dayBound(value string) string {
	if len(value) >= 10 {
		return value[:10]
	}
	if value == "" {
		return "0000-00-00"
	}
	return value
}

func downsampleSeries(raw map[string]*models.AnalyticsSeriesPoint, order []string, step string) []models.AnalyticsSeriesPoint {
	if step == "" || step == "1" {
		series := make([]models.AnalyticsSeriesPoint, 0, len(order))
		for _, minute := range order {
			series = append(series, *raw[minute])
		}
		return series
	}
	var width int
	_, _ = fmt.Sscanf(step, "%d", &width)
	if width <= 1 {
		width = 5
	}
	var series []models.AnalyticsSeriesPoint
	var current *models.AnalyticsSeriesPoint
	count := 0
	for _, minute := range order {
		if current == nil {
			current = &models.AnalyticsSeriesPoint{Time: minute}
		}
		point := raw[minute]
		current.Requests += point.Requests
		current.Bytes += point.Bytes
		current.Errors += point.Errors
		count++
		if count >= width {
			series = append(series, *current)
			current = nil
			count = 0
		}
	}
	if current != nil {
		series = append(series, *current)
	}
	return series
}
