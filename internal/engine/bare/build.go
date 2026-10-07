package bare

import (
	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/models"
	"fmt"
)

func toolchainScript(target models.RuntimeTarget) string {
	ensure := ""
	switch target.BareToolchain {
	case "go":
		ensure = "command -v go >/dev/null || { command -v apt-get >/dev/null && apt-get update -qq && apt-get install -y -qq golang-go; } || exit 1\ncommand -v go >/dev/null\n"
	case "node":
		ensure = "command -v node >/dev/null && command -v npm >/dev/null || { command -v apt-get >/dev/null && apt-get update -qq && apt-get install -y -qq nodejs npm; } || exit 1\ncommand -v node >/dev/null\ncommand -v npm >/dev/null\n"
	case "python":
		ensure = "command -v python3 >/dev/null && command -v pip3 >/dev/null || { command -v apt-get >/dev/null && apt-get update -qq && apt-get install -y -qq python3 python3-pip python3-venv; } || exit 1\n"
	case "static":
		ensure = "command -v python3 >/dev/null || { command -v apt-get >/dev/null && apt-get update -qq && apt-get install -y -qq python3; } || exit 1\n"
	}
	return ensure
}

func sourceBuildStage(target models.RuntimeTarget) string {
	branch := target.BareBranch
	if branch == "" {
		branch = "main"
	}
	output := target.BareOutput
	if output == "" {
		output = "app"
	}
	stage := "command -v git >/dev/null\nrm -rf \"$base/source\"\ngit clone --depth 1 --branch " + ssh.ShellQuote(branch) + " " + ssh.ShellQuote(target.BareRepoURL) + " \"$base/source\"\n"
	stage += toolchainScript(target)
	stage += "cd \"$base/source\"\n"
	if target.BareInstallCommand != "" {
		stage += "sh -ec " + ssh.ShellQuote(target.BareInstallCommand) + "\n"
	}
	if target.BareBuildCommand != "" {
		stage += "sh -ec " + ssh.ShellQuote(target.BareBuildCommand) + "\n"
	} else {
		switch target.BareToolchain {
		case "go":
			stage += "go build -o app .\n"
		case "node":
			stage += "test -f package.json\nnpm run build\n"
		case "python":
			stage += "python3 -m compileall -q .\n"
		}
	}
	stage += "test -e " + ssh.ShellQuote(output) + "\nmkdir -p \"$base/releases/$release\"\ncp -r " + ssh.ShellQuote(output) + " \"$base/releases/$release/app\"\n"
	if target.BareToolchain == "static" {
		stage += "test -d \"$base/releases/$release/app\"\n"
	} else {
		stage += "test -x \"$base/releases/$release/app\" || chmod +x \"$base/releases/$release/app\"\ntest -x \"$base/releases/$release/app\"\n"
	}
	stage += "cd \"$base\"\nrm -rf \"$base/source\"\n"
	return stage
}

func routingFile(app *models.AppService) string {
	return "/etc/codedock/traefik-dynamic/codedock-app-" + app.ID + ".yaml"
}

func routingStage(app *models.AppService) string {
	if app.Domain == "" || app.RuntimeMode == models.RuntimeModeWorker {
		return "rm -f " + ssh.ShellQuote(routingFile(app)) + "\n"
	}
	document := fmt.Sprintf("http:\n  routers:\n    codedock-app-%s:\n      rule: Host(`%s`)\n      entryPoints: [web, websecure]\n      service: codedock-app-%s\n  services:\n    codedock-app-%s:\n      loadBalancer:\n        servers:\n        - url: 'http://host.docker.internal:%d'\n", app.ID, app.Domain, app.ID, app.ID, app.InternalPort)
	return "mkdir -p /etc/codedock/traefik-dynamic\nprintf %s " + ssh.ShellQuote(document) + " > " + ssh.ShellQuote(routingFile(app)) + "\nchmod 600 " + ssh.ShellQuote(routingFile(app)) + "\n"
}
