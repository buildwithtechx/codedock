import { CodedockClientOptions, HttpClient } from './client';
import type {
  ActionRun,
  ActionRunner,
  ActionWorkflow,
  ActiveSubscriptionOverview,
  ClusterScaling,
  ComputeCluster,
  IngressConnection,
  IngressDiscovery,
  IngressHostname,
  NetworkWatchdog,
  PrepaidTopupPackage,
  ServerNetwork,
} from './types';

export * from './client';
export * from './types';
export * from 'zod';

export class CodedockClient {
  private http: HttpClient;

  constructor(options: CodedockClientOptions = {}) {
    this.http = new HttpClient(options);
  }

  readonly ingress = {
    list: (): Promise<IngressConnection[]> => this.http.get<IngressConnection[]>('/dns/ingress'),
    discover: (credentialId: string, accountId?: string): Promise<IngressDiscovery> =>
      this.http.post<IngressDiscovery>('/dns/ingress/discover', { credentialId, accountId }),
    connect: (input: { serverId: string; name: string; connector: string }): Promise<IngressConnection> =>
      this.http.post<IngressConnection>('/dns/ingress', input),
    update: (id: string, input: { defaultForDomains: boolean }): Promise<IngressConnection> =>
      this.http.put<IngressConnection>(`/dns/ingress/${encodeURIComponent(id)}`, input),
    retry: (id: string): Promise<IngressConnection> =>
      this.http.post<IngressConnection>(`/dns/ingress/${encodeURIComponent(id)}/retry`),
    disconnect: (id: string): Promise<{ success: boolean }> =>
      this.http.delete<{ success: boolean }>(`/dns/ingress/${encodeURIComponent(id)}`),
    getDomain: (domainId: string): Promise<IngressHostname> =>
      this.http.get<IngressHostname>(`/dns/ingress/domains/${encodeURIComponent(domainId)}`),
    setDomain: (domainId: string, input: { connectionId?: string }): Promise<IngressHostname> =>
      this.http.put<IngressHostname>(`/dns/ingress/domains/${encodeURIComponent(domainId)}`, input),
  };

  readonly networks = {
    list: (): Promise<ServerNetwork[]> => this.http.get<ServerNetwork[]>('/system/networks'),
    get: (id: string): Promise<ServerNetwork> =>
      this.http.get<ServerNetwork>(`/system/networks/${encodeURIComponent(id)}`),
    create: (input: { name: string; mode?: string }): Promise<ServerNetwork> =>
      this.http.post<ServerNetwork>('/system/networks', input),
    reviseAccess: (id: string, accessPolicy: { mode: string; allowedPairs?: [string, string][] }): Promise<ServerNetwork> =>
      this.http.put<ServerNetwork>(`/system/networks/${encodeURIComponent(id)}/access`, { accessPolicy }),
    delete: (id: string): Promise<{ success: boolean }> =>
      this.http.delete<{ success: boolean }>(`/system/networks/${encodeURIComponent(id)}`),
    watchdog: (operationId: string): Promise<NetworkWatchdog> =>
      this.http.get<NetworkWatchdog>(`/system/networks/operations/${encodeURIComponent(operationId)}/watchdog`),
    commitWatchdog: (operationId: string): Promise<NetworkWatchdog> =>
      this.http.post<NetworkWatchdog>(`/system/networks/operations/${encodeURIComponent(operationId)}/watchdog/commit`),
    rollbackWatchdog: (operationId: string): Promise<NetworkWatchdog> =>
      this.http.post<NetworkWatchdog>(`/system/networks/operations/${encodeURIComponent(operationId)}/watchdog/rollback`),
  };

  readonly clusters = {
    list: (projectId: string): Promise<ComputeCluster[]> =>
      this.http.get<ComputeCluster[]>(`/projects/${encodeURIComponent(projectId)}/clusters`),
    getScaling: (clusterId: string): Promise<ClusterScaling> =>
      this.http.get<ClusterScaling>(`/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling`),
    setupScaling: (clusterId: string): Promise<ClusterScaling> =>
      this.http.post<ClusterScaling>(`/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling/setup`),
    retryScaling: (clusterId: string): Promise<ClusterScaling> =>
      this.http.post<ClusterScaling>(`/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling/retry`),
    removeScaling: (clusterId: string): Promise<ClusterScaling> =>
      this.http.delete<ClusterScaling>(`/system/compute-clusters/${encodeURIComponent(clusterId)}/scaling`),
  };

  readonly actions = {
    listWorkflows: (): Promise<ActionWorkflow[]> => this.http.get<ActionWorkflow[]>('/actions/workflows'),
    createWorkflow: (input: { name: string; description?: string; trigger?: string; yaml?: string }): Promise<ActionWorkflow> =>
      this.http.post<ActionWorkflow>('/actions/workflows', input),
    listRuns: (): Promise<ActionRun[]> => this.http.get<ActionRun[]>('/actions/runs'),
    getRun: (id: string): Promise<ActionRun> =>
      this.http.get<ActionRun>(`/actions/runs/${encodeURIComponent(id)}`),
    retryRun: (id: string): Promise<ActionRun> =>
      this.http.post<ActionRun>(`/actions/runs/${encodeURIComponent(id)}/retry`),
    cancelRun: (id: string): Promise<ActionRun> =>
      this.http.post<ActionRun>(`/actions/runs/${encodeURIComponent(id)}/cancel`),
    listRunners: (): Promise<ActionRunner[]> => this.http.get<ActionRunner[]>('/actions/runners'),
    createRunner: (input: { name: string; mode?: string; image?: string; labels?: string[]; enabled: boolean }): Promise<ActionRunner> =>
      this.http.post<ActionRunner>('/actions/runners', input),
  };

  readonly billing = {
    getActiveSubscription: (): Promise<ActiveSubscriptionOverview> =>
      this.http.get<ActiveSubscriptionOverview>('/billing/active-subscription'),
    getTopupPackages: (): Promise<PrepaidTopupPackage[]> =>
      this.http.get<PrepaidTopupPackage[]>('/billing/prepaid/packages'),
    topup: (packageId: string): Promise<{ checkoutUrl: string }> =>
      this.http.post<{ checkoutUrl: string }>('/billing/prepaid/topup', { packageId }),
  };
}

export function createCodedockClient(options?: CodedockClientOptions): CodedockClient {
  return new CodedockClient(options);
}
