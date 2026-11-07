/** Freshness status of a runtime container image. */
export type RuntimeImageStatus = 'current' | 'update_available' | 'pinned' | 'unknown';

/** Version metadata and status for a service or base container image. */
export interface RuntimeImageInfo {
  serviceId: string | null;
  serviceName: string | null;
  role: 'image' | 'base';
  ref: string;
  currentDigest: string | null;
  latestDigest: string | null;
  status: RuntimeImageStatus;
}

/** Result returned by a runtime image freshness inspection. */
export interface RuntimeImageCheckResult {
  checkedAt: string;
  images: RuntimeImageInfo[];
}

/** Input to trigger runtime base image updates and redeployments. */
export interface UpdateRuntimeInput {
  serviceIds?: string[];
}

/** Security advisory severity level. */
export type SecurityAdvisorySeverity = 'critical' | 'high' | 'medium' | 'low' | 'info';

/** Signed runtime security advisory from upstream feed. */
export interface SecurityAdvisory {
  id: string;
  revision: number;
  title: string;
  severity: SecurityAdvisorySeverity;
  message: string;
  url: string;
  action?: string;
}

/** Result of non-invasive security advisory match against a project/service. */
export interface AdvisoryFinding {
  advisoryId: string;
  title: string;
  severity: string;
  action: string;
  confidence: 'confirmed' | 'possible';
  serviceName?: string;
  evidence?: string;
  canRefresh?: boolean;
  runtimeUpgradeRequired?: boolean;
  imageChangeRequired?: boolean;
}

/** Security audit findings and policies for a project. */
export interface ProjectSecurityResult {
  projectId: string;
  deploymentId?: string;
  checkedAt?: string;
  autoRefresh: boolean;
  findings: AdvisoryFinding[];
}

/** Signed security advisories feed response including affected projects. */
export interface AdvisoriesResponse {
  revision: number;
  checkedAt?: string;
  expiresAt?: string;
  stale: boolean;
  advisories: SecurityAdvisory[];
  affectedProjects: ProjectSecurityResult[];
}

/** Queued or completed security remediation attempt. */
export interface RefreshAttempt {
  id: string;
  projectId: string;
  deploymentId?: string;
  advisoryId: string;
  status: string;
  error?: string;
  createdAt: string;
}

/** Policy update for automatic security patch refreshing. */
export interface SetSecurityPolicyInput {
  autoRefresh: boolean;
}

/** Input to apply a verified security advisory patch/remediation. */
export interface ApplySecurityInput {
  advisoryId: string;
  deploymentId?: string;
}
