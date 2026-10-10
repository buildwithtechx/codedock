import { Cloud, Pencil, RotateCcw, Zap } from 'lucide-react';

interface DeployTopBarProps {
  environmentName?: string;
  onEditTarget: () => void;
}

export function DeployTopBar({ environmentName = 'Production', onEditTarget }: DeployTopBarProps) {
  return (
    <button
      type="button"
      onClick={onEditTarget}
      className="group flex w-full items-center justify-between gap-3 rounded-2xl border border-border/60 bg-card p-4 text-left transition-all hover:border-primary/40 hover:bg-card/90 sm:p-5"
      aria-label="Edit deployment target"
    >
      <div className="flex items-center gap-2.5">
        <div className="flex size-7 items-center justify-center rounded-lg bg-muted text-muted-foreground transition-colors group-hover:bg-primary/10 group-hover:text-primary">
          <Cloud className="size-4" />
        </div>
        <div className="flex items-center gap-1.5 text-xs sm:text-sm">
          <span className="text-muted-foreground">Deploy:</span>
          <span className="font-semibold text-foreground transition-colors group-hover:text-primary">
            {environmentName}
          </span>
        </div>
      </div>
      <div className="flex items-center gap-3 text-muted-foreground text-xs sm:gap-4">
        <span className="flex items-center gap-1.5">
          <Zap className="size-3.5" />
          Full capacity
        </span>
        <span className="flex items-center gap-1.5">
          <RotateCcw className="size-3.5" />
          Rollback: auto
        </span>
        <div className="flex size-7 items-center justify-center rounded-lg text-muted-foreground transition-colors group-hover:bg-muted group-hover:text-foreground">
          <Pencil className="size-3.5" />
        </div>
      </div>
    </button>
  );
}
