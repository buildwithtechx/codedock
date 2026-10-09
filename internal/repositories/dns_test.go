package repositories

import (
	"context"
	"testing"

	"codedock/internal/models"
)

func TestDNSRecordRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewDNSRepo(db)
	record := &models.DNSRecord{DomainName: "example.com", RecordType: "A", RecordName: "@", RecordValue: "10.0.0.1"}
	if err := repo.Create(ctx, record); err != nil {
		t.Fatalf("create record: %v", err)
	}
	if record.ID == "" || record.TTL != 3600 || record.CreatedAt.IsZero() {
		t.Fatalf("unexpected created record: %+v", record)
	}
	loaded, err := repo.GetByID(ctx, record.ID)
	if err != nil {
		t.Fatalf("get record: %v", err)
	}
	if loaded.RecordValue != "10.0.0.1" || loaded.DomainName != "example.com" {
		t.Fatalf("unexpected record row: %+v", loaded)
	}
	other := &models.DNSRecord{DomainName: "other.com", RecordType: "A", RecordName: "@", RecordValue: "10.0.0.2", TTL: 60}
	if err := repo.Create(ctx, other); err != nil {
		t.Fatalf("create other record: %v", err)
	}
	filtered, err := repo.ListByDomain(ctx, "example.com")
	if err != nil {
		t.Fatalf("list by domain: %v", err)
	}
	if len(filtered) != 1 || filtered[0].ID != record.ID {
		t.Fatalf("unexpected filtered list: %+v", filtered)
	}
	all, err := repo.ListByDomain(ctx, "")
	if err != nil {
		t.Fatalf("list all records: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 records, got %d", len(all))
	}
	record.RecordValue = "10.0.0.9"
	record.TTL = 120
	if err := repo.Update(ctx, record); err != nil {
		t.Fatalf("update record: %v", err)
	}
	updated, err := repo.GetByID(ctx, record.ID)
	if err != nil {
		t.Fatalf("get updated record: %v", err)
	}
	if updated.RecordValue != "10.0.0.9" || updated.TTL != 120 {
		t.Fatalf("update did not persist: %+v", updated)
	}
	if err := repo.Delete(ctx, record.ID); err != nil {
		t.Fatalf("delete record: %v", err)
	}
	if _, err := repo.GetByID(ctx, record.ID); err == nil {
		t.Fatal("expected deleted record to be gone")
	}
}
