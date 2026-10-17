import { Cpu, HardDrive, Layers } from 'lucide-react';
import { useId, useMemo, useState } from 'react';
import { Button } from '#/components/ui/button';
import {
  type ComputePricing,
  DEFAULT_COMPUTE_PRICING,
  DEFAULT_CREDIT_PACKAGES_CENTS,
  DEFAULT_PAYG_TIERS,
  type PaygTier,
} from '#/constants/billing';
import type { CustomServerResources } from '#/interfaces/billing';
import {
  estimatePrepaidUsage,
  formatMoney,
  formatNumber,
  previewPaygPool,
  previewPaygTier,
} from '#/services/billing-estimate';
import { PaygEstimateCard } from './payg-estimate-card';
import { PaygResourceTiers } from './payg-resource-tiers';

export interface PayAsYouGoPlanProps {
  pricing?: ComputePricing;
  tiers?: PaygTier[];
  creditPackagesCents?: number[];
  busy?: boolean;
  onAddCredits?: (packageCents: number) => void;
}

export function PayAsYouGoPlan({
  pricing = DEFAULT_COMPUTE_PRICING,
  tiers = DEFAULT_PAYG_TIERS,
  creditPackagesCents = DEFAULT_CREDIT_PACKAGES_CENTS,
  busy = false,
  onAddCredits,
}: PayAsYouGoPlanProps) {
  const id = useId();
  const [packageCents, setPackageCents] = useState<number>(
    creditPackagesCents[1] ?? creditPackagesCents[0] ?? 2500
  );
  const [tierId, setTierId] = useState<string>(
    () => previewPaygTier(tiers, packageCents)?.id ?? tiers[0]?.id ?? 'tier_1'
  );

  const selectedTier = tiers.find((t) => t.id === tierId) ?? tiers[0];
  const packageTier = previewPaygTier(tiers, packageCents);

  const [cores, setCores] = useState<number>(2);
  const [ramGb, setRamGb] = useState<number>(4);
  const [diskGb, setDiskGb] = useState<number>(40);

  const resources: CustomServerResources = useMemo(
    () => ({ cpuCores: cores, memoryMb: ramGb * 1024, diskGb }),
    [cores, ramGb, diskGb]
  );

  const pool = previewPaygPool(selectedTier, resources, 1);
  const largerTier =
    pool && !pool.fits ? tiers.find((t) => previewPaygPool(t, resources, 1)?.fits) : null;

  const estimate = pool?.fits
    ? estimatePrepaidUsage({ pricing, resources, packageCents, serverCount: 1, cpuPercent: 100 })
    : null;

  const balanceCredits = (packageCents / 100) * pricing.creditsPerDollar;

  return (
    <div className="space-y-6">
      <div className="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
        <div className="space-y-6">
          <section
            aria-labelledby={`${id}-sizing`}
            className="rounded-2xl border border-border bg-card p-6 shadow-sm"
          >
            <div className="flex flex-wrap items-center justify-between gap-2">
              <h3 id={`${id}-sizing`} className="font-semibold text-foreground text-base">
                Host resource sizing & pool limits
              </h3>
              <span className="rounded-lg bg-muted/60 px-2.5 py-1 font-medium text-muted-foreground text-xs">
                Prepaid PAYG
              </span>
            </div>
            <p className="mt-1 text-muted-foreground text-sm">
              Configure your server footprint. Usage is billed to prepaid credits second-by-second.
            </p>

            <div className="mt-5 border-border/50 border-b pb-5">
              <PaygResourceTiers
                tiers={tiers}
                selected={selectedTier}
                onChange={setTierId}
                formatMoney={formatMoney}
              />
            </div>

            <div className="mt-5 grid gap-5 sm:grid-cols-3">
              <div>
                <label className="flex items-center justify-between font-medium text-foreground text-xs">
                  <span className="flex items-center gap-1.5">
                    <Cpu className="size-3.5 text-primary" />
                    Dedicated vCPU
                  </span>
                  <span className="font-semibold">{cores} Cores</span>
                </label>
                <input
                  type="range"
                  aria-label="vCPU Cores"
                  min={1}
                  max={selectedTier.pool.cpuCores}
                  value={cores}
                  onChange={(e) => setCores(Number(e.target.value))}
                  className="mt-2 h-2 w-full cursor-pointer appearance-none rounded-lg bg-muted accent-primary"
                />
                <span className="mt-1 block text-muted-foreground text-xs">
                  Max {selectedTier.pool.cpuCores} vCPU
                </span>
              </div>

              <div>
                <label className="flex items-center justify-between font-medium text-foreground text-xs">
                  <span className="flex items-center gap-1.5">
                    <Layers className="size-3.5 text-primary" />
                    RAM Memory
                  </span>
                  <span className="font-semibold">{ramGb} GB</span>
                </label>
                <input
                  type="range"
                  aria-label="RAM Memory"
                  min={1}
                  max={Math.floor(selectedTier.pool.memoryMb / 1024)}
                  value={ramGb}
                  onChange={(e) => setRamGb(Number(e.target.value))}
                  className="mt-2 h-2 w-full cursor-pointer appearance-none rounded-lg bg-muted accent-primary"
                />
                <span className="mt-1 block text-muted-foreground text-xs">
                  Max {Math.floor(selectedTier.pool.memoryMb / 1024)} GB
                </span>
              </div>

              <div>
                <label className="flex items-center justify-between font-medium text-foreground text-xs">
                  <span className="flex items-center gap-1.5">
                    <HardDrive className="size-3.5 text-primary" />
                    Fast NVMe
                  </span>
                  <span className="font-semibold">{diskGb} GB</span>
                </label>
                <input
                  type="range"
                  aria-label="Fast NVMe"
                  min={10}
                  max={selectedTier.pool.diskGb}
                  step={5}
                  value={diskGb}
                  onChange={(e) => setDiskGb(Number(e.target.value))}
                  className="mt-2 h-2 w-full cursor-pointer appearance-none rounded-lg bg-muted accent-primary"
                />
                <span className="mt-1 block text-muted-foreground text-xs">
                  Max {selectedTier.pool.diskGb} GB
                </span>
              </div>
            </div>

            {pool && !pool.fits && (
              <div
                role="alert"
                className="mt-4 flex flex-wrap items-center justify-between gap-3 rounded-xl bg-amber-500/10 p-3 text-amber-600 text-xs dark:text-amber-400"
              >
                <span>Selected allocation exceeds Tier {selectedTier.level} limits.</span>
                {largerTier && (
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() => setTierId(largerTier.id)}
                    className="h-7 text-xs"
                  >
                    Switch to Tier {largerTier.level}
                  </Button>
                )}
              </div>
            )}
          </section>

          <section
            aria-labelledby={`${id}-packages`}
            className="rounded-2xl border border-border bg-card p-6 shadow-sm"
          >
            <h3 id={`${id}-packages`} className="font-semibold text-foreground text-base">
              Credit deposit packages
            </h3>
            <p className="mt-1 text-muted-foreground text-xs">
              Credits never expire. Top up when balance runs low. $1.00 = {pricing.creditsPerDollar}{' '}
              credits.
            </p>

            <div className="mt-4 grid gap-3 sm:grid-cols-5">
              {creditPackagesCents.map((amount) => {
                const isSelected = packageCents === amount;
                const creditsAmt = (amount / 100) * pricing.creditsPerDollar;
                return (
                  <button
                    key={amount}
                    type="button"
                    aria-pressed={isSelected}
                    onClick={() => setPackageCents(amount)}
                    className={`flex flex-col items-center justify-center rounded-xl border p-4 transition-all ${
                      isSelected
                        ? 'border-primary bg-primary/5 shadow-xs ring-1 ring-primary'
                        : 'border-border/60 bg-muted/20 hover:bg-muted/40'
                    }`}
                  >
                    <span className="font-bold text-foreground text-lg">
                      {formatMoney(amount, 0)}
                    </span>
                    <span className="mt-0.5 text-muted-foreground text-xs">
                      {formatNumber(creditsAmt)} pts
                    </span>
                  </button>
                );
              })}
            </div>

            {packageTier && (
              <p className="mt-4 font-medium text-foreground text-xs">
                Package automatically unlocks prospective Tier {packageTier.level} capacity.
              </p>
            )}
            {packageCents < selectedTier.minimumFundingCents && (
              <p className="mt-2 text-amber-600 text-xs dark:text-amber-400">
                Notice: Tier {selectedTier.level} typically requires a minimum funding deposit of{' '}
                {formatMoney(selectedTier.minimumFundingCents)}.
              </p>
            )}
          </section>
        </div>

        <PaygEstimateCard
          id={id}
          estimate={estimate}
          cores={cores}
          ramGb={ramGb}
          diskGb={diskGb}
          packageCents={packageCents}
          balanceCredits={balanceCredits}
          fits={Boolean(pool?.fits)}
          busy={busy}
          onAddCredits={onAddCredits}
        />
      </div>
    </div>
  );
}
