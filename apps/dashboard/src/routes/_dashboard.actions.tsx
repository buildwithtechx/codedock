import { useQuery } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router';
import { Loader2 } from 'lucide-react';
import { ActionsHome } from '#/features/actions';
import { actionsService } from '#/services/actions';

export const Route = createFileRoute('/_dashboard/actions')({
  component: ActionsPage,
});

function ActionsPage() {
  const { data: workflows = [], isLoading: isLoadingWorkflows } = useQuery({
    queryKey: ['actions', 'workflows'],
    queryFn: () => actionsService.listWorkflows(),
  });

  const { data: runs = [], isLoading: isLoadingRuns } = useQuery({
    queryKey: ['actions', 'runs'],
    queryFn: () => actionsService.listRuns(),
  });

  const { data: runners = [], isLoading: isLoadingRunners } = useQuery({
    queryKey: ['actions', 'runners'],
    queryFn: () => actionsService.listRunners(),
  });

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
    />
  );
}
