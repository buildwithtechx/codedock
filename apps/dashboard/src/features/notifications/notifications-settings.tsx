import { Check } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Skeleton } from '#/components/ui/skeleton';
import { Switch } from '#/components/ui/switch';
import {
  useGetNotificationSettings,
  useTestNotification,
  useUpdateNotificationSettings,
} from '#/features/settings';
import { NotificationChannelsList, type NotifSettingsForm } from './notification-channels-list';

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

  useEffect(() => {
    if (data?.data) {
      const s = data.data as Record<string, unknown>;
      setForm({
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
      });
    }
  }, [data]);

  const set = (k: keyof NotifSettingsForm, v: unknown) => setForm((f) => ({ ...f, [k]: v }));

  const handleSave = async (provider: string) => {
    const fields = notificationFields[provider];
    if (!fields) return;

    setSavingProvider(provider);
    try {
      await update(Object.fromEntries(fields.map((field) => [field, form[field]])));
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
      <section className="flex flex-col gap-4 rounded-xl border border-border/80 bg-card p-5 shadow-sm sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h3 className="font-semibold text-sm">Alert behavior</h3>
          <p className="mt-1 text-muted-foreground text-sm">
            Enable or pause delivery across every configured channel.
          </p>
        </div>
        <div className="flex items-center gap-3">
          <Switch
            checked={form.notificationAlerts}
            onCheckedChange={(v: boolean) => set('notificationAlerts', v)}
          />
          <Button size="sm" onClick={() => handleSave('alerts')} disabled={isPending}>
            <Check className="mr-2 h-4 w-4" />
            {savingProvider === 'alerts' ? 'Saving...' : 'Save'}
          </Button>
        </div>
      </section>

      <NotificationChannelsList
        form={form}
        set={set}
        handleSave={handleSave}
        handleTest={handleTest}
        savingProvider={savingProvider}
        testingProvider={testingProvider}
        testing={testing}
      />
    </div>
  );
};
