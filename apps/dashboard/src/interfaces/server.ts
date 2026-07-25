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
  metrics: ServerMetrics | null;
  createdAt: string;
  updatedAt: string;
}

export interface CreateServerRequest {
  name: string;
  ipAddress: string;
}
