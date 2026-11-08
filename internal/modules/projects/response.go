package projects

import "time"

type CreateTokenResponse struct {
	Token        string        `json:"token"`
	ProjectToken *ProjectToken `json:"projectToken"`
}

type TokenCreatedPayload struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Token       string     `json:"token"`
	Scopes      []string   `json:"scopes"`
	IPAllowlist []string   `json:"ipAllowlist"`
	ExpiresAt   *time.Time `json:"expiresAt"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type StatusResponse struct {
	Status string `json:"status"`
}
