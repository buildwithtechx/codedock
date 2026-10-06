package ssh

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
)

func (c *Client) RunWithInput(ctx context.Context, command string, input io.Reader) (string, error) {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("create SSH session: %w", err)
	}
	defer session.Close()
	stop := context.AfterFunc(ctx, func() { session.Close() })
	defer stop()
	var output boundedOutput
	session.Stdin = input
	session.Stdout = &output
	session.Stderr = &output
	if err := session.Run(command); err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", fmt.Errorf("remote command failed: %w", err)
	}
	return output.String(), nil
}

type boundedOutput struct {
	bytes.Buffer
	mu sync.Mutex
}

func (b *boundedOutput) Write(value []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.Len()+len(value) > 2*1024*1024 {
		return 0, fmt.Errorf("remote output exceeds limit")
	}
	return b.Buffer.Write(value)
}
