import { Link } from '@tanstack/react-router';
import { ArrowLeft, Ban, Bot, Loader2, RotateCcw } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { PageHeader } from '#/components/layout/page-header';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import { useListProjects } from '#/features/projects';
import type { Deployment } from '#/features/services';
import { DeploymentFailureAi } from '#/features/services/deployment-failure-ai';
import { useGetApp } from '#/hooks/use-apps';
import { useCancelDeployment, useDeployment, useRollback } from '#/hooks/use-deployments';
import { DeploymentLogs } from './deployment-logs';
import { DeploymentRollbackDialog } from './deployment-rollback-dialog';

const inFlightStatuses = new Set(['PENDING', 'CLONING', 'PULLING', 'BUILDING']);
const rollbackStatuses = new Set(['READY', 'ACTIVE', 'SUCCESS']);

const statusTone = (status: string) => {
  const normalized = status.toUpperCase();
  if (rollbackStatuses.has(normalized)) return 'bg-emerald-500';
  if (normalized === 'FAILED') return 'bg-rose-500';
  if (inFlightStatuses.has(normalized)) return 'bg-amber-400';
  return 'bg-muted-foreground/50';
};

const statusLabel = (status: string) =>
  status
    .toLowerCase()
    .replace(/_/g, ' ')
    .replace(/\b\w/g, (letter) => letter.toUpperCase());

const formatDateTime = (value?: string) => (value ? new Date(value).toLocaleString() : '–');

export function DeploymentDetail({ deploymentId }: { deploymentId: string }) {
  const [rollbackOpen, setRollbackOpen] = useState(false);
  const { data, isLoading, isError, refetch } = useDeployment(deploymentId);
  const deployment = data?.data;
  const { data: appResponse } = useGetApp(deployment?.serviceId ?? '');
  const { data: projectsResponse } = useListProjects({ limit: 100 });
  const cancelDeployment = useCancelDeployment();
  const rollbackDeployment = useRollback();

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (isError || !deployment) {
    return (
      <QueryErrorState
        title="Deployment is unavailable"
        description="Codedock could not load this deployment."
        onRetry={() => void refetch()}
      />
    );
  }

  const serviceName = appResponse?.data?.name || 'Unknown app';
  const projectName =
    projectsResponse?.data?.records.find((project) => project.id === deployment.projectId)?.name ||
    'Unknown project';
  const normalizedStatus = deployment.status.toUpperCase();
  const isInFlight = inFlightStatuses.has(normalizedStatus);
  const canRollback = rollbackStatuses.has(normalizedStatus);
  const shortHash = deployment.commitHash?.slice(0, 7);
  const releaseLabel = [
    deployment.version != null ? `v${deployment.version}` : null,
    shortHash ? `commit ${shortHash}` : null,
    deployment.branch ? `on ${deployment.branch}` : null,
  ]
    .filter(Boolean)
    .join(' · ');

  const handleCancel = () => {
    cancelDeployment.mutate(
      { deploymentId: deployment.id },
      {
        onSuccess: () => toast.success('Deployment cancel requested'),
        onError: () => toast.error('Could not cancel this deployment'),
      }
    );
  };

  const handleRollback = () => {
    rollbackDeployment.mutate(
      { deploymentId: deployment.id },
      {
        onSuccess: () => {
          toast.success('Rollback started');
          setRollbackOpen(false);
        },
        onError: () => toast.error('Could not roll back this deployment'),
      }
    );
  };

  return (
    <div className="space-y-6">
      <Link
        to="/deployments"
        className="inline-flex items-center gap-1.5 text-muted-foreground text-sm transition-colors hover:text-foreground"
      >
        <ArrowLeft className="h-4 w-4" />
        Deployments
      </Link>
      <PageHeader
        title={serviceName}
        description={deployment.commitMessage || `Deployment ${deployment.id.slice(0, 7)}`}
        action={
          <div className="flex items-center gap-2">
            {isInFlight && (
              <Button
                variant="outline"
                onClick={handleCancel}
                disabled={cancelDeployment.isPending}
                className="gap-2"
              >
                <Ban className="h-4 w-4" />
                {cancelDeployment.isPending ? 'Cancelling...' : 'Cancel'}
              </Button>
            )}
            {canRollback && (
              <Button onClick={() => setRollbackOpen(true)} className="gap-2">
                <RotateCcw className="h-4 w-4" />
                Roll back
              </Button>
            )}
          </div>
        }
      />
      <div className="flex flex-wrap items-center gap-2">
        <Badge variant="secondary" className="gap-1.5">
          <span className={`h-1.5 w-1.5 rounded-full ${statusTone(deployment.status)}`} />
          {statusLabel(deployment.status)}
        </Badge>
        {deployment.version != null && <Badge variant="outline">v{deployment.version}</Badge>}
        {deployment.trigger && <Badge variant="outline">{deployment.trigger}</Badge>}
      </div>

      <Tabs defaultValue="overview">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="logs">Logs</TabsTrigger>
          <TabsTrigger value="explain">Explain</TabsTrigger>
        </TabsList>
        <TabsContent value="overview">
          <DeploymentOverview
            deployment={deployment}
            serviceName={serviceName}
            projectName={projectName}
          />
        </TabsContent>
        <TabsContent value="logs">
          <DeploymentLogs deploymentId={deployment.id} fallbackLogs={deployment.buildLogs} />
        </TabsContent>
        <TabsContent value="explain">
          {normalizedStatus === 'FAILED' ? (
            <DeploymentFailureAi deploymentId={deployment.id} />
          ) : (
            <div className="flex min-h-[16rem] flex-col items-center justify-center rounded-2xl bg-card px-6 py-12 text-center">
              <span className="flex h-12 w-12 items-center justify-center rounded-2xl bg-muted text-muted-foreground">
                <Bot className="h-5 w-5" />
              </span>
              <h3 className="mt-4 font-semibold">No failure to explain</h3>
              <p className="mt-1 max-w-sm text-muted-foreground text-sm">
                AI failure analysis is available when a deployment fails.
              </p>
            </div>
          )}
        </TabsContent>
      </Tabs>

      <DeploymentRollbackDialog
        open={rollbackOpen}
        releaseLabel={releaseLabel || `Deployment ${deployment.id.slice(0, 7)}`}
        busy={rollbackDeployment.isPending}
        onOpenChange={setRollbackOpen}
        onConfirm={handleRollback}
      />
    </div>
  );
}

