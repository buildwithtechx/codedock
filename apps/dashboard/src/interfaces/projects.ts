import type {
  DomainConfig,
  EnvironmentConfig,
  ProjectConfig,
  ProjectMember,
  ProjectToken,
} from '@codedock/contracts';
import type { BaseResponse, PaginatedData } from '#/interfaces/base';
import type { Database } from '#/interfaces/databases';
import type { AppService } from '#/interfaces/services';

export interface ServerlessFunctionCode {
  id: string;
  serviceId: string;
  runtime: string;
  codeContent: string;
  createdAt: string;
  updatedAt: string;
}

export interface CanvasSummary {
  id: string;
  name: string;
  description?: string;
  environmentsCount: number;
  appsCount: number;
  databasesCount: number;

  onlineServices: number;
  totalServices: number;
  serviceIcons: string[];
  stack?: string;
  deployTarget: 'local' | 'server';
  serverName?: string;
  defaultEnvironment?: EnvironmentConfig;
  createdAt: string;
  updatedAt: string;
}

export interface EnvironmentCanvas {
  revision: string;
  nodes: import('@xyflow/react').Node[];
  edges: import('#/features/canvas/topology-types').TopologyEdge[];
  environment: EnvironmentConfig;
  apps: AppService[];
  databases: Database[];
}

export type ListProjectsResponse = BaseResponse<PaginatedData<ProjectConfig>>;
export type GetProjectResponse = BaseResponse<ProjectConfig>;
export type CreateProjectResponse = BaseResponse<ProjectConfig>;
export type ListEnvironmentsResponse = BaseResponse<EnvironmentConfig[]>;
export type CreateEnvironmentResponse = BaseResponse<EnvironmentConfig>;
export type GetCanvasResponse = BaseResponse<CanvasSummary[]>;
export type GetEnvironmentCanvasResponse = BaseResponse<EnvironmentCanvas>;
export type ListDomainsResponse = BaseResponse<DomainConfig[]>;
export type CreateDomainResponse = BaseResponse<DomainConfig>;
export type ListMembersResponse = BaseResponse<ProjectMember[]>;
export type ListTokensResponse = BaseResponse<ProjectToken[]>;
