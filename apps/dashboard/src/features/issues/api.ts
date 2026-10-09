import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';
import type { ActIssueOutcome, ActIssuePayload, AttentionIssue } from './interfaces';

export const issuesService = {
  list: async (organizationId: string): Promise<BaseResponse<AttentionIssue[]>> => {
    try {
      return await apiClient.get<BaseResponse<AttentionIssue[]>>(
        `/organizations/${organizationId}/attention`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  evaluate: async (organizationId: string): Promise<BaseResponse<null>> => {
    try {
      return await apiClient.post<BaseResponse<null>>(
        `/organizations/${organizationId}/attention/evaluate`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  acknowledge: async (organizationId: string, issueId: string): Promise<BaseResponse<null>> => {
    try {
      return await apiClient.post<BaseResponse<null>>(
        `/organizations/${organizationId}/attention/${issueId}/ack`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  resolve: async (organizationId: string, issueId: string): Promise<BaseResponse<null>> => {
    try {
      return await apiClient.post<BaseResponse<null>>(
        `/organizations/${organizationId}/attention/${issueId}/resolve`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  act: async (
    organizationId: string,
    issueId: string,
    payload: ActIssuePayload
  ): Promise<BaseResponse<ActIssueOutcome>> => {
    try {
      return await apiClient.post<BaseResponse<ActIssueOutcome>>(
        `/organizations/${organizationId}/attention/${issueId}/act`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
