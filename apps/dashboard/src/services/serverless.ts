import type { ServerlessFunctionCode } from '#/features/projects/interfaces';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const serverlessService = {
  getCode: async (
    projectId: string,
    serviceId: string
  ): Promise<BaseResponse<ServerlessFunctionCode>> => {
    try {
      return await apiClient.get<BaseResponse<ServerlessFunctionCode>>(
        `/projects/${projectId}/services/${serviceId}/serverless/code`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  saveCode: async (
    projectId: string,
    serviceId: string,
    payload: { codeContent: string; runtime?: string }
  ): Promise<BaseResponse<ServerlessFunctionCode>> => {
    try {
      return await apiClient.post<BaseResponse<ServerlessFunctionCode>>(
        `/projects/${projectId}/services/${serviceId}/serverless/code`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
