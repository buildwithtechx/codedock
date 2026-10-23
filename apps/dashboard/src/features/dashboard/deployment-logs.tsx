import { Terminal } from 'lucide-react';
import { QueryErrorState } from '#/components/query-error-state';
import { Skeleton } from '@codedock/ui/skeleton';
import { useGetLogs } from '#/hooks/use-deployments';

export function DeploymentLogs({
  deploymentId,
  fallbackLogs,
}: {
  deploymentId: string;
  fallbackLogs?: string;
}) {
  const { data, isLoading, isError, refetch } = useGetLogs(deploymentId);
  const rawData = data?.data;
  const logs =
    typeof rawData === 'string'
      ? rawData
      : (rawData as { buildLogs?: string } | undefined)?.buildLogs || fallbackLogs || '';

  if (isLoading) {
    return (
      <div className="space-y-2 rounded-2xl border border-border/50 bg-card p-5 shadow-xs">
        <Skeleton className="h-4 w-3/4" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-5/6" />
        <Skeleton className="h-4 w-2/3" />
      </div>
    );
  }

  if (isError && !fallbackLogs) {
    return (
      <QueryErrorState
        title="Build logs are unavailable"
        description="Codedock could not load logs for this deployment."
        onRetry={() => void refetch()}
        className="min-h-[16rem]"
      />
    );
  }

  if (!logs) {
    return (
      <div className="flex min-h-[16rem] flex-col items-center justify-center rounded-2xl border border-border/50 bg-card px-6 py-12 text-center shadow-xs">
        <span className="flex h-12 w-12 items-center justify-center rounded-2xl bg-muted text-muted-foreground">
          <Terminal className="h-5 w-5" />
        </span>
        <h3 className="mt-4 font-semibold">No logs yet</h3>
        <p className="mt-1 max-w-sm text-muted-foreground text-sm">
          This deployment has not produced any build output.
        </p>
      </div>
    );
  }

  return (
    <section className="overflow-hidden rounded-2xl border border-border/50 bg-card shadow-xs">
      <pre className="max-h-[36rem] overflow-auto whitespace-pre-wrap p-5 font-mono text-foreground/85 text-xs leading-6">
        {logs}
      </pre>
    </section>
  );
}
