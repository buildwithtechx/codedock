package models

import (
	"time"

	"codedock.run/codedock/pkg/types"
)

type ProjectConfig = types.ProjectConfig
type EnvironmentConfig = types.EnvironmentConfig
type CreateProjectRequest = types.CreateProjectRequest

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectToken struct {
	ID            string     `json:"id" db:"id"`
	ProjectID     string     `json:"projectId" db:"project_id"`
	EnvironmentID string     `json:"environmentId,omitempty" db:"environment_id"`
	Name          string     `json:"name" db:"name"`
	TokenHash     string     `json:"-" db:"token_hash"`
	TokenPrefix   string     `json:"tokenPrefix,omitempty" db:"token_prefix"`
	Prefix        string     `json:"prefix" db:"prefix"`
	Role          string     `json:"role" db:"role"`
	Scopes        []string   `json:"scopes,omitempty" db:"scopes"`
	IPAllowlist   []string   `json:"ipAllowlist,omitempty" db:"ip_allowlist"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty" db:"expires_at"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
}

type CreateTokenRequest struct {
	Name      string     `json:"name"`
	Role      string     `json:"role"`
	ExpiresAt *time.Time `json:"expiresAt"`
}

type CreateTokenResponse struct {
	Token        string        `json:"token"`
	ProjectToken *ProjectToken `json:"projectToken"`
}

type ProjectMember struct {
	ID        string    `json:"id" db:"id"`
	ProjectID string    `json:"projectId" db:"project_id"`
	UserID    string    `json:"userId" db:"user_id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type AddMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type ServerlessFunctionCode struct {
	ID          string    `json:"id" db:"id"`
	ServiceID   string    `json:"serviceId" db:"service_id"`
	Runtime     string    `json:"runtime" db:"runtime"`
	CodeContent string    `json:"codeContent" db:"code_content"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type CanvasSummary struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
	EnvironmentsCount  int                `json:"environmentsCount"`
	AppsCount          int                `json:"appsCount"`
	DatabasesCount     int                `json:"databasesCount"`
	TotalServices      int                `json:"totalServices"`
	OnlineServices     int                `json:"onlineServices"`
	DefaultEnvironment *EnvironmentConfig `json:"defaultEnvironment,omitempty"`
	ServiceIcons       []string           `json:"serviceIcons"`
	Nodes              []string           `json:"nodes,omitempty"`
	NodeCount          int                `json:"nodeCount,omitempty"`
	EdgeCount          int                `json:"edgeCount,omitempty"`
}

type CanvasPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type CanvasNode struct {
	ID   string         `json:"id"`
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
	Pos  CanvasPosition `json:"position"`
}

type CanvasEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type EnvironmentCanvas struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Environment *EnvironmentConfig `json:"environment"`
	Apps        []*AppService      `json:"apps"`
	Databases   []*Database        `json:"databases"`
	Nodes       []CanvasNode       `json:"nodes"`
	Edges       []CanvasEdge       `json:"edges"`
}
