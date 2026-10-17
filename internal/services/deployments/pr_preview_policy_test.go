package deployments

import (
	"testing"

	"codedock/internal/models"
)

func TestIsPreviewEligible(t *testing.T) {
	validPR := PullRequestInfo{
		Number:   42,
		State:    "open",
		Draft:    false,
		HeadSHA:  "a1b2c3d4e5f60718293a4b5c6d7e8f9012345678",
		HeadRef:  "feature-branch",
		HeadRepo: "myorg/myrepo",
		BaseRef:  "main",
		BaseRepo: "myorg/myrepo",
	}

	eligible, reason := IsPreviewEligible(validPR, "MyOrg", "MyRepo", "main")
	if !eligible {
		t.Fatalf("expected PR to be eligible, failed with: %s", reason)
	}

	// Closed PR
	closedPR := validPR
	closedPR.State = "closed"
	if eligible, _ := IsPreviewEligible(closedPR, "myorg", "myrepo", "main"); eligible {
		t.Fatalf("expected closed PR to be ineligible")
	}

	// Draft PR
	draftPR := validPR
	draftPR.Draft = true
	if eligible, _ := IsPreviewEligible(draftPR, "myorg", "myrepo", "main"); eligible {
		t.Fatalf("expected draft PR to be ineligible")
	}

	// Branch mismatch
	branchMismatchPR := validPR
	branchMismatchPR.BaseRef = "develop"
	if eligible, _ := IsPreviewEligible(branchMismatchPR, "myorg", "myrepo", "main"); eligible {
		t.Fatalf("expected branch mismatch PR to be ineligible")
	}

	// Invalid SHA
	badShaPR := validPR
	badShaPR.HeadSHA = "not-a-valid-sha"
	if eligible, _ := IsPreviewEligible(badShaPR, "myorg", "myrepo", "main"); eligible {
		t.Fatalf("expected invalid SHA PR to be ineligible")
	}
}

func TestValidatePreviewVolumes(t *testing.T) {
	validVolumes := []models.ServiceVolume{
		{HostPath: "codedock-data", ContainerPath: "/data"},
		{HostPath: "app_cache", ContainerPath: "/root/.cache"},
	}
	if err := ValidatePreviewVolumes(validVolumes); err != nil {
		t.Fatalf("expected valid volumes, got error: %v", err)
	}

	invalidBind := []models.ServiceVolume{
		{HostPath: "/host/data", ContainerPath: "/app/data"},
	}
	if err := ValidatePreviewVolumes(invalidBind); err == nil {
		t.Fatalf("expected error on host bind mount")
	}

	invalidRelative := []models.ServiceVolume{
		{HostPath: "./relative", ContainerPath: "/app/data"},
	}
	if err := ValidatePreviewVolumes(invalidRelative); err == nil {
		t.Fatalf("expected error on relative host bind mount")
	}
}

func TestValidatePreviewCommands(t *testing.T) {
	if err := ValidatePreviewCommands("npm run build", "node server.js"); err != nil {
		t.Fatalf("expected plain commands to pass: %v", err)
	}

	if err := ValidatePreviewCommands("NODE_ENV=production npm run build"); err == nil {
		t.Fatalf("expected error on command embedding env vars")
	}
}

func TestNormalizeRepoIdentity(t *testing.T) {
	id := NormalizeRepoIdentity("github", "MyOrg", "MyRepo.git")
	if id != "github:myorg/myrepo" {
		t.Fatalf("expected 'github:myorg/myrepo', got %q", id)
	}
}

func TestGeneratePreviewDomain(t *testing.T) {
	domain := GeneratePreviewDomain("preview.example.com", "", "web-app", "proj-123", 10)
	if domain == "" || !testing.Short() {
		// Verify structure: web-app-pr-10-<suffix>.preview.example.com
		if domain[:14] != "web-app-pr-10-" {
			t.Fatalf("unexpected domain prefix: %s", domain)
		}
	}
}
