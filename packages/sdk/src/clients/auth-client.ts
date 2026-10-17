import type { HttpClient } from '../http';
import type { AccountInfo, PairingClaimResult, PairingHandshake, TenantCorsConfig } from '../types';

export function createAuthOperations(http: HttpClient) {
  return {
    getAccount: (): Promise<AccountInfo> => http.get<AccountInfo>('/auth/account'),
    createPairingCode: (): Promise<PairingHandshake> =>
      http.post<PairingHandshake>('/auth/pairing/code'),
    pair: (pairingCode: string): Promise<PairingClaimResult> =>
      http.post<PairingClaimResult>('/auth/pair', { code: pairingCode }),
    getCors: (tenantId?: string): Promise<TenantCorsConfig> =>
      http.get<TenantCorsConfig>(
        tenantId ? `/auth/cors?tenantId=${encodeURIComponent(tenantId)}` : '/auth/cors'
      ),
  };
}
