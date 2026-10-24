import { Clock, RotateCcw, StopCircle, User } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@codedock/ui/card';
import type { ActionRun } from '#/interfaces/actions';
import { useCancelRun, useRetryRun } from '#/hooks/use-actions';
import { ActionLogs } from './action-logs';
import { ActionStatus } from './action-status';

interface RunWorkspaceProps {
  run: ActionRun;
  onRunUpdated?: (run: ActionRun) => void;
}

export function RunWorkspace({ run: initialRun, onRunUpdated }: RunWorkspaceProps) {
  const [run, setRun] = useState<ActionRun>(initialRun);
  const retryMutation = useRetryRun();
  const cancelMutation = useCancelRun();

  const handleRetry = async () => {
    try {
      const updated = await retryMutation.mutateAsync(run.id);
      setRun(updated);
      onRunUpdated?.(updated);
      toast.success(`Run #${updated.number} triggered for workflow ${updated.name}`);
    } catch {
      toast.error('Failed to re-run action');
    }
  };

  const handleCancel = async () => {
    try {
      const updated = await cancelMutation.mutateAsync(run.id);
      setRun(updated);
      onRunUpdated?.(updated);
      toast.info(`Run #${updated.number} cancelled`);
    } catch {
      toast.error('Failed to cancel action run');
    }
  };

  const isRetrying = retryMutation.isPending;
  const isCancelling = cancelMutation.isPending;

  const isPending = run.status === 'running' || run.status === 'queued';

  return (
    <div className="space-y-6">
      <Card className="border-border/70 shadow-sm">
        <CardHeader className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div className="space-y-1">
            <div className="flex items-center gap-3">
              <CardTitle className="font-bold text-xl tracking-tight">
                {run.name} <span className="text-muted-foreground">#{run.number}</span>
              </CardTitle>
              <ActionStatus status={run.status} />
            </div>
            <CardDescription className="flex flex-wrap items-center gap-4 text-xs">
              <span className="flex items-center gap-1 text-muted-foreground">
                <User className="size-3.5" />
                Triggered by <span className="font-medium text-foreground">{run.triggeredBy}</span>
              </span>
              <span className="flex items-center gap-1 text-muted-foreground">
                <Clock className="size-3.5" />
                Duration <span className="font-medium text-foreground">{run.durationSeconds}s</span>
              </span>
              {run.runnerId && (
                <span className="flex items-center gap-1 text-muted-foreground">
                  Runner <span className="font-medium text-foreground">{run.runnerId}</span>
                </span>
              )}
            </CardDescription>
          </div>

          <div className="flex items-center gap-2">
            {isPending && (
              <Button
                variant="destructive"
                size="sm"
                onClick={handleCancel}
                disabled={isCancelling}
                className="gap-1.5"
              >
                <StopCircle className="size-4" />
                {isCancelling ? 'Cancelling...' : 'Cancel Run'}
              </Button>
            )}
            <Button
              variant="outline"
              size="sm"
              onClick={handleRetry}
              disabled={isRetrying}
              className="gap-1.5"
            >
              <RotateCcw className="size-4" />
              {isRetrying ? 'Retrying...' : 'Re-run'}
            </Button>
          </div>
        </CardHeader>

        <CardContent className="space-y-6">
          {run.jobs && run.jobs.length > 0 && (
            <div className="space-y-3">
              <h4 className="font-semibold text-muted-foreground text-xs uppercase tracking-wider">
                Jobs & Steps
              </h4>
              <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
                {run.jobs.map((job) => (
                  <div
                    key={job.id}
                    className="rounded-lg border border-border/70 bg-card/60 p-3 text-sm"
                  >
                    <div className="mb-2 flex items-center justify-between">
                      <span className="font-medium">{job.name}</span>
                      <ActionStatus status={job.status} />
                    </div>
                    {job.steps && job.steps.length > 0 && (
                      <ul className="space-y-1.5 text-muted-foreground text-xs">
                        {job.steps.map((step) => (
                          <li key={step.id} className="flex items-center justify-between">
                            <span className="truncate">{step.name}</span>
                            <span className="shrink-0 font-mono text-[11px]">
                              {step.durationSeconds ?? 0}s
                            </span>
                          </li>
                        ))}
                      </ul>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}

          <div className="space-y-2">
            <h4 className="font-semibold text-muted-foreground text-xs uppercase tracking-wider">
              Execution Log
            </h4>
            <ActionLogs logs={run.logs} isRunning={run.status === 'running'} />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
