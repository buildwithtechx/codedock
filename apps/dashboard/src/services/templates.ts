import type { OneClickApp } from '#/features/databases';
import type {
  ArchiveDeployResponse,
  ComposeDeployResponse,
  ExampleApp,
  InstallAppInput,
  InstallPreview,
  OneClickAppDetails,
  OneClickDeployRequest,
  OneClickDeployResponse,
} from '#/interfaces/templates';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const templatesService = {
  listOneClickApps: async (): Promise<OneClickApp[]> => {
    try {
      const response = await apiClient.get<{ data: OneClickApp[] }>('/one-click');
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listExampleApps: async (): Promise<ExampleApp[]> => {
    try {
      const response = await apiClient.get<{ data: ExampleApp[] }>('/examples');
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getOneClickApp: async (appId: string): Promise<OneClickAppDetails> => {
    try {
      const response = await apiClient.get<{ data: OneClickAppDetails }>(`/one-click/${appId}`);
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  reviewOneClickInstall: async (payload: InstallAppInput): Promise<InstallPreview> => {
    try {
      const response = await apiClient.post<{ data: InstallPreview }>('/one-click/review', payload);
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deployOneClickApp: async (payload: OneClickDeployRequest): Promise<OneClickDeployResponse> => {
    try {
      return await apiClient.post<OneClickDeployResponse>('/one-click/deploy', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deployCompose: async (projectId: string, file: File): Promise<ComposeDeployResponse> => {
    try {
      const formData = new FormData();
      formData.append('projectId', projectId);
      formData.append('file', file);

      return await apiClient.post<ComposeDeployResponse>('/compose/deploy', formData);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deployArchive: async (projectId: string, file: File): Promise<ArchiveDeployResponse> => {
    try {
      const formData = new FormData();
      if (projectId) {
        formData.append('projectId', projectId);
      }
      formData.append('file', file);

      return await apiClient.post<ArchiveDeployResponse>('/deploy/archive', formData);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
