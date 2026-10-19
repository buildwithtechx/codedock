import type {
  Deployment,
  DeploymentFailureExplanation,
  ListOrganizationDeploymentsParams,
  OrganizationDeployment,
  PRPreview,
  TriggerDeploymentRequest,
} from '@codedock/contracts';
import type { HttpClient } from '../http';

export function createDeploymentsOperations(http: HttpClient) {
  return {
    listByOrganization: (
      organizationId: string,
      params: ListOrganizationDeploymentsParams = {}
    ): Promise<OrganizationDeployment[]> => {
      const query = new URLSearchParams({ organizationId });
      if (params.projectId) query.set('projectId', params.projectId);
      if (params.serviceId) query.set('serviceId', params.serviceId);
      if (params.status) query.set('status', params.status);
      if (params.search) query.set('search', params.search);
      if (params.days) query.set('days', String(params.days));
      if (params.page) query.set('page', String(params.page));
      if (params.limit) query.set('limit', String(params.limit));
      return http.get<OrganizationDeployment[]>(`/deployments?${query.toString()}`);
    },

    listByService: (serviceId: string): Promise<Deployment[]> =>
      http.get<Deployment[]>(`/services/${encodeURIComponent(serviceId)}/deployments`),

    listPRPreviews: (serviceId: string): Promise<PRPreview[]> =>
      http.get<PRPreview[]>(`/services/${encodeURIComponent(serviceId)}/previews`),

    get: (deploymentId: string): Promise<Deployment> =>
      http.get<Deployment>(`/deployments/${encodeURIComponent(deploymentId)}`),

    trigger: (serviceId: string, payload?: TriggerDeploymentRequest): Promise<Deployment> =>
      http.post<Deployment>(`/services/${encodeURIComponent(serviceId)}/deploy`, payload),

    triggerProject: (projectId: string, environmentId?: string): Promise<void> => {
      const url = environmentId
        ? `/projects/${encodeURIComponent(projectId)}/deploy?environmentId=${encodeURIComponent(environmentId)}`
        : `/projects/${encodeURIComponent(projectId)}/deploy`;
      return http.post<void>(url);
    },

    cancel: (deploymentId: string): Promise<Deployment> =>
      http.post<Deployment>(`/deployments/${encodeURIComponent(deploymentId)}/cancel`),

    rollback: (deploymentId: string): Promise<Deployment> =>
      http.post<Deployment>(`/deployments/${encodeURIComponent(deploymentId)}/rollback`),

    getLogs: (deploymentId: string): Promise<string> =>
      http.get<string>(`/deployments/${encodeURIComponent(deploymentId)}/logs`),

    togglePin: (deploymentId: string, pinned: boolean): Promise<Deployment> =>
      http.post<Deployment>(`/deployments/${encodeURIComponent(deploymentId)}/pin`, { pinned }),

    explainFailure: (deploymentId: string): Promise<DeploymentFailureExplanation> =>
      http.get<DeploymentFailureExplanation>(
        `/deployments/${encodeURIComponent(deploymentId)}/explain-failure`
      ),
  };
}
