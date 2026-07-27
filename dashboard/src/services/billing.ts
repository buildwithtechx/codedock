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
    const res = await apiClient.get<any>('/billing/config');
    return res.data;
  },
  createCheckoutSession: async (payload: {
    priceId: string;
    successUrl: string;
    cancelUrl: string;
  }): Promise<any> => {
    const res = await apiClient.post<any>('/billing/checkout', payload);
    return res.data;
  },
};
