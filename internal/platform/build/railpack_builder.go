package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"

	"codedock/internal/config"
)

type RailpackBuilder struct {
	dockerClient *client.Client
}

func NewRailpackBuilder(dockerClient *client.Client) *RailpackBuilder {
	return &RailpackBuilder{dockerClient: dockerClient}
}

func (r *RailpackBuilder) Build(ctx context.Context, opts BuildOptions, engineName string) (string, error) {
	imageTag := fmt.Sprintf("codedock-app-%s:latest", strings.ToLower(opts.ProjectID))
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "🌟 [Railpack/Nixpacks] Auto-detecting language & framework in %s...\n", opts.SourceDir)
	}
	stack := r.detectLanguageStack(opts.SourceDir)
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "🛠️ [Railpack/Nixpacks] Stack detected: %s\n", stack)
	}
	versionEnv := versionBuildEnv(stack, appRuntimeVersion(opts))
	if opts.LogWriter != nil && len(versionEnv) > 0 {
		fmt.Fprintf(opts.LogWriter, "🧩 [Railpack/Nixpacks] Pinning runtime version: %s\n", appRuntimeVersion(opts))
	}
	absSourceDir, err := filepath.Abs(opts.SourceDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute source dir: %w", err)
	}

	dockerRunArgs := r.dockerRunArgs()

	if engineName == "buildpacks" {
		if opts.LogWriter != nil {
			fmt.Fprintf(opts.LogWriter, "⚙️ [Buildpacks] Executing pack builder engine via Docker container...\n")
		}
		packImage := config.Get().Builder.PackImage
		if packImage == "" {
			packImage = "buildpacksio/pack:latest"
		}
		builderImage := config.Get().Builder.BuilderImage
		if builderImage == "" {
			builderImage = "paketobuildpacks/builder:base"
		}
		args := append([]string{}, dockerRunArgs...)
		args = append(args,
			"-v", fmt.Sprintf("%s:/app", absSourceDir),
			packImage, "build", imageTag, "--path", "/app", "--builder", builderImage)
		for k, v := range opts.EnvVars {
			args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
		}
		for k, v := range versionEnv {
			args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
		}
		cmd := exec.CommandContext(ctx, "docker", args...)
		cmd.Stdout = opts.LogWriter
		cmd.Stderr = opts.LogWriter
		if err := cmd.Run(); err == nil {
			return imageTag, nil
		} else {
			if opts.LogWriter != nil {
				fmt.Fprintf(opts.LogWriter, "⚠️ [Buildpacks] Container execution failed, falling back to synthesized Dockerfile. Error: %v\n", err)
			}
		}
	} else {
		if opts.LogWriter != nil {
			fmt.Fprintf(opts.LogWriter, "⚙️ [Nixpacks] Executing Nixpacks builder engine via Docker container...\n")
		}
		nixpacksImage := config.Get().Builder.NixpacksImage
		if nixpacksImage == "" {
			nixpacksImage = "ghcr.io/railwayapp/nixpacks:latest"
		}
		args := append([]string{}, dockerRunArgs...)
		args = append(args,
			"-v", fmt.Sprintf("%s:/app", absSourceDir),
			nixpacksImage, "build", "/app", "--name", imageTag)
		for k, v := range opts.EnvVars {
			args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
		}
		for k, v := range versionEnv {
			args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
		}
		if opts.AppConfig != nil {
			if opts.AppConfig.InstallCommand != "" {
				args = append(args, "--install-cmd", opts.AppConfig.InstallCommand)
			}
			if opts.AppConfig.BuildCommand != "" {
				args = append(args, "--build-cmd", opts.AppConfig.BuildCommand)
			}
			if opts.AppConfig.StartCommand != "" {
				args = append(args, "--start-cmd", opts.AppConfig.StartCommand)
			}
		}
		cmd := exec.CommandContext(ctx, "docker", args...)
		cmd.Stdout = opts.LogWriter
		cmd.Stderr = opts.LogWriter
		if err := cmd.Run(); err == nil {
			return imageTag, nil
		} else {
			if opts.LogWriter != nil {
				fmt.Fprintf(opts.LogWriter, "⚠️ [Nixpacks] Container execution failed, falling back to synthesized Dockerfile. Error: %v\n", err)
			}
		}
	}
	if opts.LogWriter != nil {
		fmt.Fprintf(opts.LogWriter, "⚠️ [Railpack/Nixpacks] Native CLI not found; synthesizing zero-config OCI build plan for %s...\n", stack)
	}
	synthesizedDockerfile, err := synthesizeDockerfile(stack, opts)
	if err != nil {
		return "", fmt.Errorf("failed to synthesize fallback dockerfile: %w", err)
	}
	dockerfilePath := filepath.Join(opts.SourceDir, ".codedock.Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(synthesizedDockerfile), 0o644); err != nil {
		return "", fmt.Errorf("failed to write synthesized dockerfile: %w", err)
	}
	defer os.Remove(dockerfilePath)
	fallbackOpts := opts
	fallbackOpts.DockerfilePath = ".codedock.Dockerfile"
	fallbackBuilder := NewDockerfileBuilder(r.dockerClient)
	return fallbackBuilder.Build(ctx, fallbackOpts)
}

func (r *RailpackBuilder) dockerRunArgs() []string {
	args := []string{"run", "--rm"}
	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		dockerSocket := config.Get().Docker.SocketPath
		if dockerSocket == "" {
			dockerSocket = "/var/run/docker.sock"
		}
		return append(args, "-v", fmt.Sprintf("%s:/var/run/docker.sock", dockerSocket))
	}

	network := config.Get().Docker.RuntimeNetwork
	if network == "" {
		network = "codedock-network"
	}
	return append(args, "--network", network, "-e", "DOCKER_HOST="+dockerHost)
}

func (r *RailpackBuilder) detectLanguageStack(sourceDir string) string {
	if _, err := os.Stat(filepath.Join(sourceDir, "package.json")); err == nil {
		return "Node.js (npm/yarn/pnpm)"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "go.mod")); err == nil {
		return "Go (Golang)"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "requirements.txt")); err == nil {
		return "Python (pip)"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "pyproject.toml")); err == nil {
		return "Python (poetry/pyproject)"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "Cargo.toml")); err == nil {
		return "Rust (cargo)"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "composer.json")); err == nil {
		return "PHP (composer)"
	}
	if _, err := os.Stat(filepath.Join(sourceDir, "Gemfile")); err == nil {
		return "Ruby (bundler)"
	}
	return "Static / Universal HTML"
}
