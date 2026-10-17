import { useId } from 'react';
import type { PaygTier } from '#/constants/billing';

export interface PaygResourceTiersProps {
  tiers: PaygTier[];
  selected: PaygTier;
  onChange: (id: string) => void;
  formatMoney: (cents: number) => string;
}

export function PaygResourceTiers({
  tiers,
  selected,
  onChange,
  formatMoney,
}: PaygResourceTiersProps) {
  const id = useId();

  const formatNumber = (value: number) => new Intl.NumberFormat('en-US').format(value);

  const fields = [
    { key: 'cpuCores' as const, label: 'Dedicated vCPU', divisor: 1, unit: 'vCPU' },
    { key: 'memoryMb' as const, label: 'Total Memory', divisor: 1024, unit: 'GB' },
    { key: 'diskGb' as const, label: 'Fast Storage', divisor: 1, unit: 'GB' },
    { key: 'servers' as const, label: 'Max Server Hosts', divisor: 1, unit: '' },
  ];

  return (
    <section aria-label="Resource pool tier capacity" className="min-w-0 space-y-3">
      <div className="flex flex-wrap items-center justify-between gap-x-3 gap-y-1">
        <h4 className="font-medium text-foreground text-sm">Resource pool limits</h4>
        <div role="tablist" aria-label="Pool tiers" className="inline-flex rounded-lg bg-muted p-1">
          {tiers.map((tier) => {
            const isSelected = tier.id === selected.id;
            return (
              <button
                key={tier.id}
                role="tab"
                id={`${id}-tab-${tier.id}`}
                aria-selected={isSelected}
                aria-controls={`${id}-panel-${tier.id}`}
                tabIndex={isSelected ? 0 : -1}
                type="button"
                onClick={() => onChange(tier.id)}
                className={`rounded-md px-3 py-1 font-medium text-xs transition-colors ${
                  isSelected
                    ? 'bg-card text-foreground shadow-xs'
                    : 'text-muted-foreground hover:text-foreground'
                }`}
              >
                Tier {tier.level}
              </button>
            );
          })}
        </div>
      </div>

      <p className="text-muted-foreground text-xs">
        Unlocks with deposit of {formatMoney(selected.minimumFundingCents)} or greater.
      </p>

      {tiers.map((tier) => (
        <div
          key={tier.id}
          id={`${id}-panel-${tier.id}`}
          role="tabpanel"
          aria-labelledby={`${id}-tab-${tier.id}`}
          hidden={tier.id !== selected.id}
        >
          <dl className="grid grid-cols-2 gap-x-4 gap-y-3 rounded-xl bg-muted/40 p-3.5 sm:grid-cols-4">
            {fields.map(({ key, label, divisor, unit }) => (
              <div key={key} className="min-w-0">
                <dt className="text-muted-foreground text-xs">{label}</dt>
                <dd className="mt-0.5 font-semibold text-foreground text-sm tabular-nums">
                  {formatNumber(tier.pool[key] / divisor)}
                  {unit ? ` ${unit}` : ''}
                </dd>
              </div>
            ))}
          </dl>
        </div>
      ))}

      <p className="text-muted-foreground text-xs">
        Each tier defines prospective cluster limits. Hourly usage is metered continuously per host.
      </p>
    </section>
  );
}
