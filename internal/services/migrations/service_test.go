package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"codedock.run/codedock/internal/engine/dockerprobe"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	projectservices "codedock.run/codedock/internal/services/projects"

	_ "modernc.org/sqlite"
)

type migrationTestVault struct{}

func (migrationTestVault) Encrypt(plaintext string) (string, error) {
	return "enc:" + plaintext, nil
}

func (migrationTestVault) Decrypt(ciphertext string) (string, error) {
	if !strings.HasPrefix(ciphertext, "enc:") {
		return "", fmt.Errorf("not encrypted")
	}
	return strings.TrimPrefix(ciphertext, "enc:"), nil
}

type fakeRunner struct {
	mu       sync.Mutex
	ps       string
	inspect  map[string]string
	volumes  map[string]string
	commands []string
	stopped  []string
	removed  []string
	started  []string
	loaded   int
	created  []string
}

func (f *fakeRunner) Run(ctx context.Context, cmd string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.commands = append(f.commands, cmd)
	switch {
	case strings.HasPrefix(cmd, "docker ps"):
		return f.ps, nil
	case strings.HasPrefix(cmd, "docker inspect -f"):
		return "running\n", nil
	case strings.HasPrefix(cmd, "docker inspect "):
		id := strings.TrimPrefix(cmd, "docker inspect ")
		if out, ok := f.inspect[id]; ok {
			return out, nil
		}
		return "", fmt.Errorf("no such container")
	case strings.HasPrefix(cmd, "docker volume ls -q --filter name="):
		name := strings.TrimPrefix(cmd, "docker volume ls -q --filter name=")
		if _, ok := f.volumes[name]; ok {
			return name + "\n", nil
		}
		return "", nil
	case strings.HasPrefix(cmd, "docker volume ls"):
		var names []string
		for name := range f.volumes {
			names = append(names, name)
		}
		return strings.Join(names, "\n"), nil
	case strings.HasPrefix(cmd, "docker volume create "):
		name := strings.TrimPrefix(cmd, "docker volume create ")
		f.volumes[name] = ""
		f.created = append(f.created, name)
		return name + "\n", nil
	case strings.HasPrefix(cmd, "docker volume rm"):
		name := strings.TrimSuffix(strings.TrimPrefix(cmd, "docker volume rm -f "), "\n")
		delete(f.volumes, strings.TrimSpace(name))
		return "", nil
	case strings.Contains(cmd, "alpine ls -A"):
		for name, data := range f.volumes {
			if strings.Contains(cmd, "-v "+name+":") {
				return data, nil
			}
		}
		return "", nil
	case strings.HasPrefix(cmd, "docker stop "):
		f.stopped = append(f.stopped, strings.TrimPrefix(cmd, "docker stop "))
		return "", nil
	case strings.HasPrefix(cmd, "docker rm "):
		f.removed = append(f.removed, strings.TrimPrefix(cmd, "docker rm "))
		return "", nil
	case strings.HasPrefix(cmd, "docker start "):
		f.started = append(f.started, strings.TrimPrefix(cmd, "docker start "))
		return "", nil
	}
	return "", fmt.Errorf("unexpected command: %s", cmd)
}

func (f *fakeRunner) RunPipe(ctx context.Context, cmd string, stdin io.Reader, stdout io.Writer) error {
	f.mu.Lock()
	f.commands = append(f.commands, "PIPE:"+cmd)
	if strings.Contains(cmd, "docker load") {
		f.loaded++
	}
	f.mu.Unlock()
	if stdin != nil {
		_, _ = io.Copy(io.Discard, stdin)
	}
	if stdout != nil {
		_, _ = stdout.Write([]byte("stream-bytes"))
	}
	return nil
}

type fakeDeployer struct {
	mu       sync.Mutex
	apps     repositories.AppServiceRepository
	deployed []string
	removed  []string
	failNext error
}

func (f *fakeDeployer) DeployAppService(ctx context.Context, appID, sourceDir string, logWriter io.Writer) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failNext != nil {
		err := f.failNext
		f.failNext = nil
		return "", err
	}
	app, err := f.apps.GetByID(ctx, appID)
	if err != nil {
		return "", err
	}
	app.ContainerID = "target-" + appID
	if len(app.ContainerID) > 20 {
		app.ContainerID = app.ContainerID[:20]
	}
	app.Status = models.AppServiceStatusRunning
	if err := f.apps.Update(ctx, app); err != nil {
		return "", err
	}
	f.deployed = append(f.deployed, appID)
	_, _ = logWriter.Write([]byte("deployed " + appID + "\n"))
	return app.ContainerID, nil
}

