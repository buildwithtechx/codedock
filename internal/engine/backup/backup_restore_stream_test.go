package backup

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/client"
)

func TestRestoreEarlyExitUnblocksInputAndPreservesDiagnostics(t *testing.T) {
	release := make(chan struct{})
	defer close(release)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/containers/database/json"):
			fmt.Fprint(w, `{"Id":"database","State":{"Running":true}}`)
		case strings.HasSuffix(r.URL.Path, "/containers/database/exec"):
			fmt.Fprint(w, `{"Id":"restore"}`)
		case strings.HasSuffix(r.URL.Path, "/exec/restore/json"):
			fmt.Fprint(w, `{"Running":false,"ExitCode":1}`)
		case strings.HasSuffix(r.URL.Path, "/exec/restore/start"):
			conn, writer, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Error(err)
				return
			}
			defer conn.Close()
			writer.WriteString("HTTP/1.1 101 UPGRADED\r\nContent-Type: application/vnd.docker.raw-stream\r\nConnection: Upgrade\r\nUpgrade: tcp\r\n\r\n")
			diagnostic := []byte("permission denied")
			frame := make([]byte, 8)
			frame[0] = 2
			binary.BigEndian.PutUint32(frame[4:], uint32(len(diagnostic)))
			writer.Write(frame)
			writer.Write(diagnostic)
			if err := writer.Flush(); err != nil {
				t.Error(err)
				return
			}
			if err := conn.(*net.TCPConn).CloseWrite(); err != nil {
				t.Error(err)
				return
			}
			<-release
		default:
			t.Errorf("unexpected request: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	dockerClient, err := client.NewClientWithOpts(client.WithHost(strings.Replace(server.URL, "http://", "tcp://", 1)), client.WithVersion("1.47"), client.WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}
	defer dockerClient.Close()
	manager := NewBackupManager(dockerClient, nil, t.TempDir())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- manager.executeRestore(ctx, "database", []string{"restore"}, make([]byte, 16*1024*1024))
	}()
	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "restore command exited with status 1") || !strings.Contains(err.Error(), "permission denied") {
			t.Fatalf("missing restore failure diagnostics: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("restore waited for unread archive input after command exit")
	}
}
