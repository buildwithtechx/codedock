import { useRef, useState } from 'react';
import { Coins, Loader2, Zap } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { DEFAULT_TOPUP_PACKS, type TopupPack, billingService } from '#/services/billing';

export interface BillingTopupsProps {
  balanceCredits?: number;
  packs?: TopupPack[];
  disabled?: boolean;
  onBuy?: (pack: TopupPack, idempotencyKey: string) => Promise<void> | void;
}

export function BillingTopups({
  balanceCredits,
  packs = DEFAULT_TOPUP_PACKS,
  disabled = false,
  onBuy,
}: BillingTopupsProps) {
  const [buyingPackId, setBuyingPackId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const idempotencyKeysRef = useRef<Map<string, string>>(new Map());

  const handleBuy = async (pack: TopupPack) => {
    if (disabled || buyingPackId) return;

    if (!idempotencyKeysRef.current.has(pack.id)) {
      idempotencyKeysRef.current.set(pack.id, `topup_${pack.id}_${Date.now()}`);
    }
    const idempotencyKey = idempotencyKeysRef.current.get(pack.id)!;

    setBuyingPackId(pack.id);
    setError(null);

    try {
      if (onBuy) {
        await onBuy(pack, idempotencyKey);
      } else {
        const { url } = await billingService.purchaseTopup({
          packId: pack.id,
          idempotencyKey,
        });
        if (url) window.location.href = url;
      }
    } catch {
      setError(`Failed to initiate checkout for ${pack.name}. Please try again.`);
    } finally {
      setBuyingPackId(null);
    }
  };

  return (
    <div className="space-y-4 rounded-2xl border border-border bg-card p-6 shadow-sm">
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h3 className="flex items-center gap-2 font-semibold text-foreground text-lg">
            <Coins className="size-5 text-primary" />
            Prepaid Compute Credits
          </h3>
          <p className="mt-1 text-muted-foreground text-sm">
            Top up on-demand credits for managed server hours, storage growth, and build concurrency.
          </p>
        </div>

        {balanceCredits !== undefined && (
          <div className="rounded-xl border border-border bg-muted/40 px-3.5 py-2">
            <p className="font-medium text-muted-foreground text-xs">Current Balance</p>
            <p className="font-bold text-foreground text-lg">
              {balanceCredits.toLocaleString()} <span className="font-normal text-muted-foreground text-xs">credits</span>
            </p>
          </div>
        )}
      </div>

      {error && (
        <div className="rounded-lg bg-destructive/10 p-3 text-destructive text-sm" role="alert">
          {error}
        </div>
      )}

      <div className="grid gap-3.5 sm:grid-cols-2 lg:grid-cols-5">
        {packs.map((pack) => {
          const isBuying = buyingPackId === pack.id;
          const priceFormatted = `$${(pack.priceCents / 100).toFixed(0)}`;

          return (
            <div
              key={pack.id}
              className="relative flex flex-col justify-between rounded-xl border border-border bg-muted/20 p-4 transition-colors hover:border-primary/40 hover:bg-muted/40"
            >
              <div>
                <div className="flex items-center justify-between">
                  <span className="font-bold text-foreground text-xl">{priceFormatted}</span>
                  {pack.bonusPercent && (
                    <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-emerald-500 text-xs">
                      <Zap className="size-3" />
                      +{pack.bonusPercent}%
                    </span>
                  )}
                </div>
                <p className="mt-2 font-medium text-foreground text-sm">
                  {pack.credits.toLocaleString()} credits
                </p>
                <p className="text-muted-foreground text-xs">{pack.name}</p>
              </div>

              <div className="mt-4">
                <Button
                  size="sm"
                  variant="outline"
                  className="w-full text-xs"
                  disabled={disabled || Boolean(buyingPackId)}
                  onClick={() => handleBuy(pack)}
                >
                  {isBuying ? (
                    <Loader2 className="size-3.5 animate-spin" />
                  ) : (
                    `Buy ${priceFormatted}`
                  )}
                </Button>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
