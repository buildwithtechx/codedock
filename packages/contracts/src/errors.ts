/** Error codes for rollback admission and execution failures. */
export const RollbackErrorCodes = {
  NOT_READY: 'ROLLBACK_NOT_READY',
  ALREADY_ACTIVE: 'ROLLBACK_ALREADY_ACTIVE',
  ARTIFACT_GONE: 'ROLLBACK_ARTIFACT_GONE',
  CONFIGURATION_REDACTED: 'ROLLBACK_CONFIGURATION_REDACTED',
} as const;

export type RollbackErrorCode = (typeof RollbackErrorCodes)[keyof typeof RollbackErrorCodes];

/** Error codes for subscription and quota billing gates. */
export const BillingGateErrorCodes = {
  GATE_REQUIRED: 'BILLING_GATE_REQUIRED',
  PAYMENT_REQUIRED: 'PAYMENT_REQUIRED',
} as const;

export type BillingGateErrorCode =
  (typeof BillingGateErrorCodes)[keyof typeof BillingGateErrorCodes];

/** Structured API failure response returned on non-2xx status codes. */
export interface ApiErrorResponse {
  error: string;
  code?: string;
  success: false;
  details?: Record<string, unknown>;
}
