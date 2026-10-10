export interface ServerMetrics {
  cpu_usage_percentage: number;
  memory_usage_bytes: number;
  memory_limit_bytes: number;
  disk_usage_bytes: number;
  disk_total_bytes: number;
}

export interface Server {
  id: string;
  userId: string;
  name: string;
  ipAddress: string;
  isLocal?: boolean;
  sshHost?: string;
  sshPort?: number;
  sshUser?: string;
  sshAuthMethod?: 'key' | 'password' | 'agent';
  sshTransport?: 'direct' | 'cloudflare';
  sshJumpHost?: string;
  isControlPlane?: boolean;
  provider?: string;
  externalId?: string;
  region?: string;
  serverType?: string;
  status: string;
  workerToken?: string;
  lastSeenAt?: string;
  managed?: CloudWorkspaceSummary | null;
  capabilities?: ServerCapabilities | null;
  purpose?: string;
  metrics: ServerMetrics | string | null;
  createdAt: string;
  updatedAt: string;
}

function toFiniteNumber(value: unknown): number {
  return typeof value === 'number' && Number.isFinite(value) ? value : 0;
}

export function parseServerMetrics(
  metrics: ServerMetrics | string | null | undefined
): ServerMetrics | null {
  let parsed: unknown = null;
  if (!metrics) return null;
  if (typeof metrics === 'object') {
    parsed = metrics;
  } else if (typeof metrics === 'string') {
    try {
      const decoded = metrics.startsWith('{') ? metrics : atob(metrics);
      parsed = JSON.parse(decoded);
    } catch {
      return null;
    }
  }
  if (!parsed || typeof parsed !== 'object') return null;
  const fields = parsed as Partial<ServerMetrics>;
  return {
    cpu_usage_percentage: toFiniteNumber(fields.cpu_usage_percentage),
    memory_usage_bytes: toFiniteNumber(fields.memory_usage_bytes),
    memory_limit_bytes: toFiniteNumber(fields.memory_limit_bytes),
    disk_usage_bytes: toFiniteNumber(fields.disk_usage_bytes),
    disk_total_bytes: toFiniteNumber(fields.disk_total_bytes),
  };
}

export interface CreateServerRequest {
  name: string;
  ipAddress?: string;
  isLocal?: boolean;
  sshHost?: string;
  sshPort?: number;
  sshUser?: string;
  sshAuthMethod?: 'key' | 'password' | 'agent';
  sshKey?: string;
  sshPrivateKey?: string;
  sshPassword?: string;
  sshTransport?: 'direct' | 'cloudflare';
  sshJumpHost?: string;
}

export interface UpdateServerRequest {
  name?: string;
  ipAddress?: string;
  isLocal?: boolean;
  sshHost?: string;
  sshPort?: number;
  sshUser?: string;
  sshAuthMethod?: 'key' | 'password' | 'agent';
  sshKey?: string;
  sshPrivateKey?: string;
  sshPassword?: string;
  sshTransport?: 'direct' | 'cloudflare';
  sshJumpHost?: string;
}

export interface TestSSHRequest {
  sshHost: string;
  sshPort?: number;
  sshUser?: string;
  sshKey?: string;
  sshPrivateKey?: string;
  sshPassword?: string;
}

export interface ServerComponentStatus {
  name: string;
  label: string;
  description: string;
  installable: boolean;
  installed: boolean;
  version?: string;
  healthy: boolean;
  message: string;
}

export interface ServerListener {
  port: number;
  protocol: string;
  state: string;
  exposed: boolean;
  process: string;
}

export interface ServerPortScanResult {
  scanned: boolean;
  serverId: string;
  listeners: ServerListener[];
}

export interface ServerRateLimitConfig {
  rps: number;
  burst: number;
  whitelist: string[];
}

export interface ServerCapabilities {
  monitor?: boolean;
  terminal?: boolean;
  exec?: boolean;
  hostConfiguration?: boolean;
  ssh?: boolean;
  networkSettings?: boolean;
}

export interface CloudWorkspaceOperation {
  id: string;
  kind: string;
  status: string;
  requestedAt: string;
  nextAttemptAt?: string | null;
  error?: string | null;
  logs: string[];
}

export interface CloudWorkspaceSummary {
  id: string;
  serverId: string;
  name: string;
  planTierId?: string;
  subscriptionStatus?: string;
  projectCount?: number;
  state: string;
  resources?: {
    cpuCores: number;
    memoryMb: number;
    diskMb: number;
  };
  operation?: CloudWorkspaceOperation | null;
  createdAt?: string;
}

export interface ServerNetworkSettings {
  internetAccess?: boolean | null;
  ingressPorts: number[];
  ingressAll?: boolean;
  egress?: string[] | null;
  privateIp?: string | null;
  outboundIp?: string | null;
  outboundMode?: string | null;
  revision?: string;
}

export interface ManagedWorkload {
  id: string;
  name: string;
  state: string;
  restartPolicy?: string | null;
  source: 'manual' | 'project' | 'system';
  projectId?: string | null;
  manageable: boolean;
}

export interface ManagedSshStatus {
  enabled?: boolean | null;
  keyConfigured?: boolean | null;
  passwordConfigured?: boolean | null;
  requiresIdentityAccess: boolean;
  connection?: {
    user: string;
    host: string;
    bastion: string;
    command: string;
  } | null;
}

export interface ManagedRuntimeStatus {
  enabled?: boolean | null;
  running?: boolean | null;
}

export interface ManagedRuntimeCredential {
  endpoint: string;
  token: string;
  revision: string;
  expiresAt?: string | null;
}

export interface ManagedServerInfo {
  workspaceId: string;
  image: string;
  state: string;
  mode?: string | null;
  operatingSystem?: string | null;
  restartPolicy?: string | null;
  resources: {
    cpuCores: number;
    memoryMb: number;
    diskMb: number;
  };
}
