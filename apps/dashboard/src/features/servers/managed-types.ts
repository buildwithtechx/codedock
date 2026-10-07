export type ManagedTier = {
  name: string;
  cores: number;
  memoryGB: number;
  diskGB: number;
  monthlyPrice: number;
};
export type ManagedCatalog = {
  provider: string;
  tiers: ManagedTier[];
  regions: string[];
  images: string[];
};
export type ManagedCredential = {
  id: string;
  organizationId: string;
  provider: string;
  label: string;
  createdAt: string;
  updatedAt: string;
};
export type ManagedQuota = {
  organizationId: string;
  maxServers: number;
  maxMemoryGB: number;
  usedServers: number;
  usedMemoryGB: number;
  updatedAt: string;
};
export type ManagedServerPlan = {
  credentialId: string;
  organizationId: string;
  userId: string;
  serverId: string;
  externalId: string;
  name: string;
  region: string;
  serverType: string;
  resizeTo?: string;
  image: string;
  sshKeyName: string;
  provider: string;
  tier: ManagedTier;
  monthlyPrice: number;
  action: string;
};
export type ManagedOperation = {
  id: string;
  status: string;
  phase: string;
  effects: string;
  logs: string;
  error: string;
  expiresAt: number;
};
export type ManagedReview = {
  review: { operation: ManagedOperation; confirmation: string };
  plan: ManagedServerPlan;
};
export type CreateManagedCredentialRequest = {
  provider: string;
  label: string;
  token: string;
};
export type ReviewManagedProvisionRequest = {
  credentialId: string;
  serverId?: string;
  name: string;
  region: string;
  serverType: string;
  image: string;
};
