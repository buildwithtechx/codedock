import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { BaseResponse } from '#/interfaces/base';
import type { Server, CreateServerRequest } from '#/interfaces/server';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

const serverService = {
  list: async (): Promise<Server[]> => {
    try {
      const res = await apiClient.get<BaseResponse<Server[]>>('/servers');
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  create: async (payload: CreateServerRequest): Promise<Server> => {
    try {
      const res = await apiClient.post<BaseResponse<Server>>('/servers', payload);
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
};

export const useListServers = () =>
  useQuery({ queryKey: ['servers'], queryFn: () => serverService.list() });

export const useCreateServer = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateServerRequest) => serverService.create(payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['servers'] }),
  });
};
