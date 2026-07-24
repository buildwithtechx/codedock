import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { BaseResponse } from '#/interfaces/base';
import type { CreateOrganizationRequest, Organization } from '#/interfaces/organization';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

const orgService = {
  list: async (): Promise<Organization[]> => {
    try {
      const res = await apiClient.get<BaseResponse<Organization[]>>('/organizations');
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  create: async (payload: CreateOrganizationRequest): Promise<Organization> => {
    try {
      const res = await apiClient.post<BaseResponse<Organization>>('/organizations', payload);
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  delete: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/organizations/${id}`);
    } catch (err) {
      throw handleApiError(err);
    }
  },
};

export const useListOrganizations = () =>
  useQuery({ queryKey: ['organizations'], queryFn: () => orgService.list() });

export const useCreateOrganization = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateOrganizationRequest) => orgService.create(payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['organizations'] }),
  });
};

export const useDeleteOrganization = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => orgService.delete(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['organizations'] }),
  });
};
