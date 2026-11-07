import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { UpdateSettingsRequest } from '#/interfaces/settings';
import { settingsService, systemService } from '#/services/settings';

export const useGetSettings = () => {
  return useQuery({
    queryKey: ['settings', 'getSettings'],
    queryFn: () => settingsService.getSettings(),
  });
};

export const useGetPublicSettings = () => {
  return useQuery({
    queryKey: ['settings', 'getPublicSettings'],
    queryFn: () => settingsService.getPublicSettings(),
  });
};

export const useGetSetupStatus = () => {
  return useQuery({
    queryKey: ['settings', 'getSetupStatus'],
    queryFn: () => settingsService.getSetupStatus(),
  });
};

export const useUpdateSettings = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: UpdateSettingsRequest }) =>
      settingsService.updateSettings(payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings'] });
    },
  });
};

export const useGetAISettings = () => {
  return useQuery({
    queryKey: ['settings', 'getAISettings'],
    queryFn: () => settingsService.getAISettings(),
  });
};

export const useUpdateAISettings = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Record<string, unknown>) => settingsService.updateAISettings(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings', 'getAISettings'] });
    },
  });
};

export const useGetGitApps = () => {
  return useQuery({
    queryKey: ['settings', 'getGitApps', 'github'],
    queryFn: () => settingsService.getGitApps(),
  });
};

export const useSaveGitApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Parameters<typeof settingsService.saveGitApp>[0]) =>
      settingsService.saveGitApp(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings', 'getGitApps', 'github'] });
    },
  });
};

export const useDeleteGitApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => settingsService.deleteGitApp(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings', 'getGitApps', 'github'] });
    },
  });
};

export const useExchangeGithubManifest = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Record<string, unknown>) =>
      settingsService.exchangeGithubManifest(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['settings', 'getGitApps', 'github'] });
    },
  });
};

export const useGetUpdateStatus = () => {
  return useQuery({
    queryKey: ['settings', 'getUpdateStatus'],
    queryFn: () => settingsService.getUpdateStatus(),
  });
};

export const useCheckUpdate = () => {
  return useMutation({
    mutationFn: () => settingsService.checkUpdate(),
  });
};

export const useDeployUpdate = () => {
  return useMutation({
    mutationFn: () => settingsService.deployUpdate(),
  });
};

export const useGetSystemStats = () => {
  return useQuery({
    queryKey: ['system', 'stats'],
    queryFn: () => systemService.getSystemStats(),
    refetchInterval: 30_000,
  });
};

export const useExportSystem = () => {
  return useMutation({
    mutationFn: (passphrase: string) => systemService.exportSystem(passphrase),
  });
};

export const useImportSystem = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: FormData) => systemService.importSystem(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['system'] });
    },
  });
};

export const useRestartSystem = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => systemService.restartSystem(),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['system'] });
    },
  });
};

export const useCleanupSystem = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => systemService.cleanupSystem(),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['system'] });
    },
  });
};
