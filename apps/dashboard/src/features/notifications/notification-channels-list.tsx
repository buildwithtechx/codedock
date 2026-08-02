import { Bell, Mail, MessageSquare, Phone, Send, Webhook } from 'lucide-react';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Switch } from '#/components/ui/switch';
import { NotificationSection } from './notification-section';

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
};

export const NotificationChannelsList = ({
  form,
  set,
  handleSave,
  handleTest,
  savingProvider,
  testingProvider,
  testing,
}: Props) => {
  return (
    <>
      <NotificationSection
        icon={<MessageSquare className="h-4 w-4" />}
        title="Discord"
        provider="discord"
        enabled={form.discordEnabled ?? false}
        onToggle={(v) => set('discordEnabled', v)}
        onSave={() => handleSave('discord')}
        onTest={() => handleTest('discord')}
        saving={savingProvider === 'discord'}
        testing={testingProvider === 'discord' && testing}
      >
        <div className="space-y-2">
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
      </NotificationSection>

      <NotificationSection
        icon={<Bell className="h-4 w-4" />}
        title="Slack"
        provider="slack"
        enabled={form.slackEnabled ?? false}
        onToggle={(v) => set('slackEnabled', v)}
        onSave={() => handleSave('slack')}
        onTest={() => handleTest('slack')}
        saving={savingProvider === 'slack'}
        testing={testingProvider === 'slack' && testing}
      >
        <div className="space-y-2">
          <Label className="text-xs">Webhook URL</Label>
          <Input
            value={form.slackWebhookUrl ?? ''}
            onChange={(e) => set('slackWebhookUrl', e.target.value)}
            placeholder="https://hooks.slack.com/services/..."
            className="font-mono text-xs"
          />
        </div>
      </NotificationSection>

      <NotificationSection
        icon={<Send className="h-4 w-4" />}
        title="Telegram"
        provider="telegram"
        enabled={form.telegramEnabled ?? false}
        onToggle={(v) => set('telegramEnabled', v)}
        onSave={() => handleSave('telegram')}
        onTest={() => handleTest('telegram')}
        saving={savingProvider === 'telegram'}
        testing={testingProvider === 'telegram' && testing}
      >
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label className="text-xs">Bot Token</Label>
            <Input
              type="password"
              value={form.telegramBotToken ?? ''}
              onChange={(e) => set('telegramBotToken', e.target.value)}
              placeholder="1234567890:AAF..."
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">Chat ID</Label>
            <Input
              value={form.telegramChatId ?? ''}
              onChange={(e) => set('telegramChatId', e.target.value)}
              placeholder="-100..."
              className="font-mono text-xs"
            />
          </div>
        </div>
      </NotificationSection>

      <NotificationSection
        icon={<Mail className="h-4 w-4" />}
        title="SMTP Email"
        provider="smtp"
        enabled={form.smtpEnabled ?? false}
        onToggle={(v) => set('smtpEnabled', v)}
        onSave={() => handleSave('smtp')}
        onTest={() => handleTest('smtp')}
        saving={savingProvider === 'smtp'}
        testing={testingProvider === 'smtp' && testing}
      >
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label className="text-xs">SMTP Host</Label>
            <Input
              value={form.smtpHost ?? ''}
              onChange={(e) => set('smtpHost', e.target.value)}
              placeholder="smtp.example.com"
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">Port</Label>
            <Input
              type="number"
              value={form.smtpPort ?? 587}
              onChange={(e) => set('smtpPort', Number(e.target.value))}
              placeholder="587"
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">Username</Label>
            <Input
              value={form.smtpUser ?? ''}
              onChange={(e) => set('smtpUser', e.target.value)}
              placeholder="user@example.com"
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">Password</Label>
            <Input
              type="password"
              value={form.smtpPassword ?? ''}
              onChange={(e) => set('smtpPassword', e.target.value)}
              placeholder="••••••••"
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">From Name</Label>
            <Input
              value={form.smtpFromName ?? ''}
              onChange={(e) => set('smtpFromName', e.target.value)}
              placeholder="Codedock"
              className="text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">From Address</Label>
            <Input
              value={form.smtpFromAddress ?? ''}
              onChange={(e) => set('smtpFromAddress', e.target.value)}
              placeholder="noreply@example.com"
              className="font-mono text-xs"
            />
          </div>
        </div>
      </NotificationSection>

      <NotificationSection
        icon={<Mail className="h-4 w-4" />}
        title="Resend"
        provider="resend"
        enabled={form.resendEnabled ?? false}
        onToggle={(v) => set('resendEnabled', v)}
        onSave={() => handleSave('resend')}
        onTest={() => handleTest('resend')}
        saving={savingProvider === 'resend'}
        testing={testingProvider === 'resend' && testing}
      >
        <div className="space-y-2">
          <Label className="text-xs">API Key</Label>
          <Input
            type="password"
            value={form.resendApiKey ?? ''}
            onChange={(e) => set('resendApiKey', e.target.value)}
            placeholder="re_..."
            className="font-mono text-xs"
          />
        </div>
      </NotificationSection>

      <NotificationSection
        icon={<Phone className="h-4 w-4" />}
        title="Pushover"
        provider="pushover"
        enabled={form.pushoverEnabled ?? false}
        onToggle={(v) => set('pushoverEnabled', v)}
        onSave={() => handleSave('pushover')}
        onTest={() => handleTest('pushover')}
        saving={savingProvider === 'pushover'}
        testing={testingProvider === 'pushover' && testing}
      >
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label className="text-xs">User Key</Label>
            <Input
              value={form.pushoverUserKey ?? ''}
              onChange={(e) => set('pushoverUserKey', e.target.value)}
              placeholder="u..."
              className="font-mono text-xs"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-xs">API Token</Label>
            <Input
              type="password"
              value={form.pushoverApiToken ?? ''}
              onChange={(e) => set('pushoverApiToken', e.target.value)}
              placeholder="a..."
              className="font-mono text-xs"
            />
          </div>
        </div>
      </NotificationSection>

      <NotificationSection
        icon={<Webhook className="h-4 w-4" />}
        title="Generic Webhook"
        provider="webhook"
        enabled={form.genericWebhookEnabled ?? false}
        onToggle={(v) => set('genericWebhookEnabled', v)}
        onSave={() => handleSave('webhook')}
        onTest={() => handleTest('webhook')}
        saving={savingProvider === 'webhook'}
        testing={testingProvider === 'webhook' && testing}
      >
        <div className="space-y-2">
          <Label className="text-xs">Webhook URL</Label>
          <Input
            value={form.genericWebhookUrl ?? ''}
            onChange={(e) => set('genericWebhookUrl', e.target.value)}
            placeholder="https://example.com/webhook"
            className="font-mono text-xs"
          />
        </div>
      </NotificationSection>
    </>
  );
};
