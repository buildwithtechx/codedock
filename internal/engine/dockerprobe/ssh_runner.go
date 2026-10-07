package dockerprobe

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Host        string
	Port        int
	User        string
	PrivateKey  string
	Password    string
	Fingerprint string
}

type SSHRunner struct {
	config SSHConfig
}

func NewSSHRunner(config SSHConfig) *SSHRunner {
	return &SSHRunner{config: config}
}

func FingerprintFor(host string, port int, timeout time.Duration) (string, error) {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()
	var fingerprint string
	callback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		fingerprint = ssh.FingerprintSHA256(key)
		return fmt.Errorf("capture only")
	}
	config := &ssh.ClientConfig{User: "probe", HostKeyCallback: callback, Timeout: timeout}
	_, _, _, _ = ssh.NewClientConn(conn, addr, config)
	if fingerprint == "" {
		return "", fmt.Errorf("no host key presented")
	}
	return fingerprint, nil
}

func (r *SSHRunner) dial(ctx context.Context) (*ssh.Client, error) {
	if r.config.Fingerprint == "" {
		return nil, fmt.Errorf("ssh host fingerprint is required")
	}
	var methods []ssh.AuthMethod
	if r.config.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(r.config.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("parse private key: %w", err)
		}
		methods = append(methods, ssh.PublicKeys(signer))
	}
	if r.config.Password != "" {
		methods = append(methods, ssh.Password(r.config.Password))
	}
	if len(methods) == 0 {
		return nil, fmt.Errorf("ssh key or password is required")
	}
	expected := r.config.Fingerprint
	config := &ssh.ClientConfig{
		User:            r.config.User,
		Auth:            methods,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			actual := ssh.FingerprintSHA256(key)
			if actual != expected {
				return fmt.Errorf("host key fingerprint mismatch: expected %s, got %s", expected, actual)
			}
			return nil
		},
		Timeout: 15 * time.Second,
	}
	port := r.config.Port
	if port <= 0 {
		port = 22
	}
	addr := net.JoinHostPort(r.config.Host, fmt.Sprintf("%d", port))
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
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
	sshConn, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	close(handshakeDone)
	if err != nil {
		return nil, fmt.Errorf("ssh handshake: %w", err)
	}
	return ssh.NewClient(sshConn, chans, reqs), nil
}

func (r *SSHRunner) Run(ctx context.Context, cmd string) (string, error) {
	var stdout bytes.Buffer
	if err := r.RunPipe(ctx, cmd, nil, &stdout); err != nil {
		return stdout.String(), err
	}
	return stdout.String(), nil
}

func (r *SSHRunner) RunPipe(ctx context.Context, cmd string, stdin io.Reader, stdout io.Writer) error {
	client, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	var stderr bytes.Buffer
	session.Stderr = &stderr
	if stdin != nil {
		session.Stdin = stdin
	}
	if stdout != nil {
		session.Stdout = stdout
	}
	if err := session.Start(cmd); err != nil {
		return fmt.Errorf("start remote command: %w", err)
	}
	waitDone := make(chan error, 1)
	go func() { waitDone <- session.Wait() }()
	select {
	case <-ctx.Done():
		_ = session.Signal(ssh.SIGKILL)
		<-waitDone
		return ctx.Err()
	case err := <-waitDone:
		if err != nil && stderr.Len() > 0 {
			return fmt.Errorf("%w: %s", err, stderr.String())
		}
		return err
	}
}

type LocalRunner struct{}

func (LocalRunner) Run(ctx context.Context, cmd string) (string, error) {
	var stdout bytes.Buffer
	if err := (LocalRunner{}).RunPipe(ctx, cmd, nil, &stdout); err != nil {
		return stdout.String(), err
	}
	return stdout.String(), nil
}

func (LocalRunner) RunPipe(ctx context.Context, cmd string, stdin io.Reader, stdout io.Writer) error {
	command := exec.CommandContext(ctx, "sh", "-c", cmd)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if stdin != nil {
		command.Stdin = stdin
	}
	if stdout != nil {
		command.Stdout = stdout
	}
	if err := command.Run(); err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("%w: %s", err, stderr.String())
		}
		return err
	}
	return nil
}
