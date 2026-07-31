package system

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"codedock.run/codedock/internal/models"
)

type TakeoverScanner struct{}

func NewTakeoverScanner() *TakeoverScanner {
	return &TakeoverScanner{}
}

type dockerPsRow struct {
	ID    string `json:"ID"`
	Names string `json:"Names"`
	Image string `json:"Image"`
	State string `json:"State"`
	Ports string `json:"Ports"`
}

type dockerInspectMount struct {
	Source string `json:"Source"`
	Type   string `json:"Type"`
}

type dockerInspectResult struct {
	Name   string `json:"Name"`
	Config struct {
		Image  string            `json:"Image"`
		Env    []string          `json:"Env"`
		Labels map[string]string `json:"Labels"`
	} `json:"Config"`
	HostConfig struct {
		Binds        []string               `json:"Binds"`
		PortBindings map[string]interface{} `json:"PortBindings"`
	} `json:"HostConfig"`
	Mounts []dockerInspectMount `json:"Mounts"`
}

func (s *TakeoverScanner) Scan(ctx context.Context, req models.TakeoverScanRequest) (*models.DiscoveredStack, error) {
	host, port, err := net.SplitHostPort(req.Host)
	if err != nil {
		host = req.Host
		port = "22"
	}
	addr := net.JoinHostPort(host, port)
	client, err := dialSSH(ctx, addr, req.SSHUser, req.SSHKey, req.SSHFingerprint)
	if err != nil {
		return nil, fmt.Errorf("ssh connect: %w", err)
	}
	defer client.Close()

	psOut, err := runSSHCommand(client, `docker ps -a --format '{"ID":"{{.ID}}","Names":"{{.Names}}","Image":"{{.Image}}","State":"{{.State}}","Ports":"{{.Ports}}"}'`)
	if err != nil {
		return nil, fmt.Errorf("docker ps: %w", err)
	}

	var rows []dockerPsRow
	for _, line := range strings.Split(strings.TrimSpace(psOut), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var row dockerPsRow
		if err := json.Unmarshal([]byte(line), &row); err == nil {
			rows = append(rows, row)
		}
	}

	var containers []models.DiscoveredContainer
	composeSet := map[string]struct{}{}
	var platform models.TakeoverPlatform

	for _, row := range rows {
		ids := strings.TrimPrefix(row.ID, "/")
		inspectOut, err := runSSHCommand(client, "docker inspect "+ids)
		if err != nil {
			continue
		}
		var results []dockerInspectResult
		if err := json.Unmarshal([]byte(inspectOut), &results); err != nil || len(results) == 0 {
			continue
		}
		res := results[0]

		env := parseEnv(res.Config.Env)
		volumes := make([]string, 0, len(res.Mounts))
		for _, m := range res.Mounts {
			if m.Source != "" {
				volumes = append(volumes, m.Source)
			}
		}
		ports := parsePorts(res.HostConfig.PortBindings)
		name := strings.TrimPrefix(res.Name, "/")
		labels := res.Config.Labels
		p := detectPlatform(labels)
		if p != models.TakeoverPlatformDocker {
			platform = p
		}
		composeProject := labels["com.docker.compose.project"]
		if composeProject != "" {
			composeSet[composeProject] = struct{}{}
		}
		containers = append(containers, models.DiscoveredContainer{
			Name:           name,
			Image:          res.Config.Image,
			Ports:          ports,
			Env:            env,
			Volumes:        volumes,
			Labels:         labels,
			Status:         row.State,
			ComposeProject: composeProject,
		})
	}

	if platform == "" {
		platform = models.TakeoverPlatformDocker
	}
	composeProjects := make([]string, 0, len(composeSet))
	for k := range composeSet {
		composeProjects = append(composeProjects, k)
	}

	return &models.DiscoveredStack{
		Containers:      containers,
		ComposeProjects: composeProjects,
		Platform:        platform,
		Host:            req.Host,
	}, nil
}

func dialSSH(ctx context.Context, addr, user, privateKey, fingerprint string) (*ssh.Client, error) {
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	var hostKeyCallback ssh.HostKeyCallback
	if fingerprint != "" {
		hostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			fp := ssh.FingerprintSHA256(key)
			if fp != fingerprint {
				return fmt.Errorf("host key fingerprint mismatch: expected %s, got %s", fingerprint, fp)
			}
			return nil
		}
	} else {
		hostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			log.Printf("WARNING: accepting unknown host key for %s", addr)
			return nil
		}
	}

	cfg := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hostKeyCallback,
		Timeout:         15 * time.Second,
	}
	d := net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	handshakeDone := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		case <-handshakeDone:
		}
	}()

	sshConn, chans, reqs, err := ssh.NewClientConn(conn, addr, cfg)
	close(handshakeDone)
	if err != nil {
		return nil, fmt.Errorf("ssh handshake: %w", err)
	}
	return ssh.NewClient(sshConn, chans, reqs), nil
}

func runSSHCommand(client *ssh.Client, cmd string) (string, error) {
	sess, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()
	out, err := sess.Output(cmd)
	return string(out), err
}

func detectPlatform(labels map[string]string) models.TakeoverPlatform {
	for k := range labels {
		switch {
		case strings.HasPrefix(k, "dokploy."):
			return models.TakeoverPlatformDokploy
		case strings.HasPrefix(k, "coolify."):
			return models.TakeoverPlatformCoolify
		case strings.HasPrefix(k, "com.dokku."):
			return models.TakeoverPlatformDokku
		}
	}
	return models.TakeoverPlatformDocker
}

func parseEnv(envSlice []string) map[string]string {
	out := make(map[string]string, len(envSlice))
	for _, e := range envSlice {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			out[parts[0]] = parts[1]
		}
	}
	return out
}

func parsePorts(bindings map[string]interface{}) []string {
	var ports []string
	for containerPort := range bindings {
		ports = append(ports, containerPort)
	}
	return ports
}
