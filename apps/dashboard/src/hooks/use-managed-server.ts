import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  controlManagedWorkload,
  createManagedWorkload,
  enableManagedRuntime,
  getManagedBootLogs,
  getManagedInfo,
  getManagedRuntimeCredential,
  getManagedRuntimeStatus,
  getManagedSshStatus,
  getManagedWorkloadLogs,
  getManagedWorkloads,
  getServerNetworkSettings,
  rotateManagedRuntimeCredential,
  setManagedSsh,
  setManagedSshKey,
  setManagedSshPassword,
  updateServerNetworkSettings,
} from '#/services/managed-server';

export const useServerNetworkSettings = (serverId: string) =>
  useQuery({
    queryKey: ['servers', serverId, 'network-settings'],
    queryFn: () => getServerNetworkSettings(serverId),
    enabled: Boolean(serverId),
  });

export const useUpdateServerNetworkSettings = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      serverId,
      payload,
    }: {
      serverId: string;
      payload: Parameters<typeof updateServerNetworkSettings>[1];
    }) => updateServerNetworkSettings(serverId, payload),
    onSuccess: (_, vars) => {
      qc.invalidateQueries({ queryKey: ['servers', vars.serverId, 'network-settings'] });
    },
  });
};

export const useManagedInfo = (serverId: string) =>
  useQuery({
    queryKey: ['servers', serverId, 'managed-info'],
    queryFn: () => getManagedInfo(serverId),
    enabled: Boolean(serverId),
  });

export const useManagedBootLogs = (serverId: string, tail = 200) =>
  useQuery({
    queryKey: ['servers', serverId, 'managed-boot-logs', tail],
    queryFn: () => getManagedBootLogs(serverId, tail),
    enabled: Boolean(serverId),
  });

export const useFetchManagedBootLogs = () =>
  useMutation({
    mutationFn: ({ serverId, tail = 200 }: { serverId: string; tail?: number }) =>
      getManagedBootLogs(serverId, tail),
  });

export const useManagedSshStatus = (serverId: string) =>
  useQuery({
    queryKey: ['servers', serverId, 'managed-ssh'],
    queryFn: () => getManagedSshStatus(serverId),
    enabled: Boolean(serverId),
  });

export const useSetManagedSsh = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      serverId,
      payload,
    }: {
      serverId: string;
      payload: Parameters<typeof setManagedSsh>[1];
    }) => setManagedSsh(serverId, payload),
    onSuccess: (_, vars) => {
      qc.invalidateQueries({ queryKey: ['servers', vars.serverId, 'managed-ssh'] });
    },
  });
};

export const useSetManagedSshKey = () =>
  useMutation({
    mutationFn: ({ serverId, publicKey }: { serverId: string; publicKey: string }) =>
      setManagedSshKey(serverId, publicKey),
  });

export const useSetManagedSshPassword = () =>
  useMutation({
    mutationFn: ({ serverId, password }: { serverId: string; password: string }) =>
      setManagedSshPassword(serverId, password),
  });

export const useManagedRuntimeStatus = (serverId: string) =>
  useQuery({
    queryKey: ['servers', serverId, 'managed-runtime'],
    queryFn: () => getManagedRuntimeStatus(serverId),
    enabled: Boolean(serverId),
  });

export const useEnableManagedRuntime = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (serverId: string) => enableManagedRuntime(serverId),
    onSuccess: (_, serverId) => {
      qc.invalidateQueries({ queryKey: ['servers', serverId, 'managed-runtime'] });
    },
  });
};

export const useManagedRuntimeCredential = () =>
  useMutation({
    mutationFn: (serverId: string) => getManagedRuntimeCredential(serverId),
  });

export const useRotateManagedRuntimeCredential = () =>
  useMutation({
    mutationFn: (serverId: string) => rotateManagedRuntimeCredential(serverId),
  });

export const useManagedWorkloads = (serverId: string) =>
  useQuery({
    queryKey: ['servers', serverId, 'managed-workloads'],
    queryFn: () => getManagedWorkloads(serverId),
    enabled: Boolean(serverId),
  });

export const useCreateManagedWorkload = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      serverId,
      payload,
    }: {
      serverId: string;
      payload: Parameters<typeof createManagedWorkload>[1];
    }) => createManagedWorkload(serverId, payload),
    onSuccess: (_, vars) => {
      qc.invalidateQueries({ queryKey: ['servers', vars.serverId, 'managed-workloads'] });
    },
  });
};

export const useControlManagedWorkload = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      serverId,
      payload,
    }: {
      serverId: string;
      payload: Parameters<typeof controlManagedWorkload>[1];
    }) => controlManagedWorkload(serverId, payload),
    onSuccess: (_, vars) => {
      qc.invalidateQueries({ queryKey: ['servers', vars.serverId, 'managed-workloads'] });
    },
  });
};

export const useManagedWorkloadLogs = (serverId: string, workloadId?: string, tail = 100) =>
  useQuery({
    queryKey: ['servers', serverId, 'managed-workloads', workloadId, 'logs', tail],
    queryFn: () => (workloadId ? getManagedWorkloadLogs(serverId, workloadId, tail) : null),
    enabled: Boolean(serverId && workloadId),
  });

export const useFetchManagedWorkloadLogs = () =>
  useMutation({
    mutationFn: ({
      serverId,
      workloadId,
      tail = 100,
    }: {
      serverId: string;
      workloadId: string;
      tail?: number;
    }) => getManagedWorkloadLogs(serverId, workloadId, tail),
  });
