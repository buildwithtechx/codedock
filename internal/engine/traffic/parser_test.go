package traffic

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"codedock.run/codedock/internal/models"
)

func TestParseAccessLine(t *testing.T) {
	line := `{"ClientAddr":"203.0.113.7:41234","DownstreamStatus":200,"Duration":1500000,"RequestHost":"shop.example.com","RequestMethod":"GET","RequestPath":"/api/cart?x=1","DownstreamContentSize":512,"OriginContentSize":64,"StartUTC":"2026-10-07T08:01:02Z"}`
	entry, err := ParseAccessLine(line)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if entry.ClientIP() != "203.0.113.7" {
		t.Fatalf("expected client ip, got %q", entry.ClientIP())
	}
	if entry.Path() != "/api/cart" {
		t.Fatalf("expected stripped path, got %q", entry.Path())
	}
	if entry.Bytes() != 576 || entry.DurationMs() != 1 {
		t.Fatalf("expected bytes/duration, got %+v", entry)
	}
	if _, err := ParseAccessLine(`{"RequestHost":""}`); err == nil {
		t.Fatal("expected hostless line to fail")
	}
	if _, err := ParseAccessLine(`not json`); err == nil {
		t.Fatal("expected garbage to fail")
	}
}

func TestClassifyCountry(t *testing.T) {
	for _, check := range []struct {
		ip    string
		want  string
	}{
		{"127.0.0.1", "local"},
		{"10.1.2.3", "private"},
		{"192.168.0.9", "private"},
		{"203.0.113.7", "unknown"},
		{"not-an-ip", "unknown"},
	} {
		if got := ClassifyCountry(check.ip); got != check.want {
			t.Fatalf("classify %s: expected %s, got %s", check.ip, check.want, got)
		}
	}
}

type recordingStore struct {
	samples []models.TrafficSample
	paths   map[string]bool
	cutoffs []string
}

func (s *recordingStore) RecordBatch(ctx context.Context, samples []models.TrafficSample) error {
	s.samples = append(s.samples, samples...)
	return nil
}

func (s *recordingStore) PathsEnabled(ctx context.Context, projectID string) (bool, error) {
	return s.paths[projectID], nil
}

func (s *recordingStore) DeleteBefore(ctx context.Context, cutoff string) error {
	s.cutoffs = append(s.cutoffs, cutoff)
	return nil
}

func (s *recordingStore) RetentionCutoff(days int) string {
	return time.Now().UTC().AddDate(0, 0, -days).Format("2006-01-02T15:04")
}

func TestCollectorAttributesAndMasksPaths(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "access.json")
	line := `{"ClientAddr":"10.0.0.2:1000","DownstreamStatus":200,"Duration":2000000,"RequestHost":"shop.example.com","RequestPath":"/secret","DownstreamContentSize":10,"StartUTC":"2026-10-07T08:01:02Z"}` + "\n"
	line += `{"ClientAddr":"10.0.0.3:1000","DownstreamStatus":500,"Duration":1000000,"RequestHost":"unknown.example.com","RequestPath":"/","DownstreamContentSize":5,"StartUTC":"2026-10-07T08:01:03Z"}` + "\n"
	if err := os.WriteFile(path, []byte(line), 0o644); err != nil {
		t.Fatalf("seed log: %v", err)
	}
	store := &recordingStore{paths: map[string]bool{}}
	collector := NewCollector(path, store, func(host string) string {
		if host == "shop.example.com" {
			return "proj-1"
		}
		return ""
	})
	if err := collector.Collect(context.Background()); err != nil {
		t.Fatalf("collect: %v", err)
	}
	if len(store.samples) != 1 {
		t.Fatalf("expected 1 attributed sample, got %d", len(store.samples))
	}
	if store.samples[0].Path != "" {
		t.Fatalf("expected masked path, got %q", store.samples[0].Path)
	}
	if err := collector.Collect(context.Background()); err != nil {
		t.Fatalf("recollect: %v", err)
	}
	if len(store.samples) != 1 {
		t.Fatalf("expected offset to prevent reread, got %d", len(store.samples))
	}
	store.paths["proj-1"] = true
	collector.pathsAt = time.Time{}
	more := `{"ClientAddr":"10.0.0.2:1000","DownstreamStatus":200,"Duration":1000000,"RequestHost":"shop.example.com","RequestPath":"/visible","DownstreamContentSize":5,"StartUTC":"2026-10-07T08:02:00Z"}` + "\n"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("append log: %v", err)
	}
	if _, err := file.WriteString(more); err != nil {
		t.Fatalf("write log: %v", err)
	}
	file.Close()
	if err := collector.Collect(context.Background()); err != nil {
		t.Fatalf("collect more: %v", err)
	}
	if len(store.samples) != 2 || store.samples[1].Path != "/visible" {
		t.Fatalf("expected visible path after toggle, got %#v", store.samples)
	}
}
