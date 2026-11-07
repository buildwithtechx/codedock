export interface TeamNotificationChannel {
  id: string;
  provider: string;
  config: Record<string, unknown>;
  events: Record<string, unknown>;
  isEnabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface NotificationEvent {
  title: string;
  message: string;
  level: string;
  eventType: string;
  projectId?: string;
  url?: string;
}

export interface NotificationSubscriptionDto {
  id: string;
  userId: string;
  organizationId: string;
  category: string;
  channels: string[];
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface NotificationDefaultDto {
  organizationId: string;
  category: string;
  channels: string[];
  enabled: boolean;
  updatedAt: string;
}

export interface NotificationDeliveryDto {
  id: string;
  eventType: string;
  channel: string;
  title: string;
  success: boolean;
  error: string;
  createdAt: string;
}

export interface DeviceTokenDto {
  id: string;
  userId: string;
  token: string;
  platform: string;
  createdAt: string;
  updatedAt: string;
}

export interface NotificationSettingsDto {
  id?: string;
  discordWebhookUrl?: string;
  discordPingEnabled?: boolean;
  discordEnabled?: boolean;
  slackWebhookUrl?: string;
  slackEnabled?: boolean;
  telegramBotToken?: string;
  telegramChatId?: string;
  telegramEnabled?: boolean;
  smtpHost?: string;
  smtpPort?: number;
  smtpUser?: string;
  smtpPassword?: string;
  smtpFromName?: string;
  smtpFromAddress?: string;
  smtpEnabled?: boolean;
  resendApiKey?: string;
  resendEnabled?: boolean;
  pushoverUserKey?: string;
  pushoverApiToken?: string;
  pushoverEnabled?: boolean;
  genericWebhookUrl?: string;
  genericWebhookEnabled?: boolean;
  expoPushToken?: string;
  expoEnabled?: boolean;
  notificationAlerts?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export type NotifSettingsForm = {
  discordWebhookUrl: string;
  discordPingEnabled: boolean;
  discordEnabled: boolean;
  slackWebhookUrl: string;
  slackEnabled: boolean;
  telegramBotToken: string;
  telegramChatId: string;
  telegramEnabled: boolean;
  smtpHost: string;
  smtpPort: number;
  smtpUser: string;
  smtpPassword: string;
  smtpFromName: string;
  smtpFromAddress: string;
  smtpEnabled: boolean;
  resendApiKey: string;
  resendEnabled: boolean;
  pushoverUserKey: string;
  pushoverApiToken: string;
  pushoverEnabled: boolean;
  genericWebhookUrl: string;
  genericWebhookEnabled: boolean;
  expoPushToken: string;
  expoEnabled: boolean;
  notificationAlerts: boolean;
};
