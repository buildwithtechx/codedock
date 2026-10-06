export type ClusterDatabaseSpec = {
  id: string;
  name: string;
  environmentId: string;
  engine: 'postgres' | 'redis';
  image: string;
  instances: number;
  storageGiB: number;
  storageClass: string;
  s3DestinationId: string;
};
export type ClusterDatabase = {
  id: string;
  clusterId: string;
  projectId: string;
  spec: ClusterDatabaseSpec;
  status: string;
  error: string;
};
export type DatabaseAction = 'operators' | 'create' | 'restore' | 'recover' | 'backup';
export type DatabaseOperation = {
  id: string;
  status: string;
  phase: string;
  effects: string;
  logs: string;
  error: string;
};
export type DatabaseReview = { operation: DatabaseOperation; confirmation: string };
