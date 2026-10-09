package bare

import (
	"codedock/internal/engine/ssh"
	"codedock/internal/models"
	"context"
	"strings"
)

func recoveryScript(app *models.AppService) (string, error) {
	base, unit, _, err := identity(app)
	if err != nil {
		return "", err
	}
	script, err := ownedScript(app)
	if err != nil {
		return "", err
	}
	guard := "if ! test -e " + ssh.ShellQuote(base) + "; then test ! -e " + ssh.ShellQuote("/etc/systemd/system/"+unit+".service") + "; exit 0; fi\n"
	return guard + script + `if ! test -f "$base/pending"; then exit 0; fi
old=$(sed -n '1p' "$base/pending")
if test "$old" = -; then
 systemctl stop "$unit" || true
 systemctl disable "$unit" || true
 rm -f "/etc/systemd/system/$unit.service" "$base/current"
 systemctl daemon-reload
else
 printf %s "$old" | grep -E '^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$' >/dev/null
 test -f "$base/releases/$old/unit"
 test ! -L "$base/releases/$old"
 rm -f "$base/next"
 ln -s "$base/releases/$old" "$base/next"
 mv -Tf "$base/next" "$base/current"
 install -m 600 "$base/releases/$old/unit" "/etc/systemd/system/$unit.service"
 systemctl daemon-reload
 systemctl reset-failed "$unit" || true
 systemctl restart "$unit"
 systemctl is-active --quiet "$unit"
fi
rm -f "$base/pending"
`, nil
}

func (r *Runtime) CurrentRelease(ctx context.Context, app *models.AppService, target models.RuntimeTarget) (string, error) {
	base, _, _, err := identity(app)
	if err != nil {
		return "", err
	}
	script, err := ownedScript(app)
	if err != nil {
		return "", err
	}
	guard := "if ! test -e " + ssh.ShellQuote(base) + "; then exit 0; fi\n"
	raw, err := r.runner.Host(ctx, target.BareNode, guard+script+"if test -L \"$base/current\"; then basename \"$(readlink \"$base/current\")\"; fi\n")
	return strings.TrimSpace(raw), err
}
