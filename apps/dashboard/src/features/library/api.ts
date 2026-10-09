import type { BaseResponse } from '#/interfaces/base';
import type { ArchiveDeployResponse } from '#/interfaces/templates';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';
import type { LibraryRepo, ProviderConnection } from './types';

export const libraryApi = {
  async gitStatus(): Promise<ProviderConnection[]> {
    try {
      const response = await apiClient.get<BaseResponse<ProviderConnection[]>>('/git/status');
      return response.data ?? [];
    } catch (error) {
      throw handleApiError(error);
    }
  },

  async connectProvider(payload: {
    provider: string;
    accessToken: string;
    accountName: string;
  }): Promise<void> {
    try {
      await apiClient.post('/git/connect', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  async listRepos(provider: string): Promise<LibraryRepo[]> {
    try {
      const response = await apiClient.get<BaseResponse<LibraryRepo[]>>(
        `/git/repos?provider=${encodeURIComponent(provider)}`
      );
      return (response.data ?? []).map((repo) => ({ ...repo, id: String(repo.id) }));
    } catch (error) {
      throw handleApiError(error);
    }
  },

  async deployArchive(projectId: string, name: string, file: File): Promise<ArchiveDeployResponse> {
    try {
      const formData = new FormData();
      formData.append('projectId', projectId);
      formData.append('name', name);
      formData.append('file', file);
      return await apiClient.post<ArchiveDeployResponse>('/deploy/archive', formData);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
