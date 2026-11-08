package dns

import (
	"context"
	"errors"
	"sync"
	"testing"
)

type memoryDNSRepo struct {
	mu      sync.Mutex
	records map[string]*DNSRecord
}

func newMemoryDNSRepo() *memoryDNSRepo {
	return &memoryDNSRepo{records: make(map[string]*DNSRecord)}
}

func (m *memoryDNSRepo) Create(ctx context.Context, record *DNSRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if record.ID == "" {
		record.ID = "rec_test_1"
	}
	m.records[record.ID] = record
	return nil
}

func (m *memoryDNSRepo) GetByID(ctx context.Context, id string) (*DNSRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.records[id]
	if !ok {
		return nil, ErrRecordNotFound
	}
	return r, nil
}

func (m *memoryDNSRepo) ListByDomain(ctx context.Context, domainName string) ([]*DNSRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []*DNSRecord
	for _, r := range m.records {
		if domainName == "" || r.DomainName == domainName {
			list = append(list, r)
		}
	}
	return list, nil
}

func (m *memoryDNSRepo) Update(ctx context.Context, record *DNSRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records[record.ID] = record
	return nil
}

func (m *memoryDNSRepo) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.records, id)
	return nil
}

var _ Repository = (*memoryDNSRepo)(nil)

func TestDNSService_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := newMemoryDNSRepo()
	svc := NewService(repo, nil)

	// Validation error on missing fields
	_, err := svc.CreateRecord(ctx, &CreateDNSRecordRequest{DomainName: "example.com"})
	if !errors.Is(err, ErrFieldsRequired) {
		t.Fatalf("expected ErrFieldsRequired, got: %v", err)
	}

	// Successful creation
	rec, err := svc.CreateRecord(ctx, &CreateDNSRecordRequest{
		DomainName:  "example.com",
		RecordType:  "A",
		RecordName:  "@",
		RecordValue: "1.2.3.4",
		TTL:         300,
	})
	if err != nil {
		t.Fatalf("create record: %v", err)
	}
	if rec.ID == "" || rec.DomainName != "example.com" {
		t.Fatalf("unexpected record: %+v", rec)
	}

	// List
	list, err := svc.ListByDomain(ctx, "example.com")
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 record, got %d (err: %v)", len(list), err)
	}

	// Update
	updated, err := svc.UpdateRecord(ctx, rec.ID, &UpdateDNSRecordRequest{
		RecordValue: "5.6.7.8",
	})
	if err != nil {
		t.Fatalf("update record: %v", err)
	}
	if updated.RecordValue != "5.6.7.8" {
		t.Fatalf("expected updated value, got %s", updated.RecordValue)
	}

	// Delete
	if err := svc.DeleteRecord(ctx, rec.ID); err != nil {
		t.Fatalf("delete record: %v", err)
	}
	afterDelete, _ := svc.ListByDomain(ctx, "example.com")
	if len(afterDelete) != 0 {
		t.Fatalf("expected 0 records after delete, got %d", len(afterDelete))
	}
}
