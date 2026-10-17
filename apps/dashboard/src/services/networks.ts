import type {
  CreateNetworkInput,
  NetworkAccessPolicy,
  NetworkMember,
  NetworkRollbackWatchdog,
  ServerNetwork,
} from '#/interfaces/networks';
import { apiClient } from '#/lib/api-client';

export const networksService = {
  list: async (): Promise<ServerNetwork[]> => {
    const res = await apiClient.get<{ data: ServerNetwork[] }>('/system/networks');
    return res.data ?? [];
  },

  get: async (id: string): Promise<ServerNetwork> => {
    const res = await apiClient.get<{ data: ServerNetwork }>(
      `/system/networks/${encodeURIComponent(id)}`
    );
    return res.data;
  },

  create: async (input: CreateNetworkInput): Promise<ServerNetwork> => {
    const res = await apiClient.post<{ data: ServerNetwork }>('/system/networks', input);
    return res.data;
  },

  reviseAccess: async (id: string, accessPolicy: NetworkAccessPolicy): Promise<ServerNetwork> => {
    const res = await apiClient.put<{ data: ServerNetwork }>(
      `/system/networks/${encodeURIComponent(id)}/access`,
      {
        networkId: id,
        accessPolicy,
      }
    );
    return res.data;
  },

  getWatchdog: async (operationId: string): Promise<NetworkRollbackWatchdog> => {
    const res = await apiClient.get<{ data: NetworkRollbackWatchdog }>(
      `/system/networks/operations/${encodeURIComponent(operationId)}/watchdog`
    );
    return res.data;
  },

  commitWatchdog: async (operationId: string): Promise<NetworkRollbackWatchdog> => {
    const res = await apiClient.post<{ data: NetworkRollbackWatchdog }>(
      `/system/networks/operations/${encodeURIComponent(operationId)}/watchdog/commit`,
      {}
    );
    return res.data;
  },

  rollbackWatchdog: async (operationId: string): Promise<NetworkRollbackWatchdog> => {
    const res = await apiClient.post<{ data: NetworkRollbackWatchdog }>(
      `/system/networks/operations/${encodeURIComponent(operationId)}/watchdog/rollback`,
      {}
    );
    return res.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/system/networks/${encodeURIComponent(id)}`);
  },
};

export type {
  ServerNetwork,
  NetworkMember,
  CreateNetworkInput,
  NetworkAccessPolicy,
  NetworkRollbackWatchdog,
};
