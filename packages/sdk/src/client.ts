import { createActionsOperations } from './clients/actions-client';
import { createAuthOperations } from './clients/auth-client';
import { createBillingOperations } from './clients/billing-client';
import { createClustersOperations } from './clients/clusters-client';
import { createIngressOperations } from './clients/ingress-client';
import { createNetworksOperations } from './clients/networks-client';
import { type CodedockClientOptions, HttpClient } from './http';

export class CodedockClient {
  private http: HttpClient;
  private options: CodedockClientOptions;

  readonly auth: ReturnType<typeof createAuthOperations>;
  readonly ingress: ReturnType<typeof createIngressOperations>;
  readonly networks: ReturnType<typeof createNetworksOperations>;
  readonly clusters: ReturnType<typeof createClustersOperations>;
  readonly actions: ReturnType<typeof createActionsOperations>;
  readonly billing: ReturnType<typeof createBillingOperations>;

  constructor(options: CodedockClientOptions = {}) {
    this.options = options;
    this.http = new HttpClient(options);
    this.auth = createAuthOperations(this.http);
    this.ingress = createIngressOperations(this.http);
    this.networks = createNetworksOperations(this.http);
    this.clusters = createClustersOperations(this.http);
    this.actions = createActionsOperations(this.http);
    this.billing = createBillingOperations(this.http);
  }

  scope(tenantId: string): CodedockClient {
    return new CodedockClient({ ...this.options, tenantId });
  }
}

export function createCodedockClient(options?: CodedockClientOptions): CodedockClient {
  return new CodedockClient(options);
}
