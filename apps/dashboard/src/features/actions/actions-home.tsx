import {
  ArrowLeft,
  Check,
  ChevronRight,
  Clock,
  FileCode2,
  GitBranch,
  PlayCircle,
  Plus,
  RefreshCw,
  Search,
  Server,
  Workflow as WorkflowIcon,
} from 'lucide-react';
import { useMemo, useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import { Textarea } from '#/components/ui/textarea';
import {
  actionsService,
  type ActionRun,
  type ActionRunner,
  type ActionWorkflow,
} from '#/services/actions';
import { ActionStatus } from './action-status';
import { RunWorkspace } from './run-workspace';
import { RunnerCatalog } from './runner-catalog';

interface ActionsHomeProps {
  initialWorkflows: ActionWorkflow[];
  initialRuns: ActionRun[];
  initialRunners: ActionRunner[];
}

export function ActionsHome({
  initialWorkflows,
  initialRuns,
  initialRunners,
}: ActionsHomeProps) {
  const [workflows, setWorkflows] = useState<ActionWorkflow[]>(initialWorkflows);
  const [runs, setRuns] = useState<ActionRun[]>(initialRuns);
  const [runners, setRunners] = useState<ActionRunner[]>(initialRunners);
  const [selectedRun, setSelectedRun] = useState<ActionRun | null>(null);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [activeTab, setActiveTab] = useState<'runs' | 'workflows' | 'runners'>('runs');
  const [workflowSearch, setWorkflowSearch] = useState('');

  // Workflow Dialog State
  const [workflowDialogOpen, setWorkflowDialogOpen] = useState(false);
  const [wfName, setWfName] = useState('');
  const [wfDescription, setWfDescription] = useState('');
  const [wfTrigger, setWfTrigger] = useState('push');
  const [wfYaml, setWfYaml] = useState(
    'name: CI\non: [push]\njobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n      - run: npm test\n'
  );
  const [isCreatingWf, setIsCreatingWf] = useState(false);

  const readyRunnersCount = useMemo(
    () => runners.filter((r) => r.enabled && r.status === 'online').length,
    [runners]
  );

  const latestRun = runs[0] ?? null;

  const filteredWorkflows = useMemo(() => {
    const q = workflowSearch.trim().toLowerCase();
    if (!q) return workflows;
    return workflows.filter(
      (w) =>
        w.name.toLowerCase().includes(q) ||
        (w.description && w.description.toLowerCase().includes(q)) ||
        w.trigger.toLowerCase().includes(q)
    );
  }, [workflows, workflowSearch]);

  const handleRefresh = async () => {
    setIsRefreshing(true);
    try {
      const [wfList, runList, runnerList] = await Promise.all([
        actionsService.listWorkflows(),
        actionsService.listRuns(),
        actionsService.listRunners(),
      ]);
      setWorkflows(wfList);
      setRuns(runList);
      setRunners(runnerList);
      if (selectedRun) {
        const updated = runList.find((r) => r.id === selectedRun.id);
        if (updated) setSelectedRun(updated);
      }
      toast.success('Actions data refreshed');
    } catch {
      toast.error('Failed to refresh actions');
    } finally {
      setIsRefreshing(false);
    }
  };

  const handleCreateWorkflow = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!wfName.trim()) {
      toast.error('Workflow name is required');
      return;
    }

    setIsCreatingWf(true);
    try {
      const created = await actionsService.createWorkflow({
        name: wfName.trim(),
        description: wfDescription.trim(),
        trigger: wfTrigger,
        yaml: wfYaml,
      });
      setWorkflows([created, ...workflows]);
      setWorkflowDialogOpen(false);
      setWfName('');
      setWfDescription('');
      toast.success(`Workflow ${created.name} created`);
    } catch {
      toast.error('Failed to create workflow');
    } finally {
      setIsCreatingWf(false);
    }
  };

  if (selectedRun) {
    return (
      <div className="space-y-6">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setSelectedRun(null)}
          className="gap-1.5 text-xs text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="size-4" />
          Back to Actions Overview
        </Button>
        <RunWorkspace
          run={selectedRun}
          onRunUpdated={(updated) => {
            setSelectedRun(updated);
            setRuns(runs.map((r) => (r.id === updated.id ? updated : r)));
          }}
        />
      </div>
    );
  }

  return (
    <div className="@container/actions-home space-y-6">
      {/* Header matching OpenShip */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <span className="flex size-11 shrink-0 items-center justify-center rounded-2xl bg-card text-primary shadow-xs">
            <WorkflowIcon className="size-6" />
          </span>
          <div>
            <h1 className="text-2xl font-semibold tracking-tight text-foreground">Actions</h1>
            <p className="text-sm text-muted-foreground">
              Continuous integration workflows, runner execution, and build history.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            disabled={isRefreshing}
            className="gap-1.5 rounded-xl"
          >
            <RefreshCw className={`size-3.5 ${isRefreshing ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
          <Button
            size="sm"
            onClick={() => {
              if (activeTab === 'runners') {
                setActiveTab('runners');
              } else {
                setWorkflowDialogOpen(true);
              }
            }}
            className="gap-1.5 rounded-xl"
          >
            <Plus className="size-4" />
            {activeTab === 'runners' ? 'New Runner' : 'New Workflow'}
          </Button>
        </div>
      </div>

      {/* Main 2-column Grid matching OpenShip */}
      <div className="grid items-start gap-6 @min-[980px]:grid-cols-[minmax(0,1fr)_340px]">
        {/* Left Column: Tabs & Content */}
        <div className="min-w-0 space-y-6">
          <Tabs
            value={activeTab}
            onValueChange={(val) => setActiveTab(val as 'runs' | 'workflows' | 'runners')}
            className="space-y-5"
          >
            <TabsList className="rounded-xl">
              <TabsTrigger value="runs" className="gap-2 rounded-lg">
                <Clock className="size-4" />
                Recent Runs ({runs.length})
              </TabsTrigger>
              <TabsTrigger value="workflows" className="gap-2 rounded-lg">
                <FileCode2 className="size-4" />
                Workflows ({workflows.length})
              </TabsTrigger>
              <TabsTrigger value="runners" className="gap-2 rounded-lg">
                <Server className="size-4" />
                Runners ({runners.length})
              </TabsTrigger>
            </TabsList>

            {/* Runs Tab */}
            <TabsContent value="runs" className="space-y-4">
              {runs.length === 0 ? (
                <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-border/70 bg-card p-12 text-center">
                  <Clock className="mb-2 size-8 text-muted-foreground/50" />
                  <h4 className="font-semibold text-sm">No action runs yet</h4>
                  <p className="mt-1 max-w-sm text-xs text-muted-foreground">
                    Trigger a workflow to view execution progress, logs, and artifacts here.
                  </p>
                </div>
              ) : (
                <div className="divide-y divide-border/40 overflow-hidden rounded-2xl border border-border/50 bg-card shadow-xs">
                  {runs.map((run) => (
                    <div
                      key={run.id}
                      onClick={() => setSelectedRun(run)}
                      className="group flex flex-wrap items-center gap-x-5 gap-y-2 px-5 py-4 transition-colors hover:bg-muted/40 cursor-pointer"
                    >
                      <GitBranch className="size-4 shrink-0 text-muted-foreground group-hover:text-foreground" />
                      <div className="min-w-0 flex-1">
                        <p className="truncate text-sm font-medium text-foreground">
                          {run.name}{' '}
                          <span className="text-muted-foreground">#{run.number}</span>
                        </p>
                        <p className="mt-0.5 truncate text-xs text-muted-foreground">
                          Triggered by {run.triggeredBy}
                          {run.runnerId && ` · Runner: ${run.runnerId}`}
                        </p>
                      </div>
                      <ActionStatus status={run.status} />
                      <time className="text-xs text-muted-foreground tabular-nums">
                        {run.durationSeconds}s
                      </time>
                      <ChevronRight className="size-4 text-muted-foreground/60 transition-transform group-hover:translate-x-0.5 group-hover:text-foreground" />
                    </div>
                  ))}
                </div>
              )}
            </TabsContent>

            {/* Workflows Tab */}
            <TabsContent value="workflows" className="space-y-4">
              <div className="relative w-full max-w-sm">
                <Search className="pointer-events-none absolute top-1/2 left-3.5 size-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  type="search"
                  placeholder="Search workflows…"
                  aria-label="Search workflows"
                  value={workflowSearch}
                  onChange={(e) => setWorkflowSearch(e.target.value)}
                  className="rounded-xl pl-10"
                />
              </div>

              {filteredWorkflows.length === 0 ? (
                <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-border/70 bg-card p-12 text-center">
                  <FileCode2 className="mb-2 size-8 text-muted-foreground/50" />
                  <h4 className="font-semibold text-sm">
                    {workflows.length === 0 ? 'No workflows configured' : 'No matching workflows'}
                  </h4>
                  <p className="mt-1 max-w-sm text-xs text-muted-foreground mb-4">
                    Define automated pipelines with jobs, matrix runners, and custom triggers.
                  </p>
                  <Button
                    onClick={() => setWorkflowDialogOpen(true)}
                    size="sm"
                    className="gap-1.5 rounded-xl"
                  >
                    <Plus className="size-4" />
                    Create Workflow
                  </Button>
                </div>
              ) : (
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                  {filteredWorkflows.map((wf) => (
                    <div
                      key={wf.id}
                      className="group flex flex-col justify-between rounded-2xl border border-border/60 bg-card p-5 transition-all hover:border-primary/40 hover:shadow-xs"
                    >
                      <div>
                        <div className="flex items-start justify-between gap-3">
                          <span className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
                            <PlayCircle className="size-5" />
                          </span>
                          <Badge variant="outline" className="font-mono text-[11px] uppercase">
                            {wf.trigger}
                          </Badge>
                        </div>
                        <h3 className="mt-3 truncate font-semibold text-foreground text-sm">
                          {wf.name}
                        </h3>
                        {wf.description && (
                          <p className="mt-1 line-clamp-2 text-xs text-muted-foreground">
                            {wf.description}
                          </p>
                        )}
                      </div>

                      <div className="mt-4">
                        <pre className="max-h-24 overflow-y-auto rounded-lg bg-muted/60 p-2 font-mono text-[10px] text-muted-foreground">
                          {wf.yaml || '# empty'}
                        </pre>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </TabsContent>

            {/* Runners Tab */}
            <TabsContent value="runners">
              <RunnerCatalog
                runners={runners}
                onRunnerRegistered={(newRunner) => setRunners([...runners, newRunner])}
              />
            </TabsContent>
          </Tabs>
        </div>

        {/* Right Sidebar Column: Overview & Quick Context */}
        <aside className="space-y-5 sticky top-6">
          <section className="rounded-2xl border border-border/60 bg-card p-5 shadow-xs">
            <h2 className="text-base font-semibold text-foreground">Actions Overview</h2>

            <dl className="mt-4 space-y-3 text-sm">
              <div className="flex items-center justify-between gap-3">
                <dt className="text-muted-foreground">Workflows</dt>
                <dd className="font-semibold tabular-nums text-foreground">{workflows.length}</dd>
              </div>
              <div className="flex items-center justify-between gap-3">
                <dt className="text-muted-foreground">Ready Runners</dt>
                <dd className="font-semibold tabular-nums text-foreground">
                  {readyRunnersCount}{' '}
                  <span className="font-normal text-muted-foreground">/ {runners.length}</span>
                </dd>
              </div>
              <div className="flex items-center justify-between gap-3">
                <dt className="text-muted-foreground">Total Runs</dt>
                <dd className="font-semibold tabular-nums text-foreground">{runs.length}</dd>
              </div>
            </dl>

            {latestRun && (
              <div className="mt-5 border-t border-border/50 pt-4">
                <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground">
                  Latest Run
                </p>
                <div
                  onClick={() => setSelectedRun(latestRun)}
                  className="mt-2 flex items-center justify-between gap-3 rounded-xl border border-border/40 bg-muted/30 p-2.5 transition-colors hover:bg-muted/60 cursor-pointer"
                >
                  <div className="min-w-0">
                    <p className="truncate text-xs font-semibold text-foreground">
                      Run #{latestRun.number}
                    </p>
                    <p className="text-[11px] text-muted-foreground">{latestRun.durationSeconds}s duration</p>
                  </div>
                  <ActionStatus status={latestRun.status} />
                </div>
              </div>
            )}
          </section>

          {/* Quick Automation Guide Card */}
          <section className="rounded-2xl border border-border/60 bg-card p-5 shadow-xs">
            <h3 className="text-sm font-semibold text-foreground">Workflow Guide</h3>
            <ol className="mt-3 space-y-3 text-xs text-muted-foreground">
              <li className="flex items-start gap-2.5">
                <span className="flex size-5 shrink-0 items-center justify-center rounded-full bg-primary/10 font-medium text-primary">
                  1
                </span>
                <span>Configure target runners in docker or isolated environment.</span>
              </li>
              <li className="flex items-start gap-2.5">
                <span className="flex size-5 shrink-0 items-center justify-center rounded-full bg-primary/10 font-medium text-primary">
                  2
                </span>
                <span>Define YAML workflow triggers, matrix tasks, and secrets.</span>
              </li>
              <li className="flex items-start gap-2.5">
                <span className="flex size-5 shrink-0 items-center justify-center rounded-full bg-primary/10 font-medium text-primary">
                  3
                </span>
                <span>Inspect execution logs, artifacts, and test summaries live.</span>
              </li>
            </ol>
          </section>
        </aside>
      </div>

      {/* New Workflow Dialog */}
      <Dialog open={workflowDialogOpen} onOpenChange={setWorkflowDialogOpen}>
        <DialogContent className="sm:max-w-[540px]">
          <form onSubmit={handleCreateWorkflow}>
            <DialogHeader>
              <DialogTitle>Create Action Workflow</DialogTitle>
              <DialogDescription>
                Define an automated pipeline with jobs, triggers, and execution steps.
              </DialogDescription>
            </DialogHeader>

            <div className="grid gap-4 py-4">
              <div className="space-y-1.5">
                <Label htmlFor="wf-name">Workflow Name</Label>
                <Input
                  id="wf-name"
                  value={wfName}
                  onChange={(e) => setWfName(e.target.value)}
                  placeholder="e.g. Build and Verify"
                  required
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="wf-desc">Description</Label>
                <Input
                  id="wf-desc"
                  value={wfDescription}
                  onChange={(e) => setWfDescription(e.target.value)}
                  placeholder="e.g. Automated CI checks on pull request"
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="wf-trigger">Trigger</Label>
                <Input
                  id="wf-trigger"
                  value={wfTrigger}
                  onChange={(e) => setWfTrigger(e.target.value)}
                  placeholder="push, pull_request, manual"
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="wf-yaml">YAML Definition</Label>
                <Textarea
                  id="wf-yaml"
                  rows={6}
                  value={wfYaml}
                  onChange={(e) => setWfYaml(e.target.value)}
                  className="font-mono text-xs"
                />
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setWorkflowDialogOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isCreatingWf}>
                {isCreatingWf ? 'Saving...' : 'Save Workflow'}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
