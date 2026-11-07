import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { ChevronDown, ChevronUp, History, Loader2, StopCircle, Terminal } from 'lucide-react';
import { useState } from 'react';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';

type Operation = {
  id: string;
  kind: string;
  status: string;
  phase: string;
  logs: string;
  error: string;
};

function OperationRow({
  operation,
  onCancel,
  isCancelling,
}: {
  operation: Operation;
  onCancel: (id: string) => void;
  isCancelling: boolean;
}) {
  const [showLogs, setShowLogs] = useState(operation.status === 'RUNNING');
  const isRunning = operation.status === 'RUNNING';
  const isFailed = operation.status === 'FAILED' || Boolean(operation.error);

  return (
    <div className="rounded-xl border border-border/50 bg-muted/20 p-4 transition-colors hover:bg-muted/30">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-2.5 min-w-0">
          <Badge
            variant={isRunning ? 'default' : isFailed ? 'destructive' : 'secondary'}
            className="capitalize text-xs font-medium"
          >
            {isRunning && <Loader2 className="mr-1.5 size-3 animate-spin" />}
            {operation.kind}
          </Badge>
          <span className="font-mono text-xs text-muted-foreground truncate">
            {operation.id.slice(0, 8)}
          </span>
          <span className="text-xs text-foreground font-medium">
            {operation.phase || operation.status}
          </span>
        </div>

        <div className="flex items-center gap-2">
          {operation.logs && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setShowLogs((prev) => !prev)}
              className="h-8 gap-1.5 text-xs text-muted-foreground hover:text-foreground"
            >
              <Terminal className="size-3.5" />
              <span>{showLogs ? 'Hide logs' : 'View logs'}</span>
              {showLogs ? <ChevronUp className="size-3.5" /> : <ChevronDown className="size-3.5" />}
            </Button>
          )}

          {isRunning && (
            <Button
              size="sm"
              variant="destructive"
              disabled={isCancelling}
              onClick={() => onCancel(operation.id)}
              className="h-8 gap-1.5 text-xs"
            >
              <StopCircle className="size-3.5" />
              <span>Cancel</span>
            </Button>
          )}
        </div>
      </div>

      {operation.error && (
        <p role="alert" className="mt-2 text-xs text-destructive font-mono">
          {operation.error}
        </p>
      )}

      {showLogs && operation.logs && (
        <pre className="mt-3 max-h-52 overflow-auto rounded-lg bg-black/60 p-3 font-mono text-[11px] text-zinc-300 leading-relaxed">
          {operation.logs}
        </pre>
      )}
    </div>
  );
}

export function BackupOperationHistory() {
  const client = useQueryClient();
  const query = useQuery({
    queryKey: ['backup-operations'],
    queryFn: () => apiClient.get<BaseResponse<Operation[]>>('/operations'),
    refetchInterval: 3000,
  });

  const cancel = useMutation({
    mutationFn: (id: string) => apiClient.post(`/operations/${id}/cancel`),
    onSettled: () => client.invalidateQueries({ queryKey: ['backup-operations'] }),
  });

  const operations =
    query.data?.data?.filter(
      (operation) => operation.kind === 'backup' || operation.kind === 'restore'
    ) ?? [];

  if (operations.length === 0 && !query.error) {
    return null;
  }

  return (
    <section
      aria-label="Backup and restore operations"
      className="min-w-0 overflow-hidden rounded-2xl border border-border/50 bg-card p-5"
    >
      <div className="flex items-center justify-between gap-4 border-b border-border/40 pb-4">
        <div className="flex items-center gap-3">
          <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <History className="size-4" />
          </div>
          <div>
            <h2 className="font-semibold text-base text-foreground">Backup & Restore Runs</h2>
            <p className="text-muted-foreground text-xs">
              Live execution status and runtime logs for background operations.
            </p>
          </div>
        </div>
        <span className="text-muted-foreground text-xs font-mono">
          {operations.length} {operations.length === 1 ? 'operation' : 'operations'}
        </span>
      </div>

      {query.error && (
        <p role="alert" className="mt-3 text-xs text-destructive">
          {query.error.message}
        </p>
      )}

      <div className="mt-4 space-y-3">
        {operations.map((op) => (
          <OperationRow
            key={op.id}
            operation={op}
            onCancel={(id) => cancel.mutate(id)}
            isCancelling={cancel.isPending}
          />
        ))}
      </div>
    </section>
  );
}
