export interface BillingPlan {
  id: string;
  name: string;
  description: string;
  priceMonthly: number;
  priceId?: string;
  features: string[];
  limits: {
    cpuCores: number;
    ramGb: number;
    bandwidthGb: number;
    teamMembers: number;
  };
}

export interface BillingConfig {
  publishableKey?: string;
  plans: BillingPlan[];
  topupPacks: TopupPack[];
  portalUrl?: string;
}

export interface TopupPack {
  id: string;
  name: string;
  priceCents: number;
  credits: number;
  creditCents?: number;
  bonusPercent?: number;
  label?: string;
}

export interface PaygResourceAllocation {
  cpuCores: number;
  ramGb: number;
  bandwidthGb: number;
}

export interface PaygHourlyBurnRate {
  baseHourly: number;
  maxHourly: number;
  dailyEstimate: number;
  monthlyEstimate: number;
}

export interface PaygTierQuota {
  id: string;
  name: string;
  minCentsPerHour: number;
  maxCores: number;
  maxRamGb: number;
  networkPriority: 'standard' | 'high' | 'dedicated';
  recommended?: boolean;
}

export interface CustomPlanQuote {
  cores: number;
  memoryGb: number;
  bandwidthTb: number;
  monthlyPrice: number;
  hourlyPrice: number;
}

export interface PendingCheckout {
  id?: string;
  checkoutId: string;
  sessionId?: string;
  priceId?: string;
  url?: string;
  idempotencyKey?: string;
  amountCents?: number;
  status: 'open' | 'pending' | 'processing' | 'failed' | 'canceled' | 'complete';
  canCancel?: boolean;
  cancelUnavailableReason?: 'provider_unavailable' | 'unsupported';
  createdAt?: string;
  expiresAt?: string;
}

export interface CustomServerResources {
  cpuCores: number;
  memoryMb: number;
  diskGb: number;
}

export interface PrepaidUsageEstimate {
  balanceCredits: number;
  hourlyUsage: {
    cpuVcpuHours: number;
    memoryGiBHours: number;
    storageGiBHours: number;
  };
  hourlyCents: {
    cpu: number;
    memory: number;
    disk: number;
  };
  totalHourlyCents: number;
  hourlyCredits: number;
  coveredElapsedHours: number;
  durationRange: {
    minElapsedHours: number;
    maxElapsedHours: number;
  };
}
