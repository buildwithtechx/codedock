import { Link } from '@tanstack/react-router';
import { Plus, Server } from 'lucide-react';
import { Button } from '#/components/ui/button';
import type { ActionRunner } from '#/services/actions';
import { ActionsEmptyState } from './actions-empty-state';

interface ActionRunnerCardsProps {
  runners: ActionRunner[];
}

export function ActionRunnerCards({ runners }: ActionRunnerCardsProps) {
  if (runners.length === 0) {
    return (
      <ActionsEmptyState
        kind="runner"
        title="No runners configured"
        description="Runners execute your workflow jobs in isolated containers or directly on host machines."
      >
        <Button asChild className="rounded-xl">
          <Link to={'/actions/runners/new' as never}>
            <Plus className="size-4" />
            Add runner
          </Link>
        </Button>
      </ActionsEmptyState>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
      {runners.map((runner) => (
        <article key={runner.id} className="rounded-2xl bg-card p-5">
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
            {runner.config?.labels?.map((lbl) => (
              <span key={lbl} className="rounded-md bg-muted px-2 py-0.5 text-muted-foreground text-xs">
                {lbl}
              </span>
            ))}
          </div>

          <div className="mt-4 border-border/40 border-t pt-3 text-muted-foreground text-xs">
            {runner.config?.cpu ?? 1} vCPU · {(runner.config?.memoryMb ?? 2048) / 1024} GiB · Max {runner.config?.maxParallel ?? 1} parallel jobs
          </div>
        </article>
      ))}
    </div>
  );
}
