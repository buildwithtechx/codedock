import { useMutation, useQuery } from '@tanstack/react-query';
import { type BillingConfig, billingService } from '#/services/billing';

export type { BillingConfig };

export function useBillingConfig() {
  return useQuery<BillingConfig>({
    queryKey: ['billing-config'],
    queryFn: () => billingService.getConfig(),
  });
}

export function useCreateCheckoutSession() {
  return useMutation({
    mutationFn: (payload: { priceId: string; successUrl: string; cancelUrl: string }) =>
      billingService.createCheckoutSession(payload),
  });
}
