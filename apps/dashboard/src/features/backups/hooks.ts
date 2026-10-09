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

export const useUpdateS3Destination = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      id: string;
      payload: Parameters<typeof backupsService.updateS3Destination>[1];
    }) => backupsService.updateS3Destination(payload.id, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['s3-destinations'] });
    },
  });
};

export const useListSFTPDestinations = (projectId?: string) => {
  return useQuery({
    queryKey: ['sftp-destinations', projectId].filter(Boolean),
    queryFn: () => backupsService.listSFTPDestinations(projectId),
  });
};

export const useCreateSFTPDestination = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Parameters<typeof backupsService.createSFTPDestination>[0]) =>
      backupsService.createSFTPDestination(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['sftp-destinations'] });
    },
  });
};

export const useVerifySFTPDestination = () => {
  return useMutation({
    mutationFn: (id: string) => backupsService.verifySFTPDestination(id),
  });
};

export const useDeleteSFTPDestination = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => backupsService.deleteSFTPDestination(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['sftp-destinations'] });
    },
  });
};

export const usePrepareRestore = () => {
  return useMutation({
    mutationFn: (payload: { recordId: string; targetDatabaseId?: string }) =>
      backupsService.prepareRestore(payload.recordId, payload.targetDatabaseId),
  });
};

export const useRestoreOperation = (operationId: string | null, enabled = true) => {
  return useQuery({
    queryKey: ['operations', operationId].filter(Boolean),
    queryFn: () => backupsService.getOperation(operationId ?? ''),
    enabled: enabled && !!operationId,
    refetchInterval: (query) => {
      const status = query.state.data?.data?.status ?? '';
      return ['RUNNING', 'CANCELLING', 'PREPARING', 'QUEUED', 'APPLYING'].includes(status)
        ? 1500
        : false;
    },
  });
};

export const useApplyRestore = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { operationId: string; confirmation: string }) =>
      backupsService.applyRestore(payload.operationId, payload.confirmation),
    onSettled: async () => {
      await queryClient.invalidateQueries({ queryKey: ['operations'] });
    },
  });
};

export const useCancelRestore = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (operationId: string) => backupsService.cancelRestore(operationId),
    onSettled: async () => {
      await queryClient.invalidateQueries({ queryKey: ['operations'] });
    },
  });
};
