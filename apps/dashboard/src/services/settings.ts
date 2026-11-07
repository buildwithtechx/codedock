import type { BaseResponse } from '#/interfaces/base';
import type {
  CheckUpdateResponse,
  DeployUpdateResponse,
  ExchangeGithubManifestResponse,
  GetGithubAppsResponse,
  GetServerSettingsResponse,
  GetUpdateStatusResponse,
  GitAppsManifestRequest,
  SaveGithubAppPayload,
  SaveGithubAppResponse,
  SystemStats,
  UpdateServerSettingsResponse,
  UpdateSettingsRequest,
} from '#/interfaces/settings';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const settingsService = {
  getSettings: async (): Promise<GetServerSettingsResponse> => {
    try {
      return await apiClient.get<GetServerSettingsResponse>('/settings');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getSetupStatus: async (): Promise<BaseResponse<{ setupRequired: boolean; cwd?: string }>> => {
    try {
      return await apiClient.get<BaseResponse<{ setupRequired: boolean; cwd?: string }>>(
        '/system/setup-status'
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getPublicSettings: async (): Promise<
    BaseResponse<{
      registrationEnabled: boolean;
      siteName?: string;
      emailEnabled: boolean;
      cloudMode: boolean;
      selfHosted?: boolean;
      version?: string;
      hostDomain?: string;
      defaultWildcardDomain?: string;
      serverIp?: string;
    }>
  > => {
    try {
      return await apiClient.get<
        BaseResponse<{
          registrationEnabled: boolean;
          siteName?: string;
          emailEnabled: boolean;
          cloudMode: boolean;
          selfHosted?: boolean;
          version?: string;
          hostDomain?: string;
          defaultWildcardDomain?: string;
          serverIp?: string;
        }>
      >('/system/public');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateSettings: async (payload: UpdateSettingsRequest): Promise<UpdateServerSettingsResponse> => {
    try {
      return await apiClient.put<UpdateServerSettingsResponse>('/settings', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getAISettings: async (): Promise<BaseResponse<Record<string, unknown>>> => {
    try {
      return await apiClient.get<BaseResponse<Record<string, unknown>>>('/ai');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateAISettings: async (
    payload: Record<string, unknown>
  ): Promise<BaseResponse<Record<string, unknown>>> => {
    try {
      return await apiClient.put<BaseResponse<Record<string, unknown>>>('/ai', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getGitApps: async (): Promise<GetGithubAppsResponse> => {
    return await settingsService.getGithubApps();
  },

  saveGitApp: async (payload: SaveGithubAppPayload): Promise<SaveGithubAppResponse> => {
    return await settingsService.saveGithubApp(payload);
  },

  deleteGitApp: async (id: string): Promise<void> => {
    return await settingsService.deleteGithubApp(id);
  },

  getGithubApps: async (): Promise<GetGithubAppsResponse> => {
    try {
      return await apiClient.get<GetGithubAppsResponse>('/settings/git_apps/github');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  saveGithubApp: async (payload: SaveGithubAppPayload): Promise<SaveGithubAppResponse> => {
    try {
      return await apiClient.put<SaveGithubAppResponse>('/settings/git_apps/github', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteGithubApp: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/settings/git_apps/github/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  exchangeGithubManifest: async (
    payload: GitAppsManifestRequest | Record<string, unknown>
  ): Promise<ExchangeGithubManifestResponse> => {
    try {
      return await apiClient.post<ExchangeGithubManifestResponse>(
        '/settings/git_apps/github/manifest-callback',
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getUpdateStatus: async (): Promise<GetUpdateStatusResponse> => {
    try {
      return await apiClient.get<GetUpdateStatusResponse>('/settings/updates/status');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  checkUpdate: async (): Promise<CheckUpdateResponse> => {
    try {
      return await apiClient.post<CheckUpdateResponse>('/settings/updates/check');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deployUpdate: async (): Promise<DeployUpdateResponse> => {
    try {
      return await apiClient.post<DeployUpdateResponse>('/settings/updates/deploy');
    } catch (error) {
      throw handleApiError(error);
    }
  },
};

export const systemService = {
  getSystemStats: async (): Promise<BaseResponse<SystemStats>> => {
    try {
      return await apiClient.get<BaseResponse<SystemStats>>('/system/stats');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  exportSystem: async (passphrase: string): Promise<Blob> => {
    try {
      return await apiClient.postBlob('/system/export', { passphrase });
    } catch (error) {
      throw handleApiError(error);
    }
  },

  importSystem: async (payload: FormData): Promise<void> => {
    try {
      await apiClient.post('/system/import', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  restartSystem: async (): Promise<void> => {
    try {
      await apiClient.post('/system/restart');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  cleanupSystem: async (): Promise<void> => {
    try {
      await apiClient.post('/system/maintenance/cleanup');
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
