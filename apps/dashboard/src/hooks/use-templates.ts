import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { InstallAppInput, OneClickDeployRequest } from '#/interfaces/templates';
import { templatesService } from '#/services/templates';

export const useListOneClickApps = () => {
  return useQuery({
    queryKey: ['oneClickApps'],
    queryFn: () => templatesService.listOneClickApps(),
  });
};

export const useListExampleApps = () => {
  return useQuery({
    queryKey: ['exampleApps'],
    queryFn: () => templatesService.listExampleApps(),
  });
};

export const useOneClickApp = (appId: string) => {
  return useQuery({
    queryKey: ['oneClickApp', appId],
    queryFn: () => templatesService.getOneClickApp(appId),
    enabled: appId.length > 0,
  });
};

export const useReviewOneClickInstall = () => {
  return useMutation({
    mutationFn: (payload: InstallAppInput) => templatesService.reviewOneClickInstall(payload),
  });
};

export const useDeployOneClickApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: OneClickDeployRequest) => templatesService.deployOneClickApp(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['apps'] });
      await queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
};

export const useDeployCompose = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, file }: { projectId: string; file: File }) =>
      templatesService.deployCompose(projectId, file),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['apps'] });
    },
  });
};

export const useDeployArchive = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, file }: { projectId: string; file: File }) =>
      templatesService.deployArchive(projectId, file),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['apps'] });
    },
  });
};
