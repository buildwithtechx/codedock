import { CheckCircle2, Loader2, MapPin } from 'lucide-react';
import type { ManagedOffers } from '#/interfaces/managed';
import { cn } from '#/lib/utils';

function priceFor(offer: ManagedOffers['offers'][number], location: string): number | null {
  const live = offer.prices[location];
  if (typeof live === 'number') return live;
  return null;
}

export function HetznerPlanStep({
  offers,
  loading,
  loadError,
  location,
  tierName,
  onLocationChange,
  onTierChange,
}: {
  offers: ManagedOffers | null;
  loading: boolean;
  loadError: string | null;
  location: string;
  tierName: string;
  onLocationChange: (location: string) => void;
  onTierChange: (tierName: string) => void;
}) {
  if (loading) {
    return (
      <div className="flex items-center gap-3 rounded-xl border border-border/60 p-6 text-muted-foreground text-sm">
        <Loader2 className="size-4 animate-spin" />
        Loading live Hetzner locations and server types…
      </div>
    );
  }

  if (loadError || !offers) {
    return (
      <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-4 text-destructive text-xs">
        {loadError ?? 'Hetzner offers are unavailable for this credential.'}
      </div>
    );
  }

  return (
    <div className="space-y-5">
      <div>
        <h3 className="font-medium text-foreground text-sm">Location</h3>
        <div className="mt-2 grid grid-cols-2 gap-2 sm:grid-cols-3">
          {offers.locations.map((entry) => {
            const selected = entry.name === location;
            return (
              <button
                key={entry.name}
                type="button"
                onClick={() => onLocationChange(entry.name)}
                className={cn(
                  'flex items-center gap-2.5 rounded-xl border px-3 py-2.5 text-left transition-colors',
                  selected
                    ? 'border-primary bg-primary/5'
                    : 'border-border/70 hover:border-primary/30 hover:bg-muted/30'
                )}
              >
                <MapPin
                  className={cn(
                    'size-4 shrink-0',
                    selected ? 'text-primary' : 'text-muted-foreground'
                  )}
                />
                <span className="min-w-0">
                  <span className="block font-medium font-mono text-foreground text-xs">
                    {entry.name}
                  </span>
                  <span className="block truncate text-[11px] text-muted-foreground">
                    {entry.city}, {entry.country}
                  </span>
                </span>
              </button>
            );
          })}
        </div>
      </div>

      <div>
        <h3 className="font-medium text-foreground text-sm">Server type</h3>
        <p className="mt-1 text-muted-foreground text-xs">
          Prices are live Hetzner Cloud list prices for {location || 'the selected location'}.
        </p>
        <div className="mt-2 space-y-2">
          {offers.offers.map((offer) => {
            const selected = offer.tier.name === tierName;
            const live = location ? priceFor(offer, location) : null;
            return (
              <button
                key={offer.tier.name}
                type="button"
                onClick={() => onTierChange(offer.tier.name)}
                className={cn(
                  'flex w-full items-center gap-4 rounded-xl border px-4 py-3 text-left transition-colors',
                  selected
                    ? 'border-primary bg-primary/5'
                    : 'border-border/70 hover:border-primary/30 hover:bg-muted/30'
                )}
              >
                <span className="min-w-0 flex-1">
                  <span className="block font-medium font-mono text-foreground text-sm">
                    {offer.tier.name}
                  </span>
                  <span className="block text-muted-foreground text-xs">
                    {offer.tier.cores} vCPU · {offer.tier.memoryGB} GB RAM · {offer.tier.diskGB} GB
                    SSD
                  </span>
                </span>
                <span className="text-right">
                  <span className="block font-medium text-foreground text-sm">
                    €{(live ?? offer.tier.monthlyPrice).toFixed(2)}/mo
                  </span>
                  <span className="block text-[11px] text-muted-foreground">
                    {live === null ? 'estimated' : 'live price'}
                  </span>
                </span>
                {selected && <CheckCircle2 className="size-4 shrink-0 text-primary" />}
              </button>
            );
          })}
        </div>
      </div>
    </div>
  );
}
