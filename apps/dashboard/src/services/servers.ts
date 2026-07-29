import type { BaseResponse } from '#/interfaces/base';
import type { CreateServerRequest, Server } from '#/interfaces/server';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const serverService = {
  list: async (): Promise<Server[]> => {
    try {
      const res = await apiClient.get<BaseResponse<Server[]>>('/servers');
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  create: async (payload: CreateServerRequest): Promise<Server> => {
    try {
      const res = await apiClient.post<BaseResponse<Server>>('/servers', payload);
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
};
