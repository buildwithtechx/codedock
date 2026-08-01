import type { CanvasSummary, EnvironmentCanvas } from '#/features/projects';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const canvasService = {
  listCanvasSummaries: async (
    organizationId?: string | null
  ): Promise<BaseResponse<CanvasSummary[]>> => {
    try {
      const query = organizationId ? `?organizationId=${encodeURIComponent(organizationId)}` : '';
      return await apiClient.get<BaseResponse<CanvasSummary[]>>(`/canvas/projects${query}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getCanvasSummary: async (projectId: string): Promise<BaseResponse<CanvasSummary>> => {
    try {
      return await apiClient.get<BaseResponse<CanvasSummary>>(`/projects/${projectId}/summary`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getEnvironmentCanvas: async (envId: string): Promise<BaseResponse<EnvironmentCanvas>> => {
    try {
      return await apiClient.get<BaseResponse<EnvironmentCanvas>>(`/environments/${envId}/canvas`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
