import type {
  AdvisoriesResponse,
  ApplySecurityInput,
  ProjectSecurityResult,
  RefreshAttempt,
} from '@codedock/contracts';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const securityService = {
  getAdvisories: async (): Promise<BaseResponse<AdvisoriesResponse>> => {
    try {
      return await apiClient.get<BaseResponse<AdvisoriesResponse>>('/updates/advisories');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getProjectSecurity: async (projectId: string): Promise<BaseResponse<ProjectSecurityResult>> => {
    try {
      return await apiClient.get<BaseResponse<ProjectSecurityResult>>(
        `/updates/${projectId}/security`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  checkProject: async (projectId: string): Promise<BaseResponse<ProjectSecurityResult>> => {
    try {
      return await apiClient.post<BaseResponse<ProjectSecurityResult>>(
        `/updates/${projectId}/security/check`,
        {}
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  setProjectPolicy: async (
    projectId: string,
    autoRefresh: boolean
  ): Promise<BaseResponse<null>> => {
    try {
      return await apiClient.put<BaseResponse<null>>(`/updates/${projectId}/security/policy`, {
        autoRefresh,
      });
    } catch (error) {
      throw handleApiError(error);
    }
  },

  applyRemediation: async (
    projectId: string,
    payload: ApplySecurityInput
  ): Promise<BaseResponse<RefreshAttempt>> => {
    try {
      return await apiClient.post<BaseResponse<RefreshAttempt>>(
        `/updates/${projectId}/security/apply`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
