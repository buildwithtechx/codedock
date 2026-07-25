import { useMutation, useQuery } from '@tanstack/react-query';
import { api } from '#/lib/api';

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
      const res = await api.get('/billing/config');
      return res.data;
    },
  });
}

export function useCreateCheckoutSession() {
  return useMutation({
    mutationFn: async (payload: { priceId: string; successUrl: string; cancelUrl: string }) => {
      const res = await api.post<{ url: string }>('/billing/checkout', payload);
      return res.data;
    },
  });
}
