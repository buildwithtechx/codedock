/** Hosted management control plane pricing definition. */
export interface ManagementPricing {
  version: number;
  monthlyCents: number;
  includedServers: number;
  extraServerMonthlyCents: number;
  maxServers: number;
  includedWithCloud: boolean;
  dedicatedContact: string;
}

/** Active customer subscription for managed server slots. */
export interface ManagementSubscription {
  servers: number;
  baseIncluded: boolean;
  priceCents: number;
  periodEnd: string | null;
  cancelAtPeriodEnd: boolean;
  status: string;
}

/** Pending plan or slot count change. */
export interface ManagementChange {
  status: string;
  effectiveAt: string;
  paymentUrl: string | null;
}

/** Comprehensive management state for an organization. */
export interface ManagementState {
  pricing: ManagementPricing;
  purchasesAvailable: boolean;
  cloudIncluded: boolean;
  active: boolean;
  serverLimit: number;
  connectedServers: number;
  pendingCheckout: boolean;
  subscription: ManagementSubscription | null;
  change: ManagementChange | null;
}

/** Request payload to purchase or adjust managed server capacity. */
export interface ManagementPurchase {
  servers: number;
  idempotencyKey: string;
}

/** Response result from a management purchase attempt. */
export interface ManagementPurchaseResult {
  kind: 'checkout' | 'quote';
  url: string | null;
  quoteId: string | null;
  amountDueNow: number | null;
  monthlyCents: number;
  effectiveAt: string | null;
}
