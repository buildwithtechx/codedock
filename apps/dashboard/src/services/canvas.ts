import type { CanvasSummary, EnvironmentCanvas } from '#/interfaces/projects';
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

  observeResource: async (
    envId: string,
    nodeId: string
  ): Promise<BaseResponse<{ logs: string; metrics: Record<string, unknown> }>> => {
    try {
      return await apiClient.get<BaseResponse<{ logs: string; metrics: Record<string, unknown> }>>(
        `/environments/${envId}/canvas/observe?node=${encodeURIComponent(nodeId)}`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  executeAction: async (path: string, operation: string): Promise<void> => {
    try {
      await apiClient.post(`${path}/${operation}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
