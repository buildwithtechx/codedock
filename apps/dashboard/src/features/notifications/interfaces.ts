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
