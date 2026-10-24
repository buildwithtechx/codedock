export interface ActiveSubscriptionOverview {
  status: string;
  planId: string;
  planName: string;
  prepaidBalanceUsd: number;
  autoRechargeEnabled: boolean;
  currentPeriodEnd?: string;
}

export interface PrepaidTopupPackage {
  id: string;
  amountUsd: number;
  bonusUsd: number;
  totalCreditsUsd: number;
  badge?: string;
}

export interface TenantCorsConfig {
  tenantId: string;
  allowedOrigins: string[];
  allowedMethods?: string[];
  maxAgeSeconds?: number;
}

export interface RateLimitDetails {
  retryAfterSeconds: number;
  limit?: number;
  remaining?: number;
  resetAt?: string;
  planTier?: string;
}

export interface BillingGateDetails {
  code: 'BILLING_GATE_REQUIRED' | 'PAYMENT_REQUIRED';
  message: string;
  upgradeUrl?: string;
  requiredPlan?: string;
}
