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
      const response = await apiClient.post<
        | PublicKeyCredentialCreationOptionsJSON
        | { publicKey: PublicKeyCredentialCreationOptionsJSON }
      >('/auth/passkeys/register/options');
      if (response && 'publicKey' in response && response.publicKey) {
        return response.publicKey;
      }
      return response as PublicKeyCredentialCreationOptionsJSON;
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
      const response = await apiClient.post<
        PublicKeyCredentialRequestOptionsJSON | { publicKey: PublicKeyCredentialRequestOptionsJSON }
      >('/auth/passkeys/authenticate/options');
      if (response && 'publicKey' in response && response.publicKey) {
        return response.publicKey;
      }
      return response as PublicKeyCredentialRequestOptionsJSON;
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
      return await apiClient.put<BaseResponse<void>>(`/auth/passkeys/${id}`, { name });
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
