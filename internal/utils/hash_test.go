package utils

import "testing"

func TestHashPasswordAndCheck(t *testing.T) {
	hash, err := HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if !CheckPasswordHash("correct-horse-battery-staple", hash) {
		t.Fatal("expected correct password to match")
	}
	if CheckPasswordHash("wrong-password", hash) {
		t.Fatal("expected wrong password to not match")
	}
}

func TestHashPasswordIsUnique(t *testing.T) {
	h1, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("first HashPassword failed: %v", err)
	}
	h2, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("second HashPassword failed: %v", err)
	}
	if h1 == h2 {
		t.Fatal("expected two hashes of the same password to differ (random salt)")
	}
}

func TestCheckPasswordHashInvalidFormat(t *testing.T) {
	if CheckPasswordHash("any", "not-a-valid-hash") {
		t.Fatal("expected invalid hash format to return false")
	}
	if CheckPasswordHash("any", "") {
		t.Fatal("expected empty hash to return false")
	}
}

func TestDecodeArgon2RoundTrip(t *testing.T) {
	hash, err := HashPassword("roundtrip")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	salt, stored, time, memory, threads, err := decodeArgon2(hash)
	if err != nil {
		t.Fatalf("decodeArgon2 failed: %v", err)
	}
	if len(salt) != argonSaltLen {
		t.Fatalf("expected salt length %d, got %d", argonSaltLen, len(salt))
	}
	if len(stored) != int(argonKeyLen) {
		t.Fatalf("expected hash length %d, got %d", argonKeyLen, len(stored))
	}
	if time != argonTime || memory != argonMemory || threads != argonThreads {
		t.Fatal("decoded params do not match encoding params")
	}
}
