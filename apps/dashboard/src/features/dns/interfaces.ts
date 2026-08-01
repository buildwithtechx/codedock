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

export type DNSProvisionStatus = 'provisioned' | 'failed' | 'pending' | 'manual';
export type DNSVerificationStatus =
  | 'resolves_to_server'
  | 'resolves_to_different_ip'
  | 'unresolved';

export interface DomainVerifyResult {
  domainId: string;
  domainName: string;
  verified: boolean;
  status: DNSVerificationStatus;
  resolvedIps?: string[];
  serverIp?: string;
  message: string;
}

export interface Domain {
  id: string;
  projectId?: string;
  serviceId?: string;
  domainName: string;
  sslStatus?: string;
  sslCertStatus?: string;
  dnsProvisionStatus?: DNSProvisionStatus;
  verified?: boolean;
  createdAt: string;
}

export interface CreateDomainRequest {
  projectId?: string;
  serviceId?: string;
  domainName: string;
}
