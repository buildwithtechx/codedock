import type { HttpClient } from '../http';
import type { IngressConnection, IngressDiscovery, IngressHostname } from '../types';

export function createIngressOperations(http: HttpClient) {
  return {
    list: (): Promise<IngressConnection[]> => http.get<IngressConnection[]>('/dns/ingress'),
    discover: (credentialId: string, accountId?: string): Promise<IngressDiscovery> =>
      http.post<IngressDiscovery>('/dns/ingress/discover', { credentialId, accountId }),
    connect: (input: {
      serverId: string;
      name: string;
      connector: string;
    }): Promise<IngressConnection> => http.post<IngressConnection>('/dns/ingress', input),
    update: (id: string, input: { defaultForDomains: boolean }): Promise<IngressConnection> =>
      http.put<IngressConnection>(`/dns/ingress/${encodeURIComponent(id)}`, input),
    retry: (id: string): Promise<IngressConnection> =>
      http.post<IngressConnection>(`/dns/ingress/${encodeURIComponent(id)}/retry`),
    disconnect: (id: string): Promise<{ success: boolean }> =>
      http.delete<{ success: boolean }>(`/dns/ingress/${encodeURIComponent(id)}`),
    getDomain: (domainId: string): Promise<IngressHostname> =>
      http.get<IngressHostname>(`/dns/ingress/domains/${encodeURIComponent(domainId)}`),
    setDomain: (domainId: string, input: { connectionId?: string }): Promise<IngressHostname> =>
      http.put<IngressHostname>(`/dns/ingress/domains/${encodeURIComponent(domainId)}`, input),
  };
}
