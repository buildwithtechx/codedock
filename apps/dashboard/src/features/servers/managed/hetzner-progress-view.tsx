import { CheckCircle2, CircleDashed, Loader2, XCircle } from 'lucide-react';
import { useEffect, useRef } from 'react';
import { Button } from '@codedock/ui/button';
import { cn } from '#/lib/utils';
import { hetznerFailed, hetznerStageIndex, hetznerStages } from './hetzner-progress-stages';
import { useOperationStream } from './use-operation-stream';

export function HetznerProgressView({
  operationId,
  serverName,
  serverId,
  onDone,
}: {
  operationId: string;
  serverName: string;
  serverId: string;
  onDone: (serverId: string) => void;
}) {
  const { snapshot, streaming, error } = useOperationStream(operationId);
  const terminalRef = useRef<HTMLPreElement>(null);
  const active = hetznerStageIndex(snapshot);
  const failed = hetznerFailed(snapshot);
  const done = snapshot !== null && snapshot.status === 'COMPLETED';

  useEffect(() => {
    const node = terminalRef.current;
    if (node) node.scrollTop = node.scrollHeight;
  }, [snapshot?.logs]);

  return (
    <div className="space-y-4">
      <div className="rounded-xl border border-border/70">
        <div className="flex items-center justify-between gap-3 border-border/60 border-b px-4 py-3">
          <div>
            <h3 className="font-medium text-foreground text-sm">{serverName}</h3>
            <p className="font-mono text-[11px] text-muted-foreground">
              {snapshot ? `${snapshot.status} · ${snapshot.phase || 'QUEUED'}` : 'Connecting…'}
            </p>
          </div>
          <span
            className={cn(
              'rounded-full px-2.5 py-1 font-medium text-xs',
              done && 'bg-success/10 text-success',
              failed && 'bg-destructive/10 text-destructive',
              !done && !failed && 'bg-muted text-muted-foreground'
            )}
          >
            {done ? 'Ready' : failed ? 'Failed' : streaming ? 'Live' : 'Polling'}
          </span>
        </div>
        <div className="space-y-1 p-3">
          {hetznerStages.map((stage, index) => {
            const state =
              failed && index === Math.min(active, hetznerStages.length - 1)
                ? 'failed'
                : index < active || done
                  ? 'passed'
                  : index === active
                    ? 'running'
                    : 'waiting';
            return (
              <div key={stage.id} className="flex items-center gap-3 rounded-lg px-2 py-1.5">
                {state === 'running' ? (
                  <Loader2 className="size-4 shrink-0 animate-spin text-primary" />
                ) : state === 'passed' ? (
                  <CheckCircle2 className="size-4 shrink-0 text-success" />
                ) : state === 'failed' ? (
                  <XCircle className="size-4 shrink-0 text-destructive" />
                ) : (
                  <CircleDashed className="size-4 shrink-0 text-muted-foreground/40" />
                )}
                <div className="min-w-0">
                  <p className="font-medium text-foreground text-xs">{stage.label}</p>
                  <p className="truncate text-[11px] text-muted-foreground">{stage.detail}</p>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      <div className="overflow-hidden rounded-xl border border-border/70 bg-zinc-950">
        <div className="border-white/10 border-b px-4 py-2">
          <p className="font-mono text-[11px] text-zinc-400">cloud-init · live terminal</p>
        </div>
        <pre
          ref={terminalRef}
          className="h-56 overflow-y-auto whitespace-pre-wrap p-4 font-mono text-[11px] text-zinc-200 leading-relaxed"
        >
          {snapshot?.logs?.trim() ? snapshot.logs : 'Waiting for provision logs…'}
        </pre>
      </div>

      {(error || snapshot?.error) && (
        <p className="text-destructive text-xs">{error || snapshot?.error}</p>
      )}

      <div className="flex items-center justify-end gap-2">
        <Button type="button" onClick={() => onDone(serverId)} disabled={!done} className="text-xs">
          {done ? 'Done, open server' : 'Waiting for server…'}
        </Button>
      </div>
    </div>
  );
}
