export interface Server {
  id: string;
  userId: string;
  name: string;
  ipAddress: string;
  status: string;
  workerToken: string;
  lastSeenAt: string;
  metrics: any;
  createdAt: string;
  updatedAt: string;
}

export interface CreateServerRequest {
  name: string;
  ipAddress: string;
}
