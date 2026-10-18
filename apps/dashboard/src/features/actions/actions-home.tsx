import { Link } from '@tanstack/react-router';
import { ArrowLeft, Cloud, Plus } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Tabs, TabsContent } from '@codedock/ui/tabs';
import type { ActionRun, ActionRunner, ActionWorkflow } from '#/interfaces/actions';
import {
  useDeleteRunner,
  useDeleteWorkflow,
  useRunners,
  useRuns,
  useTriggerWorkflow,
  useWorkflows,
} from '#/hooks/use-actions';
import { ActionRunList } from './action-run-list';
import { ActionRunnerCards } from './action-runner-cards';
import { ActionWorkflowList } from './action-workflow-list';
import { ActionsEmptyState } from './actions-empty-state';
import { ActionsSidebar } from './actions-sidebar';
import { RunWorkspace } from './run-workspace';

export interface ActionsHomeProps {
  initialWorkflows?: ActionWorkflow[];
  initialRuns?: ActionRun[];
  initialRunners?: ActionRunner[];
  isCloud?: boolean;
}

export function ActionsHome({
  initialWorkflows,
  initialRuns,
  initialRunners,
  isCloud = false,
}: ActionsHomeProps) {
  const [tab, setTab] = useState<'workflows' | 'runs' | 'runners'>('workflows');
  const [selectedRun, setSelectedRun] = useState<ActionRun | null>(null);

  const workflowsQuery = useWorkflows();
  const runsQuery = useRuns();
  const runnersQuery = useRunners();
  const triggerWorkflowMutation = useTriggerWorkflow();
  const deleteWorkflowMutation = useDeleteWorkflow();
  const deleteRunnerMutation = useDeleteRunner();

  const workflows = workflowsQuery.data ?? initialWorkflows ?? [];
  const runs = runsQuery.data ?? initialRuns ?? [];
  const runners = runnersQuery.data ?? initialRunners ?? [];

  const handleTriggerWorkflow = async (workflowId: string) => {
    try {
      const run = await triggerWorkflowMutation.mutateAsync(workflowId);
      setSelectedRun(run);
      toast.success(`Workflow run #${run.number} started`);
    } catch {
      toast.error('Failed to trigger workflow');
    }
  };

  const handleDeleteWorkflow = async (workflowId: string) => {
    try {
      await deleteWorkflowMutation.mutateAsync(workflowId);
      toast.success('Workflow deleted');
    } catch {
      toast.error('Failed to delete workflow');
    }
  };

  const handleDeleteRunner = async (runnerId: string) => {
    try {
      await deleteRunnerMutation.mutateAsync(runnerId);
      toast.success('Runner deleted');
    } catch {
      toast.error('Failed to delete runner');
    }
  };

  if (selectedRun) {
    return (
      <div className="space-y-4">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setSelectedRun(null)}
          className="-ms-3 gap-1.5 rounded-lg text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="size-4 rtl:rotate-180" />
          Back to Actions Overview
        </Button>
        <RunWorkspace run={selectedRun} onRunUpdated={(updated) => setSelectedRun(updated)} />
      </div>
    );
  }

  const readyRunnerCount = runners.filter((r) => r.enabled && !r.error).length;
  const firstWorkflow = workflows[0];
  const latestRun = runs[0];

  const hasItemsForTab = tab === 'runners' ? runners.length > 0 : workflows.length > 0;

  return (
    <div className="space-y-6">
      <header className="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
        <div>
          <h1 className="font-medium text-2xl text-foreground tracking-tight">Actions</h1>
          <p className="mt-1 text-muted-foreground text-sm">
            Your workflows. Your runners. One place to follow every job.
          </p>
        </div>

        {hasItemsForTab && (
          <Button asChild className="rounded-xl">
            <Link
              to={tab === 'runners' ? ('/actions/runners/new' as never) : ('/actions/new' as never)}
            >
              {tab === 'runners' ? (
                isCloud ? (
                  <Cloud className="size-4" />
                ) : (
                  <Plus className="size-4" />
                )
              ) : (
                <Plus className="size-4" />
              )}
              {tab === 'runners' ? 'Add runner' : 'New workflow'}
            </Link>
          </Button>
        )}
      </header>

      <Tabs value={tab} onValueChange={(val) => setTab(val as typeof tab)} className="space-y-6">
        <div className="flex items-center gap-1 overflow-x-auto border-border/50 border-b">
          <button
            type="button"
            onClick={() => setTab('workflows')}
            className={`relative inline-flex items-center gap-2 px-4 py-2.5 font-medium text-sm transition-colors ${
              tab === 'workflows'
                ? 'text-foreground'
                : 'text-muted-foreground hover:text-foreground/70'
            }`}
          >
            Workflows
            <span
              className={`rounded-full px-1.5 py-0.5 font-medium text-xs tabular-nums ${
                tab === 'workflows'
                  ? 'bg-muted text-foreground'
                  : 'bg-muted/60 text-muted-foreground'
              }`}
            >
              {workflows.length}
            </span>
            {tab === 'workflows' && (
              <span className="absolute inset-x-4 bottom-0 h-0.5 rounded-full bg-foreground" />
            )}
          </button>
          <button
            type="button"
            onClick={() => setTab('runs')}
            className={`relative inline-flex items-center gap-2 px-4 py-2.5 font-medium text-sm transition-colors ${
              tab === 'runs' ? 'text-foreground' : 'text-muted-foreground hover:text-foreground/70'
            }`}
          >
            Runs
            {runs.length > 0 && (
              <span
                className={`rounded-full px-1.5 py-0.5 font-medium text-xs tabular-nums ${
                  tab === 'runs' ? 'bg-muted text-foreground' : 'bg-muted/60 text-muted-foreground'
                }`}
              >
                {runs.length}
              </span>
            )}
            {tab === 'runs' && (
              <span className="absolute inset-x-4 bottom-0 h-0.5 rounded-full bg-foreground" />
            )}
          </button>
          <button
            type="button"
            onClick={() => setTab('runners')}
            className={`relative inline-flex items-center gap-2 px-4 py-2.5 font-medium text-sm transition-colors ${
              tab === 'runners'
                ? 'text-foreground'
                : 'text-muted-foreground hover:text-foreground/70'
            }`}
          >
            Runners
            <span
              className={`rounded-full px-1.5 py-0.5 font-medium text-xs tabular-nums ${
                tab === 'runners' ? 'bg-muted text-foreground' : 'bg-muted/60 text-muted-foreground'
              }`}
            >
              {runners.length}
            </span>
            {tab === 'runners' && (
              <span className="absolute inset-x-4 bottom-0 h-0.5 rounded-full bg-foreground" />
            )}
          </button>
        </div>

        <div className="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
          <div className="min-w-0 space-y-4">
            <TabsContent value="workflows" className="m-0 space-y-4">
              <ActionWorkflowList
                workflows={workflows}
                runs={runs}
                onTriggerWorkflow={handleTriggerWorkflow}
                onDeleteWorkflow={handleDeleteWorkflow}
                readyRunnerCount={readyRunnerCount}
                runnerCount={runners.length}
                onGoToRunners={() => setTab('runners')}
                isCloud={isCloud}
              />
            </TabsContent>

            <TabsContent value="runs" className="m-0 space-y-4">
              <ActionRunList
                runs={runs}
                onSelectRun={(run) => setSelectedRun(run)}
                empty={
                  <ActionsEmptyState
                    kind="history"
                    title="No workflow runs yet"
                    description="Trigger a workflow to begin running jobs on your runners."
                  >
                    {firstWorkflow ? (
                      <Button
                        asChild
                        className="h-10 rounded-full bg-primary px-5 font-medium text-primary-foreground hover:bg-primary/90"
                      >
                        <Link to={'/actions/new' as never}>
                          <Plus className="size-4" />
                          New workflow
                        </Link>
                      </Button>
                    ) : (
                      <Button
                        asChild
                        className="h-10 rounded-full bg-primary px-5 font-medium text-primary-foreground hover:bg-primary/90"
                      >
                        <Link to={'/actions/runners/new' as never}>
                          {isCloud ? <Cloud className="size-4" /> : <Plus className="size-4" />}
                          Add runner
                        </Link>
                      </Button>
                    )}
                  </ActionsEmptyState>
                }
              />
            </TabsContent>

            <TabsContent value="runners" className="m-0 space-y-4">
              <ActionRunnerCards
                runners={runners}
                onDeleteRunner={handleDeleteRunner}
                isCloud={isCloud}
              />
            </TabsContent>
          </div>

          <ActionsSidebar
            runnerCount={runners.length}
            readyRunnerCount={readyRunnerCount}
            workflowCount={workflows.length}
            firstWorkflow={firstWorkflow}
            latestRun={latestRun}
            onRunners={() => setTab('runners')}
            isCloud={isCloud}
          />
        </div>
      </Tabs>
    </div>
  );
}
