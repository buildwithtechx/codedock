import { Link, useNavigate } from '@tanstack/react-router';
import { ArrowLeft, Check, Code2, PlayCircle, Sparkles, Terminal } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Textarea } from '#/components/ui/textarea';
import { actionsService } from '#/services/actions';

const STARTER_CI = `name: CI
on: [push, workflow_dispatch]
jobs:
  check:
    runs-on: [self-hosted, codedock, linux]
    steps:
      - uses: actions/checkout@v4
      - run: echo 'Building and verifying application...'
      - run: echo 'CI checks passed successfully'
`;

const STANDALONE_AUTOMATION = `name: Automation
on: workflow_dispatch
jobs:
  run:
    runs-on: [self-hosted, linux]
    steps:
      - run: echo 'Executing standalone automated task...'
      - run: echo 'Automation completed'
`;

export function WorkflowBuilder() {
  const navigate = useNavigate();
  const [name, setName] = useState('CI');
  const [description, setDescription] = useState('Continuous integration pipeline');
  const [mode, setMode] = useState<'repository' | 'standalone'>('repository');
  const [repo, setRepo] = useState('my-org/my-app');
  const [branch, setBranch] = useState('main');
  const [path, setPath] = useState('.github/workflows/ci.yml');
  const [trigger, setTrigger] = useState('push, workflow_dispatch');
  const [yaml, setYaml] = useState(STARTER_CI);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleModeChange = (newMode: 'repository' | 'standalone') => {
    setMode(newMode);
    if (newMode === 'standalone') {
      setName('Automation');
      setPath('.codedock/workflows/automation.yml');
      setYaml(STANDALONE_AUTOMATION);
      setTrigger('workflow_dispatch');
    } else {
      setName('CI');
      setPath('.github/workflows/ci.yml');
      setYaml(STARTER_CI);
      setTrigger('push, workflow_dispatch');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Workflow name is required');
      return;
    }
    setIsSubmitting(true);
    try {
      const parts = repo.split('/');
      const owner = mode === 'repository' && parts.length > 1 ? parts[0] : '';
      const repoName = mode === 'repository' && parts.length > 1 ? parts.slice(1).join('/') : repo;

      await actionsService.createWorkflow({
        name: name.trim(),
        description: description.trim(),
        path: path.trim(),
        ref: branch.trim(),
        owner,
        repo: repoName,
        trigger: trigger.trim(),
        yamlContent: yaml,
      });

      toast.success('Workflow created successfully');
      navigate({ to: '/actions' as never });
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to create workflow');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      <header className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <Button
            asChild
            variant="ghost"
            size="sm"
            className="-ms-3 mb-2 gap-1.5 rounded-lg text-muted-foreground hover:text-foreground"
          >
            <Link to={'/actions' as never}>
              <ArrowLeft className="size-4 rtl:rotate-180" />
              Back to Actions
            </Link>
          </Button>
          <h1 className="font-medium text-2xl text-foreground tracking-tight">New workflow</h1>
          <p className="mt-1 text-muted-foreground text-sm">
            Define an automated workflow with jobs, steps and triggers.
          </p>
        </div>
        <Button
          type="button"
          variant="ghost"
          onClick={() => navigate({ to: '/actions' as never })}
          className="rounded-xl"
        >
          Cancel
        </Button>
      </header>

      <form onSubmit={handleSubmit} className="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
        <fieldset disabled={isSubmitting} className="m-0 min-w-0 space-y-5 rounded-2xl border-0 bg-card p-5">
          <div className="flex gap-3">
            <button
              type="button"
              onClick={() => handleModeChange('repository')}
              className={`flex-1 rounded-xl border p-3.5 text-start transition-all ${
                mode === 'repository' ? 'border-primary bg-primary/5' : 'border-border/60 bg-muted/20'
              }`}
            >
              <div className="flex items-center gap-2">
                <Code2 className="size-4 text-primary" />
                <span className="font-medium text-foreground text-sm">Repository Workflow</span>
              </div>
              <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
                Triggered on git push or pull requests in a repo.
              </p>
            </button>

            <button
              type="button"
              onClick={() => handleModeChange('standalone')}
              className={`flex-1 rounded-xl border p-3.5 text-start transition-all ${
                mode === 'standalone' ? 'border-primary bg-primary/5' : 'border-border/60 bg-muted/20'
              }`}
            >
              <div className="flex items-center gap-2">
                <Terminal className="size-4 text-primary" />
                <span className="font-medium text-foreground text-sm">Standalone Workflow</span>
              </div>
              <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
                Executed manually or on cron schedule.
              </p>
            </button>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="wf-name">Workflow name</Label>
              <Input
                id="wf-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. Build & Test"
                required
                className="h-10 rounded-xl bg-muted/40"
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="wf-path">Definition file path</Label>
              <Input
                id="wf-path"
                value={path}
                onChange={(e) => setPath(e.target.value)}
                placeholder=".github/workflows/ci.yml"
                required
                className="h-10 rounded-xl bg-muted/40 font-mono text-xs"
              />
            </div>
          </div>

          {mode === 'repository' && (
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <Label htmlFor="wf-repo">Repository (owner/repo)</Label>
                <Input
                  id="wf-repo"
                  value={repo}
                  onChange={(e) => setRepo(e.target.value)}
                  placeholder="acme/web"
                  required
                  className="h-10 rounded-xl bg-muted/40"
                />
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="wf-branch">Target branch</Label>
                <Input
                  id="wf-branch"
                  value={branch}
                  onChange={(e) => setBranch(e.target.value)}
                  placeholder="main"
                  required
                  className="h-10 rounded-xl bg-muted/40"
                />
              </div>
            </div>
          )}

          <div className="space-y-1.5">
            <Label htmlFor="wf-yaml">Workflow YAML</Label>
            <Textarea
              id="wf-yaml"
              value={yaml}
              onChange={(e) => setYaml(e.target.value)}
              rows={10}
              className="rounded-xl bg-muted/40 font-mono text-xs leading-relaxed"
            />
          </div>
        </fieldset>

        <aside className="space-y-5 rounded-2xl bg-card p-5">
          <div className="flex items-center gap-3">
            <PlayCircle className="size-6 text-primary" />
            <div className="min-w-0">
              <h2 className="truncate font-medium text-base text-foreground">{name || 'New workflow'}</h2>
              <p className="mt-0.5 text-muted-foreground text-xs">{mode === 'repository' ? repo : 'Standalone'}</p>
            </div>
          </div>

          <dl className="space-y-3 text-sm">
            <div className="flex items-center justify-between gap-3">
              <dt className="text-muted-foreground">Trigger</dt>
              <dd className="font-medium text-foreground">{trigger}</dd>
            </div>
            <div className="flex items-center justify-between gap-3">
              <dt className="text-muted-foreground">Type</dt>
              <dd className="font-medium text-foreground">{mode}</dd>
            </div>
          </dl>

          <Button type="submit" disabled={isSubmitting || !name.trim()} className="w-full rounded-xl">
            {isSubmitting ? 'Saving workflow...' : 'Save and activate workflow'}
          </Button>
        </aside>
      </form>
    </div>
  );
}
