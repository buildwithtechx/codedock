import type { CodedockDeployment } from '../interfaces';
import { apiRequest } from './client';

export const DeploymentsApi = {
  async listDeployments(
    url: string,
    token: string,
    limit: number = 20
  ): Promise<CodedockDeployment[]> {
    const res = await apiRequest<CodedockDeployment[]>(
      url,
      token,
      `/api/deployments?limit=${limit}`
    );
    return Array.isArray(res) ? res : [];
  },

  async getDeploymentLogs(url: string, token: string, deploymentId: string): Promise<string> {
    return await apiRequest<string>(
      url,
      token,
      `/api/deployments/${encodeURIComponent(deploymentId)}/logs`
    );
  },
};
