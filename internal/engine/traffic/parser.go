package traffic

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

type AccessEntry struct {
	ClientAddr         string `json:"ClientAddr"`
	ClientHost         string `json:"ClientHost"`
	DownstreamStatus   int    `json:"DownstreamStatus"`
	Duration           int64  `json:"Duration"`
	RequestHost        string `json:"RequestHost"`
	RequestMethod      string `json:"RequestMethod"`
	RequestPath        string `json:"RequestPath"`
	DownstreamContentSize int64 `json:"DownstreamContentSize"`
	OriginContentSize  int64  `json:"OriginContentSize"`
	StartUTC           string `json:"StartUTC"`
}

func ParseAccessLine(line string) (AccessEntry, error) {
	var entry AccessEntry
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return entry, fmt.Errorf("empty access line")
	}
	if err := json.Unmarshal([]byte(trimmed), &entry); err != nil {
		return entry, fmt.Errorf("decode access line: %w", err)
	}
	if entry.RequestHost == "" {
		return entry, fmt.Errorf("access line has no host")
	}
	return entry, nil
}

func (e AccessEntry) Timestamp() time.Time {
	if parsed, err := time.Parse(time.RFC3339, e.StartUTC); err == nil {
		return parsed.UTC()
	}
	return time.Now().UTC()
}

func (e AccessEntry) ClientIP() string {
	if e.ClientHost != "" {
		return e.ClientHost
	}
	host, _, err := net.SplitHostPort(e.ClientAddr)
	if err != nil {
		return e.ClientAddr
	}
	return host
}

func (e AccessEntry) Bytes() int64 {
	return e.DownstreamContentSize + e.OriginContentSize
}

func (e AccessEntry) DurationMs() int64 {
	return e.Duration / 1000000
}

func (e AccessEntry) Path() string {
	path := e.RequestPath
	if index := strings.Index(path, "?"); index >= 0 {
		path = path[:index]
	}
	if path == "" {
		return "/"
	}
	return path
}

func ClassifyCountry(ip string) string {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return "unknown"
	}
	if parsed.IsLoopback() {
		return "local"
	}
	if parsed.IsPrivate() || parsed.IsLinkLocalUnicast() || parsed.IsLinkLocalMulticast() {
		return "private"
	}
	return "unknown"
}
