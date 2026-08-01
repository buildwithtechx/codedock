import type { Database } from '#/features/databases';
import type { DNSVerificationStatus, DomainBase } from '#/features/dns/contracts';
import type { AppService } from '#/features/services';
import type { BaseResponse, PaginatedData } from '#/interfaces/base';

export type {
  DNSProvisionStatus,
  DNSVerificationStatus,
  DomainBase,
  DomainVerifyResult,
} from '#/features/dns/contracts';

export type MemberPermission = 'admin' | 'member' | 'viewer';
export type MemberStatus = 'pending' | 'accepted';

export interface ProjectConfig {
  id: string;
  name: string;
  description?: string;
  createdAt: string;
  updatedAt: string;
}

export interface EnvironmentConfig {
  id: string;
  projectId: string;
  name: string;
  isDefault: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface DomainConfig extends DomainBase {
  redirectTo?: string;
  pathPrefix?: string;
  verificationStatus?: DNSVerificationStatus;
}

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
  defaultEnvironment?: EnvironmentConfig;
  createdAt: string;
  updatedAt: string;
}

export interface EnvironmentCanvas {
  environment: EnvironmentConfig;
  apps: AppService[];
  databases: Database[];
}

export interface CreateProjectRequest {
  name: string;
  description?: string;
  repositoryUrl?: string;
  branch?: string;
  internalPort?: number;
  domain?: string;
  serverId?: string;
  organizationId?: string;
}

export interface ProjectToken {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  tokenPrefix: string;
  scopes: string[];
  ipAllowlist: string[];
  expiresAt?: string;
  createdAt: string;
}

export interface ProjectMember {
  id: string;
  projectId: string;
  userId?: string;
  email: string;
  permission: MemberPermission;
  status: MemberStatus;
  invitedAt: string;
  acceptedAt?: string;
}

export interface CreateWebhookRequest {
  url: string;
  eventTypes: string[];
  includePrEnvironments: boolean;
}

export interface CreateTokenRequest {
  name: string;
  environmentId: string;
  scopes: string[];
  ipAllowlist?: string[];
  expiresAt?: string;
}

export interface AddMemberRequest {
  email: string;
  permission: MemberPermission;
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
