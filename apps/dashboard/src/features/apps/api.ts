import type { AppService } from '#/features/services';
import type { BaseResponse } from '#/interfaces/base';
import type {
  InstallAppInput,
  InstallPreview,
  OneClickAppDetails,
  OneClickDeployRequest,
  OneClickDeployResponse,
} from '#/interfaces/templates';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const appsApi = {
  async listInstalled(organizationId: string): Promise<AppService[]> {
    try {
      const response = await apiClient.get<BaseResponse<AppService[]>>(
        `/apps?organizationId=${encodeURIComponent(organizationId)}`
      );
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  async listCatalog(): Promise<OneClickAppDetails[]> {
    try {
      const response = await apiClient.get<{ data: OneClickAppDetails[] }>('/one-click');
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  async getCatalogApp(appId: string): Promise<OneClickAppDetails> {
    try {
      const response = await apiClient.get<{ data: OneClickAppDetails }>(
        `/one-click/${encodeURIComponent(appId)}`
      );
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  async reviewInstall(payload: InstallAppInput): Promise<InstallPreview> {
    try {
      const response = await apiClient.post<{ data: InstallPreview }>('/one-click/review', payload);
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  async deployCatalogApp(payload: OneClickDeployRequest): Promise<OneClickDeployResponse> {
    try {
      return await apiClient.post<OneClickDeployResponse>('/one-click/deploy', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
