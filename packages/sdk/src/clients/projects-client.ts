import type {
  AddMemberRequest,
  CreateProjectRequest,
  CreateTokenRequest,
  CreateWebhookRequest,
  DomainConfig,
  PaginatedData,
  ProjectConfig,
  ProjectMember,
  ProjectToken,
  ServiceWebhook,
} from '@codedock/contracts';
import type { HttpClient } from '../http';

export function createProjectsOperations(http: HttpClient) {
  return {
    list: (params?: {
      organizationId?: string;
      page?: number;
      limit?: number;
    }): Promise<PaginatedData<ProjectConfig>> => {
      const search = new URLSearchParams();
      if (params?.organizationId) search.set('organizationId', params.organizationId);
      if (params?.page) search.set('page', String(params.page));
      if (params?.limit) search.set('limit', String(params.limit));
      const query = search.toString();
      return http.get<PaginatedData<ProjectConfig>>(query ? `/projects?${query}` : '/projects');
    },

    get: (id: string): Promise<ProjectConfig> =>
      http.get<ProjectConfig>(`/projects/${encodeURIComponent(id)}`),

    create: (payload: CreateProjectRequest): Promise<ProjectConfig> =>
      http.post<ProjectConfig>('/projects', payload),

    delete: (id: string): Promise<void> => http.delete<void>(`/projects/${encodeURIComponent(id)}`),

    getVars: (id: string): Promise<Record<string, string>> =>
      http.get<Record<string, string>>(`/projects/${encodeURIComponent(id)}/env`),

    setVars: (id: string, vars: Record<string, string>): Promise<void> =>
      http.post<void>(`/projects/${encodeURIComponent(id)}/env`, vars),

    listDomains: (projectId: string): Promise<DomainConfig[]> =>
      http.get<DomainConfig[]>(`/projects/${encodeURIComponent(projectId)}/domains`),

    listMembers: (projectId: string): Promise<ProjectMember[]> =>
      http.get<ProjectMember[]>(`/projects/${encodeURIComponent(projectId)}/members`),

    addMember: (projectId: string, payload: AddMemberRequest): Promise<ProjectMember> =>
      http.post<ProjectMember>(`/projects/${encodeURIComponent(projectId)}/members`, payload),

    removeMember: (projectId: string, memberId: string): Promise<void> =>
      http.delete<void>(
        `/projects/${encodeURIComponent(projectId)}/members/${encodeURIComponent(memberId)}`
      ),

    listTokens: (projectId: string): Promise<ProjectToken[]> =>
      http.get<ProjectToken[]>(`/projects/${encodeURIComponent(projectId)}/tokens`),

    createToken: (
      projectId: string,
      payload: CreateTokenRequest
    ): Promise<{ token: string; projectToken: ProjectToken }> =>
      http.post<{ token: string; projectToken: ProjectToken }>(
        `/projects/${encodeURIComponent(projectId)}/tokens`,
        payload
      ),

    deleteToken: (projectId: string, tokenId: string): Promise<void> =>
      http.delete<void>(
        `/projects/${encodeURIComponent(projectId)}/tokens/${encodeURIComponent(tokenId)}`
      ),

    listWebhooks: (serviceId: string): Promise<ServiceWebhook[]> =>
      http.get<ServiceWebhook[]>(`/apps/${encodeURIComponent(serviceId)}/webhooks`),

    createWebhook: (serviceId: string, payload: CreateWebhookRequest): Promise<ServiceWebhook> =>
      http.post<ServiceWebhook>(`/apps/${encodeURIComponent(serviceId)}/webhooks`, payload),

    deleteWebhook: (serviceId: string, webhookId: string): Promise<void> =>
      http.delete<void>(
        `/apps/${encodeURIComponent(serviceId)}/webhooks/${encodeURIComponent(webhookId)}`
      ),
  };
}
