package models

type Operation struct {
	ID        string `json:"id" db:"id"`
	UserID    string `json:"-" db:"user_id"`
	ProjectID string `json:"projectId" db:"project_id"`
	Kind      string `json:"kind" db:"kind"`
	Target    string `json:"target" db:"target"`
	Status    string `json:"status" db:"status"`
	Phase     string `json:"phase" db:"phase"`
	Effects   string `json:"effects" db:"effects"`
	Payload   string `json:"-" db:"encrypted_payload"`
	Snapshot  string `json:"-" db:"snapshot"`
	TokenHash string `json:"-" db:"token_hash"`
	Error     string `json:"error" db:"error"`
	Logs      string `json:"logs" db:"logs"`
	ExpiresAt int64  `json:"expiresAt" db:"expires_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}

type OperationReview struct {
	Operation    *Operation `json:"operation"`
	Confirmation string     `json:"confirmation"`
}

type OperationApply struct {
	Confirmation string `json:"confirmation"`
}
