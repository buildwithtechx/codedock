import { PlayCircle } from 'lucide-react';
import { Button } from '#/components/ui/button';

interface RunnerSummaryCardProps {
  title: string;
  presetTitle: string;
  isContainer: boolean;
  cpu: number;
  memoryMb: number;
  maxParallel: number;
  matchingLabels: string[];
  enabled: boolean;
  onEnabledChange: (enabled: boolean) => void;
  isPending: boolean;
  canSubmit: boolean;
  isCloud: boolean;
}

export function RunnerSummaryCard({
  title,
  presetTitle,
  isContainer,
  cpu,
  memoryMb,
  maxParallel,
  matchingLabels,
  enabled,
  onEnabledChange,
  isPending,
  canSubmit,
  isCloud,
}: RunnerSummaryCardProps) {
  return (
    <aside className="space-y-5 rounded-2xl bg-card p-5">
      <div className="flex items-center gap-3">
        <span className="flex size-8 shrink-0 items-center justify-center rounded-full bg-blue-500/15 text-blue-500">
          <PlayCircle className="size-4" />
        </span>
        <div className="min-w-0">
          <h2 className="font-medium text-base text-foreground">{title || 'Add runner'}</h2>
          <p className="mt-0.5 text-muted-foreground text-xs">Runner configuration</p>
        </div>
      </div>

      <dl className="space-y-3 text-sm">
        <div className="flex items-center justify-between gap-3">
          <dt className="text-muted-foreground">Default environment</dt>
          <dd className="font-medium text-foreground">{presetTitle}</dd>
        </div>
        {isContainer && (
          <div className="flex items-center justify-between gap-3">
            <dt className="text-muted-foreground">Per job</dt>
            <dd className="font-medium text-foreground tabular-nums">
              {cpu} vCPU · {memoryMb / 1024} GiB
            </dd>
          </div>
        )}
        <div className="flex items-center justify-between gap-3">
          <dt className="text-muted-foreground">Concurrent jobs</dt>
          <dd className="font-medium text-foreground tabular-nums">{maxParallel}</dd>
        </div>
      </dl>

      {matchingLabels.length > 0 && (
        <div className="space-y-2 border-border/50 border-t pt-4">
          <p className="text-muted-foreground text-xs">Matching runner labels</p>
          <div className="flex flex-wrap gap-1.5">
            {matchingLabels.map((l) => (
              <span
                key={l}
                className="rounded-md bg-muted px-2 py-0.5 font-mono text-[11px] text-muted-foreground"
              >
                {l}
              </span>
            ))}
          </div>
        </div>
      )}

      <label className="flex items-center gap-2 pt-2 text-sm">
        <input
          type="checkbox"
          checked={enabled}
          onChange={(e) => onEnabledChange(e.target.checked)}
          className="size-4 rounded border-border"
        />
        <span className="font-medium text-foreground text-xs">Accept new jobs</span>
      </label>

      <Button type="submit" disabled={isPending || !canSubmit} className="w-full rounded-full">
        {isPending ? 'Verifying runner...' : 'Verify and save runner'}
      </Button>

      <p className="text-muted-foreground text-xs leading-relaxed">
        {isCloud
          ? "Jobs share this Cloud space's capacity. Only run trusted workflows."
          : "Jobs share this server's capacity. Only run trusted workflows."}
      </p>
    </aside>
  );
}
