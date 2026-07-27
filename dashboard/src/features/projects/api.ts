import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';
import type {
  AddMemberRequest,
  CreateProjectRequest,
  CreateProjectResponse,
  CreateTokenRequest,
  CreateWebhookRequest,
  GetProjectResponse,
  ListProjectsResponse,
  ProjectMember,
  ProjectToken,
} from './interfaces';

export interface ServiceWebhook {
  id: string;
  serviceId: string;
  url: string;
  eventTypes: string[];
  includePrEnvironments: boolean;
  createdAt: string;
  updatedAt: string;
}

export const projectsService = {
  listProjects: async (organizationId?: string): Promise<ListProjectsResponse> => {
    try {
      const url = organizationId
        ? `/projects?organizationId=${encodeURIComponent(organizationId)}`
        : '/projects';
      return await apiClient.get<ListProjectsResponse>(url);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getProject: async (id: string): Promise<GetProjectResponse> => {
    try {
      return await apiClient.get<GetProjectResponse>(`/projects/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createProject: async (payload: CreateProjectRequest): Promise<CreateProjectResponse> => {
    try {
      return await apiClient.post<CreateProjectResponse>('/projects', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteProject: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/projects/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getVars: async (id: string): Promise<BaseResponse<Record<string, string>>> => {
    try {
      return await apiClient.get<BaseResponse<Record<string, string>>>(`/projects/${id}/env`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  setVars: async (
    id: string,
    payload: Record<string, string> | { variables: Record<string, string> }
  ): Promise<BaseResponse<void>> => {
    try {
      const body =
        'variables' in payload && typeof payload.variables === 'object' && payload.variables
          ? payload.variables
          : payload;
      return await apiClient.put<BaseResponse<void>>(`/projects/${id}/env`, body);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};

export const projectSettingsService = {
  listWebhooks: async (serviceId: string): Promise<BaseResponse<ServiceWebhook[]>> => {
    try {
      return await apiClient.get<BaseResponse<ServiceWebhook[]>>(`/apps/${serviceId}/webhooks`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createWebhook: async (
    serviceId: string,
    payload: CreateWebhookRequest
  ): Promise<BaseResponse<ServiceWebhook>> => {
    try {
      return await apiClient.post<BaseResponse<ServiceWebhook>>(
        `/apps/${serviceId}/webhooks`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteWebhook: async (serviceId: string, id: string): Promise<void> => {
    try {
      await apiClient.delete(`/apps/${serviceId}/webhooks/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listTokens: async (projectId: string): Promise<BaseResponse<ProjectToken[]>> => {
    try {
      return await apiClient.get<BaseResponse<ProjectToken[]>>(`/projects/${projectId}/tokens`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createToken: async (
    projectId: string,
    payload: CreateTokenRequest
  ): Promise<BaseResponse<{ token: string; projectToken: ProjectToken }>> => {
    try {
      return await apiClient.post<BaseResponse<{ token: string; projectToken: ProjectToken }>>(
        `/projects/${projectId}/tokens`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteToken: async (projectId: string, id: string): Promise<void> => {
    try {
      await apiClient.delete(`/projects/${projectId}/tokens/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listMembers: async (projectId: string): Promise<BaseResponse<ProjectMember[]>> => {
    try {
      return await apiClient.get<BaseResponse<ProjectMember[]>>(`/projects/${projectId}/members`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  addMember: async (
    projectId: string,
    payload: AddMemberRequest
  ): Promise<BaseResponse<ProjectMember>> => {
    try {
      return await apiClient.post<BaseResponse<ProjectMember>>(
        `/projects/${projectId}/members`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  removeMember: async (projectId: string, memberId: string): Promise<void> => {
    try {
      await apiClient.delete(`/projects/${projectId}/members/${memberId}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
