import type { BaseResponse } from '#/interfaces/base';
import type {
  CreateServerRequest,
  Server,
  ServerComponentStatus,
  ServerPortScanResult,
  ServerRateLimitConfig,
  TestSSHRequest,
  UpdateServerRequest,
} from '#/interfaces/server';
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

  get: async (id: string): Promise<Server> => {
    try {
      const res = await apiClient.get<BaseResponse<Server>>(`/servers/${id}`);
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

  update: async (id: string, payload: UpdateServerRequest): Promise<Server> => {
    try {
      const res = await apiClient.patch<BaseResponse<Server>>(`/servers/${id}`, payload);
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  delete: async (id: string): Promise<void> => {
    try {
      await apiClient.delete<BaseResponse<null>>(`/servers/${id}`);
    } catch (err) {
      throw handleApiError(err);
    }
  },

  testSSH: async (payload: TestSSHRequest): Promise<{ success: boolean; message: string }> => {
    try {
      const res = await apiClient.post<BaseResponse<null>>('/servers/test-ssh', payload);
      return { success: true, message: res.message || 'SSH connection successful' };
    } catch (err) {
      const error = handleApiError(err);
      return { success: false, message: error.message };
    }
  },

  getComponents: async (id: string): Promise<ServerComponentStatus[]> => {
    try {
      const res = await apiClient.get<BaseResponse<ServerComponentStatus[]>>(
        `/servers/${id}/components`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  scanPorts: async (id: string): Promise<ServerPortScanResult> => {
    try {
      const res = await apiClient.post<BaseResponse<ServerPortScanResult>>(
        `/servers/${id}/scan-ports`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  getRateLimit: async (id: string): Promise<ServerRateLimitConfig> => {
    try {
      const res = await apiClient.get<BaseResponse<ServerRateLimitConfig>>(
        `/servers/${id}/rate-limit`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  updateRateLimit: async (
    id: string,
    payload: ServerRateLimitConfig
  ): Promise<ServerRateLimitConfig> => {
    try {
      const res = await apiClient.put<BaseResponse<ServerRateLimitConfig>>(
        `/servers/${id}/rate-limit`,
        payload
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
};
