package domains

import (
	"context"
	"testing"
)

type mockDomainRepo struct {
	domains []DomainConfig
}

func (m *mockDomainRepo) ListByService(ctx context.Context, serviceID string) ([]DomainConfig, error) {
	var res []DomainConfig
	for _, d := range m.domains {
		if d.ServiceID == serviceID {
			res = append(res, d)
		}
	}
	return res, nil
}

func (m *mockDomainRepo) ListAll(ctx context.Context) ([]DomainConfig, error) {
	return m.domains, nil
}

func (m *mockDomainRepo) GetByID(ctx context.Context, id string) (*DomainConfig, error) {
	for _, d := range m.domains {
		if d.ID == id {
			return &d, nil
		}
	}
	return nil, nil
}

func (m *mockDomainRepo) Create(ctx context.Context, d *DomainConfig) error {
	m.domains = append(m.domains, *d)
	return nil
}

func (m *mockDomainRepo) UpdateDNSProvisionStatus(ctx context.Context, id string, status string, provider string, provisionedIP string) error {
	for i, d := range m.domains {
		if d.ID == id {
			m.domains[i].DNSProvisionStatus = status
			m.domains[i].DNSProvider = provider
			m.domains[i].DNSProvisionedIP = provisionedIP
			return nil
		}
	}
	return nil
}

func (m *mockDomainRepo) Delete(ctx context.Context, id string) error {
	for i, d := range m.domains {
		if d.ID == id {
			m.domains = append(m.domains[:i], m.domains[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockDomainRepo) BindUnassignedVerifiedDomainsToSingleService(ctx context.Context, projectID, singleServiceID string) error {
	return nil
}

func TestCanonicalDomainName(t *testing.T) {
	tests := []struct {
		input    string
		valid    bool
		expected string
	}{
		{"example.com", true, "example.com"},
		{"Sub.Example.com.", true, "sub.example.com"},
		{"invalid..com", false, ""},
		{"-bad.com", false, ""},
		{"bad-.com", false, ""},
		{"localhost", false, ""},
		{"foo bar.com", false, ""},
	}

	for _, tc := range tests {
		res, err := CanonicalDomainName(tc.input)
		if tc.valid && err != nil {
			t.Errorf("expected valid domain for %q, got error: %v", tc.input, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("expected invalid domain for %q, got: %s", tc.input, res)
		}
		if tc.valid && res != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, res)
		}
	}
}

func TestDomainServiceCRUD(t *testing.T) {
	repo := &mockDomainRepo{}
	svc := NewDomainService(repo, nil)
	ctx := context.Background()

	// Validation
	if _, err := svc.CreateGeneratedDomain(ctx, nil); err == nil {
		t.Fatal("expected error on nil domain")
	}

	created, err := svc.CreateGeneratedDomain(ctx, &DomainConfig{
		ServiceID:  "svc-123",
		DomainName: "api.example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error creating domain: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected generated UUID")
	}

	fetched, err := svc.GetDomain(ctx, created.ID)
	if err != nil || fetched == nil || fetched.DomainName != "api.example.com" {
		t.Fatalf("failed to get domain: %v", err)
	}

	list, err := svc.ListDomainsByService(ctx, "svc-123")
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 domain in list, got %d", len(list))
	}

	if err := svc.DeleteDomain(ctx, created.ID); err != nil {
		t.Fatalf("failed to delete domain: %v", err)
	}
	listAfter, _ := svc.ListDomainsByService(ctx, "svc-123")
	if len(listAfter) != 0 {
		t.Fatalf("expected 0 domains after delete, got %d", len(listAfter))
	}
}
