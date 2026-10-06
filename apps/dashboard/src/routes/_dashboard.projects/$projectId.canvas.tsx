import { createFileRoute } from '@tanstack/react-router';
import { Activity, Grid3X3 } from 'lucide-react';
import { useState } from 'react';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { EnvironmentCanvas } from '#/features/canvas/environment-canvas';
import { useGetCanvasSummary, useGetEnvironmentCanvas } from '#/hooks/use-canvas';
import { useListByProject } from '#/hooks/use-environments';

export const Route = createFileRoute('/_dashboard/projects/$projectId/canvas')({
  component: CanvasRouteComponent,
});

function CanvasRouteComponent() {
  const { projectId } = Route.useParams();
  const [selectedEnvironment, setEnvironment] = useState('');
  const environments = useListByProject(projectId);
  const {
    data: summaryRes,
    isLoading: summaryLoading,
    isError: summaryError,
    refetch: refetchSummary,
  } = useGetCanvasSummary(projectId);

  const envId =
    environments.data?.data?.find((environment) => environment.id === selectedEnvironment)?.id ??
    summaryRes?.data?.defaultEnvironment?.id;
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
      <section className="flex min-h-[28rem] flex-col items-center justify-center px-6 text-center">
        <span className="flex h-14 w-14 items-center justify-center rounded-2xl bg-card text-muted-foreground">
          <Grid3X3 className="h-6 w-6" />
        </span>
        <h1 className="mt-5 font-medium text-foreground/90 text-xl">No environment canvas yet</h1>
        <p className="mt-2 max-w-sm text-muted-foreground/75 text-sm leading-6">
          Create an environment or add project resources to see the deployment topology here.
        </p>
      </section>
    );
  }

  return (
    <div className="flex h-full flex-col gap-4">
      <h1 className="font-medium text-2xl text-foreground/90 tracking-[-0.02em]">
        Environment Canvas
      </h1>
      <select
        aria-label="Environment"
        className="w-fit rounded-md border bg-background p-2"
        value={envId ?? ''}
        onChange={(event) => setEnvironment(event.target.value)}
      >
        {environments.data?.data?.map((environment) => (
          <option key={environment.id} value={environment.id}>
            {environment.name}
          </option>
        ))}
      </select>
      <div className="flex-1 overflow-hidden rounded-2xl bg-card">
        <EnvironmentCanvas key={envId} envData={envRes.data} />
      </div>
    </div>
  );
}