func (f *fakeDeployer) StopAppService(ctx context.Context, app *models.AppService) error {
	return nil
}

func (f *fakeDeployer) RemoveAppService(ctx context.Context, app *models.AppService) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.removed = append(f.removed, app.ID)
	return nil
}

type migrationTestEnv struct {
	db       *sql.DB
	service  *Service
	source   *fakeRunner
	target   *fakeRunner
	deployer *fakeDeployer
	orgID    string
	adminID  string
	memberID string
	sourceID string
}

const (
	fakeWebID = "abc123def456"
	fakeDbID  = "def456abc123"
)

func webInspect(id string) string {
	return `[{"Id": "` + id + `", "Name": "/web", "Config": {"Image": "nginx:1.25", "Env": ["APP_SECRET=topsecret", "PORT=80"], "Labels": {"com.docker.compose.project": "shop"}}, "HostConfig": {"Binds": [], "PortBindings": {"80/tcp": [{"HostIp": "0.0.0.0", "HostPort": "8080"}]}}, "Mounts": [{"Source": "webdata", "Destination": "/usr/share/nginx/html", "Type": "volume", "Name": "webdata"}], "State": {"Status": "running", "Running": true, "Restarting": false}, "NetworkSettings": {"Ports": {}}}]`
}

func dbInspect(id string) string {
	return `[{"Id": "` + id + `", "Name": "/db", "Config": {"Image": "postgres:16", "Env": ["POSTGRES_PASSWORD=dbsecret"], "Labels": {}}, "HostConfig": {"Binds": [], "PortBindings": {}}, "Mounts": [{"Source": "dbdata", "Destination": "/var/lib/postgresql/data", "Type": "volume", "Name": "dbdata"}], "State": {"Status": "exited", "Running": false, "Restarting": false}, "NetworkSettings": {"Ports": {}}}]`
}

const fakePs = `{"ID":"abc123def456","Names":"web","Image":"nginx:1.25","State":"running","Ports":"0.0.0.0:8080->80/tcp"}
{"ID":"def456abc123","Names":"db","Image":"postgres:16","State":"exited","Ports":""}
`

func setupMigrationTest(t *testing.T) *migrationTestEnv {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(ON)")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	if err := repositories.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	ctx := context.Background()
	vault := migrationTestVault{}
	userRepo := repositories.NewUserRepo(db)
	orgRepo := repositories.NewOrganizationRepository(db)
	envRepo := repositories.NewEnvironmentRepo(db)
	projectRepo := repositories.NewProjectRepo(db, envRepo)
	appRepo := repositories.NewAppServiceRepo(db)
	appService := projectservices.NewAppService(appRepo, repositories.NewServiceVarRepo(db), repositories.NewServiceVolumeRepo(db))
	serverRepo := repositories.NewServerRepository(db, vault)
	migrationRepo := repositories.NewMigrationRepository(db, vault)
	for _, user := range []*models.User{
		{ID: "admin", Email: "admin@example.com", Role: models.UserRoleMember, IsActive: true, PlanType: "pro"},
		{ID: "member", Email: "member@example.com", Role: models.UserRoleMember, IsActive: true, PlanType: "free"},
		{ID: "outsider", Email: "outsider@example.com", Role: models.UserRoleMember, IsActive: true, PlanType: "free"},
	} {
		if err := userRepo.CreateUser(ctx, user); err != nil {
			t.Fatalf("create user %s: %v", user.ID, err)
		}
	}
	if err := orgRepo.Create(ctx, &models.Organization{ID: "org-1", Name: "Acme"}); err != nil {
		t.Fatalf("create org: %v", err)
	}
	now := time.Now().UTC()
	for _, member := range []*models.OrganizationMember{
		{ID: "m-admin", OrganizationID: "org-1", UserID: "admin", Email: "admin@example.com", Permission: models.MemberPermissionAdmin, Status: models.MemberStatusActive, InvitedAt: now, AcceptedAt: now},
		{ID: "m-member", OrganizationID: "org-1", UserID: "member", Email: "member@example.com", Permission: models.MemberPermissionMember, Status: models.MemberStatusActive, InvitedAt: now, AcceptedAt: now},
	} {
		if err := orgRepo.AddMember(ctx, member); err != nil {
			t.Fatalf("add member %s: %v", member.UserID, err)
		}
	}
	source := &fakeRunner{ps: fakePs, inspect: map[string]string{fakeWebID: webInspect(fakeWebID), fakeDbID: dbInspect(fakeDbID)}, volumes: map[string]string{}}
	target := &fakeRunner{volumes: map[string]string{}}
	deployer := &fakeDeployer{apps: appRepo}
	service := NewService(migrationRepo, projectRepo, appService, appRepo, repositories.NewServiceVolumeRepo(db), envRepo, serverRepo, userRepo, orgRepo, nil, projectservices.NewComposeParserService(), deployer, nil, nil,
		func(*models.MigrationSource) dockerprobe.Runner { return source },
		func(*models.Server) dockerprobe.Runner { return target })
	env := &migrationTestEnv{db: db, service: service, source: source, target: target, deployer: deployer, orgID: "org-1", adminID: "admin", memberID: "member"}
	if err := migrationRepo.CreateSource(ctx, &models.MigrationSource{ID: "src-1", OrganizationID: "org-1", Name: "legacy", SSHHost: "10.0.0.9", SSHPort: 22, SSHUser: "root", SSHAuthMethod: "key", SSHKey: "k", Fingerprint: "fp"}); err != nil {
		t.Fatalf("create source: %v", err)
	}
	env.sourceID = "src-1"
	return env
}

