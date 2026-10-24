import { Checkbox } from '@codedock/ui/checkbox';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';

export const RELAY_PROVIDERS = [
  { id: 'ses', label: 'Amazon SES', host: 'email-smtp.us-east-1.amazonaws.com', port: '587' },
  { id: 'postmark', label: 'Postmark', host: 'smtp.postmarkapp.com', port: '587' },
  { id: 'sendgrid', label: 'SendGrid', host: 'smtp.sendgrid.net', port: '587' },
  { id: 'resend', label: 'Resend', host: 'smtp.resend.com', port: '587' },
  { id: 'custom', label: 'Custom SMTP', host: '', port: '587' },
];

export interface RelaySetupState {
  enabled: boolean;
  provider: string;
  host: string;
  port: string;
  username: string;
  password: string;
}

interface RelaySetupAccordionProps {
  relay: RelaySetupState;
  onChange: (relay: RelaySetupState) => void;
}

export function RelaySetupAccordion({ relay, onChange }: RelaySetupAccordionProps) {
  const handleSelectRelay = (providerId: string) => {
    const prov = RELAY_PROVIDERS.find((p) => p.id === providerId);
    onChange({
      ...relay,
      provider: providerId,
      host: prov?.host ?? '',
      port: prov?.port ?? '587',
    });
  };

  return (
    <div className="space-y-4 rounded-xl border border-border/50 bg-muted/[0.15] p-4">
      <div
        className="flex cursor-pointer items-start gap-3"
        onClick={() => onChange({ ...relay, enabled: !relay.enabled })}
      >
        <Checkbox
          checked={relay.enabled}
          onCheckedChange={(val) => onChange({ ...relay, enabled: Boolean(val) })}
          className="mt-0.5"
        />
        <div className="min-w-0">
          <span className="block font-medium text-foreground text-sm">
            Outbound relay (recommended)
          </span>
          <span className="mt-0.5 block text-muted-foreground text-xs leading-relaxed">
            Send outgoing mail through an established email provider to avoid blocklisted server IPs
            and port 25 limits.
          </span>
        </div>
      </div>

      {relay.enabled && (
        <div className="space-y-4 pt-2">
          <div className="flex flex-wrap gap-2">
            {RELAY_PROVIDERS.map((p) => (
              <button
                key={p.id}
                type="button"
                onClick={() => handleSelectRelay(p.id)}
                className={`inline-flex items-center gap-2 rounded-xl border px-3 py-1.5 font-medium text-xs transition-all ${
                  relay.provider === p.id
                    ? 'border-primary/40 bg-primary/[0.08] text-foreground'
                    : 'border-border/60 text-muted-foreground hover:bg-muted/40'
                }`}
              >
                {p.label}
              </button>
            ))}
          </div>

          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1">
              <Label className="text-xs">SMTP Host</Label>
              <Input
                value={relay.host}
                onChange={(e) => onChange({ ...relay, host: e.target.value })}
                placeholder="smtp.example.com"
                className="rounded-lg text-xs"
                required
              />
            </div>
            <div className="space-y-1">
              <Label className="text-xs">SMTP Port</Label>
              <Input
                value={relay.port}
                onChange={(e) => onChange({ ...relay, port: e.target.value })}
                placeholder="587"
                className="rounded-lg text-xs"
                required
              />
            </div>
            <div className="space-y-1">
              <Label className="text-xs">SMTP Username</Label>
              <Input
                value={relay.username}
                onChange={(e) => onChange({ ...relay, username: e.target.value })}
                placeholder="smtp-user"
                className="rounded-lg text-xs"
                required
              />
            </div>
            <div className="space-y-1">
              <Label className="text-xs">SMTP Password</Label>
              <Input
                type="password"
                value={relay.password}
                onChange={(e) => onChange({ ...relay, password: e.target.value })}
                placeholder="••••••••••••"
                className="rounded-lg text-xs"
                required
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
