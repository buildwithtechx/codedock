import type { HttpClient } from '../http';
import type { NetworkWatchdog, ServerNetwork } from '../types';

export function createNetworksOperations(http: HttpClient) {
  return {
    list: (): Promise<ServerNetwork[]> => http.get<ServerNetwork[]>('/system/networks'),
    get: (id: string): Promise<ServerNetwork> =>
      http.get<ServerNetwork>(`/system/networks/${encodeURIComponent(id)}`),
    create: (input: { name: string; mode?: string }): Promise<ServerNetwork> =>
      http.post<ServerNetwork>('/system/networks', input),
    reviseAccess: (
      id: string,
      accessPolicy: { mode: string; allowedPairs?: [string, string][] }
    ): Promise<ServerNetwork> =>
      http.put<ServerNetwork>(`/system/networks/${encodeURIComponent(id)}/access`, {
        accessPolicy,
      }),
    delete: (id: string): Promise<{ success: boolean }> =>
      http.delete<{ success: boolean }>(`/system/networks/${encodeURIComponent(id)}`),
    watchdog: (operationId: string): Promise<NetworkWatchdog> =>
      http.get<NetworkWatchdog>(
        `/system/networks/operations/${encodeURIComponent(operationId)}/watchdog`
      ),
    commitWatchdog: (operationId: string): Promise<NetworkWatchdog> =>
      http.post<NetworkWatchdog>(
        `/system/networks/operations/${encodeURIComponent(operationId)}/watchdog/commit`
      ),
    rollbackWatchdog: (operationId: string): Promise<NetworkWatchdog> =>
      http.post<NetworkWatchdog>(
        `/system/networks/operations/${encodeURIComponent(operationId)}/watchdog/rollback`
      ),
  };
}
