package environments

type RouteRuleType string

const (
	RouteRuleTypeRateLimit   RouteRuleType = "rate_limit"
	RouteRuleTypeIPAllowlist RouteRuleType = "ip_allowlist"
	RouteRuleTypeIPBlocklist RouteRuleType = "ip_blocklist"
	RouteRuleTypeHeaders     RouteRuleType = "headers"
)

const (
	DNSProvisionStatusPending = "pending"
	DNSProvisionStatusSuccess = "provisioned"
	DNSProvisionStatusFailed  = "failed"
)

const (
	RouteAuthCookieName = "codedock_route_session"
	DefaultSessionSec   = 86400
	MinPasswordLen      = 12
	MaxPasswordLen      = 256
	PBKDF2Iterations    = 600000
)

const (
	RouteProtectionStatusApplied = "applied"
	RouteProtectionStatusPending = "pending"
	RouteProtectionStatusFailed  = "failed"
)

const (
	ActionProtect = "protect"
	ActionRevoke  = "revoke"
	ActionRetry   = "retry"
	ActionUnlock  = "unlock"
)
