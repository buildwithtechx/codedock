export const RollbackErrorCodes = {
  NOT_READY: 'ROLLBACK_NOT_READY',
  ALREADY_ACTIVE: 'ROLLBACK_ALREADY_ACTIVE',
  ARTIFACT_GONE: 'ROLLBACK_ARTIFACT_GONE',
  CONFIGURATION_REDACTED: 'ROLLBACK_CONFIGURATION_REDACTED',
} as const;

export type RollbackErrorCode = (typeof RollbackErrorCodes)[keyof typeof RollbackErrorCodes];

export const BillingGateErrorCodes = {
  GATE_REQUIRED: 'BILLING_GATE_REQUIRED',
  PAYMENT_REQUIRED: 'PAYMENT_REQUIRED',
} as const;

export type BillingGateErrorCode =
  (typeof BillingGateErrorCodes)[keyof typeof BillingGateErrorCodes];

export interface ApiErrorResponse {
  error: string;
  code?: string;
  success: false;
  details?: Record<string, unknown>;
}
