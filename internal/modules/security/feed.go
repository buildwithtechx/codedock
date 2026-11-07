package security

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	DefaultFeedURL   = "https://api.github.com/repos/oblien/openship/contents/security-advisories.v2.signed.json?ref=security-advisories"
	FeedRowID        = "security:v2"
	PollIntervalMin  = 5 * time.Minute
	PollIntervalNorm = 1 * time.Hour
)

// Embedded baseline advisories for air-gapped or offline operation.
var embeddedBaselineManifest = SecurityManifest{
	Format:    1,
	Revision:  1,
	IssuedAt:  time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC),
	ExpiresAt: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
	Advisories: []SecurityAdvisory{
		{
			ID:          "CVE-2026-111497-node",
			Revision:    1,
			Title:       "Update Node.js for the Brotli security fix",
			Message:     "Affected Node.js releases contain a Brotli decompression vulnerability. Use Node.js 22.23.4, 24.22.1 or 26.12.1, or a later patched release.",
			Severity:    SeverityHigh,
			URL:         "https://nodejs.org/en/blog/vulnerability/october-2026-security-releases",
			PublishedAt: time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC),
			Action:      ActionRefreshImages,
			Targets: []Target{
				{
					Kind:     TargetKindRuntime,
					Runtime:  "node",
					Affected: "<22.23.4 || >=23.0.0 <24.22.1 || >=25.0.0 <26.12.1",
					Fixed:    stringPtr(">=22.23.4 <23 || >=24.22.1 <25 || >=26.12.1"),
				},
			},
		},
		{
			ID:          "CVE-2026-111497-bun",
			Revision:    1,
			Title:       "Update Bun for the Brotli security fix",
			Message:     "Bun versions before 1.4.4 include the affected Brotli decoder. Upgrade to Bun 1.4.4 or later.",
			Severity:    SeverityHigh,
			URL:         "https://github.com/oven-sh/bun/issues/46352",
			PublishedAt: time.Date(2026, 10, 22, 0, 0, 0, 0, time.UTC),
			Action:      ActionRefreshImages,
			Targets: []Target{
				{
					Kind:     TargetKindRuntime,
					Runtime:  "bun",
					Affected: "<1.4.4",
					Fixed:    stringPtr(">=1.4.4"),
				},
			},
		},
		{
			ID:          "PHP-2026-10-21-rce",
			Revision:    1,
			Title:       "Update PHP for the October 21 RCE fixes",
			Message:     "PHP 8.2.35, 8.3.36, 8.4.27 and 8.5.12 fix CVE-2026-109618 and CVE-2026-109619. Update your runtime image.",
			Severity:    SeverityCritical,
			URL:         "https://www.php.net/archive/2026.php",
			PublishedAt: time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC),
			Action:      ActionRefreshImages,
			Targets: []Target{
				{
					Kind:     TargetKindRuntime,
					Runtime:  "php",
					Affected: "<8.2.35 || >=8.3.0 <8.3.36 || >=8.4.0 <8.4.27 || >=8.5.0 <8.5.12",
					Fixed:    stringPtr(">=8.2.35 <8.3.0 || >=8.3.36 <8.4.0 || >=8.4.27 <8.5.0 || >=8.5.12 <8.6.0"),
				},
			},
		},
	},
}

func stringPtr(s string) *string {
	return &s
}

type FeedManager struct {
	db       *sql.DB
	verifier *Verifier
	client   *http.Client
	feedURL  string
	mu       sync.RWMutex
	cached   *SecurityManifest
}

func NewFeedManager(db *sql.DB, verifier *Verifier, feedURL string) *FeedManager {
	if feedURL == "" {
		feedURL = DefaultFeedURL
	}
	return &FeedManager{
		db:       db,
		verifier: verifier,
		client:   &http.Client{Timeout: 15 * time.Second},
		feedURL:  feedURL,
	}
}

func (fm *FeedManager) GetManifest(ctx context.Context) (*SecurityManifest, error) {
	fm.mu.RLock()
	if fm.cached != nil {
		manifest := fm.cached
		fm.mu.RUnlock()
		return manifest, nil
	}
	fm.mu.RUnlock()

	// Try loading from database
	if fm.db != nil {
		var bodyRaw []byte
		var revision int64
		err := fm.db.QueryRowContext(ctx, "SELECT body, revision FROM advisory_feed WHERE id = $1", FeedRowID).Scan(&bodyRaw, &revision)
		if err == nil && len(bodyRaw) > 0 {
			var proof SignedProof
			if err := json.Unmarshal(bodyRaw, &proof); err == nil && proof.Payload != "" {
				manifest, err := fm.verifier.VerifyProof(proof, revision, true, time.Now())
				if err == nil {
					fm.mu.Lock()
					fm.cached = manifest
					fm.mu.Unlock()
					return manifest, nil
				}
			}
		}
	}

	// Fallback to embedded baseline
	fm.mu.Lock()
	fm.cached = &embeddedBaselineManifest
	fm.mu.Unlock()
	return &embeddedBaselineManifest, nil
}

func (fm *FeedManager) FetchRemote(ctx context.Context, githubToken string) (*SecurityManifest, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fm.feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create feed request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.raw+json")
	if githubToken != "" {
		req.Header.Set("Authorization", "Bearer "+githubToken)
	}

	var storedETag string
	var storedRevision int64
	if fm.db != nil {
		_ = fm.db.QueryRowContext(ctx, "SELECT etag, revision FROM advisory_feed WHERE id = $1", FeedRowID).Scan(&storedETag, &storedRevision)
		if storedETag != "" {
			req.Header.Set("If-None-Match", storedETag)
		}
	}

	resp, err := fm.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch advisory feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		return fm.GetManifest(ctx)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("feed request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxPayloadBytes+65536))
	if err != nil {
		return nil, fmt.Errorf("read feed response: %w", err)
	}

	var proof SignedProof
	if err := json.Unmarshal(body, &proof); err != nil {
		return nil, fmt.Errorf("parse feed envelope: %w", err)
	}

	manifest, err := fm.verifier.VerifyProof(proof, storedRevision, false, time.Now())
	if err != nil {
		return nil, fmt.Errorf("verify feed proof: %w", err)
	}

	newETag := strings.Trim(resp.Header.Get("ETag"), `"`)
	now := time.Now().UTC()
	nextCheck := now.Add(PollIntervalNorm)

	if fm.db != nil {
		query := `
			INSERT INTO advisory_feed (id, body, etag, revision, checked_at, next_check_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (id) DO UPDATE SET
				body = EXCLUDED.body,
				etag = EXCLUDED.etag,
				revision = EXCLUDED.revision,
				checked_at = EXCLUDED.checked_at,
				next_check_at = EXCLUDED.next_check_at,
				updated_at = EXCLUDED.updated_at
		`
		_, _ = fm.db.ExecContext(ctx, query, FeedRowID, body, newETag, manifest.Revision, now, nextCheck, now)
	}

	fm.mu.Lock()
	fm.cached = manifest
	fm.mu.Unlock()

	return manifest, nil
}
