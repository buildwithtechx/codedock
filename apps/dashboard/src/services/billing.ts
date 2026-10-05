import { apiClient } from '#/lib/api-client';

export interface BillingConfig {
  publishableKey: string;
  plans: {
    id: string;
    name: string;
    price: number;
    priceId?: string;
    features: string[];
  }[];
}

export const billingService = {
  getConfig: async (): Promise<BillingConfig> => {
    return await apiClient.get<BillingConfig>('/billing/config');
  },
  createCheckoutSession: async (payload: {
    priceId: string;
    successUrl: string;
    cancelUrl: string;
  }): Promise<{ url: string }> => {
    return await apiClient.post<{ url: string }>('/billing/checkout', payload);
  },
};
