package auth

import (
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"codedock/internal/models"
)

type webAuthnUser struct {
	id          []byte
	name        string
	displayName string
	credentials []webauthn.Credential
}

func (u *webAuthnUser) WebAuthnID() []byte                         { return u.id }
func (u *webAuthnUser) WebAuthnName() string                       { return u.name }
func (u *webAuthnUser) WebAuthnDisplayName() string                { return u.displayName }
func (u *webAuthnUser) WebAuthnIcon() string                       { return "" }
func (u *webAuthnUser) WebAuthnCredentials() []webauthn.Credential { return u.credentials }

func toWebAuthnUser(user *models.User, creds []*models.PasskeyCredential) *webAuthnUser {
	var webCreds []webauthn.Credential
	for _, c := range creds {
		credID, err := base64.RawURLEncoding.DecodeString(c.CredentialID)
		if err != nil {
			continue
		}
		pubKey, err := base64.RawURLEncoding.DecodeString(c.PublicKey)
		if err != nil {
			continue
		}
		var transports []protocol.AuthenticatorTransport
		if c.Transports != "" {
			for _, t := range strings.Split(c.Transports, ",") {
				transports = append(transports, protocol.AuthenticatorTransport(strings.TrimSpace(t)))
			}
		}
		aaguid, _ := hex.DecodeString(c.AAGUID)
		webCreds = append(webCreds, webauthn.Credential{
			ID:        credID,
			PublicKey: pubKey,
			Authenticator: webauthn.Authenticator{
				AAGUID:    aaguid,
				SignCount: uint32(c.Counter),
			},
			Transport: transports,
			Flags: webauthn.CredentialFlags{
				BackupState: c.BackedUp,
			},
		})
	}

	dispName := user.Name
	if dispName == "" {
		dispName = user.Email
	}
	return &webAuthnUser{
		id:          []byte(user.ID),
		name:        user.Email,
		displayName: dispName,
		credentials: webCreds,
	}
}

func toPasskeyResponse(c *models.PasskeyCredential) *models.PasskeyResponse {
	var transports []string
	if c.Transports != "" {
		for _, t := range strings.Split(c.Transports, ",") {
			transports = append(transports, strings.TrimSpace(t))
		}
	}
	return &models.PasskeyResponse{
		ID:         c.ID,
		Name:       c.Name,
		DeviceType: c.DeviceType,
		BackedUp:   c.BackedUp,
		Transports: transports,
		AAGUID:     c.AAGUID,
		CreatedAt:  c.CreatedAt,
		LastUsedAt: c.LastUsedAt,
	}
}
