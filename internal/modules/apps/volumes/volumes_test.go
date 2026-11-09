package volumes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockVolumeService struct {
	volumes map[string]*ServiceVolume
}

func newMockVolumeService() *mockVolumeService {
	return &mockVolumeService{volumes: make(map[string]*ServiceVolume)}
}

func (m *mockVolumeService) CreateVolume(ctx context.Context, vol *ServiceVolume) (*ServiceVolume, error) {
	if vol.ID == "" {
		vol.ID = "vol-1"
	}
	m.volumes[vol.ID] = vol
	return vol, nil
}

func (m *mockVolumeService) ListVolumes(ctx context.Context, serviceID string) ([]ServiceVolume, error) {
	var list []ServiceVolume
	for _, v := range m.volumes {
		if v.ServiceID == serviceID {
			list = append(list, *v)
		}
	}
	return list, nil
}

func (m *mockVolumeService) DeleteVolume(ctx context.Context, id string) error {
	delete(m.volumes, id)
	return nil
}

func (m *mockVolumeService) EmptyVolume(ctx context.Context, serviceID, volumeID string) (*ServiceVolume, string, error) {
	vol := m.volumes[volumeID]
	prev := vol.HostPath
	vol.HostPath = "codedock-fresh-123"
	return vol, prev, nil
}

func TestVolumeRoutes(t *testing.T) {
	e := echo.New()
	svc := newMockVolumeService()
	h := NewVolumeHandler(svc, nil)

	g := e.Group("")
	RegisterRoutes(g, h)

	// Create
	body := `{"hostPath":"/data","containerPath":"/app/data"}`
	req := httptest.NewRequest(http.MethodPost, "/apps/srv-1/volumes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	// List
	req = httptest.NewRequest(http.MethodGet, "/apps/srv-1/volumes", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var listResp struct {
		Data []ServiceVolume `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(listResp.Data) != 1 {
		t.Fatalf("expected 1 item, got %d", len(listResp.Data))
	}

	// Empty
	req = httptest.NewRequest(http.MethodPost, "/apps/srv-1/volumes/vol-1/empty", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/apps/srv-1/volumes/vol-1", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}
