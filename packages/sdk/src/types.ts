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

export interface NetworkMember {
  serverId: string;
  serverName?: string;
  privateIp: string;
  interfaceName: string;
  status?: string;
  latencyMs?: number;
}

export interface NetworkWatchdog {
  operationId: string;
  armed: boolean;
  timeoutSeconds: number;
  remainingSeconds: number;
  expiresAt: string;
  canCommit: boolean;
  canRollback: boolean;
  status: string;
}

export interface ServerNetwork {
  id: string;
  organizationId?: string;
  name: string;
  mode: 'wireguard' | 'native';
  subnetCidr: string;
  mtu: number;
  verificationPort: number;
  status: string;
  revision: number;
  members: NetworkMember[];
  accessPolicy: {
    mode: 'mesh' | 'isolated' | 'custom';
    allowedPairs?: [string, string][];
  };
  watchdog?: NetworkWatchdog;
  createdAt: string;
  updatedAt: string;
}

export interface ClusterNode {
  serverId: string;
  serverName?: string;
  role: 'server' | 'agent';
  privateIp: string;
  interface?: string;
  fingerprint?: string;
}

export interface ClusterScaling {
  status: 'not_enabled' | 'setting_up' | 'ready' | 'failed' | 'interrupted' | 'removing';
  k3sVersion: string;
  controlNodes: number;
  workerNodes: number;
  podCidr?: string;
  serviceCidr?: string;
  clusterDomain?: string;
  error?: string;
  updatedAt: string;
}

export interface ComputeCluster {
  id: string;
  projectId?: string;
  name: string;
  version: string;
  networkId?: string;
  controls: number;
  nodes: ClusterNode[];
  revision: number;
  status: string;
  scaling?: ClusterScaling;
  error?: string;
  updatedAt?: string;
}

export interface ActionWorkflow {
  id: string;
  projectId?: string;
  name: string;
  description?: string;
  trigger: string;
  yaml?: string;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface ActionRun {
  id: string;
  workflowId: string;
  name: string;
  number: number;
  status: 'queued' | 'running' | 'ready' | 'failed' | 'cancelled';
  triggeredBy: string;
  runnerId?: string;
  logs?: string;
  durationSeconds: number;
  createdAt: string;
  updatedAt: string;
}

export interface ActionRunner {
  id: string;
  name: string;
  mode: 'container' | 'native';
  image?: string;
  status: 'online' | 'offline' | 'busy';
  labels: string[];
  enabled: boolean;
  lastHeartbeatAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ActiveSubscriptionOverview {
  status: string;
  planId: string;
  planName: string;
  prepaidBalanceUsd: number;
  autoRechargeEnabled: boolean;
  currentPeriodEnd?: string;
}

export interface PrepaidTopupPackage {
  id: string;
  amountUsd: number;
  bonusUsd: number;
  totalCreditsUsd: number;
  badge?: string;
}

export interface AccountInfo {
  id: string;
  email: string;
  name?: string;
  organizationId?: string;
  role?: string;
  tier?: string;
  createdAt?: string;
}

export interface PairingHandshake {
  code: string;
  instanceId?: string;
  expiresAt: string;
}

export interface PairingClaimResult {
  token: string;
  account: AccountInfo;
  organizationId?: string;
}

export interface TenantCorsConfig {
  tenantId: string;
  allowedOrigins: string[];
  allowedMethods?: string[];
  maxAgeSeconds?: number;
}

export interface RateLimitDetails {
  retryAfterSeconds: number;
  limit?: number;
  remaining?: number;
  resetAt?: string;
  planTier?: string;
}

export interface BillingGateDetails {
  code: 'BILLING_GATE_REQUIRED' | 'PAYMENT_REQUIRED';
  message: string;
  upgradeUrl?: string;
  requiredPlan?: string;
}
