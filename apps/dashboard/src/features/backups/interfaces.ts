import type { BaseResponse } from '#/interfaces/base';

export type BackupConfigStatus = 'active' | 'inactive';
export type BackupRecordStatus = 'running' | 'completed' | 'failed' | 'expiring' | 'expired';
export type DestinationKind = 's3' | 'sftp';

export interface BackupConfig {
  preDeployment?: boolean;
  id: string;
  projectId: string;
  databaseId?: string;
  serviceId?: string;
  volumeName?: string;
  s3DestinationId?: string;
  sftpDestinationId?: string;
  name: string;
  description: string;
  dbUser: string;
  dbPassword?: string;
  backupEnabled: boolean;
  s3Enabled: boolean;
  sftpEnabled?: boolean;
  disableLocal: boolean;
  schedule: string;
  timezone: string;
  timeout: number;
  retentionDays: number;
  maxBackups: number;
  maxStorageGb: number;
  status: BackupConfigStatus;
  createdAt: string;
  updatedAt: string;
}

export interface BackupRecord {
  protectedUntil?: number;
  sha256?: string;
  verifiedAt?: string;
  id: string;
  backupConfigId: string;
  projectId: string;
  databaseId?: string;
  s3DestinationId?: string;
  sftpDestinationId?: string;
  status: BackupRecordStatus;
  filePath: string;
  fileSizeBytes: number;
  s3Url?: string;
  sftpUrl?: string;
  logs: string;
  startedAt: string;
  completedAt: string;
}

export interface S3Destination {
  id: string;
  projectId: string;
  name: string;
  description: string;
  provider: string;
  endpoint: string;
  bucket: string;
  region: string;
  pathPrefix?: string;
  isDefault?: boolean;
  lastVerifiedAt?: string;
  lastVerifyError?: string;
  accessKeyId: string;
  secretAccessKey: string;
  createdAt: string;
}

export interface SFTPDestination {
  id: string;
  projectId?: string;
  name: string;
  description: string;
  host: string;
  port: number;
  username: string;
  pathPrefix: string;
  lastVerifiedAt?: string;
  lastVerifyError?: string;
  createdAt: string;
}

export interface CreateBackupConfigRequest {
  preDeployment?: boolean;
  projectId: string;
  name: string;
  description: string;
  dbUser: string;
  dbPassword?: string;
  backupEnabled: boolean;
  s3Enabled: boolean;
  disableLocal: boolean;
  schedule: string;
  timezone: string;
  timeout: number;
  retentionDays: number;
  maxBackups: number;
  maxStorageGb: number;
  databaseId?: string;
  serviceId?: string;
  volumeName?: string;
  s3DestinationId?: string;
}

export interface CreateS3DestinationRequest {
  projectId: string;
  name: string;
  description: string;
  provider: string;
  endpoint: string;
  bucket: string;
  region: string;
  pathPrefix?: string;
  isDefault?: boolean;
  accessKeyId: string;
  secretAccessKey: string;
}

export interface UpdateS3DestinationRequest {
  name: string;
  description: string;
  provider: string;
  endpoint: string;
  bucket: string;
  region: string;
  pathPrefix?: string;
  isDefault?: boolean;
  accessKeyId: string;
  secretAccessKey: string;
}

export interface CreateSFTPDestinationRequest {
  projectId?: string;
  name: string;
  description?: string;
  host: string;
  port: number;
  username: string;
  password?: string;
  privateKey?: string;
  pathPrefix?: string;
}

export interface VerifyS3Response {
  ok: boolean;
  reason?: string;
}

export interface RestoreOperation {
  id: string;
  status: string;
  phase: string;
  effects: string;
  error: string;
  logs: string;
  expiresAt: number;
  cancelRequested?: boolean;
  meta?: Record<string, unknown>;
}

export interface RestoreReview {
  operation: RestoreOperation;
  confirmation: string;
}

export type ListBackupsResponse = BaseResponse<BackupConfig[]>;
export type GetBackupResponse = BaseResponse<BackupConfig>;
export type CreateBackupResponse = BaseResponse<BackupConfig>;
export type ListBackupRecordsResponse = BaseResponse<BackupRecord[]>;
export type ListS3DestinationsResponse = BaseResponse<S3Destination[]>;
export type CreateS3DestinationResponse = BaseResponse<S3Destination>;
export type UpdateS3DestinationResponse = BaseResponse<S3Destination>;
export type ListSFTPDestinationsResponse = BaseResponse<SFTPDestination[]>;
export type CreateSFTPDestinationResponse = BaseResponse<SFTPDestination>;
export type VerifyS3DestinationResponse = BaseResponse<VerifyS3Response>;
export type VerifySFTPDestinationResponse = BaseResponse<VerifyS3Response>;
export type RestoreReviewResponse = BaseResponse<RestoreReview>;
export type RestoreOperationResponse = BaseResponse<RestoreOperation>;
