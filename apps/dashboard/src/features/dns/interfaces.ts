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

export interface Domain {
  id: string;
  projectId: string;
  domainName: string;
  sslStatus: string;
  createdAt: string;
}

export interface CreateDomainRequest {
  projectId: string;
  domainName: string;
}
