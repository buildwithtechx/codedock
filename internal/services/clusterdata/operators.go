package clusterdata

import (
	"codedock.run/codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func (s *Service) operatorManifests(ctx context.Context, cluster string, request models.ClusterDataRequest) ([]models.OperatorManifest, error) {
	pattern := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	for _, version := range []string{request.OperatorVersion, request.BarmanVersion, request.CertManagerVersion} {
		if !pattern.MatchString(version) {
			return nil, fmt.Errorf("operator dependencies require exact numeric release versions")
		}
	}
	var major, minor, patch int
	if _, err := fmt.Sscanf(request.OperatorVersion, "%d.%d.%d", &major, &minor, &patch); err != nil || major != 1 || minor < 26 {
		return nil, fmt.Errorf("Barman plugin requires CloudNativePG 1.26 or newer")
	}
	releases := []struct{ name, url string }{{"cert-manager", "https://github.com/cert-manager/cert-manager/releases/download/v" + request.CertManagerVersion + "/cert-manager.yaml"}, {"CloudNativePG", fmt.Sprintf("https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/release-%d.%d/releases/cnpg-%s.yaml", major, minor, request.OperatorVersion)}, {"Barman Cloud", "https://github.com/cloudnative-pg/plugin-barman-cloud/releases/download/v" + request.BarmanVersion + "/manifest.yaml"}}
	result := []models.OperatorManifest{}
	for _, release := range releases {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, release.url, nil)
		if err != nil {
			return nil, err
		}
		response, err := s.http.Do(request)
		if err != nil {
			return nil, err
		}
		data, readErr := io.ReadAll(io.LimitReader(response.Body, 8*1024*1024+1))
		closeErr := response.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		if closeErr != nil {
			return nil, closeErr
		}
		if response.StatusCode != http.StatusOK || len(data) == 0 || len(data) > 8*1024*1024 {
			return nil, fmt.Errorf("reviewed %s release manifest unavailable", release.name)
		}
		objects := []map[string]any{}
		decoder := yaml.NewDecoder(strings.NewReader(string(data)))
		for {
			var object map[string]any
			err := decoder.Decode(&object)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			if len(object) == 0 {
				continue
			}
			metadata, ok := object["metadata"].(map[string]any)
			if !ok {
				return nil, fmt.Errorf("operator resource metadata missing")
			}
			labels, ok := metadata["labels"].(map[string]any)
			if !ok {
				labels = map[string]any{}
			}
			labels["codedock.run/cluster"] = cluster
			metadata["labels"] = labels
			objects = append(objects, object)
		}
		encoded, err := json.Marshal(map[string]any{"apiVersion": "v1", "kind": "List", "items": objects})
		if err != nil {
			return nil, err
		}
		sum := sha256.Sum256(encoded)
		result = append(result, models.OperatorManifest{Name: release.name, Manifest: string(encoded), SHA256: hex.EncodeToString(sum[:])})
	}
	return result, nil
}
func (s *Service) operatorsReady(ctx context.Context, cluster *models.Cluster) error {
	for _, deployment := range []string{"cnpg-controller-manager", "barman-cloud"} {
		if _, err := s.commands.Kubectl(ctx, cluster.Nodes[0], []string{"-n", "cnpg-system", "rollout", "status", "deployment/" + deployment, "--timeout=5s"}, ""); err != nil {
			return fmt.Errorf("install and verify PostgreSQL operators first: %w", err)
		}
	}
	return nil
}
