export type {
  DNSProvisionStatus,
  DNSVerificationStatus,
  DomainBase,
  DomainVerifyResult,
} from './contracts';

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
