import type {
  ChatAppChallenge,
  ChatAppConnection,
  ChatAppIdentity,
  CreateChatAppInput,
  UpdateChatAppInput,
  VerifyPairingInput,
} from '@codedock/contracts';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const chatAppsService = {
  listConnections: async (): Promise<BaseResponse<ChatAppConnection[]>> => {
    try {
      return await apiClient.get<BaseResponse<ChatAppConnection[]>>('/chat-apps');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getConnection: async (id: string): Promise<BaseResponse<ChatAppConnection>> => {
    try {
      return await apiClient.get<BaseResponse<ChatAppConnection>>(`/chat-apps/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createConnection: async (input: CreateChatAppInput): Promise<BaseResponse<ChatAppConnection>> => {
    try {
      return await apiClient.post<BaseResponse<ChatAppConnection>>('/chat-apps', input);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateConnection: async (
    id: string,
    input: UpdateChatAppInput
  ): Promise<BaseResponse<ChatAppConnection>> => {
    try {
      return await apiClient.patch<BaseResponse<ChatAppConnection>>(`/chat-apps/${id}`, input);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteConnection: async (id: string): Promise<BaseResponse<null>> => {
    try {
      return await apiClient.delete<BaseResponse<null>>(`/chat-apps/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  testConnection: async (
    id: string
  ): Promise<BaseResponse<{ ok: boolean; status: string; checked_at: string }>> => {
    try {
      return await apiClient.post<
        BaseResponse<{ ok: boolean; status: string; checked_at: string }>
      >(`/chat-apps/${id}/test`, {});
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listIdentities: async (id: string): Promise<BaseResponse<ChatAppIdentity[]>> => {
    try {
      return await apiClient.get<BaseResponse<ChatAppIdentity[]>>(`/chat-apps/${id}/identities`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteIdentity: async (id: string): Promise<BaseResponse<null>> => {
    try {
      return await apiClient.delete<BaseResponse<null>>(`/chat-apps/identities/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getChallenge: async (token: string): Promise<BaseResponse<ChatAppChallenge>> => {
    try {
      return await apiClient.get<BaseResponse<ChatAppChallenge>>(
        `/chat-apps/pairing/challenge/${token}`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  verifyPairing: async (input: VerifyPairingInput): Promise<BaseResponse<ChatAppIdentity>> => {
    try {
      return await apiClient.post<BaseResponse<ChatAppIdentity>>(
        '/chat-apps/pairing/verify',
        input
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
