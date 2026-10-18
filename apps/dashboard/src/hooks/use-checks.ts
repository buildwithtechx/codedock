import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  type BuildAdmissionStatus,
  type ProjectChecksConfig,
  type UpdateChecksInput,
  checksService,
} from '#/services/checks';

export const useBuildAdmission = (projectId: string) => {
  return useQuery<BuildAdmissionStatus>({
    queryKey: ['projects', projectId, 'admission'],
    queryFn: () => checksService.getAdmission(projectId),
    enabled: Boolean(projectId),
  });
};

export const useProjectChecks = (projectId: string) => {
  return useQuery<ProjectChecksConfig>({
    queryKey: ['projects', projectId, 'checks'],
    queryFn: () => checksService.getChecks(projectId),
    enabled: Boolean(projectId),
  });
};

export const useUpdateProjectChecks = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, input }: { projectId: string; input: UpdateChecksInput }) =>
      checksService.updateChecks(projectId, input),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({
        queryKey: ['projects', variables.projectId, 'checks'],
      });
    },
  });
};
