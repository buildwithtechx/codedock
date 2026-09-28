package utils

import "testing"

func TestSanitizeDomainNameCollisions(t *testing.T) {
	name1 := SanitizeDomainName("foo_bar")
	name2 := SanitizeDomainName("foobar")
	if name1 == name2 {
		t.Fatalf("expected foo_bar and foobar to produce distinct domains, got %s for both", name1)
	}

	name3 := SanitizeDomainName("foo-bar")
	if name1 == name3 {
		t.Fatalf("expected foo_bar and foo-bar to produce distinct domains, got %s for both", name1)
	}

	if name2 != "foobar" {
		t.Fatalf("expected foobar, got %s", name2)
	}
	if name3 != "foo-bar" {
		t.Fatalf("expected foo-bar, got %s", name3)
	}
}
