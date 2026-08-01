import { useQuery } from '@tanstack/react-query';
import { canvasService } from '#/services/canvas';
import { useOrganizationStore } from '#/stores/organization-store';

export const useListCanvasSummaries = () => {
  const activeOrganizationId = useOrganizationStore((state) => state.activeOrganizationId);

  return useQuery({
    queryKey: ['canvas', 'listCanvasSummaries', activeOrganizationId],
    queryFn: () => canvasService.listCanvasSummaries(activeOrganizationId),
  });
};

export const useGetCanvasSummary = (projectId: string) => {
  return useQuery({
    queryKey: ['canvas', 'getCanvasSummary', projectId].filter(Boolean),
    queryFn: () => canvasService.getCanvasSummary(projectId),
  });
};

export const useGetEnvironmentCanvas = (envId: string) => {
  return useQuery({
    queryKey: ['canvas', 'getEnvironmentCanvas', envId].filter(Boolean),
    queryFn: () => canvasService.getEnvironmentCanvas(envId),
  });
};
