package bare

import (
	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/models"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"strings"
)

func identity(app *models.AppService) (string, string, string, error) {
	if _, err := uuid.Parse(app.ID); err != nil {
		return "", "", "", err
	}
	if _, err := uuid.Parse(app.ProjectID); err != nil {
		return "", "", "", err
	}
	return "/opt/codedock/services/" + app.ID, "codedock-app-" + app.ID, "cd-" + strings.ReplaceAll(app.ID, "-", "")[:20], nil
}
func ownedScript(app *models.AppService) (string, error) {
	base, unit, user, err := identity(app)
	if err != nil {
		return "", err
	}
	return "set -eu\numask 077\nbase=" + ssh.ShellQuote(base) + "\nunit=" + ssh.ShellQuote(unit) + "\nuser=" + ssh.ShellQuote(user) + "\ntest ! -L \"$base\"\ntest \"$(cat \"$base/owner\")\" = " + ssh.ShellQuote(app.ProjectID) + "\nexec 9>\"$base/lock\"\nflock -n 9\n", nil
}
func DeploymentScript(app *models.AppService, target models.RuntimeTarget, variables map[string]string, release string) (string, error) {
	if err := Validate(app, target); err != nil {
		return "", err
	}
	if _, err := uuid.Parse(release); err != nil {
		return "", err
	}
	base, unit, user, err := identity(app)
	if err != nil {
		return "", err
	}
	launch := "#!/bin/sh\nset -eu\ncd " + ssh.ShellQuote(base+"/current") + "\nexport CODEDOCK_DATA_DIR=" + ssh.ShellQuote(base+"/data") + "\n"
	for key, value := range variables {
		if !regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(key) || strings.ContainsRune(value, 0) {
			return "", fmt.Errorf("invalid native environment variable")
		}
		launch += "export " + key + "=" + ssh.ShellQuote(value) + "\n"
	}
	args := []string{}
	if target.BareToolchain == "static" {
		staticDir := target.BareStaticDir
		if staticDir == "" {
			staticDir = "app"
		}
		launch += "exec python3 -m http.server " + fmt.Sprint(app.InternalPort) + " --directory " + ssh.ShellQuote(staticDir) + "\n"
	} else {
		for _, arg := range target.BareCommand {
			args = append(args, ssh.ShellQuote(arg))
		}
		launch += "exec " + strings.Join(args, " ") + "\n"
	}
	config := "[Unit]\nDescription=Codedock application\nAfter=network-online.target\nStartLimitIntervalSec=60\nStartLimitBurst=3\n[Service]\nType=simple\nUser=" + user + "\nGroup=" + user + "\nWorkingDirectory=" + base + "/current\nExecStart=/bin/sh " + base + "/current/launch.sh\nRestart=on-failure\nRestartSec=3\nNoNewPrivileges=true\nPrivateTmp=true\nProtectSystem=strict\nProtectHome=true\nReadWritePaths=" + base + "/data\n[Install]\nWantedBy=multi-user.target\n"
	if app.MemoryLimit > 0 {
		config += fmt.Sprintf("[Service]\nMemoryMax=%dM\n", app.MemoryLimit)
	}
	if app.CPULimit > 0 {
		config += fmt.Sprintf("[Service]\nCPUQuota=%.2f%%\n", app.CPULimit*100)
	}
	script := "set -eu\numask 077\ncommand -v systemctl >/dev/null\ncommand -v curl >/dev/null\ncommand -v sha256sum >/dev/null\ncommand -v flock >/dev/null\nbase=" + ssh.ShellQuote(base) + "\nunit=" + ssh.ShellQuote(unit) + "\nuser=" + ssh.ShellQuote(user) + "\nrelease=" + ssh.ShellQuote(release) + "\n"
	script += `test ! -L /opt/codedock
test ! -L /opt/codedock/services
test ! -L "$base"
if test -e "$base"; then test "$(cat "$base/owner")" = ` + ssh.ShellQuote(app.ProjectID) + `; else
 test ! -e "/etc/systemd/system/$unit.service"
 mkdir -p "$base"
 printf %s ` + ssh.ShellQuote(app.ProjectID) + ` > "$base/owner"
fi
exec 9>"$base/lock"
flock -n 9
test ! -e "$base/pending"
if id "$user" >/dev/null 2>&1; then test "$(getent passwd "$user" | cut -d: -f6)" = "$base"; else useradd --system --home-dir "$base" --shell /usr/sbin/nologin "$user"; fi
chown root:"$user" "$base"
chmod 750 "$base"
test ! -L "$base/releases"
test ! -L "$base/data"
mkdir -p "$base/releases/$release" "$base/data"
chown root:"$user" "$base/releases"
chown "$user":"$user" "$base/data"
chmod 750 "$base/releases" "$base/data"
`
	if target.BareRepoURL != "" {
		script += sourceBuildStage(target)
	} else {
		script += `curl --proto '=https' --proto-redir '=https' --fail --location --max-time 180 --max-filesize 536870912 --output "$base/releases/$release/app" -- ` + ssh.ShellQuote(target.BareReleaseURL) + `
printf '%s  %s\n' ` + ssh.ShellQuote(strings.ToLower(target.BareSHA256)) + ` "$base/releases/$release/app" | sha256sum -c - >/dev/null
`
	}
	script += `printf %s ` + ssh.ShellQuote(launch) + ` > "$base/releases/$release/launch.sh"
printf %s ` + ssh.ShellQuote(config) + ` > "$base/releases/$release/unit"
chown -R root:"$user" "$base/releases/$release"
chmod 750 "$base/releases/$release"
chmod 550 "$base/releases/$release/app" "$base/releases/$release/launch.sh"
chmod 600 "$base/releases/$release/unit"
old=-
if test -L "$base/current"; then old=$(basename "$(readlink "$base/current")"); test "$(readlink "$base/current")" = "$base/releases/$old"; test -f "$base/releases/$old/unit"; fi
printf '%s\n%s\n' "$old" "$release" > "$base/pending"
ln -s "$base/releases/$release" "$base/next"
mv -Tf "$base/next" "$base/current"
install -m 600 "$base/releases/$release/unit" "/etc/systemd/system/$unit.service"
systemctl daemon-reload
systemctl reset-failed "$unit" || true
systemctl enable "$unit" >/dev/null
systemctl restart "$unit"
`
	script += routingStage(app)
	script += readinessScript(app)
	return script, nil
}
func readinessScript(app *models.AppService) string {
	command := "systemctl is-active --quiet \"$unit\""
	if app.RuntimeMode != models.RuntimeModeWorker {
		path := app.HealthCheckPath
		if path == "" {
			path = "/"
		}
		command += " && curl --fail --silent --max-time 2 -- " + ssh.ShellQuote(fmt.Sprintf("http://127.0.0.1:%d%s", app.InternalPort, path)) + " >/dev/null"
	}
	return "ready=0\nfor attempt in $(seq 1 60); do if " + command + "; then ready=$((ready+1)); else ready=0; fi; if test \"$ready\" -ge 5; then exit 0; fi; sleep 1; done\nexit 1\n"
}
