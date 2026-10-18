import Constants from 'expo-constants';
import * as Notifications from 'expo-notifications';
import { Platform } from 'react-native';
import { MobileApi } from '@/services/api';

Notifications.setNotificationHandler({
  handleNotification: async () => ({
    shouldShowAlert: true,
    shouldPlaySound: true,
    shouldSetBadge: true,
    shouldShowBanner: true,
    shouldShowList: true,
  }),
});

export async function registerForPushNotificationsAsync(
  instanceUrl: string,
  token: string
): Promise<string | null> {
  if (Platform.OS === 'web') {
    return null;
  }

  const { status: existingStatus } = await Notifications.getPermissionsAsync();
  let finalStatus = existingStatus;

  if (existingStatus !== 'granted') {
    const { status } = await Notifications.requestPermissionsAsync();
    finalStatus = status;
  }

  if (finalStatus !== 'granted') {
    return null;
  }

  const projectId = Constants?.expoConfig?.extra?.eas?.projectId ?? Constants?.easConfig?.projectId;

  try {
    const pushTokenData = await Notifications.getExpoPushTokenAsync({
      projectId: projectId || undefined,
    });
    const deviceToken = pushTokenData.data;

    if (deviceToken && instanceUrl && token) {
      await MobileApi.registerDeviceToken(instanceUrl, token, deviceToken, Platform.OS);
    }

    return deviceToken;
  } catch (error) {
    console.warn('Failed to register push token:', error);
    return null;
  }
}
