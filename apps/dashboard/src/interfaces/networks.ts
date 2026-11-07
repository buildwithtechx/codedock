export type NetworkMode = 'wireguard' | 'native';

export type NetworkStatus = 'applying' | 'verified' | 'attention' | 'rolled_back' | 'ready';

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
  mode: NetworkMode;
  subnetCidr: string;
  mtu: number;
  verificationPort: number;
  status: NetworkStatus;
  revision: number;
  members: NetworkMember[];
  accessPolicy: NetworkAccessPolicy;
  watchdog?: NetworkRollbackWatchdog;
  createdAt: string;
  updatedAt: string;
}

export interface CreateNetworkInput {
  name: string;
  mode: NetworkMode;
  subnetCidr: string;
  mtu?: number;
  verificationPort?: number;
  members: NetworkMember[];
  accessPolicy: NetworkAccessPolicy;
}
