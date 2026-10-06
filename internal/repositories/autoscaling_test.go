package repositories

import (
	"context"
	"testing"
	"time"
)

func TestAutoscalingPersistenceAndReservationGuards(t *testing.T) {
	db := openTestDB(t)
	db.SetMaxOpenConns(1)
	if err := RunMigrations(db); err != nil {
		t.Fatal(err)
	}
	for _, q := range []string{
		`INSERT INTO organizations(id,name) VALUES('org','Org')`,
		`INSERT INTO project_apps(id,organization_id,name,slug) VALUES('app','org','App','app')`,
		`INSERT INTO projects(id,app_id,organization_id,name,slug) VALUES('project','app','org','Project','project')`,
		`INSERT INTO app_services(id,project_id,name,status,replicas) VALUES('service','project','Service','running',2)`,
	} {
		if _, err := db.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	ctx := context.Background()
	repo := NewAutoscalingRepo(db)
	policy, err := repo.Get(ctx, "service")
	if err != nil || policy.Enabled {
		t.Fatalf("default policy enabled: %v", err)
	}
	policy.Enabled = true
	if err := repo.Save(ctx, policy); err != nil {
		t.Fatal(err)
	}
	policies, err := repo.ListEnabled(ctx)
	if err != nil || len(policies) != 1 {
		t.Fatal("enabled policy missing")
	}
	if ok, err := repo.SetReplicas(ctx, policy, 2, 3); err != nil || !ok {
		t.Fatalf("reservation failed: %v", err)
	}
	if ok, err := repo.SetReplicas(ctx, policy, 2, 4); err != nil || ok {
		t.Fatalf("stale replica reservation accepted: %v", err)
	}
	if err := repo.Observe(ctx, policy, 95, "requested deployment", true); err != nil {
		t.Fatal(err)
	}
	saved, err := repo.Get(ctx, "service")
	if err != nil || saved.LastDecision != "requested deployment" || saved.LastScaledAt == "" {
		t.Fatal("decision was not persisted")
	}
	if _, err := time.Parse(time.RFC3339Nano, saved.LastScaledAt); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO deployments(id,organization_id,project_id,service_id,status) VALUES('ready','org','project','service','READY')`); err != nil {
		t.Fatal(err)
	}
	if ok, err := repo.SetReplicas(ctx, policy, 3, 4); err != nil || !ok {
		t.Fatalf("completed READY deployment blocked scaling: %v", err)
	}
	if err := repo.RollbackReplicas(ctx, "service", 4, 3); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO deployments(id,organization_id,project_id,service_id,status) VALUES('busy','org','project','service','BUILDING')`); err != nil {
		t.Fatal(err)
	}
	if ok, err := repo.SetReplicas(ctx, policy, 3, 4); err != nil || ok {
		t.Fatalf("active deployment allowed scaling: %v", err)
	}
	if err := repo.RollbackReplicas(ctx, "service", 3, 2); err != nil {
		t.Fatal(err)
	}
	var replicas int
	if err := db.QueryRow(`SELECT replicas FROM app_services WHERE id='service'`).Scan(&replicas); err != nil || replicas != 2 {
		t.Fatal("replicas were not rolled back")
	}
	policy.Enabled = false
	if err := repo.Save(ctx, policy); err != nil {
		t.Fatal(err)
	}
	if ok, err := repo.SetReplicas(ctx, policy, 2, 3); err != nil || ok {
		t.Fatalf("disabled policy allowed scaling: %v", err)
	}
}
