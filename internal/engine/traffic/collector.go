package traffic

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"codedock/internal/models"
)

type SampleStore interface {
	RecordBatch(ctx context.Context, samples []models.TrafficSample) error
	PathsEnabled(ctx context.Context, projectID string) (bool, error)
	DeleteBefore(ctx context.Context, cutoff string) error
	RetentionCutoff(days int) string
}

type AttributeFunc func(host string) (projectID string)

type Collector struct {
	path       string
	store      SampleStore
	attribute  AttributeFunc

	mu         sync.Mutex
	offset     int64
	paths      map[string]bool
	pathsAt    time.Time
	interval   time.Duration
	retainDays int
}

func NewCollector(path string, store SampleStore, attribute AttributeFunc) *Collector {
	return &Collector{
		path: path, store: store, attribute: attribute,
		paths: map[string]bool{}, interval: 30 * time.Second, retainDays: 30,
	}
}

func (c *Collector) SetInterval(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.interval = interval
}

func (c *Collector) Run(ctx context.Context) {
	c.skipToEnd()
	ticker := time.NewTicker(c.currentInterval())
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.Collect(ctx); err != nil {
				slog.Error("collect traffic", "error", err)
			}
			if err := c.store.DeleteBefore(ctx, c.store.RetentionCutoff(c.retainDays)); err != nil {
				slog.Error("retire traffic", "error", err)
			}
		}
	}
}

func (c *Collector) currentInterval() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.interval <= 0 {
		return 30 * time.Second
	}
	return c.interval
}

func (c *Collector) skipToEnd() {
	info, err := os.Stat(c.path)
	if err != nil {
		return
	}
	c.mu.Lock()
	c.offset = info.Size()
	c.mu.Unlock()
}

func (c *Collector) Collect(ctx context.Context) error {
	c.mu.Lock()
	offset := c.offset
	c.mu.Unlock()
	file, err := os.Open(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.Size() < offset {
		offset = 0
	}
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	var samples []models.TrafficSample
	position := offset
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			position += int64(len(line))
			if sample, ok := c.convert(ctx, line); ok {
				samples = append(samples, sample)
			}
		}
		if err != nil {
			break
		}
	}
	if len(samples) > 0 {
		if err := c.store.RecordBatch(ctx, samples); err != nil {
			return fmt.Errorf("record traffic samples: %w", err)
		}
	}
	c.mu.Lock()
	c.offset = position
	c.mu.Unlock()
	return nil
}

func (c *Collector) convert(ctx context.Context, line string) (models.TrafficSample, bool) {
	entry, err := ParseAccessLine(line)
	if err != nil {
		return models.TrafficSample{}, false
	}
	projectID := ""
	if c.attribute != nil {
		projectID = c.attribute(entry.RequestHost)
	}
	if projectID == "" {
		return models.TrafficSample{}, false
	}
	path := ""
	if c.pathsEnabled(ctx, projectID) {
		path = entry.Path()
	}
	return models.TrafficSample{
		Time:       entry.Timestamp().Format(time.RFC3339),
		ProjectID:  projectID,
		Domain:     entry.RequestHost,
		Path:       path,
		Status:     entry.DownstreamStatus,
		Bytes:      entry.Bytes(),
		DurationMs: entry.DurationMs(),
		ClientIP:   entry.ClientIP(),
	}, true
}

func (c *Collector) pathsEnabled(ctx context.Context, projectID string) bool {
	c.mu.Lock()
	cached, ok := c.paths[projectID]
	fresh := time.Since(c.pathsAt) < time.Minute
	c.mu.Unlock()
	if ok && fresh {
		return cached
	}
	enabled, err := c.store.PathsEnabled(ctx, projectID)
	if err != nil {
		return false
	}
	c.mu.Lock()
	c.paths[projectID] = enabled
	c.pathsAt = time.Now()
	c.mu.Unlock()
	return enabled
}
