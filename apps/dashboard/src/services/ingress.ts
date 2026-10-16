import { apiClient } from '#/lib/api-client';

export interface IngressConnection {
  id: string;
  serverId: string;
  credentialId: string;
  accountId: string;
  tunnelId?: string | null;
  name: string;
  connector: 'existing' | 'managed';
  edgePort: number;
  managedTunnel: boolean;
  defaultForDomains: boolean;
  defaultAccessPolicyId?: string | null;
  status: 'ready' | 'pending' | 'removing' | 'error';
  error?: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface IngressHostname {
  id: string;
  domainId?: string | null;
  connectionId: string;
  hostname: string;
  accessPolicyId?: string | null;
  status: 'ready' | 'pending' | 'removing' | 'error';
  phase: string;
  error?: string | null;
  updatedAt: string;
}

export interface IngressDiscovery {
  accounts: Array<{ id: string; name: string }>;
  tunnels: Array<{ id: string; name: string; remote: boolean; status: string }>;
  policies: Array<{ id: string; name: string; decision: string }>;
  policiesError?: string | null;
}

export interface IngressConnectInput {
  serverId: string;
  credentialId: string;
  accountId: string;
  name: string;
  tunnelId?: string;
  connector: 'existing' | 'managed';
  edgePort?: number;
  defaultForDomains: boolean;
  defaultAccessPolicyId?: string;
}

export interface IngressUpdateInput {
  id: string;
  defaultForDomains: boolean;
  defaultAccessPolicyId?: string | null;
}

export interface IngressDomainInput {
  domainId: string;
  connectionId: string | null;
  accessPolicyId?: string | null;
}

export const ingressService = {
  list: async (serverId?: string): Promise<IngressConnection[]> => {
    const query = serverId ? `?serverId=${encodeURIComponent(serverId)}` : '';
    const res = await apiClient.get<{ data: IngressConnection[] }>(`/dns/ingress${query}`);
    return res.data ?? [];
  },

  discover: async (input: { credentialId: string; accountId?: string; includePolicies?: boolean }): Promise<IngressDiscovery> => {
    const res = await apiClient.post<{ data: IngressDiscovery }>('/dns/ingress/discover', input);
    return res.data ?? { accounts: [], tunnels: [], policies: [] };
  },

  connect: async (input: IngressConnectInput): Promise<IngressConnection> => {
    const res = await apiClient.post<{ data: IngressConnection }>('/dns/ingress', input);
    return res.data;
  },

  update: async ({ id, ...body }: IngressUpdateInput): Promise<IngressConnection> => {
    const res = await apiClient.put<{ data: IngressConnection }>(`/dns/ingress/${encodeURIComponent(id)}`, body);
    return res.data;
  },

  retry: async (id: string): Promise<IngressConnection> => {
    const res = await apiClient.post<{ data: IngressConnection }>(`/dns/ingress/${encodeURIComponent(id)}/retry`, {});
    return res.data;
  },

  disconnect: async (id: string): Promise<{ success: boolean }> => {
    const res = await apiClient.delete<{ data: { success: boolean } }>(`/dns/ingress/${encodeURIComponent(id)}`);
    return res.data ?? { success: true };
  },

  getDomain: async (domainId: string): Promise<IngressHostname | null> => {
    const res = await apiClient.get<{ data: IngressHostname | null }>(`/dns/ingress/domains/${encodeURIComponent(domainId)}`);
    return res.data ?? null;
  },

  setDomain: async ({ domainId, ...body }: IngressDomainInput): Promise<IngressHostname | null> => {
    const res = await apiClient.put<{ data: IngressHostname | null }>(`/dns/ingress/domains/${encodeURIComponent(domainId)}`, body);
    return res.data ?? null;
  },
};
