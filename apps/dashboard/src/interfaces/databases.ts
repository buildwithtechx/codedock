import type { BaseResponse } from '#/interfaces/base';

export type DatabaseEngine =
  | 'postgres'
  | 'postgresql'
  | 'mysql'
  | 'mariadb'
  | 'redis'
  | 'mongodb'
  | 'clickhouse'
  | 'dragonfly'
  | 'keydb'
  | 'timescaledb'
  | 'kafka'
  | 'rabbitmq'
  | 'nats';

export interface Database {
  id: string;
  projectId: string;
  environmentId?: string;
  name: string;
  engine: DatabaseEngine;
  version: string;
  username: string;
  password?: string;
  databaseName: string;
  internalHost: string;
  internalDns?: string;
  externalDns?: string;
  volumePath?: string;
  port: number;
  externalPort?: number;
  status: 'deploying' | 'running' | 'stopped' | 'error';
  createdAt: string;
  updatedAt: string;
}

export interface CreateDatabaseRequest {
  projectId: string;
  environmentId?: string;
  name: string;
  engine: DatabaseEngine;
  version?: string;
  databaseName?: string;
  username?: string;
  password?: string;
}

export type ListDatabasesResponse = BaseResponse<Database[]>;
export type GetDatabasesResponse = BaseResponse<Database[]>;
export type GetDatabaseResponse = BaseResponse<Database>;
export type CreateDatabaseResponse = BaseResponse<Database>;
export type DeleteDatabaseResponse = BaseResponse<void>;

export interface QueryDatabaseRequest {
  query: string;
}

export interface QueryDatabaseResponse {
  columns: string[];
  rows: Record<string, any>[];
  rowCount: number;
  executionTimeMs: number;
}

export type DatabaseQueryResponseType = BaseResponse<QueryDatabaseResponse>;

export interface DatabaseTableSchema {
  name: string;
  columns: Array<{ name: string; type?: string }>;
}

export type ListTablesResponse = BaseResponse<DatabaseTableSchema[]>;

export interface TableRowPayload {
  [key: string]: any;
}

export interface ImportDatabaseRequest {
  sql: string;
}

export type ImportDatabaseResponse = BaseResponse<void>;

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

export type DatabaseReview = {
  operation: DatabaseOperation;
  confirmation: string;
};
