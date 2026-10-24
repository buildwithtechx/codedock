import { Loader2, RefreshCw, ShieldAlert } from 'lucide-react';
import type { RuntimeImageCheckResult } from '@codedock/contracts';
import { Button } from '@codedock/ui/button';

export function RuntimeImageUpdateNotice({
  result,
  busy,
  onRefresh,
}: {
  result: RuntimeImageCheckResult | null;
  busy: boolean;
  onRefresh: () => void | Promise<void>;
}) {
  const changed = result?.images.filter((img) => img.status === 'update_available') ?? [];
  if (changed.length === 0) return null;

  return (
    <div
      role="status"
      className="flex flex-col gap-3 rounded-2xl border border-amber-500/30 bg-amber-500/10 p-4 text-amber-900 dark:text-amber-200 sm:flex-row sm:items-center sm:justify-between"
    >
      <div className="flex items-start gap-3">
        <span className="mt-0.5 flex size-8 shrink-0 items-center justify-center rounded-xl bg-amber-500/20 text-amber-600 dark:text-amber-400">
          <ShieldAlert className="size-4" />
        </span>
        <div className="space-y-1">
          <h4 className="font-semibold text-sm">
            A newer base image is available with security patches
          </h4>
          <div className="space-y-1 text-xs">
            {changed.map((img) => (
              <div
                key={`${img.serviceId}:${img.ref}`}
                className="flex flex-wrap items-center gap-1.5 font-mono"
              >
                <span className="font-medium text-foreground">
                  {img.serviceName ? `${img.serviceName}: ` : ''}
                  {img.ref}
                </span>
                {img.currentDigest && img.latestDigest && (
                  <span className="text-muted-foreground text-[11px]">
                    ({img.currentDigest.slice(0, 16)} → {img.latestDigest.slice(0, 16)})
                  </span>
                )}
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="shrink-0 self-end sm:self-center">
        <Button
          size="sm"
          disabled={busy}
          onClick={() => void onRefresh()}
          className="gap-1.5 bg-amber-600 hover:bg-amber-700 text-white shadow-xs"
        >
          {busy ? (
            <Loader2 className="size-3.5 animate-spin" />
          ) : (
            <RefreshCw className="size-3.5" />
          )}
          Redeploy with latest image
        </Button>
      </div>
    </div>
  );
}
