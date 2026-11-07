import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { canvasService } from '#/services/canvas';
import { useOrganizationStore } from '#/stores/organization-store';

export const useListCanvasSummaries = () => {
  const activeOrganizationId = useOrganizationStore((state) => state.activeOrganizationId);

  return useQuery({
    queryKey: ['canvas', 'listCanvasSummaries', activeOrganizationId],
    queryFn: () => canvasService.listCanvasSummaries(activeOrganizationId),
    enabled: Boolean(activeOrganizationId),
  });
};

export const useGetCanvasSummary = (projectId: string) => {
  return useQuery({
    queryKey: ['canvas', 'getCanvasSummary', projectId].filter(Boolean),
    queryFn: () => canvasService.getCanvasSummary(projectId),
    enabled: Boolean(projectId),
  });
};

export const useGetEnvironmentCanvas = (envId: string) => {
  return useQuery({
    queryKey: ['canvas', 'getEnvironmentCanvas', envId].filter(Boolean),
    queryFn: () => canvasService.getEnvironmentCanvas(envId),
    enabled: Boolean(envId),
    refetchInterval: 5000,
  });
};

export const useObserveCanvasResource = (
  envId: string,
  nodeId: string,
  options?: { enabled?: boolean }
) => {
  return useQuery({
    queryKey: ['canvas-observation', envId, nodeId],
    queryFn: () => canvasService.observeResource(envId, nodeId),
    enabled: options?.enabled ?? (Boolean(envId) && Boolean(nodeId)),
    refetchInterval: 5000,
  });
};

export const useExecuteResourceAction = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ path, operation }: { path: string; operation: string }) =>
      canvasService.executeAction(path, operation),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};
