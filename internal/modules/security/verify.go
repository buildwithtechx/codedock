package security

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// DefaultSecurityPublicKeySPKI is the standard Ed25519 SPKI public key.
const DefaultSecurityPublicKeySPKI = "MCowBQYDK2VwAyEAc59hAZg7Jcsk9ajFBjNjqZ5JlMN3mbBm4rbXRCGDtVw="

const (
	MaxPayloadBytes = 1_000_000
	PrimaryDomain   = "openship-security-advisories-v1\n"
	SecondaryDomain = "codedock-security-advisories-v1\n"
)

type Verifier struct {
	publicKey ed25519.PublicKey
}

func NewVerifier(spkiBase64 string) (*Verifier, error) {
	if spkiBase64 == "" {
		spkiBase64 = DefaultSecurityPublicKeySPKI
	}
	der, err := base64.StdEncoding.DecodeString(spkiBase64)
	if err != nil {
		return nil, fmt.Errorf("decode public key base64: %w", err)
	}

	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse PKIX public key: %w", err)
	}

	edPub, ok := pub.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("public key is not an Ed25519 public key")
	}

	return &Verifier{publicKey: edPub}, nil
}

func (v *Verifier) VerifyProof(proof SignedProof, minRevision int64, allowExpired bool, now time.Time) (*SecurityManifest, error) {
	if proof.Payload == "" || proof.Signature == "" {
		return nil, ErrInvalidProof
	}
	if len(proof.Payload) > MaxPayloadBytes {
		return nil, ErrPayloadTooLarge
	}

	sig, err := base64.StdEncoding.DecodeString(proof.Signature)
	if err != nil || len(sig) != ed25519.SignatureSize {
		return nil, fmt.Errorf("%w: invalid signature encoding or length", ErrSignatureInvalid)
	}

	// Verify signature using primary domain separator, fallback to secondary
	verified := false
	for _, domain := range []string{PrimaryDomain, SecondaryDomain} {
		msg := append([]byte(domain), []byte(proof.Payload)...)
		if ed25519.Verify(v.publicKey, msg, sig) {
			verified = true
			break
		}
	}

	if !verified {
		return nil, ErrSignatureInvalid
	}

	var manifest SecurityManifest
	if err := json.Unmarshal([]byte(proof.Payload), &manifest); err != nil {
		return nil, fmt.Errorf("unmarshal security manifest: %w", err)
	}

	if manifest.Format != 1 {
		return nil, fmt.Errorf("%w: expected format 1, got %d", ErrInvalidFormat, manifest.Format)
	}

	if manifest.Revision < minRevision {
		return nil, fmt.Errorf("%w: revision %d < min %d", ErrRollbackDetected, manifest.Revision, minRevision)
	}

	if manifest.IssuedAt.After(now.Add(5 * time.Minute)) {
		return nil, ErrClockSkew
	}

	if !manifest.ExpiresAt.After(manifest.IssuedAt) {
		return nil, errors.New("manifest expiration must be strictly after issuance")
	}

	if manifest.ExpiresAt.Sub(manifest.IssuedAt) > 31*24*time.Hour {
		return nil, errors.New("manifest expiration exceeds maximum 31 days from issuance")
	}

	if !allowExpired && now.After(manifest.ExpiresAt) {
		return nil, ErrManifestExpired
	}

	return &manifest, nil
}
