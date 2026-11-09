package storage

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type testVolumeGate struct {
	mu     sync.Mutex
	active map[string]bool
}

func newTestVolumeGate() *testVolumeGate {
	return &testVolumeGate{active: make(map[string]bool)}
}

func (g *testVolumeGate) AcquireVolume(name string) (func(), error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.active[name] {
		return nil, fmt.Errorf("volume %s is busy; wait for its current operation to finish", name)
	}
	g.active[name] = true
	var once sync.Once
	return func() {
		once.Do(func() {
			g.mu.Lock()
			delete(g.active, name)
			g.mu.Unlock()
		})
	}, nil
}

func TestVolumeRestoreStreamsArchiveAndRemovesHelper(t *testing.T) {
	gate := newTestVolumeGate()
	archiveReceived := make(chan string, 1)
	inputDone := make(chan struct{})
	removed := make(chan struct{}, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/images/alpine/json"):
			w.WriteHeader(404)
			fmt.Fprint(w, `{"message":"missing image"}`)
		case strings.HasSuffix(r.URL.Path, "/images/create"):
			fmt.Fprint(w, `{"status":"pulled"}`)
		case strings.HasSuffix(r.URL.Path, "/volumes/owned"):
			fmt.Fprint(w, `{"Name":"owned","CreatedAt":"original"}`)
		case strings.HasSuffix(r.URL.Path, "/containers/owner/json"):
			fmt.Fprint(w, `{"Id":"owner","State":{"Running":false},"Mounts":[{"Type":"volume","Name":"owned"}]}`)
		case strings.HasSuffix(r.URL.Path, "/containers/json"):
			fmt.Fprint(w, `[]`)
		case strings.HasSuffix(r.URL.Path, "/containers/create"):
			fmt.Fprint(w, `{"Id":"helper","Warnings":[]}`)
		case strings.HasSuffix(r.URL.Path, "/containers/helper/attach"):
			conn, writer, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Error(err)
				return
			}
			defer conn.Close()
			if _, err := writer.WriteString("HTTP/1.1 101 UPGRADED\r\nContent-Type: application/vnd.docker.raw-stream\r\nConnection: Upgrade\r\nUpgrade: tcp\r\n\r\n"); err != nil {
				t.Error(err)
				return
			}
			if err := writer.Flush(); err != nil {
				t.Error(err)
				return
			}
			data, err := io.ReadAll(conn)
			if err != nil {
				t.Error(err)
			}
			archiveReceived <- string(data)
			close(inputDone)
		case strings.HasSuffix(r.URL.Path, "/containers/helper/start"):
			for _, key := range []string{"owned", "service:service"} {
				if release, err := gate.AcquireVolume(key); err == nil {
					release()
					t.Errorf("restore did not fence %s", key)
				}
			}
			w.WriteHeader(http.StatusNoContent)
		case strings.HasSuffix(r.URL.Path, "/containers/helper/wait"):
			select {
			case <-inputDone:
				fmt.Fprint(w, `{"StatusCode":0}`)
			case <-r.Context().Done():
				return
			}
		case strings.HasSuffix(r.URL.Path, "/containers/helper") && r.Method == http.MethodDelete:
			removed <- struct{}{}
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	dockerClient, err := client.NewClientWithOpts(client.WithHost(strings.Replace(server.URL, "http://", "tcp://", 1)), client.WithVersion("1.47"), client.WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}
	defer dockerClient.Close()
	filename := filepath.Join(t.TempDir(), "archive.tar.gz")
	if err := os.WriteFile(filename, []byte("archive input"), 0600); err != nil {
		t.Fatal(err)
	}
	store := &volumeTestStore{mockStore: newMockStore()}
	store.configs["config"] = &BackupConfig{ID: "config", ServiceID: "service", VolumeName: "owned", Timeout: 5}
	store.records["record"] = &BackupRecord{ID: "record", BackupConfigID: "config", Status: BackupRecordStatusCompleted, FilePath: filename}
	manager := NewBackupManager(dockerClient, store, t.TempDir())
	manager.SetVolumeOperations(gate)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := manager.RestoreVolume(ctx, "record", "owned"); err != nil {
		t.Fatal(err)
	}
	if received := <-archiveReceived; received != "archive input" {
		t.Fatalf("wrong archive input %q", received)
	}
	select {
	case <-removed:
	default:
		t.Fatal("restore helper was not removed")
	}
	for _, key := range []string{"owned", "service:service"} {
		release, err := gate.AcquireVolume(key)
		if err != nil {
			t.Fatal(err)
		}
		release()
	}
	if manager.CancelVolumeRestore("record") {
		t.Fatal("completed restore remained active")
	}
}

func TestVolumeBackupCannotReadDuringRestore(t *testing.T) {
	gate := newTestVolumeGate()
	release, err := gate.AcquireVolume("owned")
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	cli, err := client.NewClientWithOpts(client.WithHost("http://127.0.0.1:1"), client.WithVersion("1.47"))
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	manager := NewBackupManager(cli, newMockStore(), t.TempDir())
	manager.SetVolumeOperations(gate)
	if _, _, err := manager.executeVolumeBackup(context.Background(), "owned"); err == nil || !strings.Contains(err.Error(), "busy") {
		t.Fatalf("backup escaped restore fence: %v", err)
	}
}
