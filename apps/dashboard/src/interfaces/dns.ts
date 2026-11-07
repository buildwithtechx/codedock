export type DNSProvisionStatus = 'provisioned' | 'failed' | 'pending';

export type DNSVerificationStatus =
  | 'resolves_to_server'
  | 'resolves_to_different_ip'
  | 'server_ip_not_configured'
  | 'unresolved';

export interface DomainBase {
  id: string;
  serviceId?: string;
  projectId?: string;
  domainName: string;
  redirectTo?: string;
  pathPrefix?: string;
  sslCertStatus?: string;
  dnsProvisionStatus?: DNSProvisionStatus;
  dnsProvider?: string;
  verified?: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface DomainVerifyResult {
  domainId: string;
  domainName: string;
  verified: boolean;
  status: DNSVerificationStatus;
  resolvedIp?: string;
  serverIp?: string;
  message: string;
}

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
