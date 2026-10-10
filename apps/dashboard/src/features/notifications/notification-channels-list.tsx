import { Mail, Webhook } from 'lucide-react';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Switch } from '#/components/ui/switch';
import { NotificationChannelRow } from './notification-channel-row';

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
  notificationAlerts: boolean;
};

type Props = {
  form: NotifSettingsForm;
  set: (k: keyof NotifSettingsForm, v: unknown) => void;
  handleSave: (provider: string) => void;
  handleTest: (provider: string) => void;
  savingProvider: string | null;
  testingProvider: string | null;
  testing: boolean;
  disabled?: boolean;
};

export const NotificationChannelsList = ({
  form,
  set,
  handleSave,
  handleTest,
  savingProvider,
  testingProvider,
  testing,
  disabled = false,
}: Props) => {
  return (
    <div className="space-y-4">
      <NotificationChannelRow
        icon={
          <img
            src="/notification-providers/discord.svg"
            alt="Discord"
            className="h-5 w-5 object-contain"
          />
        }
        name="Discord"
        description="Discord webhook channel dispatch"
        isConfigured={Boolean(form.discordWebhookUrl?.trim())}
        enabled={form.discordEnabled ?? false}
        onToggle={(v) => set('discordEnabled', v)}
        onSave={() => handleSave('discord')}
        onTest={() => handleTest('discord')}
        saving={savingProvider === 'discord'}
        testing={testingProvider === 'discord' && testing}
        disabled={disabled}
      >
        <div className="space-y-1.5">
          <Label className="text-xs">Webhook URL</Label>
          <Input
            value={form.discordWebhookUrl ?? ''}
            onChange={(e) => set('discordWebhookUrl', e.target.value)}
            placeholder="https://discord.com/api/webhooks/..."
            className="font-mono text-xs"
          />
        </div>
        <div className="flex items-center gap-2">
          <Switch
            checked={form.discordPingEnabled ?? false}
            onCheckedChange={(v) => set('discordPingEnabled', v)}
          />
          <Label className="text-muted-foreground text-xs">@here ping on critical alerts</Label>
        </div>
      </NotificationChannelRow>

      <NotificationChannelRow
        icon={
          <img
            src="/notification-providers/slack.svg"
            alt="Slack"
            className="h-5 w-5 object-contain"
          />
        }
        name="Slack"
        description="Slack incoming webhook channel dispatch"
        isConfigured={Boolean(form.slackWebhookUrl?.trim())}
        enabled={form.slackEnabled ?? false}
        onToggle={(v) => set('slackEnabled', v)}
        onSave={() => handleSave('slack')}
        onTest={() => handleTest('slack')}
        saving={savingProvider === 'slack'}
        testing={testingProvider === 'slack' && testing}
        disabled={disabled}
      >
        <div className="space-y-1.5">
          <Label className="text-xs">Webhook URL</Label>
          <Input
            value={form.slackWebhookUrl ?? ''}
            onChange={(e) => set('slackWebhookUrl', e.target.value)}
            placeholder="https://hooks.slack.com/services/..."
            className="font-mono text-xs"
          />
        </div>
      </NotificationChannelRow>

      <NotificationChannelRow
        icon={
          <img
            src="/notification-providers/telegram.svg"
            alt="Telegram"
            className="h-5 w-5 object-contain"
          />
        }
        name="Telegram"
        description="Telegram bot token and destination chat ID"
        isConfigured={Boolean(form.telegramBotToken?.trim() && form.telegramChatId?.trim())}
        enabled={form.telegramEnabled ?? false}
        onToggle={(v) => set('telegramEnabled', v)}
        onSave={() => handleSave('telegram')}
        onTest={() => handleTest('telegram')}
        saving={savingProvider === 'telegram'}
        testing={testingProvider === 'telegram' && testing}
        disabled={disabled}
      >
        <div className="space-y-1.5">
          <Label className="text-xs">Bot Token</Label>
          <Input
            type="password"
            value={form.telegramBotToken ?? ''}
            onChange={(e) => set('telegramBotToken', e.target.value)}
            placeholder="1234567890:AAF..."
            className="font-mono text-xs"
          />
        </div>
        <div className="space-y-1.5">
          <Label className="text-xs">Chat ID</Label>
          <Input
            value={form.telegramChatId ?? ''}
            onChange={(e) => set('telegramChatId', e.target.value)}
            placeholder="-100..."
            className="font-mono text-xs"
          />
        </div>
      </NotificationChannelRow>

      <NotificationChannelRow
        icon={<Mail className="h-5 w-5" />}
        name="SMTP Email"
        description="Standard SMTP mail server credentials for alerts"
        isConfigured={Boolean(form.smtpHost?.trim())}
        enabled={form.smtpEnabled ?? false}
        onToggle={(v) => set('smtpEnabled', v)}
        onSave={() => handleSave('smtp')}
        onTest={() => handleTest('smtp')}
        saving={savingProvider === 'smtp'}
        testing={testingProvider === 'smtp' && testing}
        disabled={disabled}
      >
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <div className="space-y-1.5">
            <Label className="text-xs">Host</Label>
            <Input
              value={form.smtpHost ?? ''}
              onChange={(e) => set('smtpHost', e.target.value)}
              placeholder="smtp.example.com"
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-1.5">
            <Label className="text-xs">Port</Label>
            <Input
              type="number"
              value={form.smtpPort ?? 587}
              onChange={(e) => set('smtpPort', Number(e.target.value))}
              className="font-mono text-xs"
            />
          </div>
        </div>
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <div className="space-y-1.5">
            <Label className="text-xs">Username</Label>
            <Input
              value={form.smtpUser ?? ''}
              onChange={(e) => set('smtpUser', e.target.value)}
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-1.5">
            <Label className="text-xs">Password</Label>
            <Input
              type="password"
              value={form.smtpPassword ?? ''}
              onChange={(e) => set('smtpPassword', e.target.value)}
              className="font-mono text-xs"
            />
          </div>
        </div>
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <div className="space-y-1.5">
            <Label className="text-xs">From Name</Label>
            <Input
              value={form.smtpFromName ?? ''}
              onChange={(e) => set('smtpFromName', e.target.value)}
              placeholder="Codedock Alerts"
              className="text-xs"
            />
          </div>
          <div className="space-y-1.5">
            <Label className="text-xs">From Address</Label>
            <Input
              value={form.smtpFromAddress ?? ''}
              onChange={(e) => set('smtpFromAddress', e.target.value)}
              placeholder="alerts@example.com"
              className="font-mono text-xs"
            />
          </div>
        </div>
      </NotificationChannelRow>

      <NotificationChannelRow
        icon={<Mail className="h-5 w-5" />}
        name="Resend"
        description="Modern email API service via Resend API key"
        isConfigured={Boolean(form.resendApiKey?.trim())}
        enabled={form.resendEnabled ?? false}
        onToggle={(v) => set('resendEnabled', v)}
        onSave={() => handleSave('resend')}
        onTest={() => handleTest('resend')}
        saving={savingProvider === 'resend'}
        testing={testingProvider === 'resend' && testing}
        disabled={disabled}
      >
        <div className="space-y-1.5">
          <Label className="text-xs">API Key</Label>
          <Input
            type="password"
            value={form.resendApiKey ?? ''}
            onChange={(e) => set('resendApiKey', e.target.value)}
            placeholder="re_..."
            className="font-mono text-xs"
          />
        </div>
      </NotificationChannelRow>

      <NotificationChannelRow
        icon={
          <img
            src="/notification-providers/pushover.svg"
            alt="Pushover"
            className="h-5 w-5 object-contain"
          />
        }
        name="Pushover"
        description="Real-time push notifications to iOS and Android devices"
        isConfigured={Boolean(form.pushoverUserKey?.trim() && form.pushoverApiToken?.trim())}
        enabled={form.pushoverEnabled ?? false}
        onToggle={(v) => set('pushoverEnabled', v)}
        onSave={() => handleSave('pushover')}
        onTest={() => handleTest('pushover')}
        saving={savingProvider === 'pushover'}
        testing={testingProvider === 'pushover' && testing}
        disabled={disabled}
      >
        <div className="space-y-1.5">
          <Label className="text-xs">User Key</Label>
          <Input
            value={form.pushoverUserKey ?? ''}
            onChange={(e) => set('pushoverUserKey', e.target.value)}
            className="font-mono text-xs"
          />
        </div>
        <div className="space-y-1.5">
          <Label className="text-xs">API Token</Label>
          <Input
            type="password"
            value={form.pushoverApiToken ?? ''}
            onChange={(e) => set('pushoverApiToken', e.target.value)}
            className="font-mono text-xs"
          />
        </div>
      </NotificationChannelRow>

      <NotificationChannelRow
        icon={<Webhook className="h-5 w-5" />}
        name="Generic Webhook"
        description="HTTP POST JSON payload dispatch to custom HTTP endpoints"
        isConfigured={Boolean(form.genericWebhookUrl?.trim())}
        enabled={form.genericWebhookEnabled ?? false}
        onToggle={(v) => set('genericWebhookEnabled', v)}
        onSave={() => handleSave('webhook')}
        onTest={() => handleTest('webhook')}
        saving={savingProvider === 'webhook'}
        testing={testingProvider === 'webhook' && testing}
        disabled={disabled}
      >
        <div className="space-y-1.5">
          <Label className="text-xs">Webhook URL</Label>
          <Input
            value={form.genericWebhookUrl ?? ''}
            onChange={(e) => set('genericWebhookUrl', e.target.value)}
            placeholder="https://api.example.com/webhook"
            className="font-mono text-xs"
          />
        </div>
      </NotificationChannelRow>
    </div>
  );
};
