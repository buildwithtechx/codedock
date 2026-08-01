export type DNSProvisionStatus = 'provisioned' | 'failed' | 'pending' | 'manual';

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
