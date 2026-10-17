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
