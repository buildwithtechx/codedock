import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { CreateNetworkInput, NetworkAccessPolicy } from '#/interfaces/networks';
import { networksService } from '#/services/networks';

export function useNetworks() {
  return useQuery({
    queryKey: ['system', 'networks'],
    queryFn: () => networksService.list(),
  });
}

export function useNetwork(id?: string) {
  return useQuery({
    queryKey: ['system', 'networks', id],
    queryFn: () => (id ? networksService.get(id) : null),
    enabled: Boolean(id),
  });
}

export function useNetworkWatchdog(operationId?: string) {
  return useQuery({
    queryKey: ['system', 'networks', 'watchdog', operationId],
    queryFn: () => (operationId ? networksService.getWatchdog(operationId) : null),
    enabled: Boolean(operationId),
    refetchInterval: (query) => (query.state.data?.armed ? 2000 : false),
  });
}

export function useCreateNetwork() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateNetworkInput) => networksService.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['system', 'networks'] });
    },
  });
}

export function useReviseNetworkAccess() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, accessPolicy }: { id: string; accessPolicy: NetworkAccessPolicy }) =>
      networksService.reviseAccess(id, accessPolicy),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['system', 'networks', variables.id] });
      queryClient.invalidateQueries({ queryKey: ['system', 'networks'] });
    },
  });
}

export function useCommitNetworkWatchdog() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (operationId: string) => networksService.commitWatchdog(operationId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['system', 'networks'] });
    },
  });
}

export function useRollbackNetworkWatchdog() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (operationId: string) => networksService.rollbackWatchdog(operationId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['system', 'networks'] });
    },
  });
}
