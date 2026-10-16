import { apiClient } from '#/lib/api-client';

export interface NetworkMember {
  serverId: string;
  serverName?: string;
  privateIp: string;
  interfaceName: string;
  endpoint?: string;
  publicKey?: string;
  status: string;
  latencyMs?: number;
  packetLoss?: number;
}

export interface NetworkAccessPolicy {
  mode: 'full_mesh' | 'custom_allowlist' | 'isolated';
  allowedPairs?: [string, string][];
}

export interface NetworkRollbackWatchdog {
  operationId: string;
  armed: boolean;
  timeoutSeconds: number;
  remainingSeconds: number;
  expiresAt: string;
  canCommit: boolean;
  canRollback: boolean;
  status: 'armed' | 'committed' | 'rolled_back' | 'expired';
}

export interface ServerNetwork {
  id: string;
  organizationId?: string;
  name: string;
  mode: 'wireguard' | 'native';
  subnetCidr: string;
  mtu: number;
  verificationPort: number;
  status: 'applying' | 'verified' | 'attention' | 'rolled_back';
  revision: number;
  members: NetworkMember[];
  accessPolicy: NetworkAccessPolicy;
  watchdog?: NetworkRollbackWatchdog;
  createdAt: string;
  updatedAt: string;
}

export interface CreateNetworkInput {
  name: string;
  mode: 'wireguard' | 'native';
  subnetCidr: string;
  mtu?: number;
  verificationPort?: number;
  members: NetworkMember[];
  accessPolicy: NetworkAccessPolicy;
}

export const networksService = {
  list: async (): Promise<ServerNetwork[]> => {
    const res = await apiClient.get<{ data: ServerNetwork[] }>('/system/networks');
    return res.data ?? [];
  },

  get: async (id: string): Promise<ServerNetwork> => {
    const res = await apiClient.get<{ data: ServerNetwork }>(`/system/networks/${encodeURIComponent(id)}`);
    return res.data;
  },

  create: async (input: CreateNetworkInput): Promise<ServerNetwork> => {
    const res = await apiClient.post<{ data: ServerNetwork }>('/system/networks', input);
    return res.data;
  },

  reviseAccess: async (id: string, accessPolicy: NetworkAccessPolicy): Promise<ServerNetwork> => {
    const res = await apiClient.put<{ data: ServerNetwork }>(`/system/networks/${encodeURIComponent(id)}/access`, {
      networkId: id,
      accessPolicy,
    });
    return res.data;
  },

  getWatchdog: async (operationId: string): Promise<NetworkRollbackWatchdog> => {
    const res = await apiClient.get<{ data: NetworkRollbackWatchdog }>(`/system/networks/operations/${encodeURIComponent(operationId)}/watchdog`);
    return res.data;
  },

  commitWatchdog: async (operationId: string): Promise<NetworkRollbackWatchdog> => {
    const res = await apiClient.post<{ data: NetworkRollbackWatchdog }>(`/system/networks/operations/${encodeURIComponent(operationId)}/watchdog/commit`, {});
    return res.data;
  },

  rollbackWatchdog: async (operationId: string): Promise<NetworkRollbackWatchdog> => {
    const res = await apiClient.post<{ data: NetworkRollbackWatchdog }>(`/system/networks/operations/${encodeURIComponent(operationId)}/watchdog/rollback`, {});
    return res.data;
  },

  delete: async (id: string): Promise<{ success: boolean }> => {
    const res = await apiClient.delete<{ data: { success: boolean } }>(`/system/networks/${encodeURIComponent(id)}`);
    return res.data ?? { success: true };
  },
};
