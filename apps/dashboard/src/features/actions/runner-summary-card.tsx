import { PlayCircle } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import { Checkbox } from '@codedock/ui/checkbox';

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
  isCloud?: boolean;
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
  isCloud = false,
}: RunnerSummaryCardProps) {
  return (
    <aside className="space-y-5 rounded-2xl bg-card p-5">
      <div className="flex items-center gap-3">
        <span className="flex size-7 shrink-0 items-center justify-center rounded-full text-sky-400">
          <PlayCircle className="size-6" />
        </span>
        <div className="min-w-0">
          <h2 className="truncate font-medium text-base text-foreground">
            {title || 'Add runner'}
          </h2>
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
          <div className="flex flex-wrap gap-1.5" dir="ltr">
            {matchingLabels.map((l) => (
              <span
                key={l}
                className="rounded-md bg-muted px-2 py-1 font-mono text-xs text-muted-foreground"
              >
                {l}
              </span>
            ))}
          </div>
        </div>
      )}

      <label className="flex cursor-pointer items-center gap-2 pt-1 text-sm">
        <Checkbox checked={enabled} onCheckedChange={(val) => onEnabledChange(Boolean(val))} />
        <span className="font-medium text-foreground text-sm">Accept new jobs</span>
      </label>

      <Button
        type="submit"
        disabled={isPending || !canSubmit}
        className="h-10 w-full rounded-xl bg-muted/80 text-foreground/80 hover:bg-muted hover:text-foreground disabled:cursor-not-allowed disabled:opacity-50"
      >
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
