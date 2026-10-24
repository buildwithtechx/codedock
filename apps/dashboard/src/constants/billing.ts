export interface ComputePricing {
  tariffId: string;
  currency: 'usd';
  creditsPerDollar: number;
  paygCapPercent?: number;
  usage: {
    activeVcpuHourCents: number;
    reservedGiBHourCents: number;
    retainedGiBMonthCents: number;
    monthHours: number;
  };
  network?: {
    managedProxyGiBCents: number;
    minimumTopupCents: number;
  };
  retentionDays?: number;
  paygCheckoutAvailable?: boolean;
}

export interface PaygTier {
  id: string;
  level: number;
  minimumFundingCents: number;
  pool: {
    cpuCores: number;
    memoryMb: number;
    diskGb: number;
    servers: number;
  };
}

export const DEFAULT_COMPUTE_PRICING: ComputePricing = {
  tariffId: 'codedock-standard-compute',
  currency: 'usd',
  creditsPerDollar: 100,
  paygCapPercent: 125,
  usage: {
    activeVcpuHourCents: 0.7,
    reservedGiBHourCents: 0.35,
    retainedGiBMonthCents: 10,
    monthHours: 730,
  },
  network: {
    managedProxyGiBCents: 5,
    minimumTopupCents: 1000,
  },
  retentionDays: 30,
  paygCheckoutAvailable: true,
};

export const DEFAULT_PAYG_TIERS: PaygTier[] = [
  {
    id: 'tier_1',
    level: 1,
    minimumFundingCents: 1000,
    pool: { cpuCores: 2, memoryMb: 4096, diskGb: 60, servers: 2 },
  },
  {
    id: 'tier_2',
    level: 2,
    minimumFundingCents: 2500,
    pool: { cpuCores: 4, memoryMb: 12288, diskGb: 150, servers: 4 },
  },
  {
    id: 'tier_3',
    level: 3,
    minimumFundingCents: 5000,
    pool: { cpuCores: 8, memoryMb: 24576, diskGb: 350, servers: 8 },
  },
  {
    id: 'tier_4',
    level: 4,
    minimumFundingCents: 10000,
    pool: { cpuCores: 16, memoryMb: 65536, diskGb: 800, servers: 16 },
  },
];

export const DEFAULT_CREDIT_PACKAGES_CENTS = [1000, 2500, 5000, 10000, 25000];

export const DEFAULT_TOPUP_PACKS = [
  { id: 'topup_10', name: '$10 Starter Pack', credits: 1000, priceCents: 1000 },
  { id: 'topup_25', name: '$25 Growth Pack', credits: 2500, priceCents: 2500 },
  { id: 'topup_50', name: '$50 Scale Pack', credits: 5500, priceCents: 5000, bonusPercent: 10 },
  { id: 'topup_100', name: '$100 Pro Pack', credits: 11500, priceCents: 10000, bonusPercent: 15 },
  { id: 'topup_250', name: '$250 Max Pack', credits: 30000, priceCents: 25000, bonusPercent: 20 },
];

export function isCloudPurchaseIntent(value: unknown): value is string {
  return (
    typeof value === 'string' &&
    /^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$/i.test(value)
  );
}
