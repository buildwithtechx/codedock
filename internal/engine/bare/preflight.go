package bare

import (
	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/models"
	"context"
)

func (r *Runtime) Preflight(ctx context.Context, app *models.AppService, target models.RuntimeTarget) error {
	base, unit, _, err := identity(app)
	if err != nil {
		return err
	}
	script := "set -eu\ntest \"$(uname -s)\" = Linux\ncommand -v systemctl >/dev/null\ncommand -v curl >/dev/null\ncommand -v sha256sum >/dev/null\ncommand -v flock >/dev/null\ncommand -v useradd >/dev/null\ncommand -v runuser >/dev/null\nbase=" + ssh.ShellQuote(base) + "\ntest ! -L /opt/codedock\ntest ! -L /opt/codedock/services\ntest ! -L \"$base\"\nif test -e \"$base\"; then test \"$(cat \"$base/owner\")\" = " + ssh.ShellQuote(app.ProjectID) + "; test ! -e \"$base/pending\"; else test ! -e " + ssh.ShellQuote("/etc/systemd/system/"+unit+".service") + "; fi\n"
	_, err = r.runner.Host(ctx, target.BareNode, script)
	return err
}
