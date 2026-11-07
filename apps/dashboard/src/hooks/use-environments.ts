import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { environmentsService } from '#/services/environments';

export const projectEnvironmentsQueryOptions = (projectId: string) => ({
  queryKey: ['environments', 'listByProject', projectId].filter(Boolean),
  queryFn: () => environmentsService.listByProject(projectId),
});

export const useListByProject = (projectId: string) => {
  return useQuery(projectEnvironmentsQueryOptions(projectId));
};

export const useCreateEnvironment = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { projectId: string; name: string }) =>
      environmentsService.createEnvironment(payload.projectId, payload.name),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['environments'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useDeleteEnvironment = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { environmentId: string }) =>
      environmentsService.deleteEnvironment(payload.environmentId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['environments'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};
