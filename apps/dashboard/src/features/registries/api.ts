import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';
import type { CreateRegistryRequest, Registry } from './interfaces';

export const registryService = {
  listByProject: async (projectId: string): Promise<Registry[]> => {
    try {
      const res = await apiClient.get<Registry[] | BaseResponse<Registry[]>>(
        `/projects/${projectId}/registries`
      );
      if (Array.isArray(res)) return res;
      return (res as BaseResponse<Registry[]>)?.data || [];
    } catch (error) {
      throw handleApiError(error);
    }
  },

  create: async (projectId: string, payload: CreateRegistryRequest): Promise<Registry> => {
    try {
      const res = await apiClient.post<Registry | BaseResponse<Registry>>(
        `/projects/${projectId}/registries`,
        payload
      );
      if (res && 'data' in res && res.data) return res.data as Registry;
      return res as Registry;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  delete: async (projectId: string, id: string): Promise<void> => {
    try {
      await apiClient.delete(`/projects/${projectId}/registries/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
