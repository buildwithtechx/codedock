import type {
  IngressConnectInput,
  IngressConnection,
  IngressDiscovery,
  IngressDomainInput,
  IngressHostname,
  IngressUpdateInput,
} from '#/interfaces/ingress';
import { apiClient } from '#/lib/api-client';

export const ingressService = {
  list: async (serverId?: string): Promise<IngressConnection[]> => {
    const query = serverId ? `?serverId=${encodeURIComponent(serverId)}` : '';
    const res = await apiClient.get<{ data: IngressConnection[] }>(`/dns/ingress${query}`);
    return res.data ?? [];
  },

  discover: async (input: {
    credentialId: string;
    accountId?: string;
    includePolicies?: boolean;
  }): Promise<IngressDiscovery> => {
    const res = await apiClient.post<{ data: IngressDiscovery }>('/dns/ingress/discover', input);
    return res.data ?? { accounts: [], tunnels: [], policies: [] };
  },

  connect: async (input: IngressConnectInput): Promise<IngressConnection> => {
    const res = await apiClient.post<{ data: IngressConnection }>('/dns/ingress', input);
    return res.data;
  },

  update: async ({ id, ...body }: IngressUpdateInput): Promise<IngressConnection> => {
    const res = await apiClient.put<{ data: IngressConnection }>(
      `/dns/ingress/${encodeURIComponent(id)}`,
      body
    );
    return res.data;
  },

  retry: async (id: string): Promise<IngressConnection> => {
    const res = await apiClient.post<{ data: IngressConnection }>(
      `/dns/ingress/${encodeURIComponent(id)}/retry`,
      {}
    );
    return res.data;
  },

  disconnect: async (id: string): Promise<void> => {
    await apiClient.delete(`/dns/ingress/${encodeURIComponent(id)}`);
  },

  listHostnames: async (domainId?: string): Promise<IngressHostname[]> => {
    const query = domainId ? `?domainId=${encodeURIComponent(domainId)}` : '';
    const res = await apiClient.get<{ data: IngressHostname[] }>(`/dns/ingress/hostnames${query}`);
    return res.data ?? [];
  },

  getDomain: async (domainId: string): Promise<IngressHostname | null> => {
    const query = `?domainId=${encodeURIComponent(domainId)}`;
    const res = await apiClient.get<{ data: IngressHostname[] }>(`/dns/ingress/hostnames${query}`);
    return res.data?.[0] ?? null;
  },

  bindDomain: async (input: IngressDomainInput): Promise<IngressHostname> => {
    const res = await apiClient.put<{ data: IngressHostname }>('/dns/ingress/domain', input);
    return res.data;
  },

  setDomain: async (input: IngressDomainInput): Promise<IngressHostname> => {
    const res = await apiClient.put<{ data: IngressHostname }>('/dns/ingress/domain', input);
    return res.data;
  },
};
