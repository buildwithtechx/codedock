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

export interface OneClickApp {
  id: string;
  name: string;
  description: string;
  icon: string;
  category: string;
  dockerImage: string;
  defaultPort: number;
  envVariables: Array<{ key: string; label: string; defaultValue?: string }>;
}
