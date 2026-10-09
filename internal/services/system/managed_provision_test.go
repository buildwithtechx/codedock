package system

import (
	"context"
	"fmt"
	"testing"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

func provisionManaged(t *testing.T, env *managedTestEnv, name string) *models.ManagedServerPlan {
	t.Helper()
	req := provisionRequest(env.credID)
	req.Name = name
	review, err := env.service.ReviewProvision(context.Background(), env.adminID, env.orgID, req)
	if err != nil {
		t.Fatalf("review %s: %v", name, err)
	}
	if err := env.service.ApplyOperation(context.Background(), env.adminID, review.Review.Operation.ID, review.Review.Confirmation); err != nil {
		t.Fatalf("apply %s: %v", name, err)
	}
	if op := waitManagedOperation(t, env.ops, review.Review.Operation.ID); op.Status != "COMPLETED" {
		t.Fatalf("expected %s to complete, got %s", name, op.Status)
	}
	return review.Plan
}

func TestManagedProvisionHappyPath(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	plan := provisionManaged(t, env, "web-1")
	if plan.MonthlyPrice != 4.51 {
		t.Fatalf("expected cx22 price 4.51, got %v", plan.MonthlyPrice)
	}
	server, err := env.serverRepo.GetByID(ctx, plan.ServerID)
	if err != nil {
		t.Fatalf("load server: %v", err)
	}
	if server.Status != models.ServerStatusOnline || server.IPAddress != "203.0.113.10" || server.ExternalID == "" {
		t.Fatalf("expected ready server, got %#v", server)
	}
	if server.Provider != "hetzner" || server.Region != "fsn1" || server.ServerType != "cx22" {
		t.Fatalf("expected hetzner placement, got %#v", server)
	}
	if server.SSHHost != server.IPAddress || server.SSHPrivateKey == "" || server.SSHUser != "root" {
		t.Fatalf("expected ssh target credentials, got %#v", server)
	}
	if server.UserID != env.adminID {
		t.Fatalf("expected owner %s, got %s", env.adminID, server.UserID)
	}
	link, err := env.managedRepo.GetLink(ctx, plan.ServerID)
	if err != nil {
		t.Fatalf("load link: %v", err)
	}
	if link.CredentialID != env.credID || link.SSHKeyName != plan.SSHKeyName {
		t.Fatalf("expected managed link, got %#v", link)
	}
	visible, err := env.serverRepo.ListByUser(ctx, env.adminID)
	if err != nil {
		t.Fatalf("list servers: %v", err)
	}
	found := false
	for _, candidate := range visible {
		if candidate.ID == plan.ServerID {
			found = true
		}
	}
	if !found {
		t.Fatal("expected managed server in owner target list")
	}
	if env.fake.creates != 1 {
		t.Fatalf("expected 1 provider server, got %d", env.fake.creates)
	}
	if _, ok := env.fake.keys[plan.SSHKeyName]; !ok {
		t.Fatal("expected provider ssh key to be registered")
	}
}

func TestManagedProvisionRetryReusesRow(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	env.fake.createErr = fmt.Errorf("provider outage")
	req := provisionRequest(env.credID)
	review, err := env.service.ReviewProvision(ctx, env.adminID, env.orgID, req)
	if err != nil {
		t.Fatalf("review: %v", err)
	}
	if err := env.service.ApplyOperation(ctx, env.adminID, review.Review.Operation.ID, review.Review.Confirmation); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if op := waitManagedOperation(t, env.ops, review.Review.Operation.ID); op.Status != "FAILED" {
		t.Fatalf("expected FAILED, got %s", op.Status)
	}
	failed, err := env.serverRepo.GetByID(ctx, review.Plan.ServerID)
	if err != nil {
		t.Fatalf("load failed server: %v", err)
	}
	if failed.Status != models.ServerStatusOffline {
		t.Fatalf("expected offline failure state, got %s", failed.Status)
	}
	env.fake.createErr = nil
	retry := provisionRequest(env.credID)
	retry.ServerID = review.Plan.ServerID
	second, err := env.service.ReviewProvision(ctx, env.adminID, env.orgID, retry)
	if err != nil {
		t.Fatalf("retry review: %v", err)
	}
	if err := env.service.ApplyOperation(ctx, env.adminID, second.Review.Operation.ID, second.Review.Confirmation); err != nil {
		t.Fatalf("retry apply: %v", err)
	}
	if op := waitManagedOperation(t, env.ops, second.Review.Operation.ID); op.Status != "COMPLETED" {
		t.Fatalf("expected retry to complete, got %s", op.Status)
	}
	if len(env.fake.servers) != 1 {
		t.Fatalf("expected 1 provider server after retry, got %d", len(env.fake.servers))
	}
	online := provisionRequest(env.credID)
	online.ServerID = review.Plan.ServerID
	if _, err := env.service.ReviewProvision(ctx, env.adminID, env.orgID, online); err == nil {
		t.Fatal("expected provision review of running server to fail")
	}
}

func TestManagedProvisionAdoptsLabelledServer(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	review, err := env.service.ReviewProvision(ctx, env.adminID, env.orgID, provisionRequest(env.credID))
	if err != nil {
		t.Fatalf("review: %v", err)
	}
	env.fake.servers["2001"] = &models.ManagedInstance{
		ExternalID: "2001", Name: "web-1", Status: "running", PublicIP: "203.0.113.20",
		ServerType: "cx22", Region: "fsn1", Labels: map[string]string{"codedock-server": review.Plan.ServerID},
	}
	if err := env.service.ApplyOperation(ctx, env.adminID, review.Review.Operation.ID, review.Review.Confirmation); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if op := waitManagedOperation(t, env.ops, review.Review.Operation.ID); op.Status != "COMPLETED" {
		t.Fatalf("expected adopt to complete, got %s", op.Status)
	}
	if env.fake.creates != 0 {
		t.Fatalf("expected adoption without create, got %d creates", env.fake.creates)
	}
	server, err := env.serverRepo.GetByID(ctx, review.Plan.ServerID)
	if err != nil {
		t.Fatalf("load server: %v", err)
	}
	if server.ExternalID != "2001" || server.IPAddress != "203.0.113.20" {
		t.Fatalf("expected adopted address, got %#v", server)
	}
}

func TestManagedApplyRejectsBadConfirmation(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	review, err := env.service.ReviewProvision(ctx, env.adminID, env.orgID, provisionRequest(env.credID))
	if err != nil {
		t.Fatalf("review: %v", err)
	}
	if err := env.service.ApplyOperation(ctx, env.adminID, review.Review.Operation.ID, "wrong"); err == nil {
		t.Fatal("expected wrong confirmation to fail")
	}
	if err := env.service.ApplyOperation(ctx, env.memberID, review.Review.Operation.ID, review.Review.Confirmation); err == nil {
		t.Fatal("expected cross-user apply to fail")
	}
	op, err := env.ops.Get(ctx, review.Review.Operation.ID)
	if err != nil {
		t.Fatalf("load operation: %v", err)
	}
	if op.Status != "REVIEWED" {
		t.Fatalf("expected operation to stay REVIEWED, got %s", op.Status)
	}
}

func TestManagedResizeLifecycle(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	plan := provisionManaged(t, env, "web-1")
	review, err := env.service.ReviewResize(ctx, env.adminID, env.orgID, plan.ServerID, models.ReviewManagedResizeRequest{ServerType: "cx32"})
	if err != nil {
		t.Fatalf("review resize: %v", err)
	}
	if review.Plan.MonthlyPrice != 8.75 {
		t.Fatalf("expected cx32 price 8.75, got %v", review.Plan.MonthlyPrice)
	}
	if err := env.service.ApplyOperation(ctx, env.adminID, review.Review.Operation.ID, review.Review.Confirmation); err != nil {
		t.Fatalf("apply resize: %v", err)
	}
	if op := waitManagedOperation(t, env.ops, review.Review.Operation.ID); op.Status != "COMPLETED" {
		t.Fatalf("expected resize to complete, got %s", op.Status)
	}
	server, err := env.serverRepo.GetByID(ctx, plan.ServerID)
	if err != nil {
		t.Fatalf("load server: %v", err)
	}
	if server.ServerType != "cx32" || server.Status != models.ServerStatusOnline {
		t.Fatalf("expected cx32 online, got %#v", server)
	}
	if len(env.fake.changes) != 1 || env.fake.powerOff != 1 || env.fake.powerOn != 1 {
		t.Fatalf("expected stop/change/start, got %+v", env.fake)
	}
	if _, err := env.service.ReviewResize(ctx, env.adminID, env.orgID, plan.ServerID, models.ReviewManagedResizeRequest{ServerType: "cx32"}); err == nil {
		t.Fatal("expected same-tier resize to fail")
	}
	if _, err := env.service.ReviewResize(ctx, env.adminID, env.orgID, plan.ServerID, models.ReviewManagedResizeRequest{ServerType: "cx22"}); err == nil {
		t.Fatal("expected smaller-disk resize to fail")
	}
	if err := env.serverRepo.UpdateStatus(ctx, plan.ServerID, models.ServerStatusProvisioning); err != nil {
		t.Fatalf("mark provisioning: %v", err)
	}
	if _, err := env.service.ReviewResize(ctx, env.adminID, env.orgID, plan.ServerID, models.ReviewManagedResizeRequest{ServerType: "cx42"}); err == nil {
		t.Fatal("expected provisioning resize to fail")
	}
}

func TestManagedDeleteLifecycle(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	plan := provisionManaged(t, env, "web-1")
	plainService := NewServerService(env.serverRepo, repositories.NewUserRepo(env.db), nil)
	if err := plainService.DeleteServer(ctx, plan.ServerID, env.adminID); err == nil {
		t.Fatal("expected plain delete of managed server to fail")
	}
	review, err := env.service.ReviewDelete(ctx, env.adminID, env.orgID, plan.ServerID)
	if err != nil {
		t.Fatalf("review delete: %v", err)
	}
	if err := env.service.ApplyOperation(ctx, env.adminID, review.Review.Operation.ID, review.Review.Confirmation); err != nil {
		t.Fatalf("apply delete: %v", err)
	}
	if op := waitManagedOperation(t, env.ops, review.Review.Operation.ID); op.Status != "COMPLETED" {
		t.Fatalf("expected delete to complete, got %s", op.Status)
	}
	server, err := env.serverRepo.GetByID(ctx, plan.ServerID)
	if err != nil {
		t.Fatalf("load server: %v", err)
	}
	if server != nil {
		t.Fatalf("expected server row removed, got %#v", server)
	}
	if link, err := env.managedRepo.GetLink(ctx, plan.ServerID); err != nil || link != nil {
		t.Fatalf("expected link removed, got %#v %v", link, err)
	}
	if len(env.fake.servers) != 0 || len(env.fake.keys) != 0 {
		t.Fatal("expected provider server and key removed")
	}
}

func TestManagedRefreshAndRecover(t *testing.T) {
	env := setupManagedTest(t)
	ctx := context.Background()
	plan := provisionManaged(t, env, "web-1")
	server, err := env.serverRepo.GetByID(ctx, plan.ServerID)
	if err != nil {
		t.Fatalf("load server: %v", err)
	}
	if err := env.fake.DeleteServer(ctx, server.ExternalID); err != nil {
		t.Fatalf("remove provider server: %v", err)
	}
	refreshed, err := env.service.RefreshServer(ctx, env.memberID, env.orgID, plan.ServerID)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if refreshed.Status != models.ServerStatusOffline {
		t.Fatalf("expected missing provider server to read offline, got %s", refreshed.Status)
	}
	if err := env.serverRepo.UpdateStatus(ctx, plan.ServerID, models.ServerStatusProvisioning); err != nil {
		t.Fatalf("mark provisioning: %v", err)
	}
	if err := env.service.Recover(ctx); err != nil {
		t.Fatalf("recover: %v", err)
	}
	recovered, err := env.serverRepo.GetByID(ctx, plan.ServerID)
	if err != nil {
		t.Fatalf("load recovered: %v", err)
	}
	if recovered.Status != models.ServerStatusOffline {
		t.Fatalf("expected stuck provisioning without provider to recover offline, got %s", recovered.Status)
	}
}
