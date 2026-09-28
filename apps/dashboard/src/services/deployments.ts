import type {
  ExplainDeploymentResponse,
  GetDeploymentLogsResponse,
  GetDiagnosticsResponse,
  GetServiceMetricsResponse,
  ListDeploymentsResponse,
  ListOrganizationDeploymentsParams,
  ListOrganizationDeploymentsResponse,
  ListPRPreviewsResponse,
  RollbackDeploymentResponse,
  TriggerDeploymentRequest,
  TriggerDeploymentResponse,
} from '#/features/services';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const deploymentsService = {
  listByOrganization: async (
    organizationId: string,
    params: ListOrganizationDeploymentsParams = {}
  ): Promise<ListOrganizationDeploymentsResponse> => {
    try {
      const query = new URLSearchParams({ organizationId });
      if (params.projectId) query.set('projectId', params.projectId);
      if (params.serviceId) query.set('serviceId', params.serviceId);
      if (params.status) query.set('status', params.status);
      if (params.search) query.set('search', params.search);
      if (params.page) query.set('page', String(params.page));
      if (params.limit) query.set('limit', String(params.limit));
      return await apiClient.get<ListOrganizationDeploymentsResponse>(`/deployments?${query}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listByService: async (serviceId: string): Promise<ListDeploymentsResponse> => {
    try {
      return await apiClient.get<ListDeploymentsResponse>(`/services/${serviceId}/deployments`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listPRPreviews: async (serviceId: string): Promise<ListPRPreviewsResponse> => {
    try {
      return await apiClient.get<ListPRPreviewsResponse>(`/services/${serviceId}/previews`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  trigger: async (
    serviceId: string,
    payload?: TriggerDeploymentRequest
  ): Promise<TriggerDeploymentResponse> => {
    try {
      return await apiClient.post<TriggerDeploymentResponse>(
        `/services/${serviceId}/deploy`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  rollback: async (deploymentId: string): Promise<RollbackDeploymentResponse> => {
    try {
      return await apiClient.post<RollbackDeploymentResponse>(
        `/deployments/${deploymentId}/rollback`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getLogs: async (deploymentId: string): Promise<GetDeploymentLogsResponse> => {
    try {
      return await apiClient.get<GetDeploymentLogsResponse>(`/deployments/${deploymentId}/logs`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getMetrics: async (serviceId: string): Promise<GetServiceMetricsResponse> => {
    try {
      return await apiClient.get<GetServiceMetricsResponse>(`/services/${serviceId}/metrics`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  diagnostics: async (deploymentId: string): Promise<GetDiagnosticsResponse> => {
    try {
      return await apiClient.post<GetDiagnosticsResponse>(
        `/deployments/${deploymentId}/diagnostics`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  explainFailure: async (deploymentId: string): Promise<ExplainDeploymentResponse> => {
    try {
      return await apiClient.get<ExplainDeploymentResponse>(`/deployments/${deploymentId}/explain`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
