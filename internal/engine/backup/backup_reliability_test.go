package backup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"codedock.run/codedock/internal/models"
	"github.com/docker/docker/client"
)

type databaseBackupStore struct {
	*mockStore
	database *models.Database
}

func (s *databaseBackupStore) GetDatabase(id string) (*models.Database, error) {
	return s.database, nil
}

type failingBackupStore struct{ *mockStore }

func (s *failingBackupStore) UpdateBackupRecord(opts models.UpdateBackupRecordOpts) error {
	return errors.New("storage unavailable")
}

func TestScheduleTimezoneAndInvalidReplacement(t *testing.T) {
	now := time.Date(2026, 10, 6, 0, 0, 0, 0, time.UTC)
	for _, expression := range []string{"0 2 * * *", "0 0 2 * * *", "@daily", "CRON_TZ=Africa/Lagos 0 2 * * *"} {
		schedule, err := ParseSchedule(expression, "Africa/Lagos")
		if err != nil {
			t.Fatal(err)
		}
		next := schedule.Next(now)
		expectedHour := 1
		if expression == "@daily" {
			expectedHour = 23
		}
		if next.UTC().Hour() != expectedHour {
			t.Fatalf("%s scheduled at %s", expression, next)
		}
	}
	for _, zone := range []string{"Invalid/Zone", "Africa/Lagos"} {
		expression := "0 2 * * *"
		if zone == "Africa/Lagos" {
			expression = "CRON_TZ=UTC 0 2 * * *"
		}
		if _, err := ParseSchedule(expression, zone); err == nil {
			t.Fatalf("invalid timezone combination accepted: %s", zone)
		}
	}
	manager := NewBackupManager(nil, nil, t.TempDir())
	cfg := &models.BackupConfig{ID: "backup", Schedule: "0 2 * * *", Timezone: "Africa/Lagos", BackupEnabled: true, Status: models.BackupConfigStatusActive}
	if err := manager.RegisterBackup(cfg); err != nil {
		t.Fatal(err)
	}
	entry := manager.entries[cfg.ID]
	cfg.Timezone = "Invalid/Zone"
	if err := manager.RegisterBackup(cfg); err == nil {
		t.Fatal("invalid replacement accepted")
	}
	if manager.entries[cfg.ID] != entry {
		t.Fatal("invalid replacement removed the valid schedule")
	}
}

func TestUploadUsesDestinationPrefix(t *testing.T) {
	requestPaths := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPaths <- r.URL.Path
		if r.Method != http.MethodPut {
			t.Errorf("unexpected method %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	manager := NewBackupManager(nil, nil, t.TempDir())
	dest := &models.S3Destination{Endpoint: server.URL, Bucket: "archives", PathPrefix: "/team/database backups/"}
	uri, err := manager.uploadToS3(context.Background(), dest, "backup.sql", []byte("dump"))
	if err != nil {
		t.Fatal(err)
	}
	requestPath := <-requestPaths
	if requestPath != "/archives/team/database backups/backup.sql" {
		t.Fatalf("wrong object path: %s", requestPath)
	}
	if uri != "s3://archives/team/database backups/backup.sql" {
		t.Fatalf("wrong archive URI: %s", uri)
	}
}

func TestUnavailableBackupRuntimeAndUnsupportedEngineFail(t *testing.T) {
	store := &databaseBackupStore{mockStore: newMockStore(), database: &models.Database{ID: "database", Engine: models.DatabaseEngine("nats")}}
	manager := NewBackupManager(nil, store, t.TempDir())
	cfg := &models.BackupConfig{ID: "backup", DatabaseID: "database", BackupEnabled: true}
	if _, _, _, err := manager.buildDumpCommand(cfg); err == nil {
		t.Fatal("engine without dump metadata accepted")
	}
	if _, _, err := manager.buildRestoreCommand(cfg); err == nil {
		t.Fatal("engine without restore metadata accepted")
	}
	if _, _, err := manager.executeVolumeBackup(context.Background(), "data"); err == nil {
		t.Fatal("missing volume runtime reported success")
	}
	if _, _, err := manager.executeDump(context.Background(), "container", nil, "database"); err == nil {
		t.Fatal("missing database runtime reported success")
	}
	if err := manager.executeRestore(context.Background(), "container", nil, []byte("dump")); err == nil {
		t.Fatal("missing restore runtime reported success")
	}
	cfg.DatabaseID = ""
	cfg.VolumeName = "data"
	store.configs[cfg.ID] = cfg
	if _, err := manager.TriggerBackup(context.Background(), cfg.ID); err == nil {
		t.Fatal("unavailable volume backup reported success")
	}
	for _, record := range store.records {
		if record.Status != models.BackupRecordStatusFailed {
			t.Fatalf("wrong failed run status: %s", record.Status)
		}
	}
}

func TestCompletionPersistenceFailureIsReturned(t *testing.T) {
	manager := NewBackupManager(nil, &failingBackupStore{newMockStore()}, t.TempDir())
	record, err := manager.finalizeBackupRecord(FinalizeBackupOpts{Record: &models.BackupRecord{ID: "record"}})
	if err == nil || record != nil {
		t.Fatal("unpersisted completion reported success")
	}
}

func TestCommandCompletionUsesExitStatus(t *testing.T) {
	for _, exitCode := range []int{0, 2} {
		t.Run(fmt.Sprint(exitCode), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !strings.HasSuffix(r.URL.Path, "/exec/command/json") {
					t.Errorf("unexpected request: %s", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"Running":false,"ExitCode":%d}`, exitCode)
			}))
			defer server.Close()
			dockerClient, err := client.NewClientWithOpts(client.WithHost(server.URL), client.WithVersion("1.47"), client.WithHTTPClient(server.Client()))
			if err != nil {
				t.Fatal(err)
			}
			defer dockerClient.Close()
			manager := NewBackupManager(dockerClient, nil, t.TempDir())
			err = manager.waitExecSuccess(context.Background(), "command", "backup")
			if (err != nil) != (exitCode != 0) {
				t.Fatalf("exit %d returned %v", exitCode, err)
			}
		})
	}
}
