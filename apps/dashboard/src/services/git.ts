import type {
  ConnectGitRequest,
  ConnectGitResponse,
  GetGitStatusResponse,
  ListGitBranchesResponse,
  ListGitReposResponse,
  RepoTreeEntry,
  RepositoryInspection,
} from '#/interfaces/git';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
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

  getRepoTree: async (payload: {
    repositoryUrl: string;
    branch?: string;
  }): Promise<BaseResponse<RepoTreeEntry[]>> => {
    try {
      return await apiClient.post<BaseResponse<RepoTreeEntry[]>>('/git/tree', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  inspectRepository: async (payload: {
    repositoryUrl: string;
    branch?: string;
    rootDirectory?: string;
  }): Promise<BaseResponse<RepositoryInspection>> => {
    try {
      return await apiClient.post<BaseResponse<RepositoryInspection>>('/git/inspect', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
