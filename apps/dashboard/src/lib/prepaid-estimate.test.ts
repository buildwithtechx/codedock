import { describe, expect, it } from 'vitest';
import {
  type ComputePricing,
  DEFAULT_COMPUTE_PRICING,
  DEFAULT_PAYG_TIERS,
  estimatePrepaidUsage,
  hasPrepaidRates,
  previewPaygPool,
  previewPaygTier,
} from './prepaid-estimate';

const testPricing: ComputePricing = {
  tariffId: 'test-tariff',
  currency: 'usd',
  creditsPerDollar: 100,
  paygCapPercent: 125,
  usage: {
    activeVcpuHourCents: 3,
    reservedGiBHourCents: 0.8,
    retainedGiBMonthCents: 5,
    monthHours: 720,
  },
  network: { managedProxyGiBCents: 10, minimumTopupCents: 500 },
  retentionDays: 30,
  paygCheckoutAvailable: true,
};

const input = {
  pricing: testPricing,
  resources: { cpuCores: 1, memoryMb: 4096, diskGb: 25 },
  packageCents: 500,
  serverCount: 1,
  cpuPercent: 100,
};

describe('prepaid pricing preview', () => {
  it('validates prepaid pricing rates', () => {
    expect(hasPrepaidRates(testPricing)).toBe(true);
    expect(hasPrepaidRates(DEFAULT_COMPUTE_PRICING)).toBe(true);
    expect(hasPrepaidRates(null)).toBe(false);
    expect(hasPrepaidRates({ currency: 'eur' })).toBe(false);
  });

  it('compares configured server resources with a shared tier pool', () => {
    const tier = DEFAULT_PAYG_TIERS[1]; // Tier 2: 4 cores, 12288 MB, 150 GB, 4 servers
    expect(previewPaygPool(tier, input.resources, 2)).toEqual({
      selected: { cpuCores: 2, memoryMb: 8192, diskGb: 50, servers: 2 },
      exceeded: [],
      fits: true,
    });

    const exceeded = previewPaygPool(tier, input.resources, 5);
    expect(exceeded?.fits).toBe(false);
    expect(exceeded?.exceeded).toContain('cpuCores');
    expect(exceeded?.exceeded).toContain('servers');

    expect(previewPaygPool(tier, { ...input.resources, diskGb: 200 }, 1)?.exceeded).toEqual([
      'diskGb',
    ]);
    expect(previewPaygPool(tier, { ...input.resources, cpuCores: -1 }, 1)).toBeNull();
  });

  it('previews prospective tier unlocked by purchase amount', () => {
    expect(previewPaygTier(DEFAULT_PAYG_TIERS, 500)).toBeNull();
    expect(previewPaygTier(DEFAULT_PAYG_TIERS, 1000)?.id).toBe('tier_1');
    expect(previewPaygTier(DEFAULT_PAYG_TIERS, 2499)?.id).toBe('tier_1');
    expect(previewPaygTier(DEFAULT_PAYG_TIERS, 2500)?.id).toBe('tier_2');
    expect(previewPaygTier(DEFAULT_PAYG_TIERS, 5000)?.id).toBe('tier_3');
    expect(previewPaygTier(DEFAULT_PAYG_TIERS, 10000)?.id).toBe('tier_4');
    expect(previewPaygTier(DEFAULT_PAYG_TIERS, -50)).toBeNull();
  });

  it('estimates hourly usage and calculates credits and runtime duration', () => {
    const result = estimatePrepaidUsage(input);
    expect(result).not.toBeNull();
    if (!result) return;

    expect(result.balanceCredits).toBe(500);
    expect(result.hourlyCents.cpu).toBe(3);
    expect(result.hourlyCents.memory).toBe(3.2);
    expect(result.hourlyCents.disk * 720).toBeCloseTo(125);
    expect(result.totalHourlyCents).toBeCloseTo(6.3736, 3);
    expect(result.hourlyCredits).toBeCloseTo(6.3736, 3);
    expect(result.coveredElapsedHours).toBeCloseTo(500 / result.totalHourlyCents);
    expect(result.durationRange.minElapsedHours).toBeLessThan(result.durationRange.maxElapsedHours);
  });

  it('scales covered duration linearly with package amount', () => {
    const baseline = estimatePrepaidUsage(input)!;
    const doubled = estimatePrepaidUsage({ ...input, packageCents: 1000 })!;
    expect(doubled.coveredElapsedHours).toBeCloseTo(baseline.coveredElapsedHours * 2);
  });

  it('maintains memory and disk costs when CPU is at idle (0%)', () => {
    const idle = estimatePrepaidUsage({ ...input, cpuPercent: 0 })!;
    expect(idle.hourlyCents.cpu).toBe(0);
    expect(idle.hourlyCents.memory).toBeGreaterThan(0);
    expect(idle.hourlyCents.disk).toBeGreaterThan(0);
    expect(idle.durationRange.maxElapsedHours).toBeCloseTo(
      input.packageCents / idle.totalHourlyCents
    );
  });
});
