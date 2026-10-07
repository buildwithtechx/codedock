package clusters

import (
	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"time"
)

func (s *Service) Apply(ctx context.Context, user, id, confirmation string) error {
	operation, err := s.operations.Get(ctx, id)
	if err != nil {
		return err
	}
	if operation.Kind != "cluster" || operation.UserID != user {
		return fmt.Errorf("cluster operation not found")
	}
	var plan models.ClusterPlan
	if err := json.Unmarshal([]byte(operation.Payload), &plan); err != nil {
		return err
	}
	if plan.Action != "remove" {
		digest := sha256.Sum256([]byte(plan.Installer))
		if hex.EncodeToString(digest[:]) != plan.InstallerSHA256 {
			return fmt.Errorf("reviewed installer checksum mismatch")
		}
	}
	current, err := s.store.Get(ctx, plan.Cluster.ID)
	if err != nil {
		return err
	}
	if err := s.validateNodes(ctx, user, current); err != nil {
		return err
	}
	return s.operations.Apply(ctx, id, user, confirmation, clusterSnapshot(current), func(ctx context.Context, op *models.Operation, progress func(string, string) error) (runErr error) {
		releases, err := s.acquireTargets(current)
		if err != nil {
			return err
		}
		defer func() {
			for i := len(releases) - 1; i >= 0; i-- {
				releases[i]()
			}
		}()
		defer func() {
			final, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			status, message := "READY", ""
			if runErr != nil {
				status, message = "FAILED", runErr.Error()
			}
			if plan.Action == "remove" && runErr == nil {
				status = "REMOVED"
			}
			if err := s.store.Observe(final, current.ID, status, message); err != nil {
				slog.Error("persist cluster outcome", "cluster", current.ID, "error", err)
				runErr = fmt.Errorf("persist cluster outcome: %w", err)
			}
		}()
		latest, err := s.store.Get(ctx, current.ID)
		if err != nil {
			return err
		}
		if clusterSnapshot(latest) != op.Snapshot {
			return fmt.Errorf("cluster changed before execution; prepare another review")
		}
		if plan.Action == "recover" {
			if err := s.resumeUpgrade(ctx, &plan); err != nil {
				return err
			}
			store, ok := s.store.(UpgradeStore)
			if !ok {
				return fmt.Errorf("upgrade journal storage unavailable")
			}
			if err := store.FinishUpgrade(ctx, current.ID, plan.Cluster.Version, "READY", ""); err != nil {
				return err
			}
			return progress("RECOVERED", "Cluster reached the reviewed release")
		}
		if err := s.store.Observe(ctx, current.ID, "APPLYING", ""); err != nil {
			return err
		}
		for i, node := range current.Nodes {
			if err := s.runner.Preflight(ctx, node, current.ID, i == 0); err != nil {
				return err
			}
		}
		order := make([]int, len(current.Nodes))
		for i := range order {
			order[i] = i
		}
		if plan.Action == "remove" {
			for i := range order {
				order[i] = len(order) - 1 - i
			}
		}
		if plan.Action == "upgrade" {
			store, ok := s.store.(UpgradeStore)
			if !ok {
				return fmt.Errorf("upgrade recovery storage unavailable")
			}
			plan.OperationID = op.ID
			if err := store.BeginUpgrade(ctx, &plan); err != nil {
				return err
			}
			defer func() { runErr = s.finishUpgrade(&plan, runErr) }()
		}
		for _, i := range order {
			node := current.Nodes[i]
			if plan.Action == "join" && i < plan.ExistingNodes {
				continue
			}
			phase := "INSTALLING"
			role := "agent"
			if i == 0 {
				role = "leader"
			} else if i < current.Controls {
				role = "control"
				phase = "JOINING"
			} else if i > 0 {
				phase = "JOINING"
			}
			if plan.Action == "remove" {
				phase = "REMOVING"
			}
			if err := progress(phase, "Processing node "+node.ServerID); err != nil {
				return err
			}
			script := installationScript(current, node, role, plan.Installer)
			if plan.Action == "remove" {
				script = removalScript(current.ID, i == 0)
			}
			if err := s.runner.Script(ctx, node, op.ID, script); err != nil {
				return fmt.Errorf("node %s operation failed: %w", node.ServerID, err)
			}
			if plan.Action == "upgrade" {
				if err := s.waitRelease(ctx, current, node); err != nil {
					return err
				}
			}
		}
		if plan.Action == "remove" {
			return progress("REMOVED", "Owned K3s cluster removed")
		}
		if err := progress("READINESS", "Waiting for all saved nodes to become Ready"); err != nil {
			return err
		}
		for _, node := range current.Nodes {
			name := "codedock-" + node.ServerID
			if _, err := s.runner.Kubectl(ctx, current.Nodes[0], []string{"wait", "--for=condition=Ready", "node/" + name, "--timeout=180s"}, ""); err != nil {
				return err
			}
		}
		if err := progress("VERIFYING", "Verifying forwarding, ports and cross-node connectivity"); err != nil {
			return err
		}
		for i, node := range current.Nodes {
			if err := s.runner.VerifyPorts(ctx, node, i < current.Controls); err != nil {
				return err
			}
		}
		return nil
	})
}
func (s *Service) acquireTargets(cluster *models.Cluster) ([]func(), error) {
	if s.gate == nil {
		return nil, fmt.Errorf("cluster operation locking unavailable")
	}
	names := []string{"cluster:" + cluster.ID}
	for _, node := range cluster.Nodes {
		names = append(names, "server:"+node.ServerID)
	}
	sort.Strings(names)
	releases := []func(){}
	for _, name := range names {
		release, err := s.gate.AcquireVolume(name)
		if err != nil {
			for i := len(releases) - 1; i >= 0; i-- {
				releases[i]()
			}
			return nil, err
		}
		releases = append(releases, release)
	}
	return releases, nil
}
func installationScript(cluster *models.Cluster, node models.ClusterNode, role string, installer string) string {
	extra := ""
	execution := ""
	server := ""
	switch role {
	case "leader":
		role = "server"
		extra = " --advertise-address " + node.PrivateIP + " --bind-address " + node.PrivateIP + " --tls-san " + node.PrivateIP + " --cluster-init"
		execution = role + " --node-name codedock-" + node.ServerID + " --node-ip " + node.PrivateIP + " --flannel-iface " + node.Interface + " --node-label codedock.run/cluster=" + cluster.ID + extra
	case "control":
		role = "server"
		extra = " --advertise-address " + node.PrivateIP + " --bind-address " + node.PrivateIP + " --tls-san " + node.PrivateIP
		execution = role + " --node-name codedock-" + node.ServerID + " --node-ip " + node.PrivateIP + " --flannel-iface " + node.Interface + " --node-label codedock.run/cluster=" + cluster.ID + extra
		server = "K3S_URL=" + ssh.ShellQuote("https://"+cluster.Nodes[0].PrivateIP+":6443") + " "
	default:
		role = "agent"
		execution = role + " --node-name codedock-" + node.ServerID + " --node-ip " + node.PrivateIP + " --flannel-iface " + node.Interface + " --node-label codedock.run/cluster=" + cluster.ID + extra
		server = "K3S_URL=" + ssh.ShellQuote("https://"+cluster.Nodes[0].PrivateIP+":6443") + " "
	}
	return "set -eu\n" + "umask 077\nsysctl -w net.ipv4.ip_forward=1 >/dev/null\nprintf %s " + ssh.ShellQuote(cluster.ID) + " > /etc/codedock-cluster\n" + "installer=$(mktemp)\ntrap 'rm -f \"$installer\"' EXIT\nprintf %s " + ssh.ShellQuote(installer) + " > \"$installer\"\n" + server + "INSTALL_K3S_VERSION=" + ssh.ShellQuote(cluster.Version) + " INSTALL_K3S_EXEC=" + ssh.ShellQuote(execution) + " K3S_TOKEN=" + ssh.ShellQuote(cluster.JoinToken) + " sh \"$installer\"\n"
}
func removalScript(clusterID string, leader bool) string {
	command := "/usr/local/bin/k3s-agent-uninstall.sh"
	if leader {
		command = "/usr/local/bin/k3s-uninstall.sh"
	}
	return "set -eu\nif ! test -e /etc/codedock-cluster; then if command -v k3s >/dev/null || test -e /etc/rancher/k3s; then exit 1; else exit 0; fi; fi\ntest \"$(cat /etc/codedock-cluster)\" = " + ssh.ShellQuote(clusterID) + "\n" + command + "\nrm /etc/codedock-cluster\n"
}
