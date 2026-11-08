package ssh

import (
	"context"
	"fmt"
	"io"
)

func (c *Client) Stream(ctx context.Context, command string, input io.Reader, output io.Writer) error {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("create SSH session: %w", err)
	}
	defer session.Close()
	stop := context.AfterFunc(ctx, func() { session.Close() })
	defer stop()
	if input != nil {
		stdin, err := session.StdinPipe()
		if err != nil {
			return fmt.Errorf("open SSH stdin: %w", err)
		}
		defer stdin.Close()
		go func() {
			_, _ = io.Copy(stdin, input)
			_ = stdin.Close()
		}()
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("open SSH stdout: %w", err)
	}
	if err := session.Start(command); err != nil {
		return fmt.Errorf("start remote stream: %w", err)
	}
	_, copyErr := io.Copy(output, stdout)
	waitErr := session.Wait()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if copyErr != nil {
		return copyErr
	}
	return waitErr
}

func (c *Client) StreamPTY(ctx context.Context, command string, input io.Reader, output io.Writer) error {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("create SSH session: %w", err)
	}
	defer session.Close()
	stop := context.AfterFunc(ctx, func() { session.Close() })
	defer stop()
	if err := session.RequestPty("xterm-256color", 80, 40, nil); err != nil {
		return fmt.Errorf("allocate remote terminal: %w", err)
	}
	if input != nil {
		stdin, err := session.StdinPipe()
		if err != nil {
			return fmt.Errorf("open SSH stdin: %w", err)
		}
		defer stdin.Close()
		go func() {
			_, _ = io.Copy(stdin, input)
			_ = stdin.Close()
		}()
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("open SSH stdout: %w", err)
	}
	if err := session.Start(command); err != nil {
		return fmt.Errorf("start remote terminal: %w", err)
	}
	_, copyErr := io.Copy(output, stdout)
	waitErr := session.Wait()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if copyErr != nil {
		return copyErr
	}
	return waitErr
}
