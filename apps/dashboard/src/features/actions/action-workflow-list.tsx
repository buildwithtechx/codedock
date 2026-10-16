import { Link } from '@tanstack/react-router';
import { Cloud, PlayCircle, Search } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import type { ActionRun, ActionWorkflow } from '#/services/actions';
import { ActionStatus } from './action-status';
import { ActionsEmptyState } from './actions-empty-state';

interface ActionWorkflowListProps {
  workflows: ActionWorkflow[];
  runs: ActionRun[];
  onTriggerWorkflow: (workflowId: string) => void;
}

export function ActionWorkflowList({
  workflows,
  runs,
  onTriggerWorkflow,
}: ActionWorkflowListProps) {
  const [search, setSearch] = useState('');

  if (workflows.length === 0) {
    return (
      <ActionsEmptyState
        title="From commit to done"
        description="Build, test and ship from your repositories. Follow every job and its logs in one place."
      >
        <Button asChild className="rounded-full bg-white text-zinc-950 hover:bg-zinc-100 font-medium px-5">
          <Link to={'/actions/runners/new' as never}>
            <Cloud className="size-4" />
            Add runner
          </Link>
        </Button>
      </ActionsEmptyState>
    );
  }

  const filtered = workflows.filter((w) =>
    `${w.name} ${w.owner ?? ''}/${w.repo ?? ''}`.toLowerCase().includes(search.trim().toLowerCase())
  );

  return (
    <div className="space-y-4">
      <div className="relative">
        <Search className="pointer-events-none absolute start-3.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          placeholder="Search workflows..."
          type="search"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="h-10 rounded-xl bg-muted/60 ps-10 pe-4"
        />
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        {filtered.map((workflow) => {
          const last = runs.find((r) => r.workflowId === workflow.id);
          return (
            <div
              key={workflow.id}
              className="group relative flex flex-col justify-between rounded-2xl bg-card p-5 transition-colors hover:bg-muted/40"
            >
              <div>
                <div className="flex items-center gap-3">
                  <span className="rounded-xl bg-muted/50 p-2.5 text-muted-foreground">
                    <PlayCircle className="size-5" />
                  </span>
                  <div className="min-w-0 flex-1">
                    <h2 className="truncate font-medium text-base text-foreground">
                      {workflow.name}
                    </h2>
                    <p className="mt-0.5 truncate text-muted-foreground text-xs">
                      {workflow.owner ? `${workflow.owner}/${workflow.repo}` : 'Standalone'}
                    </p>
                  </div>
                </div>
                <p className="mt-3 line-clamp-2 text-muted-foreground text-xs">
                  {workflow.description || 'Automated workflow execution pipeline'}
                </p>
              </div>

              <div className="mt-5 flex items-center justify-between gap-2 border-border/40 border-t pt-3">
                <code className="text-muted-foreground text-xs">
                  {(workflow.path || 'workflow.yml').split('/').pop()}
                </code>
                <div className="flex items-center gap-2">
                  {last ? (
                    <ActionStatus status={last.status} />
                  ) : (
                    <span className="text-muted-foreground text-xs">No runs</span>
                  )}
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => onTriggerWorkflow(workflow.id)}
                    className="h-7 rounded-lg px-2 text-xs"
                  >
                    Run
                  </Button>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
