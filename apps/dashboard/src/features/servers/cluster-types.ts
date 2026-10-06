export type ClusterNode = {
  serverId: string;
  privateIp: string;
  interface: string;
  fingerprint: string;
};
export type Cluster = {
  id: string;
  name: string;
  version: string;
  nodes: ClusterNode[];
  revision: number;
  status: string;
  error: string;
};
export type ClusterOperation = {
  id: string;
  status: string;
  phase: string;
  effects: string;
  logs: string;
  error: string;
  expiresAt: number;
};
export type ClusterReview = { operation: ClusterOperation; confirmation: string };
