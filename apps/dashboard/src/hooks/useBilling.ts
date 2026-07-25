import { useMutation, useQuery } from '@tanstack/react-query';
import { apiClient } from '#/lib/apiClient';

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

export function useBillingConfig() {
  return useQuery<BillingConfig>({
    queryKey: ['billing-config'],
    queryFn: async () => {
      const res = await apiClient.get<any>('/billing/config');
      return res.data;
    },
  });
}

export function useCreateCheckoutSession() {
  return useMutation({
    mutationFn: async (payload: { priceId: string; successUrl: string; cancelUrl: string }) => {
      const res = await apiClient.post<any>('/billing/checkout', payload);
      return res.data;
    },
  });
}
