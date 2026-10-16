import { AlertCircle, CheckCircle, Clock, Play, RotateCcw, StopCircle, User } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { actionsService, type ActionRun } from '#/services/actions';
import { ActionLogs } from './action-logs';
import { ActionStatus } from './action-status';

interface RunWorkspaceProps {
  run: ActionRun;
  onRunUpdated?: (run: ActionRun) => void;
}

export function RunWorkspace({ run: initialRun, onRunUpdated }: RunWorkspaceProps) {
  const [run, setRun] = useState<ActionRun>(initialRun);
  const [isRetrying, setIsRetrying] = useState(false);
  const [isCancelling, setIsCancelling] = useState(false);

  const handleRetry = async () => {
    setIsRetrying(true);
    try {
      const updated = await actionsService.retryRun(run.id);
      setRun(updated);
      onRunUpdated?.(updated);
      toast.success(`Run #${updated.number} triggered for workflow ${updated.name}`);
    } catch {
      toast.error('Failed to re-run action');
    } finally {
      setIsRetrying(false);
    }
  };

  const handleCancel = async () => {
    setIsCancelling(true);
    try {
      const updated = await actionsService.cancelRun(run.id);
      setRun(updated);
      onRunUpdated?.(updated);
      toast.info(`Run #${updated.number} cancelled`);
    } catch {
      toast.error('Failed to cancel action run');
    } finally {
      setIsCancelling(false);
    }
  };

  const isPending = run.status === 'running' || run.status === 'queued';

  return (
    <div className="space-y-6">
      <Card className="border-border/70 shadow-sm">
        <CardHeader className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div className="space-y-1">
            <div className="flex items-center gap-3">
              <CardTitle className="text-xl font-bold tracking-tight">
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
              <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Jobs & Steps</h4>
              <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
                {run.jobs.map((job) => (
                  <div key={job.id} className="rounded-lg border border-border/70 bg-card/60 p-3 text-sm">
                    <div className="mb-2 flex items-center justify-between">
                      <span className="font-medium">{job.name}</span>
                      <ActionStatus status={job.status} />
                    </div>
                    {job.steps && job.steps.length > 0 && (
                      <ul className="space-y-1.5 text-xs text-muted-foreground">
                        {job.steps.map((step) => (
                          <li key={step.id} className="flex items-center justify-between">
                            <span className="truncate">{step.name}</span>
                            <span className="shrink-0 text-[11px] font-mono">{step.durationSeconds ?? 0}s</span>
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
            <h4 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Execution Log</h4>
            <ActionLogs logs={run.logs} isRunning={run.status === 'running'} />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
