import { ChevronRight, GitBranch } from 'lucide-react';
import type { ReactNode } from 'react';
import type { ActionRun } from '#/interfaces/actions';
import { ActionStatus } from './action-status';
import { ActionsEmptyState } from './actions-empty-state';

export function ActionRunList({
  runs,
  empty,
  onSelectRun,
}: {
  runs: ActionRun[];
  empty?: ReactNode;
  onSelectRun?: (run: ActionRun) => void;
}) {
  if (!runs.length) {
    return (
      empty ?? (
        <ActionsEmptyState
          compact
          kind="history"
          title="No workflow runs yet"
          description="Runs will appear here once workflows are triggered manually or by repository events."
        />
      )
    );
  }

  return (
    <div className="divide-y divide-border/40 overflow-hidden rounded-2xl border border-border/50 bg-card shadow-xs">
      {runs.map((run) => (
        <button
          key={run.id}
          type="button"
          onClick={() => onSelectRun?.(run)}
          className="flex w-full flex-wrap items-center gap-x-5 gap-y-2 px-5 py-4 text-start transition-colors hover:bg-muted/40 focus-visible:outline-2 focus-visible:outline-ring"
        >
          <GitBranch className="size-4 shrink-0 text-muted-foreground" />
          <div className="min-w-0 flex-1">
            <p className="truncate font-medium text-foreground text-sm">
              {run.name}{' '}
              <span className="text-muted-foreground">
                #{run.number}
                {(run.attempt ?? 1) > 1 ? ` · attempt ${run.attempt}` : ''}
              </span>
            </p>
            <p className="mt-1 truncate text-muted-foreground text-xs">
              <bdi>
                {run.owner && run.repo
                  ? `${run.owner}/${run.repo}${run.ref ? ` · ${run.ref.replace(/^refs\/(heads|tags)\//, '')}` : ''}`
                  : 'Standalone'}
                {run.revision ? ` · ${run.revision.slice(0, 7)}` : ''}
              </bdi>
            </p>
          </div>
          <ActionStatus status={run.status} />
          <time dateTime={run.createdAt} className="text-muted-foreground text-xs">
            {new Date(run.createdAt).toLocaleString(undefined, {
              month: 'short',
              day: 'numeric',
              hour: 'numeric',
              minute: '2-digit',
            })}
          </time>
          <ChevronRight className="size-4 text-muted-foreground rtl:rotate-180" />
        </button>
      ))}
    </div>
  );
}
