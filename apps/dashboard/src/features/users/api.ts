import type { BaseResponse, PaginatedData } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';
import type { User } from './interfaces';

export const usersService = {
  list: async (): Promise<BaseResponse<PaginatedData<User>>> => {
    try {
      return await apiClient.get<BaseResponse<PaginatedData<User>>>('/users');
    } catch (err) {
      throw handleApiError(err);
    }
  },
  delete: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/users/${id}`);
    } catch (err) {
      throw handleApiError(err);
    }
  },
  invite: async (payload: { email: string; role: string }): Promise<User> => {
    try {
      const res = await apiClient.post<BaseResponse<User>>('/users/invite', payload);
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
};
