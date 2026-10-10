import { Bell, Check } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Skeleton } from '#/components/ui/skeleton';
import { Switch } from '#/components/ui/switch';
import {
  useGetNotificationSettings,
  useTestNotification,
  useUpdateNotificationSettings,
} from '#/features/settings';
import { SettingsSection } from '#/features/settings/settings-section';
import { NotificationChannelsList, type NotifSettingsForm } from './notification-channels-list';
import { NotificationSubscriptionsTable } from './notification-subscriptions-table';

const EMPTY: NotifSettingsForm = {
  discordWebhookUrl: '',
  discordPingEnabled: false,
  discordEnabled: false,
  slackWebhookUrl: '',
  slackEnabled: false,
  telegramBotToken: '',
  telegramChatId: '',
  telegramEnabled: false,
  smtpHost: '',
  smtpPort: 587,
  smtpUser: '',
  smtpPassword: '',
  smtpFromName: '',
  smtpFromAddress: '',
  smtpEnabled: false,
  resendApiKey: '',
  resendEnabled: false,
  pushoverUserKey: '',
  pushoverApiToken: '',
  pushoverEnabled: false,
  genericWebhookUrl: '',
  genericWebhookEnabled: false,
  notificationAlerts: true,
};

const notificationFields: Record<string, (keyof NotifSettingsForm)[]> = {
  alerts: ['notificationAlerts'],
  discord: ['discordWebhookUrl', 'discordPingEnabled', 'discordEnabled'],
  slack: ['slackWebhookUrl', 'slackEnabled'],
  telegram: ['telegramBotToken', 'telegramChatId', 'telegramEnabled'],
  smtp: [
    'smtpHost',
    'smtpPort',
    'smtpUser',
    'smtpPassword',
    'smtpFromName',
    'smtpFromAddress',
    'smtpEnabled',
  ],
  resend: ['resendApiKey', 'resendEnabled'],
  pushover: ['pushoverUserKey', 'pushoverApiToken', 'pushoverEnabled'],
  webhook: ['genericWebhookUrl', 'genericWebhookEnabled'],
};

export const NotificationsSettings = () => {
  const { data, isLoading } = useGetNotificationSettings();
  const { mutateAsync: update, isPending } = useUpdateNotificationSettings();
  const { mutateAsync: testNotif, isPending: testing } = useTestNotification();

  const [form, setForm] = useState<NotifSettingsForm>(EMPTY);
  const [testingProvider, setTestingProvider] = useState<string | null>(null);
  const [savingProvider, setSavingProvider] = useState<string | null>(null);

  const isInitialized = useRef(false);
  const dirtyFieldsRef = useRef<Set<keyof NotifSettingsForm>>(new Set());

  useEffect(() => {
    if (data?.data) {
      const s = data.data as Record<string, unknown>;
      const serverValues: NotifSettingsForm = {
        discordWebhookUrl: (s.discordWebhookUrl as string) ?? '',
        discordPingEnabled: (s.discordPingEnabled as boolean) ?? false,
        discordEnabled: (s.discordEnabled as boolean) ?? false,
        slackWebhookUrl: (s.slackWebhookUrl as string) ?? '',
        slackEnabled: (s.slackEnabled as boolean) ?? false,
        telegramBotToken: (s.telegramBotToken as string) ?? '',
        telegramChatId: (s.telegramChatId as string) ?? '',
        telegramEnabled: (s.telegramEnabled as boolean) ?? false,
        smtpHost: (s.smtpHost as string) ?? '',
        smtpPort: (s.smtpPort as number) ?? 587,
        smtpUser: (s.smtpUser as string) ?? '',
        smtpPassword: (s.smtpPassword as string) ?? '',
        smtpFromName: (s.smtpFromName as string) ?? '',
        smtpFromAddress: (s.smtpFromAddress as string) ?? '',
        smtpEnabled: (s.smtpEnabled as boolean) ?? false,
        resendApiKey: (s.resendApiKey as string) ?? '',
        resendEnabled: (s.resendEnabled as boolean) ?? false,
        pushoverUserKey: (s.pushoverUserKey as string) ?? '',
        pushoverApiToken: (s.pushoverApiToken as string) ?? '',
        pushoverEnabled: (s.pushoverEnabled as boolean) ?? false,
        genericWebhookUrl: (s.genericWebhookUrl as string) ?? '',
        genericWebhookEnabled: (s.genericWebhookEnabled as boolean) ?? false,
        notificationAlerts: (s.notificationAlerts as boolean) ?? true,
      };

      setForm((prev) => {
        if (!isInitialized.current) {
          isInitialized.current = true;
          return serverValues;
        }
        const updated = { ...prev };
        for (const key of Object.keys(serverValues) as Array<keyof NotifSettingsForm>) {
          if (!dirtyFieldsRef.current.has(key)) {
            (updated as Record<keyof NotifSettingsForm, unknown>)[key] = serverValues[key];
          }
        }
        return updated;
      });
    }
  }, [data]);

  const set = (k: keyof NotifSettingsForm, v: unknown) => {
    dirtyFieldsRef.current.add(k);
    setForm((f) => ({ ...f, [k]: v }));
  };

  const handleSave = async (provider: string) => {
    const fields = notificationFields[provider];
    if (!fields) return;

    setSavingProvider(provider);
    try {
      await update(Object.fromEntries(fields.map((field) => [field, form[field]])));
      for (const field of fields) {
        dirtyFieldsRef.current.delete(field);
      }
      toast.success(`${provider === 'alerts' ? 'Alert behavior' : provider} settings saved`);
    } catch {
      toast.error(`Failed to save ${provider} settings`);
    } finally {
      setSavingProvider(null);
    }
  };

  const handleTest = async (provider: string) => {
    setTestingProvider(provider);
    try {
      await testNotif({ provider });
      toast.success(`Test notification sent via ${provider}`);
    } catch {
      toast.error(`Failed to send test via ${provider}`);
    } finally {
      setTestingProvider(null);
    }
  };

  const isSavingAny = isPending || Boolean(savingProvider);

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[...Array(4)].map((_, i) => (
          <Skeleton key={i} className="h-40 w-full rounded-xl" />
        ))}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <SettingsSection
        icon={<Bell className="size-4 text-primary" />}
        title="Alert Behavior"
        description="Enable or pause delivery across every configured channel."
        action={
          <div className="flex items-center gap-3">
            <Switch
              checked={form.notificationAlerts}
              onCheckedChange={(v: boolean) => set('notificationAlerts', v)}
            />
            <Button size="sm" onClick={() => handleSave('alerts')} disabled={isSavingAny}>
              <Check className="mr-2 h-4 w-4" />
              {savingProvider === 'alerts' ? 'Saving...' : 'Save Alerts'}
            </Button>
          </div>
        }
      >
        <p className="text-muted-foreground text-xs leading-relaxed">
          When paused, notifications for failed deployments, container restarts, and runtime
          warnings will be suppressed across all configured communication channels.
        </p>
      </SettingsSection>

      <NotificationChannelsList
        form={form}
        set={set}
        handleSave={handleSave}
        handleTest={handleTest}
        savingProvider={savingProvider}
        testingProvider={testingProvider}
        testing={testing}
        disabled={isSavingAny}
      />

      <NotificationSubscriptionsTable />
    </div>
  );
};
