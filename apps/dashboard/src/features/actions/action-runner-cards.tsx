import { Link } from '@tanstack/react-router';
import { Cloud, Plus, Server, Trash2 } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import type { ActionRunner } from '#/interfaces/actions';
import { ActionsEmptyState } from './actions-empty-state';

interface ActionRunnerCardsProps {
  runners: ActionRunner[];
  onDeleteRunner?: (runnerId: string) => void;
  isCloud?: boolean;
}

export function ActionRunnerCards({
  runners,
  onDeleteRunner,
  isCloud = false,
}: ActionRunnerCardsProps) {
  if (runners.length === 0) {
    return (
      <ActionsEmptyState
        kind="runner"
        title="No runners configured"
        description="Runners execute your workflow jobs in isolated containers or directly on host machines."
      >
        <Button
          asChild
          className="h-10 rounded-full bg-primary px-5 font-medium text-primary-foreground hover:bg-primary/90"
        >
          <Link to={'/actions/runners/new' as never}>
            {isCloud ? <Cloud className="size-4" /> : <Plus className="size-4" />}
            Add runner
          </Link>
        </Button>
      </ActionsEmptyState>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
      {runners.map((runner) => (
        <article
          key={runner.id}
          className="rounded-2xl border border-border/50 bg-card p-5 shadow-xs"
        >
          <div className="flex items-center gap-3">
            <span className="flex size-10 items-center justify-center rounded-xl bg-muted/60 text-foreground">
              <Server className="size-5" />
            </span>
            <div className="min-w-0 flex-1">
              <h2 className="truncate font-medium text-base text-foreground">{runner.name}</h2>
              <p className="mt-0.5 text-muted-foreground text-xs">
                {runner.kind === 'cloud' ? 'Cloud runner' : 'Connected server runner'}
              </p>
            </div>
            <span
              className={`size-2.5 rounded-full ${
                runner.enabled && !runner.error ? 'bg-emerald-500' : 'bg-muted-foreground/40'
              }`}
            />
          </div>

          <div className="mt-4 flex flex-wrap gap-1.5">
            {runner.config?.labels?.map((lbl: string) => (
              <span
                key={lbl}
                className="rounded-md bg-muted px-2 py-0.5 text-muted-foreground text-xs"
              >
                {lbl}
              </span>
            ))}
          </div>

          <div className="mt-4 flex items-center justify-between gap-2 border-border/40 border-t pt-3 text-muted-foreground text-xs">
            <span>
              {runner.config?.cpu ?? 1} vCPU · {(runner.config?.memoryMb ?? 2048) / 1024} GiB · Max{' '}
              {runner.config?.maxParallel ?? 1} parallel jobs
            </span>
            {onDeleteRunner && (
              <Button
                size="sm"
                variant="ghost"
                onClick={() => onDeleteRunner(runner.id)}
                className="h-7 w-7 rounded-lg p-0 text-muted-foreground hover:text-destructive"
                aria-label={`Delete runner ${runner.name}`}
                title="Delete runner"
              >
                <Trash2 className="size-3.5" />
              </Button>
            )}
          </div>
        </article>
      ))}
    </div>
  );
}
