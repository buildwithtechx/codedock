import type {
  AnalyticsDashboard,
  AnalyticsDeploymentStats,
  AnalyticsGeo,
  AnalyticsOverview,
  AnalyticsSummary,
  AttentionIssue,
  ServiceUsage,
} from '#/features/servers/analytics-types';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

function withQuery(endpoint: string, params: Record<string, string | number>): string {
  const query = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== '') {
      query.set(key, String(value));
    }
  }
  const suffix = query.toString();
  return suffix ? `${endpoint}?${suffix}` : endpoint;
}

export const analyticsService = {
  summary: async (
    organizationId: string,
    projectId: string,
    params: Record<string, string> = {}
  ): Promise<AnalyticsSummary> => {
    try {
      const res = await apiClient.get<BaseResponse<AnalyticsSummary>>(
        withQuery(
          `/organizations/${organizationId}/projects/${projectId}/analytics/summary`,
          params
        )
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  overview: async (
    organizationId: string,
    projectId: string,
    params: Record<string, string> = {}
  ): Promise<AnalyticsOverview> => {
    try {
      const res = await apiClient.get<BaseResponse<AnalyticsOverview>>(
        withQuery(
          `/organizations/${organizationId}/projects/${projectId}/analytics/overview`,
          params
        )
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  geo: async (
    organizationId: string,
    projectId: string,
    params: Record<string, string> = {}
  ): Promise<AnalyticsGeo> => {
    try {
      const res = await apiClient.get<BaseResponse<AnalyticsGeo>>(
        withQuery(`/organizations/${organizationId}/projects/${projectId}/analytics/geo`, params)
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  setPaths: async (organizationId: string, projectId: string, enabled: boolean): Promise<void> => {
    try {
      await apiClient.put<BaseResponse<null>>(
        `/organizations/${organizationId}/projects/${projectId}/analytics/paths`,
        { enabled }
      );
    } catch (err) {
      throw handleApiError(err);
    }
  },
  deploymentStats: async (
    organizationId: string,
    projectId: string,
    days = 30
  ): Promise<AnalyticsDeploymentStats> => {
    try {
      const res = await apiClient.get<BaseResponse<AnalyticsDeploymentStats>>(
        withQuery(`/organizations/${organizationId}/projects/${projectId}/analytics/deployments`, {
          days,
        })
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  dashboard: async (organizationId: string): Promise<AnalyticsDashboard> => {
    try {
      const res = await apiClient.get<BaseResponse<AnalyticsDashboard>>(
        `/organizations/${organizationId}/analytics/dashboard`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  usage: async (organizationId: string, projectId: string): Promise<ServiceUsage[]> => {
    try {
      const res = await apiClient.get<BaseResponse<ServiceUsage[]>>(
        `/organizations/${organizationId}/projects/${projectId}/analytics/usage`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
};

export const attentionService = {
  list: async (organizationId: string, status = ''): Promise<AttentionIssue[]> => {
    try {
      const res = await apiClient.get<BaseResponse<AttentionIssue[]>>(
        withQuery(`/organizations/${organizationId}/attention`, status ? { status } : {})
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  evaluate: async (organizationId: string): Promise<void> => {
    try {
      await apiClient.post<BaseResponse<null>>(
        `/organizations/${organizationId}/attention/evaluate`
      );
    } catch (err) {
      throw handleApiError(err);
    }
  },
  acknowledge: async (organizationId: string, issueId: string): Promise<void> => {
    try {
      await apiClient.post<BaseResponse<null>>(
        `/organizations/${organizationId}/attention/${issueId}/ack`
      );
    } catch (err) {
      throw handleApiError(err);
    }
  },
  resolve: async (organizationId: string, issueId: string): Promise<void> => {
    try {
      await apiClient.post<BaseResponse<null>>(
        `/organizations/${organizationId}/attention/${issueId}/resolve`
      );
    } catch (err) {
      throw handleApiError(err);
    }
  },
  act: async (
    organizationId: string,
    issueId: string,
    payload: { action?: string; params?: Record<string, string> }
  ): Promise<{ outcome: string }> => {
    try {
      const res = await apiClient.post<BaseResponse<{ outcome: string }>>(
        `/organizations/${organizationId}/attention/${issueId}/act`,
        payload
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
};
