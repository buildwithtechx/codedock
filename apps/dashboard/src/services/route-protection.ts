import type { RouteProtectionInput, RouteProtectionStatus } from '@codedock/contracts';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const routeProtectionService = {
  listRouteProtection: async (
    projectId: string
  ): Promise<BaseResponse<RouteProtectionStatus[]>> => {
    try {
      return await apiClient.get<BaseResponse<RouteProtectionStatus[]>>(
        `/projects/${projectId}/route-protection`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  changeRouteProtection: async (
    projectId: string,
    payload: RouteProtectionInput
  ): Promise<BaseResponse<RouteProtectionStatus>> => {
    try {
      return await apiClient.post<BaseResponse<RouteProtectionStatus>>(
        `/projects/${projectId}/route-protection`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
