import type {
  ConnectGitRequest,
  ConnectGitResponse,
  GetGitStatusResponse,
  ListGitBranchesResponse,
  ListGitReposResponse,
} from '#/interfaces/git';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const gitService = {
  getStatus: async (): Promise<GetGitStatusResponse> => {
    try {
      return await apiClient.get<GetGitStatusResponse>('/git/status');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  connect: async (payload: ConnectGitRequest): Promise<ConnectGitResponse> => {
    try {
      return await apiClient.post<ConnectGitResponse>('/git/connect', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  disconnect: async (provider: string): Promise<void> => {
    try {
      await apiClient.delete(`/git/connect/${provider}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listRepos: async (provider: string): Promise<ListGitReposResponse> => {
    try {
      return await apiClient.get<ListGitReposResponse>(`/git/repos?provider=${provider}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listBranches: async (provider: string, repo: string): Promise<ListGitBranchesResponse> => {
    try {
      return await apiClient.get<ListGitBranchesResponse>(
        `/git/branches?provider=${provider}&repo=${encodeURIComponent(repo)}`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
