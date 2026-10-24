import type { RollbackCapacityUI } from '@codedock/contracts';
import { DEFAULT_ROLLBACK_WINDOW, MAX_ROLLBACK_WINDOW, formatBytes } from '@codedock/core';
import { Loader2, RotateCcw } from 'lucide-react';
import { InfoCard } from './info-card';

export { DEFAULT_ROLLBACK_WINDOW, MAX_ROLLBACK_WINDOW };

export interface RollbackRetentionCardsProps {
  strategy: 'git' | 'snapshot';
  artifactKind?: 'image' | 'files';
  capacity: RollbackCapacityUI | null;
  onToggleStrategy?: () => void;
  onChangeWindow?: (next: number | null) => void;
  togglingStrategy?: boolean;
  savingWindow?: boolean;
  readOnly?: boolean;
}

export function RollbackRetentionCards({
  strategy,
  capacity,
  artifactKind = 'image',
  onToggleStrategy,
  onChangeWindow,
  togglingStrategy = false,
  savingWindow = false,
  readOnly = false,
}: RollbackRetentionCardsProps) {
  const isSnapshot = strategy === 'snapshot';
  const isFiles = artifactKind === 'files';
  const windowVal = capacity?.window ?? DEFAULT_ROLLBACK_WINDOW;
  const maxWindow = capacity?.maxWindow ?? MAX_ROLLBACK_WINDOW;
  const canEdit = !readOnly && !!capacity && !!onChangeWindow;
  const busy = savingWindow || togglingStrategy;
  const automatic = capacity?.explicit === null;

  const measured: string[] = [isFiles ? 'built files on disk' : 'images on the host'];
  if (capacity?.snapshotSizeBytes) {
    measured.push(`~${formatBytes(capacity.snapshotSizeBytes)} per image, including shared layers`);
  }
  if (capacity?.diskFreeBytes) {
    measured.push(`${formatBytes(capacity.diskFreeBytes)} free`);
  }
  if (capacity?.source === 'instance-default') {
    measured.push('Automatic', 'instance default');
  }
  if (capacity?.source === 'explicit') {
    measured.push('set manually');
  }

  const historyDescription = automatic
    ? 'Keeps up to the default history limit. Under disk pressure, recoverable artifacts may be removed earlier. Active and pinned releases are extra. Commit and configuration history stays available; rollback may require a rebuild.'
    : 'Maximum number of past releases to keep. Active and pinned releases are extra. Older artifacts are removed; deployment history stays available. Set to 0 to keep only active and pinned releases.';

  const windowText = automatic
    ? `Up to ${windowVal} versions`
    : windowVal === 1
      ? '1 version'
      : `${windowVal} versions`;

  return (
    <div className="grid gap-3">
      <InfoCard
        icon={RotateCcw}
        title="Rollback strategy"
        value={
          capacity ? (isSnapshot ? 'Preserve rollback artifacts' : 'Reuse images or rebuild') : '—'
        }
        description={
          isSnapshot
            ? isFiles
              ? 'Keeps built files within the history limit. Older versions may need to be rebuilt.'
              : 'Keeps rollback artifacts within the history limit. Older versions may need to be rebuilt.'
            : isFiles
              ? 'Rebuilds the old commit, unless its files are still on the host.'
              : 'Rebuilds the old commit, unless its image is still on the host.'
        }
        action={
          readOnly || !onToggleStrategy ? undefined : (
            <button
              type="button"
              role="switch"
              aria-checked={isSnapshot}
              onClick={onToggleStrategy}
              disabled={busy || !capacity}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                isSnapshot ? 'bg-primary' : 'bg-muted'
              } ${togglingStrategy ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'}`}
              aria-label="Toggle snapshot retention"
            >
              {togglingStrategy ? (
                <span className="mx-auto">
                  <Loader2 className="size-3.5 animate-spin text-background" />
                </span>
              ) : (
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-background transition-transform ${
                    isSnapshot ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              )}
            </button>
          )
        }
      />
      <InfoCard
        icon={RotateCcw}
        title="Rollback history"
        value={capacity ? windowText : '—'}
        description={historyDescription}
        footer={
          <p className="text-[11px] text-muted-foreground/80 leading-relaxed">
            {measured.join(' · ')}
          </p>
        }
        action={
          canEdit ? (
            <div className="flex flex-wrap items-center justify-end gap-1.5">
              {!automatic && (
                <button
                  type="button"
                  disabled={busy}
                  onClick={() => onChangeWindow?.(null)}
                  className="rounded-md px-2 py-1 text-muted-foreground text-xs hover:bg-muted hover:text-foreground disabled:opacity-40"
                >
                  Automatic
                </button>
              )}
              <button
                type="button"
                onClick={() => onChangeWindow?.(windowVal - 1)}
                disabled={busy || windowVal <= 0}
                className="flex h-6 w-6 items-center justify-center rounded-md border border-border/60 text-foreground transition-colors enabled:hover:bg-muted disabled:cursor-not-allowed disabled:opacity-40"
                aria-label="Decrease rollback history"
              >
                −
              </button>
              <span className="w-5 text-center font-medium text-[13px] text-foreground tabular-nums">
                {savingWindow ? (
                  <Loader2 className="mx-auto size-3.5 animate-spin text-muted-foreground" />
                ) : (
                  windowVal
                )}
              </span>
              <button
                type="button"
                onClick={() => onChangeWindow?.(windowVal + 1)}
                disabled={busy || windowVal >= maxWindow}
                className="flex h-6 w-6 items-center justify-center rounded-md border border-border/60 text-foreground transition-colors enabled:hover:bg-muted disabled:cursor-not-allowed disabled:opacity-40"
                aria-label="Increase rollback history"
              >
                +
              </button>
            </div>
          ) : undefined
        }
      />
    </div>
  );
}

export default RollbackRetentionCards;
