import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { BillingConfig, PendingCheckout, TopupPack } from '#/interfaces/billing';
import { billingService } from '#/services/billing';

export function useBillingConfig() {
  return useQuery<BillingConfig>({
    queryKey: ['billing', 'config'],
    queryFn: () => billingService.getConfig(),
  });
}

export function useCreateCheckoutSession() {
  return useMutation({
    mutationFn: (payload: {
      priceId: string;
      successUrl: string;
      cancelUrl: string;
      idempotencyKey?: string;
      purchaseIntent?: string;
      cancelPending?: boolean;
    }) => billingService.createCheckoutSession(payload),
  });
}

export function usePendingCheckout() {
  return useQuery<{ pending: PendingCheckout | null }>({
    queryKey: ['billing', 'pending'],
    queryFn: () => billingService.getPendingCheckout(),
  });
}

export function useCancelCheckoutSession() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ checkoutId, idempotencyKey }: { checkoutId: string; idempotencyKey?: string }) =>
      billingService.cancelCheckoutSession(checkoutId, idempotencyKey),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['billing', 'pending'] });
    },
  });
}

export function useTopupPacks() {
  return useQuery<TopupPack[]>({
    queryKey: ['billing', 'topups'],
    queryFn: () => billingService.getTopupPacks(),
  });
}

export function usePurchaseTopup() {
  return useMutation({
    mutationFn: (payload: { packId: string; idempotencyKey: string }) =>
      billingService.purchaseTopup(payload),
  });
}

export function useCreatePortalSession() {
  return useMutation({
    mutationFn: () => billingService.createPortalSession(),
  });
}
