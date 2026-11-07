package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"
)

func TestVerifier(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	spkiBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatalf("marshal spki: %v", err)
	}
	spkiBase64 := base64.StdEncoding.EncodeToString(spkiBytes)

	verifier, err := NewVerifier(spkiBase64)
	if err != nil {
		t.Fatalf("new verifier: %v", err)
	}

	now := time.Now().UTC()
	manifest := SecurityManifest{
		Format:    1,
		Revision:  100,
		IssuedAt:  now.Add(-1 * time.Hour),
		ExpiresAt: now.Add(14 * 24 * time.Hour),
		Advisories: []SecurityAdvisory{
			{
				ID:       "TEST-1",
				Revision: 1,
				Title:    "Test Advisory",
				Severity: SeverityHigh,
				Action:   ActionRefreshImages,
			},
		},
	}

	payloadBytes, _ := json.Marshal(manifest)
	payloadStr := string(payloadBytes)

	sign := func(domain, text string) string {
		msg := append([]byte(domain), []byte(text)...)
		sig := ed25519.Sign(priv, msg)
		return base64.StdEncoding.EncodeToString(sig)
	}

	t.Run("verifies valid proof with primary domain separator", func(t *testing.T) {
		proof := SignedProof{
			Payload:   payloadStr,
			Signature: sign(PrimaryDomain, payloadStr),
		}
		res, err := verifier.VerifyProof(proof, 50, false, now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Revision != 100 || len(res.Advisories) != 1 {
			t.Fatalf("unexpected manifest: %+v", res)
		}
	})

	t.Run("verifies valid proof with secondary domain separator", func(t *testing.T) {
		proof := SignedProof{
			Payload:   payloadStr,
			Signature: sign(SecondaryDomain, payloadStr),
		}
		res, err := verifier.VerifyProof(proof, 50, false, now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Revision != 100 {
			t.Fatalf("unexpected manifest: %+v", res)
		}
	})

	t.Run("rejects tampered payload or signature", func(t *testing.T) {
		proof := SignedProof{
			Payload:   payloadStr + " ",
			Signature: sign(PrimaryDomain, payloadStr),
		}
		if _, err := verifier.VerifyProof(proof, 50, false, now); err != ErrSignatureInvalid {
			t.Fatalf("expected ErrSignatureInvalid, got %v", err)
		}
	})

	t.Run("rejects rollback revision", func(t *testing.T) {
		proof := SignedProof{
			Payload:   payloadStr,
			Signature: sign(PrimaryDomain, payloadStr),
		}
		if _, err := verifier.VerifyProof(proof, 200, false, now); err == nil {
			t.Fatal("expected error on revision 100 < min 200")
		}
	})

	t.Run("rejects expired manifest when allowExpired is false", func(t *testing.T) {
		expiredManifest := manifest
		expiredManifest.ExpiresAt = now.Add(-10 * time.Minute)
		expBytes, _ := json.Marshal(expiredManifest)
		expStr := string(expBytes)

		proof := SignedProof{
			Payload:   expStr,
			Signature: sign(PrimaryDomain, expStr),
		}
		if _, err := verifier.VerifyProof(proof, 50, false, now); err != ErrManifestExpired {
			t.Fatalf("expected ErrManifestExpired, got %v", err)
		}

		// Allowed when allowExpired is true
		if _, err := verifier.VerifyProof(proof, 50, true, now); err != nil {
			t.Fatalf("expected allowExpired to succeed, got %v", err)
		}
	})
}
