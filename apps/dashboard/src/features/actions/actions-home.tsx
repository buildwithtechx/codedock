import { Link } from '@tanstack/react-router';
import {
  ArrowLeft,
  Cloud,
  History,
  PlayCircle,
  Plus,
  Server,
} from 'lucide-react';
import { useCallback, useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Tabs, TabsContent } from '#/components/ui/tabs';
import {
  actionsService,
  type ActionRun,
  type ActionRunner,
  type ActionWorkflow,
} from '#/services/actions';
import { ActionsSidebar } from './actions-sidebar';
import { ActionRunList } from './action-run-list';
import { ActionWorkflowList } from './action-workflow-list';
import { ActionRunnerCards } from './action-runner-cards';
import { ActionsEmptyState } from './actions-empty-state';
import { RunWorkspace } from './run-workspace';

export interface ActionsHomeProps {
  initialWorkflows?: ActionWorkflow[];
  initialRuns?: ActionRun[];
  initialRunners?: ActionRunner[];
}

export function ActionsHome({
  initialWorkflows,
  initialRuns,
  initialRunners,
}: ActionsHomeProps) {
  const [tab, setTab] = useState<'workflows' | 'runs' | 'runners'>('workflows');
  const [selectedRun, setSelectedRun] = useState<ActionRun | null>(null);

  const [workflows, setWorkflows] = useState<ActionWorkflow[]>(initialWorkflows ?? []);
  const [runs, setRuns] = useState<ActionRun[]>(initialRuns ?? []);
  const [runners, setRunners] = useState<ActionRunner[]>(initialRunners ?? []);

  const loadData = useCallback(async () => {
    try {
      const [wfList, runList, runnerList] = await Promise.all([
        actionsService.listWorkflows(),
        actionsService.listRuns(),
        actionsService.listRunners(),
      ]);
      setWorkflows(wfList);
      setRuns(runList);
      setRunners(runnerList);
    } catch {
      // Retain existing
    }
  }, []);

  useEffect(() => {
    if (!initialWorkflows && !initialRuns && !initialRunners) {
      loadData();
    }
  }, [loadData, initialWorkflows, initialRuns, initialRunners]);

  const handleTriggerWorkflow = async (workflowId: string) => {
    try {
      const run = await actionsService.triggerWorkflow(workflowId);
      setRuns((prev) => [run, ...prev]);
      setSelectedRun(run);
      toast.success(`Workflow run #${run.number} started`);
    } catch {
      toast.error('Failed to trigger workflow');
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
        <RunWorkspace
          run={selectedRun}
          onRunUpdated={(updated) => setSelectedRun(updated)}
        />
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
            <Link to={tab === 'runners' ? ('/actions/runners/new' as never) : ('/actions/new' as never)}>
              <Plus className="size-4" />
              {tab === 'runners' ? 'Add runner' : 'New workflow'}
            </Link>
          </Button>
        )}
      </header>

      <Tabs
        value={tab}
        onValueChange={(val) => setTab(val as typeof tab)}
        className="space-y-6"
      >
        <div className="flex items-center gap-6 border-b border-border/50 text-sm">
          <button
            type="button"
            onClick={() => setTab('workflows')}
            className={`relative pb-3 font-medium transition-colors ${
              tab === 'workflows' ? 'text-foreground' : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            Workflows <span className="tabular-nums">{workflows.length}</span>
            {tab === 'workflows' && (
              <span className="-bottom-px absolute inset-x-0 h-0.5 rounded-full bg-primary" />
            )}
          </button>
          <button
            type="button"
            onClick={() => setTab('runs')}
            className={`relative pb-3 font-medium transition-colors ${
              tab === 'runs' ? 'text-foreground' : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            Runs
            {tab === 'runs' && (
              <span className="-bottom-px absolute inset-x-0 h-0.5 rounded-full bg-primary" />
            )}
          </button>
          <button
            type="button"
            onClick={() => setTab('runners')}
            className={`relative pb-3 font-medium transition-colors ${
              tab === 'runners' ? 'text-foreground' : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            Runners <span className="tabular-nums">{runners.length}</span>
            {tab === 'runners' && (
              <span className="-bottom-px absolute inset-x-0 h-0.5 rounded-full bg-primary" />
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
              />
            </TabsContent>

            <TabsContent value="runs" className="m-0 space-y-4">
              <ActionRunList
                runs={runs}
                onSelectRun={(run) => setSelectedRun(run)}
                empty={
                  <ActionsEmptyState
                    kind="history"
                    title="No runs yet"
                    description="Trigger a workflow to begin running jobs on your runners."
                  >
                    <Button asChild className="rounded-full bg-white text-zinc-950 hover:bg-zinc-100 font-medium px-5">
                      <Link to={'/actions/runners/new' as never}>
                        <Cloud className="size-4" />
                        Add runner
                      </Link>
                    </Button>
                  </ActionsEmptyState>
                }
              />
            </TabsContent>

            <TabsContent value="runners" className="m-0 space-y-4">
              <ActionRunnerCards runners={runners} />
            </TabsContent>
          </div>

          <ActionsSidebar
            runnerCount={runners.length}
            readyRunnerCount={readyRunnerCount}
            workflowCount={workflows.length}
            firstWorkflow={firstWorkflow}
            latestRun={latestRun}
            onRunners={() => setTab('runners')}
          />
        </div>
      </Tabs>
    </div>
  );
}
