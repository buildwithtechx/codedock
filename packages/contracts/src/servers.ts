export interface ServerMetric {
  timestamp: string;
  cpuPercent: number;
  memoryMB: number;
  networkRxKB: number;
  networkTxKB: number;
  status?: string;
  cpuUsagePercentage?: number;
  memoryUsageBytes?: number;
  memoryLimitBytes?: number;
}

export interface CodedockServer {
  id: string;
  name: string;
  ipAddress: string;
  status: 'online' | 'offline' | 'unreachable' | string;
  sshPort?: number;
  dockerVersion?: string;
  osInfo?: string;
  isDefault?: boolean;
}

export interface ClusterNode {
  serverId: string;
  serverName?: string;
  role: 'server' | 'agent';
  privateIp: string;
  interface?: string;
  fingerprint?: string;
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
