package bare

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/base64"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func nativeFixture() (*models.AppService, models.RuntimeTarget) {
	app := &models.AppService{ID: "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", ProjectID: "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb", RuntimeMode: models.RuntimeModeWorker, Replicas: 1}
	target := models.RuntimeTarget{Kind: "bare", BareNode: models.ClusterNode{ServerID: "cccccccc-cccc-4ccc-8ccc-cccccccccccc", PrivateIP: "192.168.1.10", Interface: "eth0", Fingerprint: "SHA256:" + base64.RawStdEncoding.EncodeToString(make([]byte, 32))}, BareReleaseURL: "https://example.com/app", BareSHA256: strings.Repeat("a", 64), BareCommand: []string{"./app"}}
	return app, target
}
func TestNativeValidationRejectsUnverifiedArtifactsAndUnsupportedScaling(t *testing.T) {
	app, target := nativeFixture()
	if err := Validate(app, target); err != nil {
		t.Fatal(err)
	}
	target.BareReleaseURL = "http://example.com/app"
	if err := Validate(app, target); err == nil {
		t.Fatal("unencrypted artifact accepted")
	}
	target.BareReleaseURL = "https://example.com/app"
	target.BareSHA256 = "missing"
	if err := Validate(app, target); err == nil {
		t.Fatal("unverified artifact accepted")
	}
	target.BareSHA256 = strings.Repeat("a", 64)
	app.Replicas = 2
	if err := Validate(app, target); err == nil {
		t.Fatal("native process advertised cluster scaling")
	}
	app.Replicas = 1
	target.BareRepoURL = "https://github.com/example/app.git"
	if err := Validate(app, target); err == nil {
		t.Fatal("simultaneous source and artifact accepted")
	}
	target.BareReleaseURL, target.BareSHA256 = "", ""
	target.BareToolchain = "go"
	if err := Validate(app, target); err != nil {
		t.Fatal("valid source build rejected", err)
	}
	target.BareToolchain = "docker"
	if err := Validate(app, target); err == nil {
		t.Fatal("unsupported toolchain accepted")
	}
}

type nativeRunner struct{ calls int }

func (r *nativeRunner) Script(context.Context, models.ClusterNode, string, string) error {
	r.calls++
	return nil
}
func (r *nativeRunner) StreamHost(context.Context, models.ClusterNode, string, io.Reader, io.Writer) error {
	return nil
}
func (r *nativeRunner) Host(context.Context, models.ClusterNode, string) (string, error) {
	return "ActiveState=failed\nSubState=failed\nNRestarts=3\nMemoryCurrent=0\n", nil
}
func TestNativeFailedServiceIsObservedAndUnsupportedLifecycleDoesNotExecute(t *testing.T) {
	app, target := nativeFixture()
	runner := &nativeRunner{}
	runtime := NewRuntime(runner)
	observed, err := runtime.Observe(context.Background(), app, target)
	if err != nil || observed.Status != "FAILED" || observed.Available != 0 {
		t.Fatal("failed native service appeared ready", err)
	}
	if err := runtime.Lifecycle(context.Background(), app, target, "scale"); err == nil || runner.calls != 0 {
		t.Fatal("unsupported scaling executed")
	}
}
func TestNativeEnvironmentRejectsShellVariableInjection(t *testing.T) {
	app, target := nativeFixture()
	if _, err := DeploymentScript(app, target, map[string]string{"KEY;bad": "value"}, "dddddddd-dddd-4ddd-8ddd-dddddddddddd"); err == nil {
		t.Fatal("invalid environment key accepted")
	}
}

func TestNativeActivationAndRecoveryScriptsHaveValidShellSyntax(t *testing.T) {
	program, err := exec.LookPath("bash")
	if runtime.GOOS == "windows" {
		program = filepath.Join(os.Getenv("ProgramFiles"), "Git", "bin", "bash.exe")
		_, err = os.Stat(program)
	}
	if err != nil {
		t.Skip("Bash syntax checker unavailable")
	}
	app, target := nativeFixture()
	deploy, err := DeploymentScript(app, target, map[string]string{"VALUE": "a'b\nnext"}, "dddddddd-dddd-4ddd-8ddd-dddddddddddd")
	if err != nil {
		t.Fatal(err)
	}
	recover, err := recoveryScript(app)
	if err != nil {
		t.Fatal(err)
	}
	for _, script := range []string{deploy, recover} {
		command := exec.Command(program, "-n")
		command.Stdin = strings.NewReader(script)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("native script syntax invalid: %s: %v", output, err)
		}
	}
}

func TestNativeIdentityRemainsUsableAfterInvalidSettingsEdit(t *testing.T) {
	app, target := nativeFixture()
	app.Replicas = 2
	app.Domain = "edited.example.com"
	if err := Validate(app, target); err == nil {
		t.Fatal("invalid deployment accepted")
	}
	if err := ValidateIdentity(app, target); err != nil {
		t.Fatal("settings edits blocked owned runtime recovery", err)
	}
}
