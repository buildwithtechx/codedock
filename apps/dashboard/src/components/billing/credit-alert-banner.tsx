import { AlertCircle, ArrowUpRight } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import { cn } from '#/lib/utils';

export interface CreditAlertInfo {
  state: 'low' | 'grace' | 'depleted' | 'ok' | string;
  percent?: number;
  remaining?: number;
  balance?: number;
  limit?: number;
}

export interface CreditAlertBannerProps {
  alert?: CreditAlertInfo | null;
  billingMode?: 'monthly' | 'payg' | string;
  onTopUp?: () => void;
  onViewPlans?: () => void;
  className?: string;
}

export function CreditAlertBanner({
  alert,
  billingMode,
  onTopUp,
  onViewPlans,
  className,
}: CreditAlertBannerProps) {
  if (!alert || billingMode === 'monthly' || !['low', 'grace', 'depleted'].includes(alert.state)) {
    return null;
  }

  const isDepleted = alert.state === 'depleted';
  const isGrace = alert.state === 'grace';

  const title = isDepleted
    ? 'Build credits exhausted'
    : isGrace
      ? 'Grace period active'
      : 'Build credits running low';

  const description = isDepleted
    ? 'You have exhausted your cloud build credits. New builds and deployments are paused until credits are added.'
    : isGrace
      ? 'Your balance is in grace credit usage. Top up credits soon to avoid deployment interruptions.'
      : `Your balance is below 15% (${alert.percent ?? 15}% remaining). Top up to keep automated builds running smoothly.`;

  return (
    <section
      aria-label="Cloud credit alert"
      className={cn(
        'flex flex-wrap items-center justify-between gap-4 rounded-2xl border p-4 text-left transition-colors',
        isDepleted
          ? 'border-destructive/30 bg-destructive/10 text-destructive-foreground'
          : 'border-amber-500/30 bg-amber-500/10 text-foreground',
        className
      )}
    >
      <div className="flex min-w-0 flex-1 items-start gap-3">
        <AlertCircle
          className={cn(
            'mt-0.5 size-5 shrink-0',
            isDepleted ? 'text-destructive' : 'text-amber-500'
          )}
          aria-hidden="true"
        />
        <div role="status">
          <h4 className="font-semibold text-foreground text-sm">{title}</h4>
          <p className="mt-0.5 text-muted-foreground text-xs leading-relaxed">{description}</p>
        </div>
      </div>
      <div className="flex items-center gap-2">
        {onTopUp && (
          <Button variant={isDepleted ? 'destructive' : 'default'} size="sm" onClick={onTopUp}>
            Add credits
          </Button>
        )}
        {onViewPlans && (
          <Button variant="outline" size="sm" onClick={onViewPlans}>
            <span>View plans</span>
            <ArrowUpRight className="size-3.5" aria-hidden="true" />
          </Button>
        )}
      </div>
    </section>
  );
}