func waitRunStatus(t *testing.T, env *migrationTestEnv, userID, runID string, wanted ...models.MigrationStatus) *RunDetail {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for {
		detail, err := env.service.GetRun(context.Background(), userID, runID)
		if err != nil {
			t.Fatalf("load run: %v", err)
		}
		for _, status := range wanted {
			if detail.Run.Status == status {
				return detail
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("run %s stuck in %s (%s): %s", runID, detail.Run.Status, detail.Run.Phase, detail.Run.Logs)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func waitRunPrompt(t *testing.T, env *migrationTestEnv, userID, runID string) *RunDetail {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for {
		detail, err := env.service.GetRun(context.Background(), userID, runID)
		if err != nil {
			t.Fatalf("load run: %v", err)
		}
		if detail.Prompt != nil {
			return detail
		}
		if detail.Run.Status == models.MigrationStatusFailed || detail.Run.Status == models.MigrationStatusCancelled {
			t.Fatalf("run ended before prompt: %s %s", detail.Run.Status, detail.Run.Error)
		}
		if time.Now().After(deadline) {
			t.Fatalf("run %s never prompted", runID)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestMigrationSourcePermissions(t *testing.T) {
	env := setupMigrationTest(t)
	ctx := context.Background()
	if _, err := env.service.ListSources(ctx, "outsider", env.orgID); err == nil {
		t.Fatal("expected outsider list to fail")
	}
	if _, err := env.service.CreateSource(ctx, "outsider", env.orgID, models.CreateMigrationSourceRequest{Name: "x"}); err == nil {
		t.Fatal("expected outsider create to fail")
	}
	if _, err := env.service.CreateSource(ctx, env.memberID, env.orgID, models.CreateMigrationSourceRequest{Name: "x"}); err == nil {
		t.Fatal("expected member create to fail")
	}
	if _, err := env.service.CreateSource(ctx, env.adminID, env.orgID, models.CreateMigrationSourceRequest{Name: "x", SSHHost: "h"}); err == nil {
		t.Fatal("expected credential validation to fail")
	}
	listed, err := env.service.ListSources(ctx, env.memberID, env.orgID)
	if err != nil {
		t.Fatalf("member list: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("expected 1 source, got %d", len(listed))
	}
}

func TestMigrationScanMasked(t *testing.T) {
	env := setupMigrationTest(t)
	stack, err := env.service.ScanSource(context.Background(), env.memberID, env.orgID, env.sourceID, nil)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(stack.Containers) != 2 {
		t.Fatalf("expected 2 containers, got %d", len(stack.Containers))
	}
	for _, container := range stack.Containers {
		for key, value := range container.Env {
			if value != dockerprobe.MaskedSecret {
				t.Fatalf("expected masked %s, got %q", key, value)
			}
		}
	}
	if len(stack.Containers[0].EnvKeys) != 2 {
		t.Fatalf("expected env keys, got %#v", stack.Containers[0].EnvKeys)
	}
	if len(stack.ComposeProjects) != 1 || stack.ComposeProjects[0] != "shop" {
		t.Fatalf("expected shop project, got %v", stack.ComposeProjects)
	}
}

func TestMigrationRevealEnv(t *testing.T) {
	env := setupMigrationTest(t)
	ctx := context.Background()
	req := models.RevealMigrationEnvRequest{SourceID: env.sourceID, ContainerID: fakeWebID}
	if _, err := env.service.RevealSourceEnv(ctx, env.memberID, env.orgID, req); err == nil {
		t.Fatal("expected member reveal to fail")
	}
	revealed, err := env.service.RevealSourceEnv(ctx, env.adminID, env.orgID, req)
	if err != nil {
		t.Fatalf("reveal: %v", err)
	}
	if revealed["APP_SECRET"] != "topsecret" {
		t.Fatalf("expected real secret, got %#v", revealed)
	}
}

func TestMigrationAdoptRecords(t *testing.T) {
	env := setupMigrationTest(t)
	result, err := env.service.AdoptSource(context.Background(), env.memberID, env.orgID, models.AdoptMigrationRequest{
		SourceID: env.sourceID, ContainerIDs: []string{fakeWebID, fakeDbID}, ProjectName: "shop", ImportEnv: true,
	})
	if err != nil {
		t.Fatalf("adopt: %v", err)
	}
	if len(result.Services) != 2 || len(result.Failures) != 0 {
		t.Fatalf("expected 2 adopted services, got %#v", result)
	}
	if _, err := env.service.AdoptSource(context.Background(), env.memberID, env.orgID, models.AdoptMigrationRequest{
		SourceID: env.sourceID, ContainerIDs: []string{"missing"}, ProjectName: "shop2",
	}); err == nil {
		t.Fatal("expected missing container to fail")
	}
}

func TestMigrationReimport(t *testing.T) {
	env := setupMigrationTest(t)
	ctx := context.Background()
	labeled := strings.Replace(webInspect(fakeWebID), `"com.docker.compose.project": "shop"`, `"com.docker.compose.project": "shop", "com.codedock.project": "proj-keep"`, 1)
	env.source.inspect[fakeWebID] = labeled
	result, err := env.service.ReimportSource(ctx, env.memberID, env.orgID, models.ReimportMigrationRequest{SourceID: env.sourceID, ContainerIDs: []string{fakeWebID}})
	if err != nil {
		t.Fatalf("reimport: %v", err)
	}
	if result.ProjectID != "proj-keep" {
		t.Fatalf("expected preserved project id, got %s", result.ProjectID)
	}
	if _, err := env.service.ReimportSource(ctx, env.memberID, env.orgID, models.ReimportMigrationRequest{SourceID: env.sourceID, ContainerIDs: []string{fakeWebID}}); err == nil {
		t.Fatal("expected double reimport to fail")
	}
	if _, err := env.service.ReimportSource(ctx, env.memberID, env.orgID, models.ReimportMigrationRequest{SourceID: env.sourceID, ContainerIDs: []string{fakeDbID}}); err == nil {
		t.Fatal("expected unlabeled reimport to fail")
	}
}

func TestMigrationPreviewConflicts(t *testing.T) {
	env := setupMigrationTest(t)
	env.target.ps = "{\"ID\":\"t1\",\"Names\":\"other\",\"Image\":\"x\",\"State\":\"running\",\"Ports\":\"0.0.0.0:8080->80/tcp\"}\n"
	env.target.volumes["webdata"] = "index.html\n"
	preview, err := env.service.PreviewMigration(context.Background(), env.memberID, env.orgID, models.PreviewMigrationRequest{
		SourceID: env.sourceID, ContainerIDs: []string{fakeWebID}, Mode: models.MigrationModeMove,
	})
	if err != nil {
		t.Fatalf("preview: %v", err)
	}
	if len(preview.Images) != 1 || len(preview.Volumes) != 1 {
		t.Fatalf("expected image and volume inventory, got %#v", preview)
	}
	if len(preview.Conflicts) != 2 {
		t.Fatalf("expected port and volume conflicts, got %#v", preview.Conflicts)
	}
	if len(preview.Warnings) == 0 || preview.Downtime == "" {
		t.Fatalf("expected warnings and downtime, got %#v", preview)
	}
}
