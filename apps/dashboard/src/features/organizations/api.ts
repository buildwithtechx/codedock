import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';
import type {
  CreateOrganizationRequest,
  InviteOrganizationMemberRequest,
  Organization,
  OrganizationMember,
  UpdateOrganizationMemberRequest,
} from './interfaces';

export const orgService = {
  list: async (): Promise<Organization[]> => {
    try {
      return await apiClient.get<Organization[]>('/organizations');
    } catch (err) {
      throw handleApiError(err);
    }
  },
  create: async (payload: CreateOrganizationRequest): Promise<Organization> => {
    try {
      return await apiClient.post<Organization>('/organizations', payload);
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
      return await apiClient.get<Organization>(`/organizations/${id}`);
    } catch (err) {
      throw handleApiError(err);
    }
  },
  listMembers: async (id: string): Promise<OrganizationMember[]> => {
    try {
      const res = await apiClient.get<OrganizationMember[] | { data: OrganizationMember[] }>(
        `/organizations/${id}/members`
      );
      if (Array.isArray(res)) return res;
      if (res && 'data' in res && Array.isArray(res.data)) return res.data;
      return [];
    } catch (err) {
      throw handleApiError(err);
    }
  },
  inviteMember: async (
    id: string,
    payload: InviteOrganizationMemberRequest
  ): Promise<OrganizationMember> => {
    try {
      return await apiClient.post<OrganizationMember>(`/organizations/${id}/members`, payload);
    } catch (err) {
      throw handleApiError(err);
    }
  },
  updateMember: async (
    orgId: string,
    memberId: string,
    payload: UpdateOrganizationMemberRequest,
    userId?: string
  ): Promise<OrganizationMember> => {
    try {
      const targetId = userId || memberId;
      return await apiClient.put<OrganizationMember>(
        `/organizations/${orgId}/members/${targetId}`,
        payload
      );
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
