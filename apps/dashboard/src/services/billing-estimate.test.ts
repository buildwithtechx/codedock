import { describe, expect, it } from 'vitest';
import {
  type ComputePricing,
  DEFAULT_COMPUTE_PRICING,
  DEFAULT_PAYG_TIERS,
} from '#/constants/billing';
import {
  estimatePrepaidUsage,
  hasPrepaidRates,
  previewPaygPool,
  previewPaygTier,
} from './billing-estimate';

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
    const tier = DEFAULT_PAYG_TIERS[1];
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

  it('identifies the prospective PAYG tier enabled by a credit top-up', () => {
    expect(previewPaygTier(DEFAULT_PAYG_TIERS, 1000)?.id).toBe('tier_1');
    expect(previewPaygTier(DEFAULT_PAYG_TIERS, 2500)?.id).toBe('tier_2');
    expect(previewPaygTier(DEFAULT_PAYG_TIERS, 5000)?.id).toBe('tier_3');
    expect(previewPaygTier(DEFAULT_PAYG_TIERS, 10000)?.id).toBe('tier_4');
    expect(previewPaygTier(DEFAULT_PAYG_TIERS, 999)).toBeNull();
  });

  it('computes hourly usage and duration ranges correctly', () => {
    const est = estimatePrepaidUsage(input);
    expect(est).not.toBeNull();
    if (!est) return;

    expect(est.hourlyUsage.cpuVcpuHours).toBe(1);
    expect(est.hourlyUsage.memoryGiBHours).toBe(4);
    expect(est.hourlyUsage.storageGiBHours).toBe(25);

    expect(est.hourlyCents.cpu).toBeCloseTo(3);
    expect(est.hourlyCents.memory).toBeCloseTo(3.2);
    expect(est.hourlyCents.disk).toBeCloseTo((25 * 5) / 720);

    expect(est.totalHourlyCents).toBeGreaterThan(6.3);
    expect(est.coveredElapsedHours).toBeCloseTo(500 / est.totalHourlyCents);
    expect(est.durationRange.minElapsedHours).toBeLessThanOrEqual(
      est.durationRange.maxElapsedHours
    );
  });
});
