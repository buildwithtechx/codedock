package hetzner

import (
	"strings"
	"testing"
)

func TestBuildOffersFiltersToPurchasableReviewedTiers(t *testing.T) {
	locations := []map[string]any{
		{"name": "fsn1", "city": "Falkenstein", "country": "DE"},
		{"name": "xyz9", "city": "Nowhere", "country": "ZZ"},
	}
	serverTypes := []map[string]any{
		{"name": "cx23", "prices": []any{
			map[string]any{"location": "fsn1", "price_monthly": map[string]any{"gross": "3.4900000000"}},
			map[string]any{"location": "xyz9", "price_monthly": map[string]any{"gross": "1.0000000000"}},
		}},
		{"name": "cx22", "deprecation": map[string]any{"announced": "2025-10-01"}, "prices": []any{
			map[string]any{"location": "fsn1", "price_monthly": map[string]any{"gross": "4.51"}},
		}},
		{"name": "ccx99", "prices": []any{
			map[string]any{"location": "fsn1", "price_monthly": map[string]any{"gross": "99"}},
		}},
	}
	offers := buildOffers(locations, serverTypes)
	if len(offers.Locations) != 1 || offers.Locations[0].City != "Falkenstein" {
		t.Fatalf("expected only reviewed locations, got %#v", offers.Locations)
	}
	if len(offers.Offers) != 1 || offers.Offers[0].Tier.Name != "cx23" {
		t.Fatalf("expected only cx23, got %#v", offers.Offers)
	}
	if price := offers.Offers[0].Prices["fsn1"]; price != 3.49 {
		t.Fatalf("expected 3.49, got %v", price)
	}
	if _, ok := offers.Offers[0].Prices["xyz9"]; ok {
		t.Fatal("unreviewed location leaked into prices")
	}
}

func TestProvisionUserDataHardensSSH(t *testing.T) {
	data := ProvisionUserData("Edge Node_1")
	for _, want := range []string{"hostname: edgenode1", "ssh_pwauth: false", "PasswordAuthentication no", "get.docker.com", "bootstrap-complete"} {
		if !strings.Contains(data, want) {
			t.Fatalf("userdata missing %q:\n%s", want, data)
		}
	}
}
