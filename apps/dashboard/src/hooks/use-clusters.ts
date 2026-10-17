import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { clustersService } from '#/services/clusters';

export function useComputeClusters(projectId: string) {
  return useQuery({
    queryKey: ['projects', projectId, 'clusters'],
    queryFn: () => clustersService.list(projectId),
    enabled: Boolean(projectId),
  });
}

export function useClusterScaling(clusterId?: string) {
  return useQuery({
    queryKey: ['system', 'compute-clusters', clusterId, 'scaling'],
    queryFn: () => (clusterId ? clustersService.getScaling(clusterId) : null),
    enabled: Boolean(clusterId),
  });
}

export function useSetupClusterScaling() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (clusterId: string) => clustersService.setupScaling(clusterId),
    onSuccess: (_, clusterId) => {
      queryClient.invalidateQueries({
        queryKey: ['system', 'compute-clusters', clusterId, 'scaling'],
      });
    },
  });
}

export function useRetryClusterScaling() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (clusterId: string) => clustersService.retryScaling(clusterId),
    onSuccess: (_, clusterId) => {
      queryClient.invalidateQueries({
        queryKey: ['system', 'compute-clusters', clusterId, 'scaling'],
      });
    },
  });
}

export function useRemoveClusterScaling() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (clusterId: string) => clustersService.removeScaling(clusterId),
    onSuccess: (_, clusterId) => {
      queryClient.invalidateQueries({
        queryKey: ['system', 'compute-clusters', clusterId, 'scaling'],
      });
    },
  });
}
