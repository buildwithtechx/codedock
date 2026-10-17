export interface CustomServerResources {
  cpuCores: number;
  memoryMb: number;
  diskGb: number;
}

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

/** Selected allocation across all proposed servers, not their metered CPU activity. */
export function previewPaygPool(
  tier: PaygTier,
  resources: CustomServerResources,
  serverCount: number
) {
  if (
    ![resources.cpuCores, resources.memoryMb, resources.diskGb, serverCount].every(
      (value) => Number.isSafeInteger(value) && value > 0
    )
  ) {
    return null;
  }
  const selected = {
    cpuCores: resources.cpuCores * serverCount,
    memoryMb: resources.memoryMb * serverCount,
    diskGb: resources.diskGb * serverCount,
    servers: serverCount,
  };
  if (!Object.values(selected).every((value) => Number.isSafeInteger(value) && value > 0)) {
    return null;
  }
  const exceeded = (Object.keys(selected) as Array<keyof typeof selected>).filter(
    (key) => selected[key] > tier.pool[key]
  );
  return { selected, exceeded, fits: exceeded.length === 0 };
}

/** Prospective access from this credit purchase, not a customer's current entitlement. */
export function previewPaygTier(tiers: PaygTier[], purchaseCents: number) {
  if (!Number.isSafeInteger(purchaseCents) || purchaseCents <= 0) return null;
  return [...tiers].reverse().find((tier) => tier.minimumFundingCents <= purchaseCents) ?? null;
}

export function hasPrepaidRates(pricing: unknown): pricing is ComputePricing {
  if (!pricing || typeof pricing !== 'object') return false;
  const p = pricing as ComputePricing;
  return Boolean(
    p.currency === 'usd' &&
      [
        p.creditsPerDollar,
        p.usage?.activeVcpuHourCents,
        p.usage?.reservedGiBHourCents,
        p.usage?.retainedGiBMonthCents,
        p.usage?.monthHours,
      ].every((val) => Number.isFinite(val) && (val ?? 0) > 0)
  );
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

/** Presentation only: published resource rates applied to prepaid usage. */
export function estimatePrepaidUsage({
  pricing,
  resources,
  packageCents,
  serverCount,
  cpuPercent,
}: {
  pricing: ComputePricing;
  resources: CustomServerResources;
  packageCents: number;
  serverCount: number;
  cpuPercent: number;
}): PrepaidUsageEstimate | null {
  if (
    !hasPrepaidRates(pricing) ||
    ![resources.cpuCores, resources.memoryMb, resources.diskGb, packageCents, serverCount].every(
      (val) => Number.isSafeInteger(val) && val > 0
    ) ||
    !Number.isFinite(cpuPercent) ||
    cpuPercent < 0 ||
    cpuPercent > 100
  ) {
    return null;
  }

  const creditsPerCent = pricing.creditsPerDollar / 100;
  const monthHours = pricing.usage.monthHours;

  const hourlyUsage = {
    cpuVcpuHours: ((resources.cpuCores * cpuPercent) / 100) * serverCount,
    memoryGiBHours: (resources.memoryMb / 1024) * serverCount,
    storageGiBHours: resources.diskGb * serverCount,
  };

  const hourlyCents = {
    cpu: hourlyUsage.cpuVcpuHours * pricing.usage.activeVcpuHourCents,
    memory: hourlyUsage.memoryGiBHours * pricing.usage.reservedGiBHourCents,
    disk: (hourlyUsage.storageGiBHours * pricing.usage.retainedGiBMonthCents) / monthHours,
  };

  const totalHourlyCents = hourlyCents.cpu + hourlyCents.memory + hourlyCents.disk;
  const idleHourlyCents = hourlyCents.memory + hourlyCents.disk;
  const fullCpuCents = resources.cpuCores * serverCount * pricing.usage.activeVcpuHourCents;
  const fullCpuHourlyCents = fullCpuCents + idleHourlyCents;

  if (
    ![totalHourlyCents, idleHourlyCents, fullCpuHourlyCents, packageCents * creditsPerCent].every(
      (value) => Number.isFinite(value) && value > 0
    )
  ) {
    return null;
  }

  const coveredElapsedHours = packageCents / totalHourlyCents;
  const durationRange = {
    minElapsedHours: packageCents / fullCpuHourlyCents,
    maxElapsedHours: packageCents / idleHourlyCents,
  };

  if (
    ![coveredElapsedHours, durationRange.minElapsedHours, durationRange.maxElapsedHours].every(
      (value) => Number.isFinite(value) && value > 0
    )
  ) {
    return null;
  }

  return {
    balanceCredits: packageCents * creditsPerCent,
    hourlyUsage,
    hourlyCents,
    totalHourlyCents,
    hourlyCredits: totalHourlyCents * creditsPerCent,
    coveredElapsedHours,
    durationRange,
  };
}

export function formatMoney(cents: number, digits = 2): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: digits,
  }).format(cents / 100);
}

export function formatNumber(val: number): string {
  return new Intl.NumberFormat('en-US').format(val);
}

export function formatDuration(hours: number): string {
  if (hours >= 24) return `${(hours / 24).toFixed(1)} days`;
  if (hours >= 1) return `${hours.toFixed(1)} hours`;
  return `${Math.round(hours * 60)} minutes`;
}