function DeploymentOverview({
  deployment,
  serviceName,
  projectName,
}: {
  deployment: Deployment;
  serviceName: string;
  projectName: string;
}) {
  const items: { label: string; value: string; mono?: boolean }[] = [
    { label: 'Service', value: serviceName },
    { label: 'Project', value: projectName },
    { label: 'Status', value: statusLabel(deployment.status) },
    { label: 'Branch', value: deployment.branch || '–', mono: true },
    { label: 'Commit', value: deployment.commitHash || deployment.id, mono: true },
    { label: 'Trigger', value: deployment.trigger || '–' },
    { label: 'Created', value: formatDateTime(deployment.createdAt) },
    { label: 'Finished', value: formatDateTime(deployment.finishedAt) },
  ];

  return (
    <div className="space-y-4">
      {deployment.commitMessage && (
        <section className="rounded-2xl bg-card p-5">
          <h2 className="font-semibold text-sm">Commit message</h2>
          <p className="mt-2 text-muted-foreground text-sm leading-6">{deployment.commitMessage}</p>
        </section>
      )}
      <section className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
        {items.map((item) => (
          <div key={item.label} className="rounded-2xl bg-card p-4">
            <p className="font-medium text-muted-foreground text-xs uppercase tracking-[0.08em]">
              {item.label}
            </p>
            <p
              className={`mt-1.5 break-words text-sm ${
                item.mono ? 'font-mono text-[13px]' : 'font-medium'
              }`}
            >
              {item.value}
            </p>
          </div>
        ))}
      </section>
    </div>
  );
}
