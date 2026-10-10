import { Cloud, Pencil, RotateCcw, Zap } from 'lucide-react';

interface DeployTopBarProps {
  environmentName?: string;
  onEditTarget: () => void;
}

export function DeployTopBar({ environmentName = 'Production', onEditTarget }: DeployTopBarProps) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-3 text-muted-foreground text-xs">
      <div className="flex items-center gap-2">
        <Cloud className="size-3.5 text-muted-foreground" />
        <span>Deploy:</span>
        <span className="font-semibold text-foreground">{environmentName}</span>
      </div>
      <div className="flex items-center gap-3 sm:gap-4">
        <span className="flex items-center gap-1.5">
          <Zap className="size-3.5" />
          Full capacity
        </span>
        <span className="flex items-center gap-1.5">
          <RotateCcw className="size-3.5" />
          Rollback: auto
        </span>
        <button
          type="button"
          onClick={onEditTarget}
          className="text-muted-foreground transition-colors hover:text-foreground"
          aria-label="Edit deployment target"
        >
          <Pencil className="size-3.5" />
        </button>
      </div>
    </div>
  );
}
