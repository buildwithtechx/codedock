package kubernetes

import (
	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net"
	"regexp"
	"strings"
)

type ClusterServers interface {
	GetByID(context.Context, string) (*models.Server, error)
}
type ClusterRunner struct{ servers ClusterServers }

func NewClusterRunner(servers ClusterServers) *ClusterRunner { return &ClusterRunner{servers: servers} }
func ValidateNode(node models.ClusterNode) error {
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
	if !regexp.MustCompile(`^[A-Za-z0-9_.:-]{1,32}$`).MatchString(node.Interface) {
		return fmt.Errorf("select a valid private network interface")
	}
	fingerprint, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(node.Fingerprint, "SHA256:"))
	if !strings.HasPrefix(node.Fingerprint, "SHA256:") || err != nil || len(fingerprint) != 32 {
		return fmt.Errorf("verify the SSH host fingerprint before preparing a cluster operation")
	}
	return nil
}
func (r *ClusterRunner) client(ctx context.Context, node models.ClusterNode) (*ssh.Client, error) {
	if err := ValidateNode(node); err != nil {
		return nil, err
	}
	server, err := r.servers.GetByID(ctx, node.ServerID)
	if err != nil {
		return nil, err
	}
	host := server.SSHHost
	if host == "" {
		host = server.IPAddress
	}
	key := server.SSHPrivateKey
	if key == "" {
		key = server.SSHKey
	}
	return ssh.NewClient(ssh.Config{Host: host, Port: server.SSHPort, User: server.SSHUser, Key: key, Password: server.SSHPassword, Fingerprint: node.Fingerprint})
}
func closeClient(client *ssh.Client) {
	if err := client.Close(); err != nil {
		slog.Warn("close cluster SSH connection", "error", err)
	}
}
func (r *ClusterRunner) Preflight(ctx context.Context, node models.ClusterNode, clusterID string, leader bool) error {
	client, err := r.client(ctx, node)
	if err != nil {
		return err
	}
	defer closeClient(client)
	cpu, memory := 1, 524288
	if leader {
		cpu, memory = 2, 2097152
	}
	command := fmt.Sprintf(`set -eu
test "$(uname -s)" = Linux
command -v systemctl >/dev/null
command -v setsid >/dev/null
command -v curl >/dev/null
test -d /sys/fs/cgroup
test "$(nproc)" -ge %d
test "$(awk '/MemTotal/ {print $2}' /proc/meminfo)" -ge %d
ip -4 -o addr show dev %s | grep -F ' %s/' >/dev/null
if test -e /etc/codedock-cluster; then test "$(cat /etc/codedock-cluster)" = %s; elif test -e /etc/rancher/k3s/config.yaml || command -v k3s >/dev/null; then exit 1; fi
`, cpu, memory, ssh.ShellQuote(node.Interface), node.PrivateIP, ssh.ShellQuote(clusterID))
	_, err = client.RunWithInput(ctx, client.RootCommand("sh -s"), strings.NewReader(command))
	if err != nil {
		return fmt.Errorf("node %s failed Linux, capacity, private networking or ownership validation", node.ServerID)
	}
	return nil
}
func (r *ClusterRunner) Script(ctx context.Context, node models.ClusterNode, id, script string) error {
	client, err := r.client(ctx, node)
	if err != nil {
		return err
	}
	defer closeClient(client)
	_, err = client.RunManagedScript(ctx, id, script)
	return err
}
func (r *ClusterRunner) Kubectl(ctx context.Context, node models.ClusterNode, args []string, input string) (string, error) {
	client, err := r.client(ctx, node)
	if err != nil {
		return "", err
	}
	defer closeClient(client)
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = ssh.ShellQuote(arg)
	}
	return client.RunWithInput(ctx, client.RootCommand("k3s kubectl "+strings.Join(quoted, " ")), strings.NewReader(input))
}

func (r *ClusterRunner) Host(ctx context.Context, node models.ClusterNode, script string) (string, error) {
	client, err := r.client(ctx, node)
	if err != nil {
		return "", err
	}
	defer closeClient(client)
	return client.RunWithInput(ctx, client.RootCommand("sh -s"), strings.NewReader(script))
}
