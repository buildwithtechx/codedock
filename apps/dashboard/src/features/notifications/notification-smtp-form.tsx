import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import type { NotifSettingsForm } from '#/interfaces/notifications';

interface Props {
  form: NotifSettingsForm;
  set: (k: keyof NotifSettingsForm, v: unknown) => void;
}

export const NotificationSmtpForm = ({ form, set }: Props) => {
  return (
    <>
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
    </>
  );
};
