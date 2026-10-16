import { apiClient } from '#/lib/api-client';

export interface BuildAdmissionStatus {
  canAdmit: boolean;
  activeBuilds: number;
  maxConcurrentBuilds: number;
  queuePosition: number;
  estimatedWaitSeconds: number;
  reason: string;
}

export interface ProjectChecksConfig {
  projectId: string;
  enabled: boolean;
  reportCommitStatus: boolean;
  requirePassingChecks: boolean;
  targetBranches: string[];
  contextName: string;
  updatedAt: string;
}

export interface UpdateChecksInput {
  enabled: boolean;
  reportCommitStatus: boolean;
  requirePassingChecks: boolean;
  targetBranches?: string[];
  contextName?: string;
}

export const checksService = {
  getAdmission: async (projectId: string): Promise<BuildAdmissionStatus> => {
    const res = await apiClient.get<{ data: BuildAdmissionStatus }>(
      `/projects/${encodeURIComponent(projectId)}/deployments/admission`
    );
    return res.data;
  },

  getChecks: async (projectId: string): Promise<ProjectChecksConfig> => {
    const res = await apiClient.get<{ data: ProjectChecksConfig }>(
      `/projects/${encodeURIComponent(projectId)}/checks`
    );
    return res.data;
  },

  updateChecks: async (projectId: string, input: UpdateChecksInput): Promise<ProjectChecksConfig> => {
    const res = await apiClient.put<{ data: ProjectChecksConfig }>(
      `/projects/${encodeURIComponent(projectId)}/checks`,
      input
    );
    return res.data;
  },
};
