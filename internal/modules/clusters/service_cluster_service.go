package clusters

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var (
	k3sVersionPattern       = regexp.MustCompile(`^v1\.[0-9]+\.[0-9]+\+k3s[0-9]+$`)
	networkInterfacePattern = regexp.MustCompile(`^[A-Za-z0-9_.:-]{1,32}$`)
)

type Store interface {
	Get(context.Context, string) (*Cluster, error)
	List(context.Context, string) ([]Cluster, error)
	Save(context.Context, *Cluster, int) error
	Observe(context.Context, string, string, string) error
}
type Projects interface {
	Get(context.Context, string) (*ProjectConfig, error)
}
type Servers interface {
	GetByID(context.Context, string) (*Server, error)
}
type Runner interface {
	Preflight(context.Context, ClusterNode, string, bool) error
	Script(context.Context, ClusterNode, string, string) error
	Kubectl(context.Context, ClusterNode, []string, string) (string, error)
	Host(context.Context, ClusterNode, string) (string, error)
	VerifyPorts(context.Context, ClusterNode, bool) error
}
type Gate interface{ AcquireVolume(string) (func(), error) }
type ClusterService struct {
	store      Store
	projects   Projects
	servers    Servers
	runner     Runner
	operations OperationsService
	gate       Gate
	http       *http.Client
}

