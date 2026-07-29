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
  status: string;
  workerToken: string;
  lastSeenAt: string;
  metrics: ServerMetrics | string | null;
  createdAt: string;
  updatedAt: string;
}

export function parseServerMetrics(
  metrics: ServerMetrics | string | null | undefined
): ServerMetrics | null {
  if (!metrics) return null;
  if (typeof metrics === 'object') return metrics;
  if (typeof metrics === 'string') {
    try {
      const decoded = metrics.startsWith('{') ? metrics : atob(metrics);
      return JSON.parse(decoded) as ServerMetrics;
    } catch {
      return null;
    }
  }
  return null;
}

export interface CreateServerRequest {
  name: string;
  ipAddress: string;
}
