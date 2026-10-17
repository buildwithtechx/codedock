import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  actionsService,
  type CreateRunnerInput,
  type CreateWorkflowInput,
} from '#/services/actions';

export function useWorkflows() {
  return useQuery({
    queryKey: ['actions', 'workflows'],
    queryFn: () => actionsService.listWorkflows(),
  });
}

export function useRuns() {
  return useQuery({
    queryKey: ['actions', 'runs'],
    queryFn: () => actionsService.listRuns(),
  });
}

export function useRunners() {
  return useQuery({
    queryKey: ['actions', 'runners'],
    queryFn: () => actionsService.listRunners(),
  });
}

export function useInspectDestination(serverId?: string) {
  return useQuery({
    queryKey: ['actions', 'destinations', serverId, 'inspect'],
    queryFn: () => (serverId ? actionsService.inspectDestination(serverId) : null),
    enabled: Boolean(serverId),
  });
}

export function useCreateRunner() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateRunnerInput) => actionsService.createRunner(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['actions', 'runners'] });
    },
  });
}

export function useCreateWorkflow() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateWorkflowInput) => actionsService.createWorkflow(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['actions', 'workflows'] });
    },
  });
}
