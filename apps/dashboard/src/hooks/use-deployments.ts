import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { ListOrganizationDeploymentsParams } from '#/features/services';
import { deploymentsService } from '#/services/deployments';
import { useOrganizationStore } from '#/stores/organization-store';

export const useListByOrganization = (params: ListOrganizationDeploymentsParams = {}) => {
  const activeOrganizationId = useOrganizationStore((state) => state.activeOrganizationId);

  return useQuery({
    queryKey: ['deployments', 'listByOrganization', activeOrganizationId, params],
    queryFn: () => deploymentsService.listByOrganization(activeOrganizationId as string, params),
    enabled: Boolean(activeOrganizationId),
  });
};

export const useListByService = (serviceId: string) => {
  return useQuery({
    queryKey: ['deployments', 'listByService', serviceId].filter(Boolean),
    queryFn: () => deploymentsService.listByService(serviceId),
    enabled: Boolean(serviceId),
  });
};

export const useTrigger = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      serviceId: string;
      payload?: Parameters<typeof deploymentsService.trigger>[1];
    }) => deploymentsService.trigger(payload.serviceId, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['deployments'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useRollback = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { deploymentId: string }) =>
      deploymentsService.rollback(payload.deploymentId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['deployments'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useGetLogs = (deploymentId: string) => {
  return useQuery({
    queryKey: ['deployments', 'getLogs', deploymentId].filter(Boolean),
    queryFn: () => deploymentsService.getLogs(deploymentId),
  });
};

export const useGetMetrics = (serviceId: string) => {
  return useQuery({
    queryKey: ['deployments', 'getMetrics', serviceId].filter(Boolean),
    queryFn: () => deploymentsService.getMetrics(serviceId),
    refetchInterval: 3000,
  });
};

export const useDiagnostics = () => {
  return useMutation({
    mutationFn: (deploymentId: string) => deploymentsService.diagnostics(deploymentId),
  });
};

export const useExplainFailure = (deploymentId: string, enabled: boolean) => {
  return useQuery({
    queryKey: ['deployments', 'explainFailure', deploymentId].filter(Boolean),
    queryFn: () => deploymentsService.explainFailure(deploymentId),
    enabled: enabled && !!deploymentId,
  });
};
