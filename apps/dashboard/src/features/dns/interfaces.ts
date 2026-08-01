export type {
  DNSProvisionStatus,
  DNSVerificationStatus,
  DomainBase,
  DomainVerifyResult,
} from './contracts';

import type { DomainBase } from './contracts';

export interface DNSRecord {
  id: string;
  domainName: string;
  recordType: string;
  recordName: string;
  recordValue: string;
  ttl: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateDNSRecordRequest {
  domainName: string;
  recordType: string;
  recordName: string;
  recordValue: string;
  ttl?: number;
}

export interface UpdateDNSRecordRequest {
  recordType: string;
  recordName: string;
  recordValue: string;
  ttl?: number;
}

export interface Domain extends DomainBase {
  sslStatus?: string;
}

export type DomainConfig = Domain;

export interface CreateDomainRequest {
  projectId?: string;
  serviceId?: string;
  domainName: string;
}
