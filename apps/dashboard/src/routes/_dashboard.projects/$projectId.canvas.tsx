import { createFileRoute } from '@tanstack/react-router';
import { Activity, Grid3X3 } from 'lucide-react';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { WorkspaceEmptyState } from '#/components/ui/workspace-empty-state';
import { EnvironmentCanvas } from '#/features/canvas/environment-canvas';
import { useGetCanvasSummary, useGetEnvironmentCanvas } from '#/hooks/use-canvas';

export const Route = createFileRoute('/_dashboard/projects/$projectId/canvas')({
  component: CanvasRouteComponent,
});

function CanvasRouteComponent() {
  const { projectId } = Route.useParams();
  const {
    data: summaryRes,
    isLoading: summaryLoading,
    isError: summaryError,
    refetch: refetchSummary,
  } = useGetCanvasSummary(projectId);

  const envId = summaryRes?.data?.defaultEnvironment?.id;
  const {
    data: envRes,
    isLoading: envLoading,
    isError: environmentError,
    refetch: refetchEnvironment,
  } = useGetEnvironmentCanvas(envId || '');

  if (summaryLoading || (envId && envLoading)) {
    return (
      <div className="flex min-h-[25rem] items-center justify-center">
        <Activity className="h-5 w-5 animate-pulse text-muted-foreground" />
      </div>
    );
  }

  if (summaryError || (envId && environmentError)) {
    return (
      <QueryErrorState
        title="Project canvas is unavailable"
        description="Codedock could not load this project's environment topology."
        onRetry={() => {
          void refetchSummary();
          if (envId) {
            void refetchEnvironment();
          }
        }}
      />
    );
  }

  if (!envRes?.data) {
    return (
      <WorkspaceEmptyState
        icon={Grid3X3}
        title="No environment canvas yet"
        description="Create an environment or add project resources to see the deployment topology here."
      />
    );
  }

  return (
    <div className="flex h-full flex-col gap-4 p-4">
      <h1 className="font-bold text-2xl">Environment Canvas</h1>
      <div className="flex-1 overflow-hidden rounded-2xl border border-border bg-card shadow-sm">
        <EnvironmentCanvas envData={envRes.data} />
      </div>
    </div>
  );
}
