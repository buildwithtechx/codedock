package backup

import (
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"
	"io"
	"net"
	"path"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func sftpDial(ctx context.Context, dest *models.SFTPDestination) (*sftp.Client, error) {
	if dest.Host == "" || dest.Username == "" {
		return nil, fmt.Errorf("sftp destination needs a host and username")
	}
	port := dest.Port
	if port <= 0 {
		port = 22
	}
	var auths []ssh.AuthMethod
	if dest.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(dest.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("parse sftp private key: %w", err)
		}
		auths = append(auths, ssh.PublicKeys(signer))
	}
	if dest.Password != "" {
		auths = append(auths, ssh.Password(dest.Password))
	}
	if len(auths) == 0 {
		return nil, fmt.Errorf("sftp destination needs a password or private key")
	}
	config := &ssh.ClientConfig{User: dest.Username, Auth: auths, HostKeyCallback: ssh.InsecureIgnoreHostKey(), Timeout: 15 * time.Second}
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", dest.Host, port))
	if err != nil {
		return nil, err
	}
	sshConn, chans, reqs, err := ssh.NewClientConn(conn, fmt.Sprintf("%s:%d", dest.Host, port), config)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	client := ssh.NewClient(sshConn, chans, reqs)
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		_ = client.Close()
		return nil, err
	}
	return sftpClient, nil
}

func VerifySFTP(ctx context.Context, dest *models.SFTPDestination) error {
	if dest == nil {
		return fmt.Errorf("sftp destination not found")
	}
	client, err := sftpDial(ctx, dest)
	if err != nil {
		return err
	}
	return client.Close()
}

func sftpJoin(prefix, name string) string {
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		return name
	}
	return path.Join(prefix, name)
}

func sftpPut(ctx context.Context, dest *models.SFTPDestination, name string, body []byte) (string, error) {
	client, err := sftpDial(ctx, dest)
	if err != nil {
		return "", err
	}
	defer client.Close()
	key := sftpJoin(dest.PathPrefix, name)
	if err := client.MkdirAll(path.Dir(key)); err != nil {
		return "", err
	}
	file, err := client.Create(key)
	if err != nil {
		return "", err
	}
	if _, err := file.Write(body); err != nil {
		_ = file.Close()
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	return "sftp://" + dest.Host + "/" + strings.TrimPrefix(key, "/"), nil
}

func sftpGet(ctx context.Context, dest *models.SFTPDestination, key string) (io.ReadCloser, error) {
	client, err := sftpDial(ctx, dest)
	if err != nil {
		return nil, err
	}
	file, err := client.Open(key)
	if err != nil {
		_ = client.Close()
		return nil, err
	}
	return &sftpReadCloser{file: file, client: client}, nil
}

type sftpReadCloser struct {
	file   *sftp.File
	client *sftp.Client
}

func (r *sftpReadCloser) Read(data []byte) (int, error) { return r.file.Read(data) }
func (r *sftpReadCloser) Close() error {
	fileErr := r.file.Close()
	clientErr := r.client.Close()
	if fileErr != nil {
		return fileErr
	}
	return clientErr
}

func sftpRemove(ctx context.Context, dest *models.SFTPDestination, key string) error {
	client, err := sftpDial(ctx, dest)
	if err != nil {
		return err
	}
	defer client.Close()
	return client.Remove(key)
}
