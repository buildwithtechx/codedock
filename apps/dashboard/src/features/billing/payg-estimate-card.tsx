import { ChevronDown, Server } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import type { PrepaidUsageEstimate } from '#/interfaces/billing';
import { formatDuration, formatMoney, formatNumber } from '#/services/billing-estimate';

export interface PaygEstimateCardProps {
  id: string;
  estimate: PrepaidUsageEstimate | null;
  cores: number;
  ramGb: number;
  diskGb: number;
  packageCents: number;
  balanceCredits: number;
  fits: boolean;
  busy?: boolean;
  onAddCredits?: (cents: number) => void;
}

export function PaygEstimateCard({
  id,
  estimate,
  cores,
  ramGb,
  diskGb,
  packageCents,
  balanceCredits,
  fits,
  busy = false,
  onAddCredits,
}: PaygEstimateCardProps) {
  return (
    <aside
      aria-labelledby={`${id}-burn`}
      className="space-y-5 rounded-2xl border border-border bg-card p-6 shadow-sm"
    >
      <div>
        <h3 id={`${id}-burn`} className="font-semibold text-base text-foreground">
          Estimated Burn Rate
        </h3>
        <div className="mt-3">
          <span className="font-bold text-3xl text-foreground tabular-nums">
            {estimate ? formatMoney(estimate.totalHourlyCents, 4) : '—'}
          </span>
          <span className="ml-1.5 text-muted-foreground text-xs">/ hour</span>
        </div>
        <p className="mt-0.5 text-muted-foreground text-xs tabular-nums">
          {estimate ? `${formatNumber(Number(estimate.hourlyCredits.toFixed(2)))} credits/hr` : '—'}
        </p>
      </div>

      <dl className="space-y-3 border-border/50 border-t pt-4 text-xs">
        <div className="flex justify-between text-muted-foreground">
          <span>Active vCPU ({cores} Cores)</span>
          <span className="font-medium text-foreground tabular-nums">
            {estimate ? formatMoney(estimate.hourlyCents.cpu, 4) : '—'}
          </span>
        </div>
        <div className="flex justify-between text-muted-foreground">
          <span>Reserved RAM ({ramGb} GB)</span>
          <span className="font-medium text-foreground tabular-nums">
            {estimate ? formatMoney(estimate.hourlyCents.memory, 4) : '—'}
          </span>
        </div>
        <div className="flex justify-between text-muted-foreground">
          <span>Retained NVMe ({diskGb} GB)</span>
          <span className="font-medium text-foreground tabular-nums">
            {estimate ? formatMoney(estimate.hourlyCents.disk, 4) : '—'}
          </span>
        </div>
      </dl>

      <details className="group border-border/50 border-t pt-3">
        <summary className="flex cursor-pointer items-center justify-between text-muted-foreground text-xs hover:text-foreground">
          <span>Estimated runtime duration</span>
          <ChevronDown className="size-3.5 transition-transform group-open:rotate-180" />
        </summary>
        <div className="mt-3 grid grid-cols-2 gap-2 text-xs">
          <div className="rounded-lg bg-muted/40 p-2.5">
            <span className="text-muted-foreground text-xs">Under 100% Load</span>
            <span className="mt-1 block font-semibold text-foreground tabular-nums">
              {estimate ? formatDuration(estimate.durationRange.minElapsedHours) : '—'}
            </span>
          </div>
          <div className="rounded-lg bg-muted/40 p-2.5">
            <span className="text-muted-foreground text-xs">At Idle (0% CPU)</span>
            <span className="mt-1 block font-semibold text-foreground tabular-nums">
              {estimate ? formatDuration(estimate.durationRange.maxElapsedHours) : '—'}
            </span>
          </div>
        </div>
      </details>

      <div className="border-border/50 border-t pt-4">
        <div className="mb-3 flex items-center justify-between">
          <div>
            <p className="font-semibold text-foreground text-sm">
              {formatMoney(packageCents, 0)} deposit
            </p>
            <p className="text-muted-foreground text-xs">{formatNumber(balanceCredits)} credits</p>
          </div>
        </div>
        <Button
          type="button"
          disabled={busy || !fits}
          onClick={() => onAddCredits?.(packageCents)}
          className="w-full gap-2 rounded-xl"
        >
          <Server className="size-4" />
          Add Credits & Deploy
        </Button>
      </div>
    </aside>
  );
}
