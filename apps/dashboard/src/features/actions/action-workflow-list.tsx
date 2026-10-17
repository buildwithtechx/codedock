import { Link } from '@tanstack/react-router';
import { ArrowRight, Cloud, PlayCircle, Plus, Search, Trash2 } from 'lucide-react';
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
  onDeleteWorkflow?: (workflowId: string) => void;
  readyRunnerCount?: number;
  runnerCount?: number;
  onGoToRunners?: () => void;
  isCloud?: boolean;
}

export function ActionWorkflowList({
  workflows,
  runs,
  onTriggerWorkflow,
  onDeleteWorkflow,
  readyRunnerCount = 0,
  runnerCount = 0,
  onGoToRunners,
  isCloud = false,
}: ActionWorkflowListProps) {
  const [search, setSearch] = useState('');

  if (workflows.length === 0) {
    const setupAction =
      readyRunnerCount > 0 ? (
        <Button asChild className="h-10 rounded-full bg-primary px-5 font-medium text-primary-foreground hover:bg-primary/90">
          <Link to={'/actions/new' as never}>
            <Plus className="size-4" />
            New workflow
          </Link>
        </Button>
      ) : runnerCount > 0 ? (
        <Button onClick={onGoToRunners} variant="secondary" className="h-10 rounded-full px-5 font-medium">
          Review runners
          <ArrowRight className="size-4 rtl:rotate-180" />
        </Button>
      ) : (
        <Button asChild className="h-10 rounded-full bg-primary px-5 font-medium text-primary-foreground hover:bg-primary/90">
          <Link to={'/actions/runners/new' as never}>
            {isCloud ? <Cloud className="size-4" /> : <Plus className="size-4" />}
            Add runner
          </Link>
        </Button>
      );

    return (
      <ActionsEmptyState
        title="From commit to done"
        description="Build, test and ship from your repositories. Follow every job and its logs in one place."
      >
        {setupAction}
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
                  {onDeleteWorkflow && (
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={() => onDeleteWorkflow(workflow.id)}
                      className="h-7 w-7 p-0 rounded-lg text-muted-foreground hover:text-destructive"
                      aria-label={`Delete workflow ${workflow.name}`}
                      title="Delete workflow"
                    >
                      <Trash2 className="size-3.5" />
                    </Button>
                  )}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
