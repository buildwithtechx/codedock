export interface AccountInfo {
  id: string;
  email: string;
  name?: string;
  organizationId?: string;
  role?: string;
  tier?: string;
  createdAt?: string;
}

export interface PairingHandshake {
  code: string;
  instanceId?: string;
  expiresAt: string;
}

export interface PairingClaimResult {
  token: string;
  account: AccountInfo;
  organizationId?: string;
}
