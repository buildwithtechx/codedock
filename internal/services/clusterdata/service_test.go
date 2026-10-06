package clusterdata

import (
	"codedock.run/codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type dataStoreFixture struct {
	Store
	plan *models.ClusterDataPlan
}

func (s dataStoreFixture) Get(context.Context, string) (*models.ClusterDataPlan, error) {
	return s.plan, nil
}
func TestCredentialsRequireMatchingProjectAndCluster(t *testing.T) {
	service := &Service{store: dataStoreFixture{plan: &models.ClusterDataPlan{Record: models.ClusterData{ProjectID: "owner", ClusterID: "cluster", Status: "READY"}, Password: "secret"}}}
	if _, err := service.Credentials(context.Background(), "foreign", "cluster", "database"); err == nil {
		t.Fatal("foreign project received credentials")
	}
	if _, err := service.Credentials(context.Background(), "owner", "foreign", "database"); err == nil {
		t.Fatal("foreign cluster received credentials")
	}
}

type manifestTransport struct{ calls int }

func (t *manifestTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	t.calls++
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("apiVersion: v1\nkind: Namespace\nmetadata:\n  name: fixture\n")), Header: make(http.Header)}, nil
}
func TestOperatorReviewPinsAndOwnsFetchedManifests(t *testing.T) {
	transport := &manifestTransport{}
	service := &Service{http: &http.Client{Transport: transport}}
	request := models.ClusterDataRequest{OperatorVersion: "1.28.4", BarmanVersion: "0.15.1", CertManagerVersion: "1.21.2"}
	manifests, err := service.operatorManifests(context.Background(), "owned-cluster", request)
	if err != nil || len(manifests) != 3 {
		t.Fatal(err)
	}
	for _, manifest := range manifests {
		sum := sha256.Sum256([]byte(manifest.Manifest))
		if manifest.SHA256 != hex.EncodeToString(sum[:]) {
			t.Fatal("review did not bind exact manifest")
		}
		var list struct {
			Items []struct {
				Metadata struct {
					Labels map[string]string `json:"labels"`
				} `json:"metadata"`
			} `json:"items"`
		}
		if err := json.Unmarshal([]byte(manifest.Manifest), &list); err != nil {
			t.Fatal(err)
		}
		if list.Items[0].Metadata.Labels["codedock.run/cluster"] != "owned-cluster" {
			t.Fatal("operator ownership missing")
		}
	}
	request.OperatorVersion = "latest"
	calls := transport.calls
	if _, err := service.operatorManifests(context.Background(), "owned-cluster", request); err == nil || transport.calls != calls {
		t.Fatal("unreviewed latest version downloaded")
	}
}

func TestDatabaseErrorsDoNotExposeCredentials(t *testing.T) {
	plan := &models.ClusterDataPlan{Password: "password-sensitive", Manifest: `{"items":[{"kind":"Secret","stringData":{"ACCESS_SECRET_KEY":"s3-sensitive"}}]}`}
	err := redactDataError(plan, fmt.Errorf("password-sensitive s3-sensitive"))
	if strings.Contains(err.Error(), "sensitive") {
		t.Fatal("credentials leaked through an operation error")
	}
}
