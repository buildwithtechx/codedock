import type { AuditLog } from '#/interfaces/audit';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export interface AuditLogRow extends AuditLog {
  category?: string;
}

export interface AuditCategoryFacet {
  id: string;
  label: string;
  description: string;
  count: number;
}

export interface AuditFacets {
  total: number;
  categories: AuditCategoryFacet[];
}

export interface AuditListParams {
  category?: string;
  limit?: number;
  offset?: number;
}

function toQuery(params: AuditListParams): string {
  const query = new URLSearchParams();
  if (params.category) query.set('category', params.category);
  if (params.limit !== undefined) query.set('limit', String(params.limit));
  if (params.offset !== undefined) query.set('offset', String(params.offset));
  const text = query.toString();
  return text ? `?${text}` : '';
}

export const auditApi = {
  list: async (params: AuditListParams = {}): Promise<BaseResponse<AuditLogRow[]>> => {
    try {
      return await apiClient.get<BaseResponse<AuditLogRow[]>>(`/audit-logs${toQuery(params)}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
  facets: async (): Promise<BaseResponse<AuditFacets>> => {
    try {
      return await apiClient.get<BaseResponse<AuditFacets>>('/audit-logs/facets');
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
