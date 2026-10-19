export type MemberPermission = 'admin' | 'member' | 'viewer';
export type MemberStatus = 'pending' | 'accepted';

export type DNSProvisionStatus = 'pending' | 'provisioned' | 'failed' | 'none';
export type DNSVerificationStatus = 'pending' | 'verified' | 'failed' | 'none';

export interface DomainBase {
  id?: string;
  domain: string;
  isCustom?: boolean;
  dnsProvisionStatus?: DNSProvisionStatus;
  dnsVerificationStatus?: DNSVerificationStatus;
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
