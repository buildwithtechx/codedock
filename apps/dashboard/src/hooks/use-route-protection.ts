import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { RouteProtectionInput } from '@codedock/contracts';
import { routeProtectionService } from '#/services/route-protection';

export const useRouteProtection = (projectId?: string) => {
  return useQuery({
    queryKey: ['projects', projectId, 'route-protection'],
    queryFn: async () => {
      if (!projectId) return [];
      const res = await routeProtectionService.listRouteProtection(projectId);
      return res.data;
    },
    enabled: Boolean(projectId),
  });
};

export const useChangeRouteProtection = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, payload }: { projectId: string; payload: RouteProtectionInput }) =>
      routeProtectionService.changeRouteProtection(projectId, payload),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({
        queryKey: ['projects', variables.projectId, 'route-protection'],
      });
    },
  });
};
