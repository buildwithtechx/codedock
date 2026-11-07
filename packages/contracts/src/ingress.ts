export interface IngressConnection {
  id: string;
  serverId: string;
  credentialId: string;
  accountId: string;
  tunnelId: string;
  name: string;
  connector: string;
  edgePort: number;
  managedTunnel: boolean;
  defaultForDomains: boolean;
  defaultAccessPolicyId?: string | null;
  status: string;
  error?: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface IngressDiscovery {
  accounts: Array<{ id: string; name: string }>;
  tunnels: Array<{ id: string; name: string; remote: boolean; status: string }>;
  policies: Array<{ id: string; name: string; decision: string }>;
  policiesError?: string | null;
}

export interface IngressHostname {
  connectionId: string;
  hostname: string;
  accessPolicyId?: string | null;
  status: string;
  phase: string;
  error?: string | null;
  updatedAt: string;
}
