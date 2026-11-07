import { Link } from '@tanstack/react-router';
import { ArrowLeft, Cloud, Plus, RefreshCw } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@codedock/ui/tabs';
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
          className="-ms-3 gap-1.5 text-muted-foreground hover:text-foreground"
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
  const isRefreshing = workflowsQuery.isFetching || runsQuery.isFetching || runnersQuery.isFetching;
  const handleRefresh = () => {
    void workflowsQuery.refetch();
    void runsQuery.refetch();
    void runnersQuery.refetch();
  };

  return (
    <div className="space-y-6">
      <header className="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
        <div>
          <h1 className="font-medium text-2xl text-foreground tracking-tight">Actions</h1>
          <p className="mt-1 text-muted-foreground text-sm">
            Your workflows. Your runners. One place to follow every job.
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="icon"
            onClick={handleRefresh}
            disabled={isRefreshing}
            title="Refresh actions"
            aria-label="Refresh actions"
            className="rounded-xl"
          >
            <RefreshCw className={`size-4 ${isRefreshing ? 'animate-spin' : ''}`} />
          </Button>
          {hasItemsForTab && (
            <Button asChild className="rounded-xl">
              <Link
                to={
                  tab === 'runners' ? ('/actions/runners/new' as never) : ('/actions/new' as never)
                }
              >
                {tab === 'runners' && isCloud ? (
                  <Cloud className="size-4" />
                ) : (
                  <Plus className="size-4" />
                )}
                {tab === 'runners' ? 'Add runner' : 'New workflow'}
              </Link>
            </Button>
          )}
        </div>
      </header>

      <Tabs value={tab} onValueChange={(val) => setTab(val as typeof tab)} className="space-y-6">
        <TabsList>
          <TabsTrigger value="workflows">
            Workflows
            <span className="rounded-full bg-muted px-1.5 py-0.5 font-medium text-xs tabular-nums">
              {workflows.length}
            </span>
          </TabsTrigger>
          <TabsTrigger value="runs">
            Runs
            {runs.length > 0 && (
              <span className="rounded-full bg-muted px-1.5 py-0.5 font-medium text-xs tabular-nums">
                {runs.length}
              </span>
            )}
          </TabsTrigger>
          <TabsTrigger value="runners">
            Runners
            <span className="rounded-full bg-muted px-1.5 py-0.5 font-medium text-xs tabular-nums">
              {runners.length}
            </span>
          </TabsTrigger>
        </TabsList>

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
                      <Button asChild className="h-10 px-5 font-medium">
                        <Link to={'/actions/new' as never}>
                          <Plus className="size-4" />
                          New workflow
                        </Link>
                      </Button>
                    ) : (
                      <Button asChild className="h-10 px-5 font-medium">
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
