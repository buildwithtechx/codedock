import { Cloud, Cpu, Pencil, Plus } from 'lucide-react';

interface DeployTopBarProps {
  environmentName?: string;
  onEditTarget: () => void;
  runtimeLabel?: string;
}

export function DeployTopBar({
  environmentName = 'Production',
  onEditTarget,
  runtimeLabel = 'Direct',
}: DeployTopBarProps) {
  return (
    <button
      type="button"
      onClick={onEditTarget}
      className="group flex w-full items-center gap-3 rounded-xl border border-border/50 bg-card px-4 py-3 text-start transition-all hover:border-primary/30"
      aria-label="Edit deployment target"
    >
      <div className="flex min-w-0 flex-1 flex-wrap items-center gap-x-4 gap-y-2">
        <div className="flex min-w-0 grow basis-64 flex-wrap items-center gap-3">
          <div className="flex min-w-0 items-center gap-2 text-xs">
            <div className="flex shrink-0 items-center gap-1 text-muted-foreground">
              <Cpu className="size-4" />
              <Plus className="size-3" />
              <Cloud className="size-4" />
            </div>
            <span className="shrink-0 text-muted-foreground">Build & deploy</span>
            <span className="truncate font-medium text-foreground" title={environmentName}>
              {environmentName}
            </span>
          </div>
        </div>
        <div className="flex flex-wrap items-center gap-x-2.5 gap-y-1.5">
          <span className="inline-flex items-center rounded-full bg-muted/60 px-2.5 py-0.5 font-medium text-muted-foreground text-xs">
            {runtimeLabel}
          </span>
          <span className="inline-flex items-center rounded-full bg-muted/60 px-2.5 py-0.5 font-medium text-muted-foreground text-xs">
            Rollback: auto
          </span>
        </div>
      </div>
      <Pencil className="size-4 shrink-0 text-muted-foreground transition-colors group-hover:text-foreground" />
    </button>
  );
}
