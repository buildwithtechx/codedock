package backup

import (
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type volumeTestStore struct {
	*mockStore
	ownerErr error
}

func (s *volumeTestStore) VolumeRestoreOwner(context.Context, *models.BackupConfig) (string, error) {
	return "owner", s.ownerErr
}

func TestVolumeRestoreRequiresOwnedExistingStoppedTarget(t *testing.T) {
	for _, test := range []struct {
		name                                  string
		missing, foreign, running, incomplete bool
		valid                                 bool
	}{
		{name: "valid", valid: true}, {name: "missing", missing: true}, {name: "foreign", foreign: true}, {name: "running", running: true}, {name: "incomplete", incomplete: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case strings.HasSuffix(r.URL.Path, "/volumes/owned"):
					if test.missing {
						w.WriteHeader(404)
						fmt.Fprint(w, `{"message":"missing"}`)
						return
					}
					fmt.Fprint(w, `{"Name":"owned"}`)
				case strings.HasSuffix(r.URL.Path, "/containers/owner/json"):
					name := "owned"
					if test.foreign {
						name = "other"
					}
					fmt.Fprintf(w, `{"Id":"owner","State":{"Running":false},"Mounts":[{"Type":"volume","Name":%q}]}`, name)
				case strings.HasSuffix(r.URL.Path, "/containers/json"):
					if test.running {
						fmt.Fprint(w, `[{"Id":"writer","Mounts":[{"Type":"volume","Name":"owned"}]}]`)
					} else {
						fmt.Fprint(w, `[]`)
					}
				default:
					t.Errorf("unexpected Docker request %s", r.URL.Path)
					w.WriteHeader(404)
				}
			}))
			defer server.Close()
			dockerClient, err := client.NewClientWithOpts(client.WithHost(server.URL), client.WithVersion("1.47"), client.WithHTTPClient(server.Client()))
			if err != nil {
				t.Fatal(err)
			}
			defer dockerClient.Close()
			store := &volumeTestStore{mockStore: newMockStore()}
			status := models.BackupRecordStatusCompleted
			if test.incomplete {
				status = models.BackupRecordStatusRunning
			}
			store.records["record"] = &models.BackupRecord{ID: "record", BackupConfigID: "config", Status: status}
			store.configs["config"] = &models.BackupConfig{ID: "config", ServiceID: "service", VolumeName: "owned", Timeout: 12}
			manager := NewBackupManager(dockerClient, store, t.TempDir())
			target, err := manager.ValidateVolumeRestore(context.Background(), "record")
			if test.valid {
				if err != nil || target.VolumeName != "owned" || target.TimeoutSeconds != 12 {
					t.Fatalf("invalid target: %v %v", target, err)
				}
				if err := manager.RestoreVolume(context.Background(), "record", "other"); err == nil {
					t.Fatal("wrong confirmation accepted")
				}
			} else if err == nil {
				t.Fatal("unsafe target accepted")
			}
		})
	}
}
func TestVolumeRestoreInterruption(t *testing.T) {
	manager := NewBackupManager(nil, newMockStore(), t.TempDir())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	manager.restores["record"] = cancel
	if !manager.CancelVolumeRestore("record") {
		t.Fatal("active restore was not interrupted")
	}
	select {
	case <-ctx.Done():
	default:
		t.Fatal("restore context not cancelled")
	}
	if manager.CancelVolumeRestore("other") {
		t.Fatal("unknown restore was cancelled")
	}
}
