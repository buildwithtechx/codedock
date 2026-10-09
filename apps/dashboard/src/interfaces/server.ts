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
