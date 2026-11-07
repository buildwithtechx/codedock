import type { ComputePricing, PaygTier } from '#/constants/billing';
import type { CustomServerResources, PrepaidUsageEstimate } from '#/interfaces/billing';

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
