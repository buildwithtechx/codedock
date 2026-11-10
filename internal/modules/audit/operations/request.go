package operations

// CancelOperationRequest defines request parameters for aborting an operation.
type CancelOperationRequest struct {
	Reason string `json:"reason,omitempty"`
}
