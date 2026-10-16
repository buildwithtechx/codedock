import {
  AlertTriangle,
  Eye,
  EyeOff,
  Globe,
  Key,
  Mail,
  Play,
  Server,
  Shield,
  Sparkles,
} from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { RelaySetupAccordion, type RelaySetupState } from './relay-setup-accordion';

function generatePassword(): string {
  const buf = new Uint8Array(18);
  crypto.getRandomValues(buf);
  return btoa(String.fromCharCode(...buf))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, '');
}

export interface MailSetupFormProps {
  onStart: (data: {
    domain: string;
    adminPassword: string;
    serverHost?: string;
    relayEnabled: boolean;
    relayProvider?: string;
    relayHost?: string;
    relayPort?: string;
    relayUsername?: string;
    relayPassword?: string;
  }) => void;
  isLoading?: boolean;
  onCancel?: () => void;
}

export function MailSetupForm({ onStart, isLoading = false, onCancel }: MailSetupFormProps) {
  const [domain, setDomain] = useState('');
  const [serverHost, setServerHost] = useState('');
  const [adminPassword, setAdminPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  const [relay, setRelay] = useState<RelaySetupState>({
    enabled: true,
    provider: 'ses',
    host: 'email-smtp.us-east-1.amazonaws.com',
    port: '587',
    username: '',
    password: '',
  });

  const handleGeneratePassword = () => {
    setAdminPassword(generatePassword());
    setShowPassword(true);
  };

  const isPasswordValid = adminPassword.length >= 12;
  const isFormValid =
    domain.trim() !== '' &&
    isPasswordValid &&
    (!relay.enabled || (relay.host.trim() !== '' && relay.username.trim() !== ''));

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!isFormValid || isLoading) return;

    onStart({
      domain: domain.trim(),
      adminPassword,
      serverHost: serverHost.trim() || undefined,
      relayEnabled: relay.enabled,
      relayProvider: relay.enabled ? relay.provider : undefined,
      relayHost: relay.enabled ? relay.host.trim() : undefined,
      relayPort: relay.enabled ? relay.port.trim() : undefined,
      relayUsername: relay.enabled ? relay.username.trim() : undefined,
      relayPassword: relay.enabled ? relay.password : undefined,
    });
  };

  return (
    <div className="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_380px]">
      <form onSubmit={handleSubmit} className="space-y-6 rounded-2xl border border-border/50 bg-card p-6">
        <div className="flex items-center gap-3">
          <div className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-violet-500/10 text-violet-500">
            <Mail className="size-5" />
          </div>
          <div>
            <h2 className="font-medium text-foreground text-lg">Install mail server</h2>
            <p className="text-muted-foreground text-sm">
              Set up a complete mail server stack with Postfix, Dovecot, and SpamAssassin.
            </p>
          </div>
        </div>

        <div className="space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="mail-domain" className="font-medium text-sm">
              Primary Mail Domain
            </Label>
            <Input
              id="mail-domain"
              type="text"
              value={domain}
              onChange={(e) => setDomain(e.target.value)}
              placeholder="example.com"
              className="rounded-xl"
              required
            />
            <p className="text-muted-foreground text-xs">
              Will be accessible at{' '}
              <strong className="text-foreground">mail.{domain.trim() || 'example.com'}</strong>
            </p>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="server-host" className="font-medium text-sm">
              Target Server Host or IP (Optional)
            </Label>
            <Input
              id="server-host"
              type="text"
              value={serverHost}
              onChange={(e) => setServerHost(e.target.value)}
              placeholder="e.g. 192.168.1.10 or default cluster"
              className="rounded-xl"
            />
          </div>

          <div className="space-y-1.5">
            <div className="flex items-baseline justify-between">
              <Label htmlFor="admin-password" className="font-medium text-sm">
                Admin Password
              </Label>
              <span className="font-mono text-muted-foreground text-xs">
                postmaster@{domain.trim() || 'your-domain.com'}
              </span>
            </div>
            <div className="flex items-stretch gap-2">
              <div className="relative flex-1">
                <Input
                  id="admin-password"
                  type={showPassword ? 'text' : 'password'}
                  value={adminPassword}
                  onChange={(e) => setAdminPassword(e.target.value)}
                  placeholder="Minimum 12 characters"
                  className="rounded-xl pe-10 font-mono text-sm"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPassword((prev) => !prev)}
                  className="absolute top-1/2 right-2.5 -translate-y-1/2 p-1 text-muted-foreground hover:text-foreground"
                >
                  {showPassword ? <EyeOff className="size-4" /> : <Eye className="size-4" />}
                </button>
              </div>
              <Button
                type="button"
                variant="secondary"
                onClick={handleGeneratePassword}
                className="gap-1.5 rounded-xl text-xs"
              >
                <Sparkles className="size-3.5" />
                Generate
              </Button>
            </div>
            <p className="text-muted-foreground text-xs">
              At least 12 characters. Protects webmail and administrative management.
            </p>
          </div>

          <RelaySetupAccordion relay={relay} onChange={setRelay} />
        </div>

        <div className="flex flex-wrap items-center gap-4 pt-2">
          <Button
            type="submit"
            disabled={!isFormValid || isLoading}
            className="gap-2 rounded-xl"
          >
            <Play className="size-4" />
            {isLoading ? 'Installing mail server...' : 'Install mail server'}
          </Button>

          {onCancel && (
            <Button type="button" variant="ghost" onClick={onCancel} className="rounded-xl">
              Cancel
            </Button>
          )}
        </div>
      </form>

      <div className="space-y-4">
        <div className="space-y-3 rounded-2xl border border-border/50 bg-card p-5">
          <div className="flex items-center gap-2 font-medium text-foreground text-sm">
            <Globe className="size-4 text-violet-500" />
            Automatic Records
          </div>
          <p className="text-muted-foreground text-xs leading-relaxed">
            During setup, MX, SPF, DKIM, and DMARC DNS records are generated automatically so your server reaches 10/10 deliverability scores.
          </p>
        </div>

        <div className="space-y-3 rounded-2xl border border-border/50 bg-card p-5">
          <div className="flex items-center gap-2 font-medium text-foreground text-sm">
            <Shield className="size-4 text-emerald-500" />
            Spam & Antivirus Protection
          </div>
          <p className="text-muted-foreground text-xs leading-relaxed">
            SpamAssassin and Rspamd run locally to filter incoming spam, phishing attempts, and unauthorized relaying.
          </p>
        </div>

        <div className="space-y-3 rounded-2xl border border-border/50 bg-card p-5">
          <div className="flex items-center gap-2 font-medium text-foreground text-sm">
            <Server className="size-4 text-blue-500" />
            Self-Hosted Isolation
          </div>
          <p className="text-muted-foreground text-xs leading-relaxed">
            Mail storage, database accounts, and encrypted credentials stay entirely within your dedicated VPS server instances.
          </p>
        </div>
      </div>
    </div>
  );
}
