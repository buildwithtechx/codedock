import { CheckCircle2, Clock, Loader2, RotateCcw, ShieldAlert } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import type { NetworkRollbackWatchdog as WatchdogType } from '#/interfaces/networks';
import { networksService } from '#/services/networks';

interface NetworkRollbackWatchdogProps {
  watchdog: WatchdogType;
  onCommitted?: () => void;
  onRolledBack?: () => void;
}

export function NetworkRollbackWatchdog({
  watchdog: initialWatchdog,
  onCommitted,
  onRolledBack,
}: NetworkRollbackWatchdogProps) {
  const [watchdog, setWatchdog] = useState<WatchdogType>(initialWatchdog);
  const [remaining, setRemaining] = useState<number>(initialWatchdog.remainingSeconds);
  const [isCommitting, setIsCommitting] = useState(false);
  const [isRollingBack, setIsRollingBack] = useState(false);

  useEffect(() => {
    setWatchdog(initialWatchdog);
    setRemaining(initialWatchdog.remainingSeconds);
  }, [initialWatchdog]);

  useEffect(() => {
    if (!watchdog.armed || remaining <= 0) return;

    const timer = setInterval(() => {
      setRemaining((prev) => {
        if (prev <= 1) {
          clearInterval(timer);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [watchdog.armed, remaining]);

  const handleCommit = async () => {
    setIsCommitting(true);
    try {
      const updated = await networksService.commitWatchdog(watchdog.operationId);
      setWatchdog(updated);
      toast.success('Network changes committed permanently');
      onCommitted?.();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to commit network changes');
    } finally {
      setIsCommitting(false);
    }
  };

  const handleRollback = async () => {
    setIsRollingBack(true);
    try {
      const updated = await networksService.rollbackWatchdog(watchdog.operationId);
      setWatchdog(updated);
      toast.warning('Emergency network rollback executed');
      onRolledBack?.();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to rollback network changes');
    } finally {
      setIsRollingBack(false);
    }
  };

  const percentage = Math.max(
    0,
    Math.min(100, Math.round((remaining / (watchdog.timeoutSeconds || 120)) * 100))
  );

  return (
    <div
      className={`rounded-2xl border p-5 ${
        watchdog.status === 'armed'
          ? 'border-amber-500/30 bg-amber-500/5'
          : watchdog.status === 'committed'
            ? 'border-emerald-500/30 bg-emerald-500/5'
            : 'border-destructive/30 bg-destructive/5'
      }`}
      data-testid="network-rollback-watchdog"
    >
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <div
            className={`flex h-10 w-10 items-center justify-center rounded-xl ${
              watchdog.status === 'armed'
                ? 'bg-amber-500/20 text-amber-500'
                : watchdog.status === 'committed'
                  ? 'bg-emerald-500/20 text-emerald-500'
                  : 'bg-destructive/20 text-destructive'
            }`}
          >
            {watchdog.status === 'armed' ? (
              <Clock className="size-5 animate-pulse" />
            ) : watchdog.status === 'committed' ? (
              <CheckCircle2 className="size-5" />
            ) : (
              <ShieldAlert className="size-5" />
            )}
          </div>
          <div>
            <h3 className="font-semibold text-base text-foreground">
              {watchdog.status === 'armed'
                ? 'Rollback Watchdog Armed (120s)'
                : watchdog.status === 'committed'
                  ? 'Network Changes Committed'
                  : 'Network Rolled Back'}
            </h3>
            <p className="text-muted-foreground text-xs">
              {watchdog.status === 'armed'
                ? 'Automatic fallback armed. If connectivity drops or verification is unconfirmed, network reverts to avoid host isolation.'
                : watchdog.status === 'committed'
                  ? 'Nodes acknowledged verified reachability. Firewall and WireGuard routes are permanent.'
                  : 'Previous host configuration restored on all member nodes.'}
            </p>
          </div>
        </div>

        {watchdog.status === 'armed' && (
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={isRollingBack || isCommitting}
              onClick={handleRollback}
              className="gap-1.5 border-destructive/40 text-destructive hover:bg-destructive/10"
            >
              {isRollingBack ? (
                <Loader2 className="size-3.5 animate-spin" />
              ) : (
                <RotateCcw className="size-3.5" />
              )}
              Emergency Rollback
            </Button>
            <Button
              size="sm"
              disabled={isCommitting || isRollingBack}
              onClick={handleCommit}
              className="gap-1.5 bg-emerald-600 text-white hover:bg-emerald-700"
            >
              {isCommitting ? (
                <Loader2 className="size-3.5 animate-spin" />
              ) : (
                <CheckCircle2 className="size-3.5" />
              )}
              Commit Changes
            </Button>
          </div>
        )}
      </div>

      {watchdog.status === 'armed' && (
        <div className="mt-4 space-y-2">
          <div className="flex justify-between font-mono text-muted-foreground text-xs">
            <span>Watchdog countdown</span>
            <span className="font-semibold text-foreground">{remaining}s remaining</span>
          </div>
          <div className="h-2 w-full overflow-hidden rounded-full bg-muted">
            <div
              className={`h-full transition-all duration-1000 ${
                remaining < 30 ? 'bg-destructive' : 'bg-amber-500'
              }`}
              style={{ width: `${percentage}%` }}
              role="progressbar"
              aria-valuenow={remaining}
              aria-valuemax={watchdog.timeoutSeconds || 120}
            />
          </div>
        </div>
      )}
    </div>
  );
}
