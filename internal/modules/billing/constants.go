package billing

type CancelUnavailableReason string

const (
	CancelUnavailableUnsupported         CancelUnavailableReason = "unsupported"
	CancelUnavailableTooClose            CancelUnavailableReason = "too_close"
	CancelUnavailableProviderUnavailable CancelUnavailableReason = "provider_unavailable"
)
