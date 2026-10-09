import type { CreateJobRequest, Job, UpdateJobRequest } from '#/features/services';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const jobsApi = {
  listByProject: async (projectId: string): Promise<BaseResponse<Job[]>> => {
    try {
      return await apiClient.get<BaseResponse<Job[]>>(
        `/scheduled-tasks?projectId=${encodeURIComponent(projectId)}`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  get: async (id: string): Promise<BaseResponse<Job>> => {
    try {
      return await apiClient.get<BaseResponse<Job>>(`/scheduled-tasks/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  create: async (payload: CreateJobRequest): Promise<BaseResponse<Job>> => {
    try {
      return await apiClient.post<BaseResponse<Job>>('/scheduled-tasks', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  update: async (id: string, payload: UpdateJobRequest): Promise<BaseResponse<Job>> => {
    try {
      return await apiClient.put<BaseResponse<Job>>(`/scheduled-tasks/${id}`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  remove: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/scheduled-tasks/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  trigger: async (id: string): Promise<void> => {
    try {
      await apiClient.post(`/scheduled-tasks/${id}/trigger`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
