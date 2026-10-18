import { createFileRoute } from '@tanstack/react-router';
import { Loader2 } from 'lucide-react';
import { ActionsHome } from '#/features/actions';
import { useRunners, useRuns, useWorkflows } from '#/hooks/use-actions';
import { useGetPublicSettings } from '#/hooks/use-settings';

export const Route = createFileRoute('/_dashboard/actions/')({
  component: ActionsIndexPage,
});

function ActionsIndexPage() {
  const { data: publicRes } = useGetPublicSettings();
  const isCloud = publicRes?.data?.cloudMode ?? false;

  const { data: workflows = [], isLoading: isLoadingWorkflows } = useWorkflows();
  const { data: runs = [], isLoading: isLoadingRuns } = useRuns();
  const { data: runners = [], isLoading: isLoadingRunners } = useRunners();

  if (isLoadingWorkflows || isLoadingRuns || isLoadingRunners) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="size-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <ActionsHome
      initialWorkflows={workflows}
      initialRuns={runs}
      initialRunners={runners}
      isCloud={isCloud}
    />
  );
}
