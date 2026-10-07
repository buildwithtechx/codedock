import type {
  CreateManagedCredentialRequest,
  ManagedCatalog,
  ManagedCredential,
  ManagedQuota,
  ManagedReview,
  ReviewManagedProvisionRequest,
} from '#/features/servers/managed-types';
import type { BaseResponse } from '#/interfaces/base';
import type { Server } from '#/interfaces/server';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const managedService = {
  catalog: async (): Promise<ManagedCatalog> => {
    try {
      const res = await apiClient.get<BaseResponse<ManagedCatalog>>('/managed/catalog');
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  listCredentials: async (organizationId: string): Promise<ManagedCredential[]> => {
    try {
      const res = await apiClient.get<BaseResponse<ManagedCredential[]>>(
        `/organizations/${organizationId}/managed/credentials`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  createCredential: async (
    organizationId: string,
    payload: CreateManagedCredentialRequest
  ): Promise<ManagedCredential> => {
    try {
      const res = await apiClient.post<BaseResponse<ManagedCredential>>(
        `/organizations/${organizationId}/managed/credentials`,
        payload
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  deleteCredential: async (organizationId: string, credentialId: string): Promise<void> => {
    try {
      await apiClient.delete<BaseResponse<null>>(
        `/organizations/${organizationId}/managed/credentials/${credentialId}`
      );
    } catch (err) {
      throw handleApiError(err);
    }
  },
  getQuota: async (organizationId: string): Promise<ManagedQuota> => {
    try {
      const res = await apiClient.get<BaseResponse<ManagedQuota>>(
        `/organizations/${organizationId}/managed/quota`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  setQuota: async (
    organizationId: string,
    payload: { maxServers: number; maxMemoryGB: number }
  ): Promise<ManagedQuota> => {
    try {
      const res = await apiClient.put<BaseResponse<ManagedQuota>>(
        `/organizations/${organizationId}/managed/quota`,
        payload
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  reviewProvision: async (
    organizationId: string,
    payload: ReviewManagedProvisionRequest
  ): Promise<ManagedReview> => {
    try {
      const res = await apiClient.post<BaseResponse<ManagedReview>>(
        `/organizations/${organizationId}/managed/servers/review`,
        payload
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  reviewResize: async (
    organizationId: string,
    serverId: string,
    serverType: string
  ): Promise<ManagedReview> => {
    try {
      const res = await apiClient.post<BaseResponse<ManagedReview>>(
        `/organizations/${organizationId}/managed/servers/${serverId}/resize/review`,
        { serverType }
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  reviewDelete: async (organizationId: string, serverId: string): Promise<ManagedReview> => {
    try {
      const res = await apiClient.post<BaseResponse<ManagedReview>>(
        `/organizations/${organizationId}/managed/servers/${serverId}/delete/review`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
  applyOperation: async (operationId: string, confirmation: string): Promise<void> => {
    try {
      await apiClient.post<BaseResponse<null>>(`/managed-operations/${operationId}/apply`, {
        confirmation,
      });
    } catch (err) {
      throw handleApiError(err);
    }
  },
  refreshServer: async (organizationId: string, serverId: string): Promise<Server> => {
    try {
      const res = await apiClient.post<BaseResponse<Server>>(
        `/organizations/${organizationId}/managed/servers/${serverId}/refresh`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
};
