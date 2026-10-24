export type DeploymentStatus =
  | 'pending'
  | 'BUILDING'
  | 'PULLING'
  | 'CLONING'
  | 'READY'
  | 'ACTIVE'
  | 'REMOVED'
  | 'SLEPT'
  | 'FAILED'
  | 'SUCCESS';

export interface Deployment {
  id: string;
  serviceId: string;
  environmentId: string;
  projectId: string;
  status: DeploymentStatus;
  branch?: string;
  commitHash?: string;
  commitMessage?: string;
  trigger?: string;
  version?: number;
  commitShaBefore?: string;
  buildLogs?: string;
  containerId?: string;
  artifactRetainedAt?: string | null;
  pinned?: boolean;
  /** Storage pressure, capacity failure, or crash diagnostics on this deployment. */
  resourceIssue?: string;
  /** Background post-deploy crash loop diagnosis with exit code and log tail. */
  restartLoop?: string;
  createdAt: string;
  updatedAt: string;
  finishedAt?: string;
}

export interface OrganizationDeployment extends Deployment {
  projectName: string;
  serviceName: string;
}

export interface ListOrganizationDeploymentsParams {
  projectId?: string;
  serviceId?: string;
  status?: string;
  search?: string;
  days?: number;
  page?: number;
  limit?: number;
}

export interface TriggerDeploymentRequest {
  commitId?: string;
}

export type PRPreviewStatus = 'active' | 'inactive' | 'building' | 'failed';

export interface PRPreview {
  id: string;
  serviceId: string;
  projectId: string;
  prNumber: number;
  branch: string;
  commitHash: string;
  status: PRPreviewStatus;
  previewDomain: string;
  containerId: string;
  createdAt: string;
  updatedAt: string;
}

export interface DeploymentFailureExplanation {
  summary: string;
  cause: string;
  suggestedFix: string;
  confidence: 'low' | 'medium' | 'high';
  commands: string[];
  relatedLogLines: string[];
}
