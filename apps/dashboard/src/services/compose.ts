import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';

export interface ComposeStackItem {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  content: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface ComposeReviewPayload {
  name: string;
  content: string;
  variables?: string;
  rootDirectory?: string;
  environmentId?: string;
}

export const composeService = {
  list: async (projectId: string): Promise<ComposeStackItem[]> => {
    const res = await apiClient.get<BaseResponse<ComposeStackItem[]>>(
      `/projects/${encodeURIComponent(projectId)}/stacks`
    );
    return res.data ?? [];
  },

  review: async (projectId: string, payload: ComposeReviewPayload): Promise<any> => {
    const res = await apiClient.post<BaseResponse<any>>(
      `/projects/${encodeURIComponent(projectId)}/stacks/review`,
      payload
    );
    return res.data;
  },

  save: async (projectId: string, payload: any): Promise<ComposeStackItem> => {
    const res = await apiClient.post<BaseResponse<ComposeStackItem>>(
      `/projects/${encodeURIComponent(projectId)}/stacks`,
      payload
    );
    return res.data;
  },

  deploy: async (projectId: string, stackId: string): Promise<any> => {
    const res = await apiClient.post<BaseResponse<any>>(
      `/projects/${encodeURIComponent(projectId)}/stacks/${encodeURIComponent(stackId)}/deploy`
    );
    return res.data;
  },

  action: async (
    projectId: string,
    stackId: string,
    action: 'restart' | 'stop' | 'cancel'
  ): Promise<any> => {
    const res = await apiClient.post<BaseResponse<any>>(
      `/projects/${encodeURIComponent(projectId)}/stacks/${encodeURIComponent(stackId)}/${encodeURIComponent(action)}`
    );
    return res.data;
  },

  config: async (projectId: string, stackId: string): Promise<string> => {
    const res = await apiClient.get<BaseResponse<{ config: string }>>(
      `/projects/${encodeURIComponent(projectId)}/stacks/${encodeURIComponent(stackId)}/config`
    );
    return res.data?.config ?? '';
  },
};
