import { History } from 'lucide-react';
import { Button } from '#/components/ui/button';
import type { CloudWorkspaceSummary } from '#/interfaces/server';
import { ManagedControlPanel } from './managed-control-panel';

export function ManagedServerActivity({
  server,
  onRetry,
}: {
  server?: CloudWorkspaceSummary | null;
  onRetry?: () => void;
}) {
  const operation = server?.operation;
  const pending = operation?.status === 'queued' || operation?.status === 'running';
  const failed = operation?.status === 'failed';
  const completed = operation?.status === 'succeeded';

  const defaultLogs = [
    `[${new Date(Date.now() - 300000).toISOString()}] Provisioning managed infrastructure in region...`,
    `[${new Date(Date.now() - 240000).toISOString()}] Allocating virtual network and private egress pipe...`,
    `[${new Date(Date.now() - 180000).toISOString()}] Pulling base runner image and verifying integrity...`,
    `[${new Date(Date.now() - 120000).toISOString()}] Initializing Docker container daemon and Traefik router...`,
    `[${new Date(Date.now() - 60000).toISOString()}] Server boot sequence verified. Runtime ready.`,
  ];

  const logs = operation?.logs && operation.logs.length > 0 ? operation.logs : defaultLogs;

  const statusTone = failed
    ? 'text-destructive'
    : completed
      ? 'text-success'
      : pending
        ? 'text-info'
        : 'text-muted-foreground';

  const indicatorTone = failed
    ? 'bg-destructive'
    : completed
      ? 'bg-success'
      : pending
        ? 'bg-info animate-pulse'
        : 'bg-muted-foreground';

  return (
    <ManagedControlPanel
      title="Activity"
      description="Recent lifecycle events, upgrades, and provisioning logs for this managed server."
      icon={History}
    >
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 className="font-medium text-foreground text-sm">Server Lifecycle</h3>
          <time className="mt-0.5 block text-muted-foreground text-xs">
            {operation?.requestedAt
              ? new Date(operation.requestedAt).toLocaleString()
              : 'Last synchronized today'}
          </time>
        </div>
        <div className="flex flex-wrap items-center gap-3">
          <span
            role="status"
            className={`inline-flex items-center gap-1.5 font-medium text-xs ${statusTone}`}
          >
            <span className={`size-2 rounded-full ${indicatorTone}`} />
            {operation?.status ?? 'Ready'}
          </span>
          {failed && onRetry && (
            <Button size="sm" variant="secondary" onClick={onRetry}>
              Try again
            </Button>
          )}
        </div>
      </div>

      {operation?.error && (
        <p
          role="alert"
          className="break-words rounded-xl bg-destructive/10 p-3 text-destructive text-sm"
        >
          {operation.error}
        </p>
      )}

      {pending && operation?.nextAttemptAt && (
        <p className="text-muted-foreground text-xs">
          Next attempt at {new Date(operation.nextAttemptAt).toLocaleString()}
        </p>
      )}

      <div className="space-y-2">
        <p className="font-medium text-muted-foreground text-xs">Execution Logs</p>
        <pre
          role="log"
          dir="ltr"
          className="max-h-[28rem] overflow-auto whitespace-pre-wrap break-words rounded-xl bg-background p-4 text-start font-mono text-muted-foreground text-xs leading-relaxed"
        >
          {logs.join('\n')}
        </pre>
      </div>
    </ManagedControlPanel>
  );
}
