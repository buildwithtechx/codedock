import type { AuditFacets, AuditListParams, AuditLogRow, AuditSettings } from '#/interfaces/audit';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

function toQuery(params: AuditListParams): string {
  const query = new URLSearchParams();
  if (params.category) query.set('category', params.category);

  const limit = params.limit ?? params.perPage;
  if (limit !== undefined) query.set('limit', String(limit));

  let offset = params.offset;
  if (offset === undefined && params.page !== undefined && limit !== undefined) {
    offset = Math.max(0, (params.page - 1) * limit);
  }
  if (offset !== undefined) query.set('offset', String(offset));

  const qParts: string[] = [];
  if (params.q) qParts.push(params.q);
  if (params.actor) qParts.push(params.actor);
  if (params.source) qParts.push(params.source);
  if (qParts.length > 0) {
    query.set('q', qParts.join(' '));
  }

  const text = query.toString();
  return text ? `?${text}` : '';
}

export const auditLogsService = {
  list: async (params: AuditListParams = {}): Promise<BaseResponse<AuditLogRow[]>> => {
    try {
      return await apiClient.get<BaseResponse<AuditLogRow[]>>(`/audit-logs${toQuery(params)}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
  facets: async (
    params: { category?: string; q?: string } = {}
  ): Promise<BaseResponse<AuditFacets>> => {
    try {
      return await apiClient.get<BaseResponse<AuditFacets>>(`/audit-logs/facets${toQuery(params)}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
  getSettings: async (organizationId?: string): Promise<BaseResponse<AuditSettings>> => {
    try {
      const query = organizationId ? `?organizationId=${organizationId}` : '';
      return await apiClient.get<BaseResponse<AuditSettings>>(`/audit-logs/settings${query}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
  updateSettings: async (
    settings: Partial<AuditSettings>
  ): Promise<BaseResponse<AuditSettings>> => {
    try {
      return await apiClient.put<BaseResponse<AuditSettings>>('/audit-logs/settings', settings);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};

export const auditApi = auditLogsService;
