/** Information for an individual Fail2ban jail. */
export interface Fail2banJail {
  jail: string;
  currentlyFailed: number;
  totalFailed: number;
  currentlyBanned: number;
  totalBanned: number;
  bannedIps: string[];
}

/** Fail2ban live status response covering all jails. */
export interface Fail2banStatus {
  jails: string[];
  details: Record<string, Fail2banJail>;
  checkedAt: string;
  canManage: boolean;
}

/** Request payload to ban or unban an IP address in Fail2ban. */
export interface Fail2banActionRequest {
  ip: string;
  jail?: string;
}

/** Result of a Fail2ban ban/unban action. */
export interface Fail2banActionResponse {
  success: boolean;
  output?: string;
}

/** Mail alias mapping an incoming address to a destination mailbox. */
export interface MailAlias {
  id: string;
  domainId: string;
  source: string;
  destination: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

/** Payload to create a new mail alias. */
export interface CreateMailAliasRequest {
  source: string;
  destination: string;
}

/** Mail server connection and metadata. */
export interface MailServer {
  id: string;
  organizationId: string;
  serverId?: string;
  domain: string;
  maildUrl: string;
  webmailUrl?: string;
  createdAt: string;
  updatedAt: string;
}

/** Mail domain record. */
export interface MailDomain {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

/** Mailbox account. */
export interface Mailbox {
  id: string;
  domainId: string;
  localPart: string;
  address: string;
  createdAt: string;
  updatedAt: string;
}
