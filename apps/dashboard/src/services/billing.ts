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
    idempotencyKey?: string
  ): Promise<{ status: string; checkout: PendingCheckout }> => {
    return await apiClient.post<{ status: string; checkout: PendingCheckout }>(
      `/billing/checkouts/${encodeURIComponent(checkoutId)}/cancel`,
      { idempotencyKey }
    );
  },

  createPortalSession: async (): Promise<{ url: string }> => {
    return await apiClient.post<{ url: string }>('/billing/portal');
  },

  getTopupPacks: async (): Promise<TopupPack[]> => {
    try {
      const res = await apiClient.get<{ packs: TopupPack[] }>('/billing/topups');
      return res.packs || DEFAULT_TOPUP_PACKS;
    } catch {
      return DEFAULT_TOPUP_PACKS;
    }
  },

  purchaseTopup: async (payload: {
    packId: string;
    idempotencyKey: string;
  }): Promise<{ url: string }> => {
    return await apiClient.post<{ url: string }>('/billing/topup', payload);
  },
};

export interface TopupPack {
  id: string;
  name: string;
  credits: number;
  priceCents: number;
  bonusPercent?: number;
}

export const DEFAULT_TOPUP_PACKS: TopupPack[] = [
  { id: 'topup_10', name: '$10 Starter Pack', credits: 10000, priceCents: 1000 },
  { id: 'topup_25', name: '$25 Growth Pack', credits: 25000, priceCents: 2500 },
  { id: 'topup_50', name: '$50 Scale Pack', credits: 55000, priceCents: 5000, bonusPercent: 10 },
  { id: 'topup_100', name: '$100 Pro Pack', credits: 115000, priceCents: 10000, bonusPercent: 15 },
  { id: 'topup_250', name: '$250 Max Pack', credits: 300000, priceCents: 25000, bonusPercent: 20 },
];
