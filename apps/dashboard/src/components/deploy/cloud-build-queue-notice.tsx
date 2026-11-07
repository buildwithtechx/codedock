import { Button } from '@codedock/ui/button';
import { cn } from '#/lib/utils';

export interface CloudBuildQueueStatus {
  position: number | null;
  reason: 'concurrency' | 'resources' | 'server';
  simultaneousBuilds: number | null;
  activeBuilds: number;
  workspaceId: string;
}

export interface CloudBuildQueueNoticeProps {
  queue: CloudBuildQueueStatus;
  onComparePlans?: (workspaceId: string, reason: string) => void;
  className?: string;
}

export function CloudBuildQueueNotice({
  queue,
  onComparePlans,
  className,
}: CloudBuildQueueNoticeProps) {
  const number = new Intl.NumberFormat('en');

  const reasonCopy = {
    concurrency:
      'Waiting for concurrency slot. Build will start automatically once a slot is freed.',
    resources:
      'Waiting for available build resources. Build will start automatically once capacity is ready.',
    server: 'Waiting for build server. Build will start automatically once a worker is assigned.',
  };

  const reason = reasonCopy[queue.reason] || reasonCopy.concurrency;

  const canComparePlans =
    queue.reason === 'resources' ||
    (queue.reason === 'concurrency' && queue.simultaneousBuilds !== null);

  const title =
    queue.position === null ? 'Build queued' : `Position ${number.format(queue.position)} in queue`;

  const allowanceText =
    queue.simultaneousBuilds === null
      ? 'No simultaneous build limit.'
      : `Simultaneous builds allowed by your plan: ${number.format(queue.simultaneousBuilds)}.`;

  const activeText = `Active builds: ${number.format(queue.activeBuilds)}.`;

  return (
    <section
      className={cn('rounded-2xl border border-border/50 bg-card p-5 text-left', className)}
      aria-label="Cloud build queue notice"
    >
      <div role="status" aria-live="polite">
        <p className="font-medium text-foreground text-sm">{title}</p>
        <p className="mt-1 text-muted-foreground text-sm">{reason}</p>
        <p className="mt-2 text-muted-foreground text-sm">
          {allowanceText} {activeText}
        </p>
      </div>
      {canComparePlans && (
        <div className="mt-3 flex flex-wrap items-center gap-3">
          <Button
            variant="secondary"
            size="sm"
            onClick={() => onComparePlans?.(queue.workspaceId, queue.reason)}
          >
            Compare plans
          </Button>
          <span className="text-muted-foreground text-xs">Upgrading is optional</span>
        </div>
      )}
    </section>
  );
}
