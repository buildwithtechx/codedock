import type { BaseResponse } from '#/interfaces/base';
import type {
  AuthenticationResponseJSON,
  PublicKeyCredentialCreationOptionsJSON,
  PublicKeyCredentialRequestOptionsJSON,
  RegistrationResponseJSON,
} from '@simplewebauthn/browser';

export interface Passkey {
  id: string;
  name: string;
  deviceType: 'singleDevice' | 'multiDevice' | string;
  backedUp: boolean;
  transports: string[];
  aaguid?: string;
  createdAt: string;
  lastUsedAt?: string;
}

export type PasskeyListResponse = BaseResponse<Passkey[]>;
export type PasskeyRegisterVerifyResponse = BaseResponse<Passkey>;

export interface PasskeyRegisterVerifyPayload {
  name: string;
  response: RegistrationResponseJSON;
}

export interface PasskeyAuthVerifyPayload {
  totpCode?: string;
  response: AuthenticationResponseJSON;
}

export type {
  PublicKeyCredentialCreationOptionsJSON,
  PublicKeyCredentialRequestOptionsJSON,
  RegistrationResponseJSON,
  AuthenticationResponseJSON,
};
