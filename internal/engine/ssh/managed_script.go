package ssh

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

func ShellQuote(value string) string { return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'" }
func (c *Client) RootCommand(command string) string {
	if c.cfg.User == "root" {
		return command
	}
	return "sudo -n " + command
}
func (c *Client) RunManagedScript(ctx context.Context, id, script string) (output string, resultErr error) {
	if _, err := uuid.Parse(id); err != nil {
		return "", fmt.Errorf("managed operation ID must be a UUID")
	}
	directory := "/run/codedock-operations/" + id
	wrapper := `set -eu
dir=` + ShellQuote(directory) + `
install -d -m 700 "$dir"
cat > "$dir/run"
chmod 700 "$dir/run"
child=''
trap 'if test -n "$child"; then kill -TERM -- "-$child" 2>/dev/null || true; fi' HUP INT TERM
setsid sh "$dir/run" &
child=$!
printf '%s' "$child" > "$dir/pid"
set +e
wait "$child"
code=$?
rm -f "$dir/pid" "$dir/run"
rmdir "$dir"
exit "$code"`
	finished := make(chan struct{})
	var cancelErr error
	stop := context.AfterFunc(ctx, func() {
		defer close(finished)
		cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		command := `dir=` + ShellQuote(directory) + `; for i in 1 2 3 4 5 6 7 8 9 10; do
if test -f "$dir/pid"; then
pid=$(cat "$dir/pid")
case "$pid" in ''|*[!0-9]*) exit 1;; esac
if test -r "/proc/$pid/cmdline" && tr '\0' '\n' < "/proc/$pid/cmdline" | grep -F -x "$dir/run" >/dev/null; then kill -TERM -- "-$pid"; fi
exit 0
fi
sleep 0.2
done`
		_, cancelErr = c.RunWithInput(cleanup, c.RootCommand("sh -c "+ShellQuote(command)), strings.NewReader(""))
	})
	defer func() {
		if !stop() {
			<-finished
			resultErr = errors.Join(resultErr, cancelErr)
		}
	}()
	return c.RunWithInput(ctx, c.RootCommand("sh -c "+ShellQuote(wrapper)), strings.NewReader(script))
}
