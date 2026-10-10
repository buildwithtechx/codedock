import { Link } from '@tanstack/react-router';
import { ArrowRight, CheckCircle2, Loader2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { Skeleton } from '#/components/ui/skeleton';
import {
  findMissingRequired,
  InstallSettingsTabs,
} from '#/features/templates/install-settings-tabs';
import type { InstallPreview } from '#/interfaces/templates';
import { useCatalogApp, useDeployCatalogApp, useReviewCatalogInstall } from './hooks';
import { InstallDomainsRouting } from './install-domains-routing';
import { WizardDestination } from './wizard-destination';
import { WizardHeader } from './wizard-header';

type Phase = 'form' | 'review' | 'done';

export function InstallWizard({ appId }: { appId: string }) {
  const { data: app, isLoading, isError, refetch } = useCatalogApp(appId);
  const reviewInstall = useReviewCatalogInstall();
  const deployApp = useDeployCatalogApp();
  const [phase, setPhase] = useState<Phase>('form');
  const [name, setName] = useState('');
  const [values, setValues] = useState<Record<string, string>>({});
  const [hostPort, setHostPort] = useState('');
  const [domain, setDomain] = useState('');
  const [projectId, setProjectId] = useState('');
  const [environmentId, setEnvironmentId] = useState('');
  const [preview, setPreview] = useState<InstallPreview | null>(null);
  const [deployMessage, setDeployMessage] = useState('');

  useEffect(() => {
    if (app) {
      setName(app.name);
      setValues({});
      setHostPort('');
      setDomain('');
      setPreview(null);
      setPhase('form');
    }
  }, [app]);

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-12 w-64" />
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
          <Skeleton className="h-96 rounded-2xl" />
          <Skeleton className="h-64 rounded-2xl" />
        </div>
      </div>
    );
  }

  if (isError || !app) {
    return (
      <div className="space-y-6">
        <Link
          to="/apps/new"
          className="inline-flex items-center gap-1.5 text-muted-foreground text-sm transition-colors hover:text-foreground"
        >
          Back to catalog
        </Link>
        <QueryErrorState
          title="App details are unavailable"
          description="Codedock could not load this catalog entry. It may have been removed."
          onRetry={() => void refetch()}
        />
      </div>
    );
  }

  const setValue = (key: string, value: string) =>
    setValues((current) => ({ ...current, [key]: value }));

  const buildPayload = (digest?: string) => {
    const secrets: Record<string, string> = {};
    const environment: Record<string, string> = {};
    for (const variable of app.envVariables) {
      const value = (values[variable.key] ?? '').trim();
      if (!value) continue;
      if (variable.secret) secrets[variable.key] = value;
      else environment[variable.key] = value;
    }
    const port = Number.parseInt(hostPort, 10);
    const resolvedProjectId = projectId === 'standalone' ? '' : projectId;
    return {
      appId,
      projectId: resolvedProjectId,
      name: (name ?? '').trim() || app.name,
      secrets,
      environment,
      hostPort: Number.isNaN(port) ? undefined : port,
      domain: (domain ?? '').trim() || undefined,
      environmentId: resolvedProjectId ? environmentId || undefined : undefined,
      digest,
    };
  };

  const configurationValid = () => {
    if (!name.trim()) {
      toast.error('Install name is required');
      return false;
    }
    const missing = findMissingRequired(app.envVariables, values);
    if (missing.length > 0) {
      toast.error(`${missing[0].label || missing[0].key} is required`);
      return false;
    }
    if (hostPort.trim()) {
      const port = Number.parseInt(hostPort, 10);
      if (Number.isNaN(port) || port < 1 || port > 65535) {
        toast.error('Host port must be between 1 and 65535');
        return false;
      }
    }
    return true;
  };

  const handleReview = async () => {
    if (!configurationValid()) return;
    try {
      setPreview(await reviewInstall.mutateAsync(buildPayload()));
      setPhase('review');
    } catch {
      toast.error('Install review failed');
    }
  };

  const handleDeploy = async () => {
    if (!preview) return;
    try {
      const response = await deployApp.mutateAsync(buildPayload(preview.digest));
      setDeployMessage(response.message || `${app.name} deployed`);
      setPhase('done');
    } catch {
      toast.error(`Failed to deploy ${app.name}`);
    }
  };

  if (phase === 'done') {
    return (
      <div className="mx-auto max-w-xl py-12 text-center">
        <CheckCircle2 className="mx-auto size-12 text-emerald-500" />
        <h2 className="mt-4 font-semibold text-2xl tracking-tight">{deployMessage}</h2>
        <p className="mt-2 text-muted-foreground text-sm">
          {name} is installing in your project. Track its progress from the project canvas.
        </p>
        <div className="mt-6 flex flex-wrap items-center justify-center gap-3">
          <Button asChild>
            <Link to="/projects/$projectId" params={{ projectId }}>
              Open project
            </Link>
          </Button>
          <Button asChild variant="outline">
            <Link to="/apps">Back to apps</Link>
          </Button>
        </div>
      </div>
    );
  }

  if (phase === 'review' && preview) {
    return (
      <div className="space-y-6">
        <WizardHeader app={app} />
        <div className="max-w-3xl space-y-4 rounded-2xl bg-card p-5">
          <h2 className="font-semibold text-foreground">Review install plan</h2>
          <div className="flex flex-wrap gap-2">
            {preview.services.map((service) => (
              <Badge key={service.service} variant="secondary">
                {service.service} · {service.image}
              </Badge>
            ))}
          </div>
          {preview.generatedSecrets.length > 0 && (
            <p className="text-muted-foreground text-sm">
              Generated secrets: {preview.generatedSecrets.join(', ')}
            </p>
          )}
          <pre className="max-h-72 overflow-auto rounded-md border bg-muted p-3 font-mono text-xs">
            {preview.composeYaml}
          </pre>
          <div className="flex justify-end gap-2">
            <Button
              variant="outline"
              onClick={() => setPhase('form')}
              disabled={deployApp.isPending}
            >
              Back
            </Button>
            <Button onClick={() => void handleDeploy()} disabled={deployApp.isPending}>
              {deployApp.isPending ? (
                <>
                  <Loader2 className="size-4 animate-spin" />
                  Deploying…
                </>
              ) : (
                'Deploy'
              )}
            </Button>
          </div>
        </div>
      </div>
    );
  }

  const busy = reviewInstall.isPending;
  return (
    <div className="space-y-6">
      <WizardHeader app={app} />
      <div className="grid grid-cols-1 items-start gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
        <div className="min-w-0 space-y-5">
          <div className="rounded-2xl bg-card p-5">
            <Label htmlFor="wizard-name">Install name</Label>
            <Input
              id="wizard-name"
              value={name}
              disabled={busy}
              onChange={(event) => setName(event.target.value)}
              placeholder={app.name}
              className="mt-2"
            />
          </div>
          {app.envVariables.length > 0 && (
            <div className="rounded-2xl bg-card p-5">
              <h2 className="mb-4 font-semibold text-foreground text-sm">Configuration</h2>
              <InstallSettingsTabs
                variables={app.envVariables}
                values={values}
                onChange={setValue}
                disabled={busy}
              />
            </div>
          )}
          <InstallDomainsRouting
            hostPort={hostPort}
            onHostPortChange={setHostPort}
            defaultPort={app.defaultPort}
            domain={domain}
            onDomainChange={setDomain}
            appName={name || app.name}
            disabled={busy}
          />
        </div>
        <div className="min-w-0 space-y-4 lg:sticky lg:top-6">
          <div className="rounded-2xl bg-card p-5">
            <h2 className="font-semibold text-foreground text-sm">Destination</h2>
            <p className="mt-0.5 text-muted-foreground text-xs">
              The project this app installs into.
            </p>
            <div className="mt-4">
              <WizardDestination
                projectId={projectId}
                environmentId={environmentId}
                onProjectChange={(next) => {
                  setProjectId(next);
                  setEnvironmentId('');
                }}
                onEnvironmentChange={setEnvironmentId}
                disabled={busy}
              />
            </div>
          </div>
          <Button className="w-full" onClick={() => void handleReview()} disabled={busy}>
            {busy ? (
              <>
                <Loader2 className="size-4 animate-spin" />
                Reviewing…
              </>
            ) : (
              <>
                Review install plan
                <ArrowRight className="size-4" />
              </>
            )}
          </Button>
        </div>
      </div>
    </div>
  );
}
