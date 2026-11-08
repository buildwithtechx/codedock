package security

import (
	"time"
)

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityModerate Severity = "moderate"
	SeverityLow      Severity = "low"
)

type Action string

const (
	ActionRefreshImages Action = "refresh-images"
	ActionManual        Action = "manual"
	ActionPatchHost     Action = "patch-host"
)

type TargetKind string

const (
	TargetKindImage       TargetKind = "image"
	TargetKindRuntime     TargetKind = "runtime"
	TargetKindPlatform    TargetKind = "platform"
	TargetKindHostPackage TargetKind = "host-package"
)

type Confidence string

const (
	ConfidenceConfirmed Confidence = "confirmed"
	ConfidencePossible  Confidence = "possible"
)

type Target struct {
	Kind TargetKind `json:"kind"`

	// TargetKindImage
	Repository      string   `json:"repository,omitempty"`
	AffectedDigests []string `json:"affectedDigests,omitempty"`
	AffectedTags    []string `json:"affectedTags,omitempty"`
	FixedDigests    []string `json:"fixedDigests,omitempty"`

	// TargetKindRuntime
	Runtime  string  `json:"runtime,omitempty"`
	Affected string  `json:"affected,omitempty"`
	Fixed    *string `json:"fixed,omitempty"`

	// TargetKindHostPackage
	Distribution string `json:"distribution,omitempty"`
	Release      string `json:"release,omitempty"`
	Package      string `json:"package,omitempty"`
	FixedVersion string `json:"fixedVersion,omitempty"`
}

type SecurityAdvisory struct {
	ID          string    `json:"id"`
	Revision    int64     `json:"revision"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	Severity    Severity  `json:"severity"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"publishedAt"`
	Action      Action    `json:"action"`
	Modes       []string  `json:"modes,omitempty"`
	Targets     []Target  `json:"targets"`
}

type SecurityManifest struct {
	Format     int                `json:"format"`
	Revision   int64              `json:"revision"`
	IssuedAt   time.Time          `json:"issuedAt"`
	ExpiresAt  time.Time          `json:"expiresAt"`
	Advisories []SecurityAdvisory `json:"advisories"`
}

type SignedProof struct {
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

type AdvisoryFinding struct {
	AdvisoryID             string                `json:"advisoryId"`
	Title                  string                `json:"title"`
	Severity               Severity              `json:"severity"`
	Message                string                `json:"message"`
	URL                    string                `json:"url"`
	Action                 Action                `json:"action"`
	TargetKind             TargetKind            `json:"targetKind"`
	Confidence             Confidence            `json:"confidence"`
	ServiceName            string                `json:"serviceName"`
	ServiceID              string                `json:"serviceId,omitempty"`
	Evidence               string                `json:"evidence"`
	CanRefresh             bool                  `json:"canRefresh"`
	RuntimeUpgradeRequired bool                  `json:"runtimeUpgradeRequired"`
	ImageChangeRequired    bool                  `json:"imageChangeRequired"`
	RuntimePatch           *SecurityRuntimePatch `json:"runtimePatch,omitempty"`
}
