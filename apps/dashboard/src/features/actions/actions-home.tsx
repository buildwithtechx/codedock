import {
  ArrowLeft,
  Clock,
  Code2,
  FileCode2,
  Loader2,
  Play,
  Plus,
  RefreshCw,
  Server,
  Workflow as WorkflowIcon,
} from 'lucide-react';
import { useState } from 'react';
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

  // Workflow Dialog State
  const [workflowDialogOpen, setWorkflowDialogOpen] = useState(false);
  const [wfName, setWfName] = useState('');
  const [wfDescription, setWfDescription] = useState('');
  const [wfTrigger, setWfTrigger] = useState('push');
  const [wfYaml, setWfYaml] = useState(
    'name: CI\non: [push]\njobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n      - run: npm test\n'
  );
  const [isCreatingWf, setIsCreatingWf] = useState(false);

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
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <div className="flex items-center gap-2">
            <WorkflowIcon className="size-6 text-primary" />
            <h1 className="text-2xl font-bold tracking-tight">Actions</h1>
          </div>
          <p className="text-sm text-muted-foreground mt-0.5">
            CI/CD workflow automation, parallel runner execution, and build verification.
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            disabled={isRefreshing}
            className="gap-1.5"
          >
            <RefreshCw className={`size-3.5 ${isRefreshing ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
          <Button
            size="sm"
            onClick={() => setWorkflowDialogOpen(true)}
            className="gap-1.5"
          >
            <Plus className="size-4" />
            New Workflow
          </Button>
        </div>
      </div>

      <Tabs defaultValue="runs" className="space-y-6">
        <TabsList>
          <TabsTrigger value="runs" className="gap-2">
            <Clock className="size-4" />
            Recent Runs ({runs.length})
          </TabsTrigger>
          <TabsTrigger value="workflows" className="gap-2">
            <FileCode2 className="size-4" />
            Workflows ({workflows.length})
          </TabsTrigger>
          <TabsTrigger value="runners" className="gap-2">
            <Server className="size-4" />
            Runners ({runners.length})
          </TabsTrigger>
        </TabsList>

        <TabsContent value="runs" className="space-y-4">
          {runs.length === 0 ? (
            <Card className="border-dashed">
              <CardContent className="flex flex-col items-center justify-center p-8 text-center">
                <Clock className="size-8 text-muted-foreground/60 mb-2" />
                <h4 className="font-semibold text-sm">No action runs yet</h4>
                <p className="text-xs text-muted-foreground max-w-sm mt-1">
                  Trigger a workflow to view build progress, logs, and artifacts here.
                </p>
              </CardContent>
            </Card>
          ) : (
            <div className="space-y-2">
              {runs.map((run) => (
                <div
                  key={run.id}
                  onClick={() => setSelectedRun(run)}
                  className="flex flex-col gap-3 rounded-lg border border-border/70 bg-card p-4 transition-all hover:border-border hover:bg-muted/40 cursor-pointer sm:flex-row sm:items-center sm:justify-between"
                >
                  <div className="flex items-start gap-3">
                    <div className="mt-0.5">
                      <ActionStatus status={run.status} />
                    </div>
                    <div>
                      <div className="flex items-center gap-2">
                        <span className="font-semibold text-sm hover:underline">{run.name}</span>
                        <span className="text-xs font-mono text-muted-foreground">#{run.number}</span>
                      </div>
                      <p className="text-xs text-muted-foreground">
                        Triggered by <span className="font-medium">{run.triggeredBy}</span>
                        {run.runnerId && ` • Runner: ${run.runnerId}`}
                      </p>
                    </div>
                  </div>

                  <div className="flex items-center gap-4 text-xs text-muted-foreground sm:text-right">
                    <span>{run.durationSeconds}s duration</span>
                    <Button variant="ghost" size="sm" className="h-8 text-xs">
                      View Logs
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="workflows" className="space-y-4">
          {workflows.length === 0 ? (
            <Card className="border-dashed">
              <CardContent className="flex flex-col items-center justify-center p-8 text-center">
                <FileCode2 className="size-8 text-muted-foreground/60 mb-2" />
                <h4 className="font-semibold text-sm">No workflows configured</h4>
                <p className="text-xs text-muted-foreground max-w-sm mt-1 mb-4">
                  Define your first YAML pipeline to automate checks and deployments.
                </p>
                <Button onClick={() => setWorkflowDialogOpen(true)} size="sm" className="gap-1.5">
                  <Plus className="size-4" />
                  Create Workflow
                </Button>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4 sm:grid-cols-2">
              {workflows.map((wf) => (
                <Card key={wf.id} className="border-border/70">
                  <CardHeader className="p-4 pb-2">
                    <div className="flex items-center justify-between">
                      <CardTitle className="text-base font-semibold">{wf.name}</CardTitle>
                      <Badge variant="outline" className="text-xs font-mono">
                        {wf.trigger}
                      </Badge>
                    </div>
                    {wf.description && (
                      <CardDescription className="text-xs">{wf.description}</CardDescription>
                    )}
                  </CardHeader>
                  <CardContent className="p-4 pt-2">
                    <pre className="max-h-32 overflow-y-auto rounded-md bg-muted/60 p-2 font-mono text-[11px] text-muted-foreground">
                      {wf.yaml || '# No definition provided'}
                    </pre>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="runners">
          <RunnerCatalog
            runners={runners}
            onRunnerRegistered={(newRunner) => setRunners([...runners, newRunner])}
          />
        </TabsContent>
      </Tabs>

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
                  placeholder="e.g. Test and Build"
                  value={wfName}
                  onChange={(e) => setWfName(e.target.value)}
                  required
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="wf-desc">Description</Label>
                <Input
                  id="wf-desc"
                  placeholder="e.g. Automated test suite and image build on push"
                  value={wfDescription}
                  onChange={(e) => setWfDescription(e.target.value)}
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="wf-trigger">Event Trigger</Label>
                <Input
                  id="wf-trigger"
                  placeholder="e.g. push, pull_request, schedule"
                  value={wfTrigger}
                  onChange={(e) => setWfTrigger(e.target.value)}
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="wf-yaml">YAML Definition</Label>
                <Textarea
                  id="wf-yaml"
                  rows={8}
                  className="font-mono text-xs"
                  value={wfYaml}
                  onChange={(e) => setWfYaml(e.target.value)}
                />
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setWorkflowDialogOpen(false)}
                disabled={isCreatingWf}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isCreatingWf}>
                {isCreatingWf && <Loader2 className="mr-2 size-4 animate-spin" />}
                Save Workflow
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
