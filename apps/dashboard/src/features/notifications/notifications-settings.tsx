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
} from '#/features/settings/hooks';
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

export const NotificationsSettings = () => {
  const { data, isLoading } = useGetNotificationSettings();
  const { mutateAsync: update, isPending } = useUpdateNotificationSettings();
  const { mutateAsync: testNotif, isPending: testing } = useTestNotification();

  const [form, setForm] = useState<NotifSettingsForm>(EMPTY);
  const [testingProvider, setTestingProvider] = useState<string | null>(null);

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

  const handleSave = async () => {
    try {
      await update(form as Record<string, unknown>);
      toast.success('Notification settings saved');
    } catch {
      toast.error('Failed to save notification settings');
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
      <div className="flex items-center justify-between">
        <div>
          <h2 className="font-semibold text-lg">Notification Channels</h2>
          <p className="text-muted-foreground text-sm">
            Configure where Codedock sends alerts for deployments, errors, and system events.
          </p>
        </div>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <span className="text-muted-foreground text-sm">Global alerts</span>
            <Switch
              checked={form.notificationAlerts}
              onCheckedChange={(v: boolean) => set('notificationAlerts', v)}
            />
          </div>
          <Button onClick={handleSave} disabled={isPending}>
            <Check className="mr-2 h-4 w-4" />
            {isPending ? 'Saving…' : 'Save Changes'}
          </Button>
        </div>
      </div>

      <NotificationChannelsList
        form={form}
        set={set}
        handleTest={handleTest}
        testingProvider={testingProvider}
        testing={testing}
      />
    </div>
  );
};
