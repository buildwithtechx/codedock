import { Link } from '@tanstack/react-router';
import { Activity, ArrowRight, Clock, GitBranch, GitCommit, RotateCw } from 'lucide-react';
import { useState } from 'react';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Card } from '#/components/ui/card';
import { useListByOrganization, useRollback } from '#/hooks/use-deployments';

interface ProjectDeploymentsTabProps {
  projectId: string;
}

export function ProjectDeploymentsTab({ projectId }: ProjectDeploymentsTabProps) {
  const [page, setPage] = useState(1);
  const {
    data: depsRes,
    isLoading,
    refetch,
  } = useListByOrganization({
    projectId,
    page,
    limit: 15,
  });
  const rollbackMutation = useRollback();

  const deployments = depsRes?.data?.records || [];
  const total = depsRes?.data?.total ?? 0;
  const totalPages = depsRes?.data?.totalPages ?? 1;

  const handleRollback = async (deploymentId: string) => {
    try {
      await rollbackMutation.mutateAsync({ deploymentId });
      refetch();
    } catch {
      // Handled
    }
  };

  if (isLoading) {
    return (
      <div className="flex h-60 items-center justify-center">
        <Activity className="h-6 w-6 animate-pulse text-primary" />
      </div>
    );
  }

  if (deployments.length === 0) {
    return (
      <Card className="flex min-h-60 flex-col items-center justify-center border-dashed p-8 text-center">
        <Clock className="mb-3 h-8 w-8 text-muted-foreground opacity-40" />
        <h4 className="font-semibold text-sm">No deployments recorded</h4>
        <p className="mt-1 max-w-sm text-muted-foreground text-xs">
          Trigger a deployment on your services to see real-time build and release logs here.
        </p>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="font-semibold text-foreground/90 text-sm">
          Deployment Activity ({total})
        </h3>
        <Button variant="ghost" size="sm" onClick={() => refetch()} className="h-8 gap-1.5 text-xs">
          <RotateCw className="h-3 w-3" />
          Refresh
        </Button>
      </div>

      <div className="divide-y divide-border/60 overflow-hidden rounded-xl border border-border/70 bg-card">
        {deployments.map((dep) => {
          const status = dep.status?.toLowerCase();
          const isSuccess = status === 'ready' || status === 'success' || status === 'running';
          const isFailed = status === 'failed' || status === 'error';
          const isPending = status === 'building' || status === 'deploying' || status === 'pending';

          return (
            <div
              key={dep.id}
              className="flex flex-col justify-between gap-4 p-4 transition-colors hover:bg-muted/30 sm:flex-row sm:items-center"
            >
              <div className="flex min-w-0 items-start gap-3">
                <div className="mt-1">
                  <div
                    className={`h-2.5 w-2.5 rounded-full ${
                      isSuccess
                        ? 'bg-emerald-500'
                        : isFailed
                          ? 'bg-red-500'
                          : isPending
                            ? 'animate-pulse bg-amber-500'
                            : 'bg-muted-foreground'
                    }`}
                  />
                </div>
                <div className="min-w-0">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="truncate font-semibold text-sm">
                      {dep.serviceName || dep.serviceId || 'Service'}
                    </span>
                    <Badge
                      variant={isSuccess ? 'default' : isFailed ? 'destructive' : 'secondary'}
                      className="px-1.5 py-0 font-mono text-[10px] uppercase"
                    >
                      {dep.status}
                    </Badge>
                    {dep.version && (
                      <span className="font-mono text-[11px] text-muted-foreground">
                        v{dep.version}
                      </span>
                    )}
                  </div>
                  <div className="mt-1 flex flex-wrap items-center gap-3 text-muted-foreground text-xs">
                    {dep.branch && (
                      <span className="flex items-center gap-1 font-mono">
                        <GitBranch className="h-3 w-3" />
                        {dep.branch}
                      </span>
                    )}
                    {dep.commitHash && (
                      <span className="flex items-center gap-1 font-mono">
                        <GitCommit className="h-3 w-3" />
                        {dep.commitHash.slice(0, 7)}
                      </span>
                    )}
                    {dep.trigger && <span className="max-w-50 truncate">{dep.trigger}</span>}
                    <span>
                      {new Date(dep.createdAt).toLocaleString(undefined, {
                        month: 'short',
                        day: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit',
                      })}
                    </span>
                  </div>
                </div>
              </div>

              <div className="flex shrink-0 items-center gap-2 self-end sm:self-center">
                {isSuccess && (
                  <Button
                    size="sm"
                    variant="ghost"
                    className="h-8 text-muted-foreground text-xs hover:text-foreground"
                    onClick={() => handleRollback(dep.id)}
                    disabled={rollbackMutation.isPending}
                  >
                    Rollback
                  </Button>
                )}
                {dep.serviceId && (
                  <Button size="sm" variant="outline" className="h-8 gap-1 text-xs" asChild>
                    <Link to="/services/$serviceId" params={{ serviceId: dep.serviceId }}>
                      Inspect
                      <ArrowRight className="h-3 w-3" />
                    </Link>
                  </Button>
                )}
              </div>
            </div>
          );
        })}
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-between pt-2 text-muted-foreground text-xs">
          <span>
            Page {page} of {totalPages}
          </span>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={page <= 1}
              onClick={() => setPage((p) => Math.max(1, p - 1))}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={page >= totalPages}
              onClick={() => setPage((p) => p + 1)}
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
