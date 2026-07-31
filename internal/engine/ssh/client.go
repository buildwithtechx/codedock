package ssh

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	dockerclient "github.com/docker/docker/client"
	"golang.org/x/crypto/ssh"

	"codedock.run/codedock/internal/models"
)

type Config struct {
	Host        string
	Port        int
	User        string
	Key         string
	Password    string
	Fingerprint string
}

type Client struct {
	cfg       Config
	sshClient *ssh.Client
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.Port <= 0 {
		cfg.Port = 22
	}
	if cfg.User == "" {
		cfg.User = "root"
	}

	var authMethods []ssh.AuthMethod
	if cfg.Key != "" {
		signer, err := ssh.ParsePrivateKey([]byte(cfg.Key))
		if err != nil {
			return nil, fmt.Errorf("failed to parse ssh key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}
	if cfg.Password != "" {
		authMethods = append(authMethods, ssh.Password(cfg.Password))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("neither ssh key nor password provided")
	}

	sshConfig := &ssh.ClientConfig{
		User: cfg.User,
		Auth: authMethods,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			if cfg.Fingerprint == "" {
				return nil
			}
			fingerprint := ssh.FingerprintSHA256(key)
			if fingerprint != cfg.Fingerprint {
				return fmt.Errorf("host key fingerprint mismatch: expected %q, got %q", cfg.Fingerprint, fingerprint)
			}
			return nil
		},
		Timeout: 10 * time.Second,
	}

	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh dial failed: %w", err)
	}

	return &Client{
		cfg:       cfg,
		sshClient: sshClient,
	}, nil
}

func (c *Client) Close() error {
	if c.sshClient != nil {
		return c.sshClient.Close()
	}
	return nil
}

func (c *Client) RunCommand(ctx context.Context, cmd string) (string, error) {
	sess, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create ssh session: %w", err)
	}

	errCh := make(chan error, 1)
	outCh := make(chan string, 1)

	go func() {
		defer sess.Close()
		out, err := sess.CombinedOutput(cmd)
		if err != nil {
			errCh <- fmt.Errorf("ssh command execution failed: %w: %s", err, string(out))
			return
		}
		outCh <- string(out)
	}()

	select {
	case <-ctx.Done():
		sess.Close()
		return "", ctx.Err()
	case err := <-errCh:
		return "", err
	case out := <-outCh:
		return out, nil
	}
}

func (c *Client) DockerClient(ctx context.Context) (*dockerclient.Client, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return c.sshClient.Dial("unix", "/var/run/docker.sock")
			},
		},
	}

	return dockerclient.NewClientWithOpts(
		dockerclient.WithHTTPClient(httpClient),
		dockerclient.WithHost("http://docker"),
		dockerclient.WithAPIVersionNegotiation(),
	)
}

func (c *Client) CollectMetrics(ctx context.Context) ([]byte, error) {
	cmd := `sh -c '
cpu=$(top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk "{print 100 - \$1}")
mem=$(free -m | awk "/Mem:/ {print \$3, \$2}")
disk=$(df -m / | awk "NR==2 {print \$3, \$2}")
echo "{\"cpu\":${cpu:-0},\"mem\":\"${mem}\",\"disk\":\"${disk}\"}"
'`
	out, err := c.RunCommand(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to collect metrics over ssh: %w", err)
	}

	out = strings.TrimSpace(out)
	var raw struct {
		CPU  float64 `json:"cpu"`
		Mem  string  `json:"mem"`
		Disk string  `json:"disk"`
	}

	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse metrics JSON: %w", err)
	}

	var memUsed, memTotal uint64
	if n, err := fmt.Sscanf(raw.Mem, "%d %d", &memUsed, &memTotal); err != nil || n != 2 {
		return nil, fmt.Errorf("failed to parse memory metrics: %w (scanned %d)", err, n)
	}

	var diskUsed, diskTotal uint64
	if n, err := fmt.Sscanf(raw.Disk, "%d %d", &diskUsed, &diskTotal); err != nil || n != 2 {
		return nil, fmt.Errorf("failed to parse disk metrics: %w (scanned %d)", err, n)
	}

	payload := models.WorkerMetricsPayload{
		CPUUsagePercentage: raw.CPU,
		MemoryUsageBytes:   memUsed * 1024 * 1024,
		MemoryLimitBytes:   memTotal * 1024 * 1024,
		DiskUsageBytes:     diskUsed * 1024 * 1024,
		DiskTotalBytes:     diskTotal * 1024 * 1024,
	}

	return json.Marshal(payload)
}
