import type { BaseResponse } from '#/interfaces/base';
import type {
  CreateOrganizationRequest,
  InviteOrganizationMemberRequest,
  Organization,
  OrganizationMember,
  UpdateOrganizationMemberRequest,
} from '#/interfaces/organization';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const orgService = {
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
  get: async (id: string): Promise<Organization> => {
    try {
      const res = await apiClient.get<BaseResponse<Organization>>(`/organizations/${id}`);
      return res.data as Organization;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  listMembers: async (id: string): Promise<OrganizationMember[]> => {
    try {
      const res = await apiClient.get<BaseResponse<OrganizationMember[]>>(
        `/organizations/${id}/members`
      );
      return res.data || [];
    } catch (err) {
      throw handleApiError(err);
    }
  },
  inviteMember: async (
    id: string,
    payload: InviteOrganizationMemberRequest
  ): Promise<OrganizationMember> => {
    try {
      const res = await apiClient.post<BaseResponse<OrganizationMember>>(
        `/organizations/${id}/members`,
        payload
      );
      return res.data as OrganizationMember;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  updateMember: async (
    orgId: string,
    memberId: string,
    payload: UpdateOrganizationMemberRequest
  ): Promise<OrganizationMember> => {
    try {
      const res = await apiClient.put<BaseResponse<OrganizationMember>>(
        `/organizations/${orgId}/members/${memberId}`,
        payload
      );
      return res.data as OrganizationMember;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  removeMember: async (orgId: string, memberId: string): Promise<void> => {
    try {
      await apiClient.delete(`/organizations/${orgId}/members/${memberId}`);
    } catch (err) {
      throw handleApiError(err);
    }
  },
};
