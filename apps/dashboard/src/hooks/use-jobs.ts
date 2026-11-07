import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { CreateJobRequest, UpdateJobRequest } from '#/interfaces/services';
import { jobsApi } from '#/services/jobs';

export const projectJobsQueryOptions = (projectId: string) => ({
  queryKey: ['jobs', 'project', projectId],
  queryFn: () => jobsApi.listByProject(projectId),
  enabled: Boolean(projectId),
});

export const useListJobsByProject = (projectId: string) => {
  return useQuery(projectJobsQueryOptions(projectId));
};

export const useGetJob = (id: string) => {
  return useQuery({
    queryKey: ['jobs', id],
    queryFn: () => jobsApi.get(id),
    enabled: Boolean(id),
  });
};

export const useCreateJob = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateJobRequest) => jobsApi.create(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['jobs'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useUpdateJob = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateJobRequest }) =>
      jobsApi.update(id, payload),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({ queryKey: ['jobs'] });
      await queryClient.invalidateQueries({ queryKey: ['jobs', variables.id] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useDeleteJob = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => jobsApi.remove(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['jobs'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useTriggerJob = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => jobsApi.trigger(id),
    onSuccess: async (_, id) => {
      await queryClient.invalidateQueries({ queryKey: ['jobs'] });
      await queryClient.invalidateQueries({ queryKey: ['jobs', id] });
    },
  });
};
