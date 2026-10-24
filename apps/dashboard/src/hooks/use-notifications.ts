import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { NotificationSettingsDto } from '#/interfaces/notifications';
import { notificationsService } from '#/services/notifications';

export const useGetNotificationSettings = () =>
  useQuery({
    queryKey: ['settings', 'getNotificationSettings'],
    queryFn: () => notificationsService.getNotificationSettings(),
  });

export const useUpdateNotificationSettings = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Partial<NotificationSettingsDto>) =>
      notificationsService.updateNotificationSettings(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ['settings', 'getNotificationSettings'],
      });
    },
  });
};

export const useTestNotification = () =>
  useMutation({
    mutationFn: (payload: { provider: string }) => notificationsService.testNotification(payload),
  });

export const useGetNotificationSubscriptions = (userId?: string) =>
  useQuery({
    queryKey: ['notification-subscriptions', userId].filter(Boolean),
    queryFn: () => notificationsService.getNotificationSubscriptions(userId),
  });

export const useUpsertNotificationSubscription = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      category: string;
      channels: string[];
      enabled: boolean;
      orgId?: string;
    }) => notificationsService.upsertNotificationSubscription(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['notification-subscriptions'] });
    },
  });
};

export const useGetNotificationDefaults = () =>
  useQuery({
    queryKey: ['notification-defaults'],
    queryFn: () => notificationsService.getNotificationDefaults(),
  });

export const useUpsertNotificationDefault = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      category: string;
      channels: string[];
      enabled: boolean;
      orgId?: string;
    }) => notificationsService.upsertNotificationDefault(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['notification-defaults'] });
    },
  });
};

export const useGetNotificationDeliveries = (limit = 50) =>
  useQuery({
    queryKey: ['notification-deliveries', limit],
    queryFn: () => notificationsService.getNotificationDeliveries(limit),
    refetchInterval: 15_000,
  });
