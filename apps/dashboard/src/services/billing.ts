import { DEFAULT_TOPUP_PACKS, isCloudPurchaseIntent } from '#/constants/billing';
import type { BillingConfig, PendingCheckout, TopupPack } from '#/interfaces/billing';
import { apiClient } from '#/lib/api-client';

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
      return (
        (res.packs as unknown as TopupPack[]) || (DEFAULT_TOPUP_PACKS as unknown as TopupPack[])
      );
    } catch {
      return DEFAULT_TOPUP_PACKS as unknown as TopupPack[];
    }
  },

  purchaseTopup: async (payload: {
    packId: string;
    idempotencyKey: string;
  }): Promise<{ url: string }> => {
    return await apiClient.post<{ url: string }>('/billing/topup', payload);
  },
};

export { isCloudPurchaseIntent };
export type { BillingConfig, PendingCheckout, TopupPack };
