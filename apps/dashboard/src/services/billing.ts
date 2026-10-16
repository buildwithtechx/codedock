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

export interface PendingCheckout {
  checkoutId: string;
  idempotencyKey: string;
  purchaseIntent?: string;
  priceId: string;
  url: string;
  status: 'open' | 'expired' | 'complete' | 'canceled';
  canCancel: boolean;
  cancelUnavailableReason?: 'unsupported' | 'provider_unavailable';
  createdAt: string;
}

export const billingService = {
  getConfig: async (): Promise<BillingConfig> => {
    return await apiClient.get<BillingConfig>('/billing/config');
  },
  createCheckoutSession: async (payload: {
    priceId: string;
    successUrl: string;
    cancelUrl: string;
    idempotencyKey?: string;
    purchaseIntent?: string;
    cancelPending?: boolean;
  }): Promise<{ url: string }> => {
    return await apiClient.post<{ url: string }>('/billing/checkout', payload);
  },
  
  getPendingCheckout: async (): Promise<{ pending: PendingCheckout | null }> => {
    return await apiClient.get<{ pending: PendingCheckout | null }>('/billing/pending');
  },

  cancelCheckoutSession: async (
    checkoutId: string,
    idempotencyKey?: string,
  ): Promise<{ status: string; checkout: PendingCheckout }> => {
    return await apiClient.post<{ status: string; checkout: PendingCheckout }>(
      `/billing/checkouts/${encodeURIComponent(checkoutId)}/cancel`,
      { idempotencyKey },
    );
  },
};
