import type { HttpClient } from '../http';
import type { ActiveSubscriptionOverview, PrepaidTopupPackage } from '../types';

export function createBillingOperations(http: HttpClient) {
  return {
    getActiveSubscription: (): Promise<ActiveSubscriptionOverview> =>
      http.get<ActiveSubscriptionOverview>('/billing/active-subscription'),
    getTopupPackages: (): Promise<PrepaidTopupPackage[]> =>
      http.get<PrepaidTopupPackage[]>('/billing/prepaid/packages'),
    topup: (packageId: string): Promise<{ checkoutUrl: string }> =>
      http.post<{ checkoutUrl: string }>('/billing/prepaid/topup', { packageId }),
  };
}
