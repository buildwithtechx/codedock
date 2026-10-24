import type { BaseResponse } from '#/interfaces/base';
import type {
  DeviceTokenDto,
  NotificationDefaultDto,
  NotificationDeliveryDto,
  NotificationSettingsDto,
  NotificationSubscriptionDto,
} from '#/interfaces/notifications';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const notificationsService = {
  getNotificationSettings: async (): Promise<BaseResponse<NotificationSettingsDto>> => {
    try {
      return await apiClient.get<BaseResponse<NotificationSettingsDto>>('/notifications');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateNotificationSettings: async (
    payload: Partial<NotificationSettingsDto>
  ): Promise<BaseResponse<NotificationSettingsDto>> => {
    try {
      return await apiClient.put<BaseResponse<NotificationSettingsDto>>('/notifications', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  testNotification: async (payload: { provider: string }): Promise<BaseResponse<null>> => {
    try {
      return await apiClient.post<BaseResponse<null>>('/settings/notifications/test', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getNotificationSubscriptions: async (
    userId?: string
  ): Promise<BaseResponse<NotificationSubscriptionDto[]>> => {
    try {
      const q = userId ? `?user_id=${encodeURIComponent(userId)}` : '';
      return await apiClient.get<BaseResponse<NotificationSubscriptionDto[]>>(
        `/notifications/subscriptions${q}`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  upsertNotificationSubscription: async (payload: {
    category: string;
    channels: string[];
    enabled: boolean;
  }): Promise<BaseResponse<NotificationSubscriptionDto>> => {
    try {
      return await apiClient.put<BaseResponse<NotificationSubscriptionDto>>(
        '/notifications/subscriptions',
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getNotificationDefaults: async (): Promise<BaseResponse<NotificationDefaultDto[]>> => {
    try {
      return await apiClient.get<BaseResponse<NotificationDefaultDto[]>>('/notifications/defaults');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  upsertNotificationDefault: async (payload: {
    category: string;
    channels: string[];
    enabled: boolean;
  }): Promise<BaseResponse<NotificationDefaultDto>> => {
    try {
      return await apiClient.put<BaseResponse<NotificationDefaultDto>>(
        '/notifications/defaults',
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getNotificationDeliveries: async (
    limit = 50
  ): Promise<BaseResponse<NotificationDeliveryDto[]>> => {
    try {
      return await apiClient.get<BaseResponse<NotificationDeliveryDto[]>>(
        `/notifications/deliveries?limit=${limit}`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getDeviceTokens: async (): Promise<BaseResponse<DeviceTokenDto[]>> => {
    try {
      return await apiClient.get<BaseResponse<DeviceTokenDto[]>>('/notifications/devices');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteDeviceToken: async (token: string): Promise<BaseResponse<null>> => {
    try {
      return await apiClient.delete<BaseResponse<null>>(
        `/notifications/devices/${encodeURIComponent(token)}`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
