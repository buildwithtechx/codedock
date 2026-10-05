import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { backupsService } from './api';

export const useList = () => {
  return useQuery({
    queryKey: ['backups'],
    queryFn: () => backupsService.list(),
  });
};

export const useGet = (id: string) => {
  return useQuery({
    queryKey: ['backups', 'get', id].filter(Boolean),
    queryFn: () => backupsService.get(id),
    enabled: !!id,
  });
};

export const useCreate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof backupsService.create>[0] }) =>
      backupsService.create(payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useUpdate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string; payload: Parameters<typeof backupsService.update>[1] }) =>
      backupsService.update(payload.id, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useDelete = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => backupsService.delete(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useTrigger = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => backupsService.trigger(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useDeleteRecord = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string; recordId: string }) =>
      backupsService.deleteRecord(payload.id, payload.recordId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useRestore = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => backupsService.restore(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useListRecords = (id: string) => {
  return useQuery({
    queryKey: ['backups', 'listRecords', id].filter(Boolean),
    queryFn: () => backupsService.listRecords(id),
    enabled: !!id,
  });
};

export const useListS3Destinations = () => {
  return useQuery({
    queryKey: ['s3-destinations'],
    queryFn: () => backupsService.listS3Destinations(),
  });
};

export const useCreateS3Destination = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof backupsService.createS3Destination>[0] }) =>
      backupsService.createS3Destination(payload.payload),
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['backups'] }),
        queryClient.invalidateQueries({ queryKey: ['s3-destinations'] }),
      ]);
    },
  });
};

export const useListAllRecords = (limit = 50) => {
  return useQuery({
    queryKey: ['backup-records', limit],
    queryFn: () => backupsService.listAllRecords(limit),
    refetchInterval: 15_000,
  });
};

export const useVerifyS3Destination = () => {
  return useMutation({
    mutationFn: (id: string) => backupsService.verifyS3Destination(id),
  });
};

export const useVerifyS3Draft = () => {
  return useMutation({
    mutationFn: (payload: Parameters<typeof backupsService.verifyS3Draft>[0]) =>
      backupsService.verifyS3Draft(payload),
  });
};

export const useDeleteS3Destination = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string; projectId: string }) =>
      backupsService.deleteS3Destination(payload.id),
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['backups'] }),
        queryClient.invalidateQueries({ queryKey: ['s3-destinations'] }),
      ]);
    },
  });
};

export const useSetDefaultS3Destination = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => backupsService.setDefaultS3Destination(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['s3-destinations'] });
    },
  });
};

export const useTriggerDatabaseBackup = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (databaseId: string) => backupsService.triggerDatabaseBackup(databaseId),
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['backups'] }),
        queryClient.invalidateQueries({ queryKey: ['backup-records'] }),
      ]);
    },
  });
};
