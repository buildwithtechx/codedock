package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"math"
	"strings"
	"testing"
)

func testWorkload() *models.KubernetesWorkload {
	return &models.KubernetesWorkload{App: models.AppService{ID: uuid.NewString(), ProjectID: uuid.NewString(), EnvironmentID: uuid.NewString(), Replicas: 1, InternalPort: 3000, RuntimeMode: models.RuntimeModeWeb}, Target: models.RuntimeTarget{Kind: "kubernetes", ClusterID: uuid.NewString()}, Image: "registry.example.com/app:version", Variables: map[string]string{"SECRET": "private"}}
}
func TestWorkersHaveNoHTTPReadinessOrService(t *testing.T) {
	workload := testWorkload()
	workload.App.RuntimeMode = models.RuntimeModeWorker
	workload.App.InternalPort = 0
	raw, err := WorkloadManifest(workload)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(raw, "readinessProbe") || strings.Contains(raw, `"kind":"Service"`) || strings.Contains(raw, "containerPort") {
		t.Fatal("worker requires an HTTP listener")
	}
	workload.App.Domain = "worker.example.com"
	if _, err := WorkloadManifest(workload); err == nil {
		t.Fatal("worker HTTP route accepted")
	}
}
func TestPersistentStorageRequiresExplicitSharedReplicaSupport(t *testing.T) {
	workload := testWorkload()
	workload.App.Replicas = 2
	workload.Target.Volumes = []models.RuntimeVolume{{Name: "data", MountPath: "/data", SizeGiB: 10}}
	if err := ValidateWorkload(workload); err == nil {
		t.Fatal("replicated workload accepted single-writer storage")
	}
	workload.Target.Volumes[0].Shared = true
	if err := ValidateWorkload(workload); err != nil {
		t.Fatal(err)
	}
	workload.Target.Volumes[0].MountPath = "/data/../etc"
	if err := ValidateWorkload(workload); err == nil {
		t.Fatal("unclean mount accepted")
	}
}

type fakeWorkloadCommands struct {
	output string
	calls  [][]string
	inputs []string
}

func (f *fakeWorkloadCommands) Kubectl(_ context.Context, _ models.ClusterNode, args []string, input string) (string, error) {
	f.calls = append(f.calls, args)
	f.inputs = append(f.inputs, input)
	return f.output, nil
}
func TestRecoveryRefusesForeignResources(t *testing.T) {
	workload := testWorkload()
	namespace, name, err := WorkloadIdentity(&workload.App)
	if err != nil {
		t.Fatal(err)
	}
	object := map[string]any{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": map[string]any{"name": name, "namespace": namespace, "labels": map[string]string{"codedock.run/project": "other", "codedock.run/service": workload.App.ID}}}
	data, err := json.Marshal(map[string]any{"items": []any{object}})
	if err != nil {
		t.Fatal(err)
	}
	commands := &fakeWorkloadCommands{}
	runtime := NewWorkloadRuntime(commands)
	if err := runtime.Rollback(context.Background(), models.ClusterNode{}, &workload.App, string(data)); err == nil {
		t.Fatal("foreign recovery journal accepted")
	}
	if len(commands.calls) != 0 {
		t.Fatal("foreign resources mutated")
	}
}
func TestMetricsParseCPUAndBinaryMemory(t *testing.T) {
	for _, test := range []struct {
		value    string
		expected float64
	}{{"123000000n", 0.123}, {"500m", 0.5}, {"256Mi", 268435456}, {"1e3", 1000}} {
		actual, err := metricQuantity(test.value)
		if err != nil || math.Abs(actual-test.expected) > 1e-10 {
			t.Fatalf("%s = %v: %v", test.value, actual, err)
		}
	}
}

func TestCleanupBindsUIDAndUsesIngressResourcePath(t *testing.T) {
	workload := testWorkload()
	_, name, err := WorkloadIdentity(&workload.App)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(map[string]any{"metadata": map[string]any{"uid": "owned-uid", "labels": map[string]string{"codedock.run/service": workload.App.ID, "codedock.run/project": workload.App.ProjectID}}})
	if err != nil {
		t.Fatal(err)
	}
	commands := &fakeWorkloadCommands{output: string(raw)}
	runtime := NewWorkloadRuntime(commands)
	if err := runtime.deleteOwned(context.Background(), models.ClusterNode{}, &workload.App, "ingress", name); err != nil {
		t.Fatal(err)
	}
	var deletion struct {
		Preconditions struct {
			UID string `json:"uid"`
		} `json:"preconditions"`
	}
	if err := json.Unmarshal([]byte(commands.inputs[len(commands.inputs)-1]), &deletion); err != nil {
		t.Fatal(err)
	}
	if deletion.Preconditions.UID != "owned-uid" {
		t.Fatal("cleanup omitted UID precondition")
	}
	args := strings.Join(commands.calls[len(commands.calls)-1], " ")
	if !strings.Contains(args, "/ingresses/"+name) || !strings.Contains(args, "--raw") {
		t.Fatal("cleanup uses an invalid resource path", args)
	}
}
