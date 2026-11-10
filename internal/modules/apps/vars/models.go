package vars

import "time"

type Variable struct {
	ID                       string    `json:"id" db:"id"`
	ServiceID                string    `json:"serviceId" db:"service_id"`
	EnvironmentID            string    `json:"environmentId,omitempty" db:"environment_id"`
	Key                      string    `json:"key" db:"key"`
	Value                    string    `json:"value" db:"value"`
	IsSecret                 bool      `json:"isSecret" db:"is_secret"`
	ExpectedTopologyRevision string    `json:"-" db:"-"`
	CreatedAt                time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt                time.Time `json:"updatedAt" db:"updated_at"`
}
