import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { CreateRunnerInput, CreateWorkflowInput } from '#/interfaces/actions';
import { actionsService } from '#/services/actions';

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

export function useTriggerWorkflow() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (workflowId: string) => actionsService.triggerWorkflow(workflowId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['actions', 'runs'] });
    },
  });
}

export function useDeleteWorkflow() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (workflowId: string) => actionsService.deleteWorkflow(workflowId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['actions', 'workflows'] });
    },
  });
}

export function useDeleteRunner() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (runnerId: string) => actionsService.deleteRunner(runnerId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['actions', 'runners'] });
    },
  });
}

export function useRetryRun() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (runId: string) => actionsService.retryRun(runId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['actions', 'runs'] });
    },
  });
}

export function useCancelRun() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (runId: string) => actionsService.cancelRun(runId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['actions', 'runs'] });
    },
  });
}
