import type { BaseResponse } from '#/interfaces/base';
import type { CreateRegistryRequest, Registry } from '#/interfaces/registry';
import { apiClient } from '#/lib/apiClient';

export const registryService = {
  listByProject: async (projectId: string): Promise<Registry[]> => {
    const res = await apiClient.get<BaseResponse<Registry[]>>(`/projects/${projectId}/registries`);
    return res.data || [];
  },

  create: async (projectId: string, payload: CreateRegistryRequest): Promise<Registry> => {
    const res = await apiClient.post<BaseResponse<Registry>>(
      `/projects/${projectId}/registries`,
      payload
    );
    return res.data as Registry;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/registries/${id}`);
  },
};
