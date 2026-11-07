package deployments

import (
	"testing"

	"codedock/internal/models"
)

func TestResolveProjectChangeImpact(t *testing.T) {
	svcWeb := &models.AppService{
		ID:            "svc-web",
		Name:          "web",
		RootDirectory: "apps/web",
	}

	svcDocs := &models.AppService{
		ID:            "svc-docs",
		Name:          "docs",
		RootDirectory: "apps/docs",
	}

	svcRoot := &models.AppService{
		ID:            "svc-root",
		Name:          "root-service",
		RootDirectory: "",
	}

	svcDocker := &models.AppService{
		ID:            "svc-docker",
		Name:          "docker-service",
		RootDirectory: "apps/worker",
		RuntimeMode:   "docker",
	}

	t.Run("skips when changes are in an unrelated workspace", func(t *testing.T) {
		files := []string{"apps/docs/src/index.ts", "apps/docs/package.json"}
		impact := ResolveProjectChangeImpact(svcWeb, files, "Update docs")
		if impact.Affected {
			t.Errorf("expected web service to be skipped for docs-only changes, got: %+v", impact)
		}
		if impact.Reason != "unrelated-workspaces" {
			t.Errorf("expected reason unrelated-workspaces, got: %s", impact.Reason)
		}
	})

	t.Run("deploys when changes are in target service directory", func(t *testing.T) {
		files := []string{"apps/web/src/App.tsx"}
		impact := ResolveProjectChangeImpact(svcWeb, files, "Fix web app")
		if !impact.Affected {
			t.Errorf("expected web service to be deployed, got: %+v", impact)
		}
		if impact.Reason != "project-files" {
			t.Errorf("expected reason project-files, got: %s", impact.Reason)
		}
	})

	t.Run("deploys when global workspace manifest changes", func(t *testing.T) {
		files := []string{"pnpm-workspace.yaml", "packages/shared/util.ts"}
		impact := ResolveProjectChangeImpact(svcWeb, files, "Add workspace")
		if !impact.Affected {
			t.Errorf("expected web service to deploy on global manifest change, got: %+v", impact)
		}
		if impact.Reason != "global-manifest" {
			t.Errorf("expected reason global-manifest, got: %s", impact.Reason)
		}
	})

	t.Run("force-deploy token overrides skipping", func(t *testing.T) {
		files := []string{"apps/docs/README.md"}
		impact := ResolveProjectChangeImpact(svcWeb, files, "Release [force-deploy]")
		if !impact.Affected {
			t.Errorf("expected force-deploy token to trigger deploy, got: %+v", impact)
		}
		if impact.Reason != "force-deploy-token" {
			t.Errorf("expected reason force-deploy-token, got: %s", impact.Reason)
		}
	})

	t.Run("repository root service always deploys", func(t *testing.T) {
		files := []string{"apps/docs/README.md"}
		impact := ResolveProjectChangeImpact(svcRoot, files, "Docs update")
		if !impact.Affected {
			t.Errorf("expected root service to always deploy, got: %+v", impact)
		}
		if impact.Reason != "repository-root" {
			t.Errorf("expected reason repository-root, got: %s", impact.Reason)
		}
	})

	t.Run("custom docker build context always deploys", func(t *testing.T) {
		files := []string{"apps/docs/README.md"}
		impact := ResolveProjectChangeImpact(svcDocker, files, "Docs update")
		if !impact.Affected {
			t.Errorf("expected docker service to always deploy, got: %+v", impact)
		}
		if impact.Reason != "custom-build-context" {
			t.Errorf("expected reason custom-build-context, got: %s", impact.Reason)
		}
	})

	t.Run("fail-open when files list is empty", func(t *testing.T) {
		impact := ResolveProjectChangeImpact(svcDocs, []string{}, "Mystery commit")
		if !impact.Affected {
			t.Errorf("expected empty files list to fail-open and deploy, got: %+v", impact)
		}
		if impact.Reason != "unknown-changes" {
			t.Errorf("expected reason unknown-changes, got: %s", impact.Reason)
		}
	})
}
