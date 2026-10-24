export type MemberPermission = 'admin' | 'member' | 'viewer';
export type MemberStatus = 'pending' | 'accepted';

export type DNSProvisionStatus = 'pending' | 'provisioned' | 'failed' | 'none';
export type DNSVerificationStatus = 'pending' | 'verified' | 'failed' | 'none';

export interface DomainBase {
  id: string;
  domain?: string;
  domainName: string;
  serviceId?: string;
  projectId?: string;
  isCustom?: boolean;
  sslCertStatus?: string;
  dnsProvisionStatus?: DNSProvisionStatus;
  dnsVerificationStatus?: DNSVerificationStatus;
  verified?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface DomainVerifyResult {
  verified: boolean;
  message?: string;
}

export interface DomainConfig extends DomainBase {
  redirectTo?: string;
  pathPrefix?: string;
}

export interface ProjectConfig {
  id: string;
  name: string;
  description?: string;
  defaultRollbackStrategy?: 'snapshot' | 'git';
  rollbackWindow?: number | null;
  snapshotSizeBytes?: number | null;
  capacityMeasuredAt?: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface RollbackCapacityUI {
  window: number;
  source: 'explicit' | 'instance-default';
  explicit: number | null;
  snapshotSizeBytes: number | null;
  measuredAt: string | null;
  diskFreeBytes: number | null;
  diskTotalBytes: number | null;
  maxWindow: number;
  diskBudgetFraction: number;
  strategy: 'snapshot' | 'git';
}

export interface UpdateRollbackSettingsRequest {
  defaultRollbackStrategy?: 'snapshot' | 'git';
  rollbackWindow?: number | null;
}

export interface EnvironmentConfig {
  id: string;
  projectId: string;
  name: string;
  isDefault: boolean;
  createdAt: string;
  updatedAt: string;
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
  gitUrl?: string;
  gitBranch?: string;
  gitProvider?: string;
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

export interface ServiceWebhook {
  id: string;
  serviceId: string;
  url: string;
  eventTypes: string[];
  includePrEnvironments: boolean;
  createdAt: string;
  updatedAt: string;
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
