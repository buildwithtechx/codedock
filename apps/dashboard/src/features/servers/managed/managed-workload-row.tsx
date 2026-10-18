import { Play, RotateCw, Square, Terminal, Trash2 } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import type { ManagedWorkload } from '#/interfaces/server';

interface ManagedWorkloadRowProps {
  workload: ManagedWorkload;
  onViewLogs: (workload: ManagedWorkload) => void;
  onControl: (workloadId: string, action: 'start' | 'stop' | 'restart' | 'delete') => void;
}

export function ManagedWorkloadRow({ workload, onViewLogs, onControl }: ManagedWorkloadRowProps) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-border/40 bg-card p-4">
      <div>
        <div className="flex items-center gap-2">
          <span className="font-medium text-foreground text-sm">{workload.name}</span>
          <span
            className={`rounded-full px-2 py-0.5 font-medium text-[10px] ${
              workload.state === 'running'
                ? 'bg-success/10 text-success'
                : 'bg-muted text-muted-foreground'
            }`}
          >
            {workload.state}
          </span>
        </div>
        <p className="mt-1 font-mono text-muted-foreground text-xs">
          Policy: {workload.restartPolicy ?? 'on-failure'} • Source: {workload.source}
        </p>
      </div>

      <div className="flex items-center gap-1.5">
        <Button
          size="icon"
          variant="ghost"
          className="size-8"
          aria-label="View logs"
          onClick={() => void onViewLogs(workload)}
        >
          <Terminal className="size-4" />
        </Button>
        {workload.state === 'running' ? (
          <Button
            size="icon"
            variant="ghost"
            className="size-8 text-warning"
            aria-label="Stop workload"
            onClick={() => void onControl(workload.id, 'stop')}
          >
            <Square className="size-4" />
          </Button>
        ) : (
          <Button
            size="icon"
            variant="ghost"
            className="size-8 text-success"
            aria-label="Start workload"
            onClick={() => void onControl(workload.id, 'start')}
          >
            <Play className="size-4" />
          </Button>
        )}
        <Button
          size="icon"
          variant="ghost"
          className="size-8"
          aria-label="Restart workload"
          onClick={() => void onControl(workload.id, 'restart')}
        >
          <RotateCw className="size-4" />
        </Button>
        <Button
          size="icon"
          variant="ghost"
          className="size-8 text-destructive"
          aria-label="Delete workload"
          onClick={() => void onControl(workload.id, 'delete')}
        >
          <Trash2 className="size-4" />
        </Button>
      </div>
    </div>
  );
}
