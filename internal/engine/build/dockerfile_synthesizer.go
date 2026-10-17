package build

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
)

//go:embed templates
var templateFS embed.FS

type dockerfileOverrides struct {
	InstallCommand string
	BuildCommand   string
	StartCommand   string
}

func synthesizeDockerfile(stack string, opts BuildOptions) (string, error) {
	var templateName string
	switch {
	case strings.HasPrefix(stack, "Node.js"):
		templateName = "templates/nodejs.Dockerfile"
	case strings.HasPrefix(stack, "Go"):
		templateName = "templates/go.Dockerfile"
	case strings.HasPrefix(stack, "Python"):
		templateName = "templates/python.Dockerfile"
	case strings.HasPrefix(stack, "Rust"):
		templateName = "templates/rust.Dockerfile"
	default:
		templateName = "templates/static.Dockerfile"
	}
	content, err := templateFS.ReadFile(templateName)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded dockerfile template %s: %w", templateName, err)
	}

	if opts.AppConfig == nil {
		return string(content), nil
	}

	overrides := dockerfileOverrides{
		InstallCommand: opts.AppConfig.InstallCommand,
		BuildCommand:   opts.AppConfig.BuildCommand,
		StartCommand:   opts.AppConfig.StartCommand,
	}

	if overrides.InstallCommand == "" && overrides.BuildCommand == "" && overrides.StartCommand == "" {
		return string(content), nil
	}

	return applyDockerfileOverrides(string(content), overrides), nil
}

func applyDockerfileOverrides(dockerfile string, o dockerfileOverrides) string {
	lines := strings.Split(dockerfile, "\n")

	var out []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if o.InstallCommand != "" && isInstallLine(trimmed) {
			out = append(out, "RUN "+o.InstallCommand)
			continue
		}
		if o.BuildCommand != "" && isBuildLine(trimmed) {
			out = append(out, "RUN "+o.BuildCommand)
			continue
		}
		if o.StartCommand != "" && strings.HasPrefix(trimmed, "CMD ") {
			out = append(out, buildCMDLine(o.StartCommand))
			continue
		}
		out = append(out, line)
	}

	var buf bytes.Buffer
	for _, l := range out {
		buf.WriteString(l)
		buf.WriteByte('\n')
	}
	return buf.String()
}

func isInstallLine(line string) bool {
	installPatterns := []string{
		"RUN npm install", "RUN npm ci", "RUN yarn install", "RUN yarn",
		"RUN pnpm install", "RUN pip install", "RUN pip3 install",
		"RUN bundle install", "RUN composer install", "RUN cargo fetch",
	}
	for _, p := range installPatterns {
		if strings.HasPrefix(line, p) {
			return true
		}
	}
	return false
}

func isBuildLine(line string) bool {
	buildPatterns := []string{
		"RUN npm run build", "RUN yarn build", "RUN pnpm build",
		"RUN go build", "RUN cargo build",
	}
	for _, p := range buildPatterns {
		if strings.HasPrefix(line, p) {
			return true
		}
	}
	return false
}

var stackVersionEnv = []struct {
	prefix   string
	language string
}{
	{"Node.js", "NODE"},
	{"Python", "PYTHON"},
	{"Go", "GO"},
	{"Rust", "RUST"},
	{"PHP", "PHP"},
	{"Ruby", "RUBY"},
}

func appRuntimeVersion(opts BuildOptions) string {
	if opts.AppConfig == nil {
		return ""
	}
	return strings.TrimSpace(opts.AppConfig.RuntimeVersion)
}

func versionBuildEnv(stack, version string) map[string]string {
	if version == "" {
		return nil
	}
	for _, entry := range stackVersionEnv {
		if strings.HasPrefix(stack, entry.prefix) {
			return map[string]string{
				"NIXPACKS_" + entry.language + "_VERSION": version,
				"RAILPACK_" + entry.language + "_VERSION": version,
			}
		}
	}
	return nil
}

func buildCMDLine(startCmd string) string {
	parts := strings.Fields(startCmd)
	if len(parts) == 0 {
		return ""
	}
	quoted := make([]string, len(parts))
	for i, p := range parts {
		quoted[i] = `"` + p + `"`
	}
	return "CMD [" + strings.Join(quoted, ", ") + "]"
}