func NewClusterService(store Store, projects Projects, servers Servers, runner Runner, ops OperationsService, gate Gate, httpClient *http.Client) *ClusterService {
	return &ClusterService{store: store, projects: projects, servers: servers, runner: runner, operations: ops, gate: gate, http: httpClient}
}
func (s *ClusterService) Get(ctx context.Context, id string) (*Cluster, error) {
	return s.store.Get(ctx, id)
}
func (s *ClusterService) List(ctx context.Context, project string) ([]Cluster, error) {
	return s.store.List(ctx, project)
}
func (s *ClusterService) Review(ctx context.Context, user, projectID string, request ClusterReviewRequest) (*OperationReview, error) {
	if request.Action == "recover" {
		return s.reviewUpgradeRecovery(ctx, user, projectID, request.Cluster.ID)
	}
	project, err := s.projects.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	cluster := request.Cluster
	if cluster.ID == "" {
		cluster.ID = uuid.NewString()
	}
	if _, err := uuid.Parse(cluster.ID); err != nil {
		return nil, fmt.Errorf("cluster ID must be a UUID")
	}
	if request.Action != "install" && request.Action != "join" && request.Action != "upgrade" && request.Action != "remove" {
		return nil, fmt.Errorf("unsupported cluster action")
	}
	if !k3sVersionPattern.MatchString(cluster.Version) {
		return nil, fmt.Errorf("select an exact K3s release version")
	}
	if strings.TrimSpace(cluster.Name) == "" || len(cluster.Name) > 100 || len(cluster.Nodes) == 0 || len(cluster.Nodes) > 32 {
		return nil, fmt.Errorf("cluster needs a name and one to thirty-two nodes")
	}
	cluster.ProjectID, cluster.OrganizationID = projectID, project.OrganizationID
	if cluster.Controls != 1 && cluster.Controls != 3 {
		if cluster.Controls != 0 {
			return nil, fmt.Errorf("select one or three control-plane servers")
		}
		cluster.Controls = 1
	}
	if cluster.Controls == 3 && len(cluster.Nodes) < 3 {
		return nil, fmt.Errorf("three control-plane servers require three nodes")
	}
	if err := s.validateNodes(ctx, user, &cluster); err != nil {
		return nil, err
	}
	previousVersion := ""
	previous := cluster.Revision
	existingNodes := 0
	if previous == 0 {
		if request.Action != "install" {
			return nil, fmt.Errorf("save a new cluster before changing its nodes")
		}
		secret := make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			return nil, err
		}
		cluster.JoinToken = hex.EncodeToString(secret)
		if err := s.store.Save(ctx, &cluster, 0); err != nil {
			return nil, err
		}
	} else {
		stored, err := s.store.Get(ctx, cluster.ID)
		if err != nil {
			return nil, err
		}
		if stored.ProjectID != projectID || stored.Revision != previous {
			return nil, fmt.Errorf("cluster changed; reload before reviewing")
		}
		previousVersion = stored.Version
		if stored.Status == "READY" {
			existingNodes = len(stored.Nodes)
		}
		if request.Action == "install" && stored.Status == "READY" {
			return nil, fmt.Errorf("ready cluster already installed; review join or upgrade instead")
		}
		if request.Action == "remove" {
			cluster = *stored
		} else {
			if err := preserveExistingNodes(stored, &cluster, request.Action); err != nil {
				return nil, err
			}
			if err := s.store.Save(ctx, &cluster, previous); err != nil {
				return nil, err
			}
		}
	}
	saved, err := s.store.Get(ctx, cluster.ID)
	if err != nil {
		return nil, err
	}
	for i, node := range saved.Nodes {
		if err := s.runner.Preflight(ctx, node, saved.ID, i == 0); err != nil {
			return nil, err
		}
	}
	plan := ClusterPlan{PreviousVersion: previousVersion, Cluster: *saved, Action: request.Action, ExistingNodes: existingNodes}
	if request.Action != "remove" {
		installer, err := s.downloadInstaller(ctx)
		if err != nil {
			return nil, err
		}
		plan.Installer = installer
		digest := sha256.Sum256([]byte(installer))
		plan.InstallerSHA256 = hex.EncodeToString(digest[:])
	}
	payload, err := json.Marshal(plan)
	if err != nil {
		return nil, err
	}
	effects := fmt.Sprintf("%s K3s %s on %d verified private-network nodes with %d control-plane server(s) and embedded etcd quorum. Install or update system services and routing/storage components. Keep private API 6443 and VXLAN 8472 restricted to these nodes. Interruptions can leave partial installation; inspect and prepare a retry. Installer SHA256: %s.", request.Action, saved.Version, len(saved.Nodes), saved.Controls, plan.InstallerSHA256)
	if request.Action == "remove" {
		effects = "Uninstall this owned K3s cluster from every saved node. Cluster workloads and local cluster storage may be destroyed; external backups remain separate."
	}
	return s.operations.Review(ctx, user, projectID, "cluster", "cluster:"+saved.ID, string(payload), clusterSnapshot(saved), effects)
}
func (s *ClusterService) validateNodes(ctx context.Context, user string, cluster *Cluster) error {
	ids, ips, fingerprints := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, node := range cluster.Nodes {
		if err := validateClusterNode(node); err != nil {
			return err
		}
		if ids[node.ServerID] || ips[node.PrivateIP] || fingerprints[node.Fingerprint] {
			return fmt.Errorf("nodes need distinct server identities, private addresses and host keys")
		}
		ids[node.ServerID], ips[node.PrivateIP], fingerprints[node.Fingerprint] = true, true, true
		server, err := s.servers.GetByID(ctx, node.ServerID)
		if err != nil {
			return err
		}
		if server.IsLocal || (server.UserID != user && server.OrganizationID != cluster.OrganizationID) {
			return fmt.Errorf("select authorized remote servers in this organization")
		}
	}
	return nil
}
func preserveExistingNodes(previous, next *Cluster, action string) error {
	if len(next.Nodes) < len(previous.Nodes) {
		return fmt.Errorf("node removal requires a reviewed cluster removal; nodes cannot be silently dropped")
	}
	for i, node := range previous.Nodes {
		if next.Nodes[i] != node {
			return fmt.Errorf("existing node identities and private network settings are immutable")
		}
	}
	if previous.Controls != 0 && next.Controls != 0 && previous.Controls != next.Controls {
		return fmt.Errorf("control-plane quorum cannot change without a reviewed cluster removal")
	}
	if next.Controls == 0 {
		next.Controls = previous.Controls
	}
	if action == "upgrade" {
		var oldMinor, oldPatch, oldBuild, newMinor, newPatch, newBuild int
		if _, err := fmt.Sscanf(previous.Version, "v1.%d.%d+k3s%d", &oldMinor, &oldPatch, &oldBuild); err != nil {
			return fmt.Errorf("stored release is invalid: %w", err)
		}
		if _, err := fmt.Sscanf(next.Version, "v1.%d.%d+k3s%d", &newMinor, &newPatch, &newBuild); err != nil {
			return fmt.Errorf("selected release is invalid: %w", err)
		}
		if newMinor < oldMinor || newMinor > oldMinor+1 || (newMinor == oldMinor && (newPatch < oldPatch || (newPatch == oldPatch && newBuild < oldBuild))) {
			return fmt.Errorf("upgrades cannot downgrade or skip a Kubernetes minor release")
		}
		if len(next.Nodes) != len(previous.Nodes) {
			return fmt.Errorf("review node joins separately from upgrades")
		}
	}
	if action == "join" && previous.Version != next.Version {
		return fmt.Errorf("join nodes using the current version; review upgrades separately")
	}
	return nil
}
func clusterSnapshot(cluster *Cluster) string {
	data := fmt.Appendf(nil, "%s|%d|%s|%d", cluster.ID, cluster.Revision, cluster.Version, cluster.Controls)
	for _, node := range cluster.Nodes {
		data = fmt.Appendf(data, "|%s|%s|%s|%s", node.ServerID, node.PrivateIP, node.Interface, node.Fingerprint)
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}
func (s *ClusterService) downloadInstaller(ctx context.Context) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://get.k3s.io", nil)
	if err != nil {
		return "", err
	}
	response, err := s.http.Do(request)
	if err != nil {
		return "", fmt.Errorf("download reviewed K3s installer: %w", err)
	}
	data, readErr := io.ReadAll(io.LimitReader(response.Body, 1024*1024+1))
	closeErr := response.Body.Close()
	if readErr != nil {
		return "", readErr
	}
	if closeErr != nil {
		return "", closeErr
	}
	if response.StatusCode != http.StatusOK || len(data) == 0 || len(data) > 1024*1024 {
		return "", fmt.Errorf("K3s installer unavailable or oversized")
	}
	return string(data), nil
}

func validateClusterNode(node ClusterNode) error {
	if _, err := uuid.Parse(node.ServerID); err != nil {
		return fmt.Errorf("select a valid registered server identity")
	}
	ip := net.ParseIP(node.PrivateIP)
	if ip == nil || !ip.IsPrivate() || ip.To4() == nil {
		return fmt.Errorf("cluster nodes need a private IPv4 address")
	}
	for _, value := range []string{"10.42.0.0/16", "10.43.0.0/16"} {
		_, network, err := net.ParseCIDR(value)
		if err != nil {
			return err
		}
		if network.Contains(ip) {
			return fmt.Errorf("private node addresses cannot overlap the default K3s pod or service network")
		}
	}
	if !networkInterfacePattern.MatchString(node.Interface) {
		return fmt.Errorf("select a valid private network interface")
	}
	fingerprint, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(node.Fingerprint, "SHA256:"))
	if !strings.HasPrefix(node.Fingerprint, "SHA256:") || err != nil || len(fingerprint) != 32 {
		return fmt.Errorf("verify the SSH host fingerprint before preparing a cluster operation")
	}
	return nil
}
