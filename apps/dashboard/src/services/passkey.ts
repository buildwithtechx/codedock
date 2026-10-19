import type { AuthResponse, LoginResponse } from '#/interfaces/auth';
import type { BaseResponse } from '#/interfaces/base';
import type {
  PasskeyAuthVerifyPayload,
  PasskeyListResponse,
  PasskeyRegisterVerifyPayload,
  PasskeyRegisterVerifyResponse,
  PublicKeyCredentialCreationOptionsJSON,
  PublicKeyCredentialRequestOptionsJSON,
} from '#/interfaces/passkey';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const passkeyService = {
  getRegisterOptions: async (): Promise<PublicKeyCredentialCreationOptionsJSON> => {
    try {
      return await apiClient.post<PublicKeyCredentialCreationOptionsJSON>(
        '/auth/passkeys/register/options'
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  verifyRegistration: async (
    payload: PasskeyRegisterVerifyPayload
  ): Promise<PasskeyRegisterVerifyResponse> => {
    try {
      return await apiClient.post<PasskeyRegisterVerifyResponse>(
        '/auth/passkeys/register/verify',
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getAuthenticateOptions: async (): Promise<PublicKeyCredentialRequestOptionsJSON> => {
    try {
      return await apiClient.post<PublicKeyCredentialRequestOptionsJSON>(
        '/auth/passkeys/authenticate/options'
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  verifyAuthentication: async (payload: PasskeyAuthVerifyPayload): Promise<AuthResponse> => {
    try {
      const response = await apiClient.post<LoginResponse>(
        '/auth/passkeys/authenticate/verify',
        payload
      );
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listPasskeys: async (): Promise<PasskeyListResponse> => {
    try {
      return await apiClient.get<PasskeyListResponse>('/auth/passkeys');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  renamePasskey: async (id: string, name: string): Promise<BaseResponse<void>> => {
    try {
      return await apiClient.patch<BaseResponse<void>>(`/auth/passkeys/${id}`, { name });
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deletePasskey: async (id: string): Promise<BaseResponse<void>> => {
    try {
      return await apiClient.delete<BaseResponse<void>>(`/auth/passkeys/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
