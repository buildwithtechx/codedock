import { Link, useNavigate } from '@tanstack/react-router';
import { ArrowLeft, Check, Globe, Mail, Plug, Server } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import {
  getMailProvider,
  MAIL_PROVIDERS,
  type MailProviderId,
} from './mail-providers-catalog';

type WizardPhase = 'choose' | 'connect' | 'done';

export function MailWizardView() {
  const navigate = useNavigate();
  const [phase, setPhase] = useState<WizardPhase>('choose');
  const [preset, setPreset] = useState<MailProviderId>('custom');
  const [hostname, setHostname] = useState('');
  const [imapHost, setImapHost] = useState('');
  const [imapPort, setImapPort] = useState(993);
  const [smtpHost, setSmtpHost] = useState('');
  const [smtpPort, setSmtpPort] = useState(465);
  const [isDeploying, setIsDeploying] = useState(false);

  const applyPreset = (id: MailProviderId) => {
    setPreset(id);
    const p = getMailProvider(id);
    setImapHost(p.imapHost);
    setImapPort(p.imapPort);
    setSmtpHost(p.smtpHost);
    setSmtpPort(p.smtpPort);
  };

  const handleDeploy = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!hostname.trim()) {
      toast.error('Hostname is required (e.g. mail.example.com)');
      return;
    }
    if (!smtpHost.trim()) {
      toast.error('SMTP host is required');
      return;
    }

    setIsDeploying(true);
    try {
      // Emulate webmail deployment setup
      await new Promise((resolve) => setTimeout(resolve, 800));
      toast.success('Webmail deployment created successfully');
      setPhase('done');
    } catch {
      toast.error('Failed to create webmail deployment');
    } finally {
      setIsDeploying(false);
    }
  };

  const currentProvider = getMailProvider(preset);

  if (phase === 'done') {
    return (
      <div className="mx-auto max-w-xl py-12 text-center">
        <div className="mx-auto flex size-14 items-center justify-center rounded-2xl bg-emerald-500/10 text-emerald-500">
          <Check className="size-7" />
        </div>
        <h1 className="mt-4 font-semibold text-2xl text-foreground">Webmail is deploying</h1>
        <p className="mt-2 text-muted-foreground text-sm">
          Your webmail client is being provisioned at <span className="font-mono text-foreground">{hostname}</span>.
        </p>
        <div className="mt-6 flex justify-center gap-3">
          <Button asChild variant="outline" className="rounded-xl">
            <Link to={'/emails' as never}>Manage Emails</Link>
          </Button>
          <Button asChild className="rounded-xl">
            <Link to={'/apps' as never}>View Apps</Link>
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-3xl pb-20 pt-6">
      <button
        type="button"
        onClick={() => (phase === 'connect' ? setPhase('choose') : navigate({ to: '/apps/new' as never }))}
        className="mb-6 inline-flex items-center gap-1.5 text-muted-foreground text-sm transition-colors hover:text-foreground"
      >
        <ArrowLeft className="size-4 rtl:rotate-180" /> Back
      </button>

      <div className="flex items-center gap-4">
        <div className="flex size-12 shrink-0 items-center justify-center rounded-2xl bg-muted/60 text-foreground">
          <Mail className="size-6 text-primary" />
        </div>
        <div className="min-w-0">
          <h1 className="font-semibold text-foreground text-xl">Codedock Mail</h1>
          <p className="text-muted-foreground text-sm">
            Your own mail on your own domain — SMTP/IMAP, mailboxes and webmail on one server. Or keep the mailboxes you already have and just run the webmail.
          </p>
        </div>
      </div>

      {phase === 'choose' && (
        <div className="mt-8 space-y-4">
          <button
            type="button"
            onClick={() => navigate({ to: '/emails' as never })}
            className="group flex w-full items-start gap-4 rounded-2xl border border-border/60 bg-card p-5 text-start transition-all hover:border-primary/40 hover:bg-muted/30"
          >
            <span className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-muted/60 text-foreground transition-colors group-hover:bg-primary/10 group-hover:text-primary">
              <Server className="size-5" />
            </span>
            <div className="min-w-0 flex-1">
              <h2 className="font-medium text-base text-foreground">Self-host mail server</h2>
              <p className="mt-1 text-muted-foreground text-sm leading-relaxed">
                Deploy the complete mail engine (maild, SMTP/IMAP daemon, DKIM, and mailbox manager) on your server.
              </p>
            </div>
          </button>

          <button
            type="button"
            onClick={() => setPhase('connect')}
            className="group flex w-full items-start gap-4 rounded-2xl border border-border/60 bg-card p-5 text-start transition-all hover:border-primary/40 hover:bg-muted/30"
          >
            <span className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-muted/60 text-foreground transition-colors group-hover:bg-primary/10 group-hover:text-primary">
              <Plug className="size-5" />
            </span>
            <div className="min-w-0 flex-1">
              <h2 className="font-medium text-base text-foreground">Connect existing mailboxes</h2>
              <p className="mt-1 text-muted-foreground text-sm leading-relaxed">
                Install a fast webmail client pointed at an existing provider (Gmail, Fastmail, or an outbound relay like SES, SendGrid, or Resend).
              </p>
            </div>
          </button>
        </div>
      )}

      {phase === 'connect' && (
        <form onSubmit={handleDeploy} className="mt-8 space-y-6">
          <div className="rounded-2xl border border-border/50 bg-card p-5">
            <h3 className="font-medium text-foreground text-sm">Provider preset</h3>
            <div className="mt-3 flex flex-wrap gap-2">
              {MAIL_PROVIDERS.map((p) => {
                const isSelected = preset === p.id;
                return (
                  <button
                    key={p.id}
                    type="button"
                    onClick={() => applyPreset(p.id)}
                    className={`inline-flex items-center gap-2 rounded-xl border px-3 py-1.5 text-xs font-medium transition-colors ${
                      isSelected
                        ? 'border-primary bg-primary/10 text-foreground'
                        : 'border-border/60 bg-muted/30 text-muted-foreground hover:bg-muted/60 hover:text-foreground'
                    }`}
                  >
                    {p.logoSrc && (
                      <img src={p.logoSrc} alt={p.label} className="size-3.5 object-contain" />
                    )}
                    {p.label}
                  </button>
                );
              })}
            </div>

            {currentProvider.hint && (
              <p className="mt-3 text-muted-foreground text-xs leading-relaxed">
                💡 {currentProvider.hint}
              </p>
            )}
          </div>

          <div className="space-y-4 rounded-2xl border border-border/50 bg-card p-5">
            <div className="space-y-1.5">
              <label htmlFor="webmail-host" className="font-medium text-foreground text-sm">
                Webmail hostname
              </label>
              <Input
                id="webmail-host"
                placeholder="mail.example.com"
                value={hostname}
                onChange={(e) => setHostname(e.target.value)}
                required
                className="h-10 rounded-xl bg-muted/40"
              />
              <p className="text-muted-foreground text-xs">
                Domain where your webmail client will be accessible over HTTPS.
              </p>
            </div>

            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <label htmlFor="imap-host" className="font-medium text-foreground text-sm">
                  IMAP host
                </label>
                <Input
                  id="imap-host"
                  placeholder="imap.example.com"
                  value={imapHost}
                  onChange={(e) => setImapHost(e.target.value)}
                  className="h-10 rounded-xl bg-muted/40"
                />
              </div>
              <div className="space-y-1.5">
                <label htmlFor="imap-port" className="font-medium text-foreground text-sm">
                  IMAP port
                </label>
                <Input
                  id="imap-port"
                  type="number"
                  value={imapPort}
                  onChange={(e) => setImapPort(Number(e.target.value))}
                  className="h-10 rounded-xl bg-muted/40"
                />
              </div>
            </div>

            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <label htmlFor="smtp-host" className="font-medium text-foreground text-sm">
                  SMTP host
                </label>
                <Input
                  id="smtp-host"
                  placeholder="smtp.example.com"
                  value={smtpHost}
                  onChange={(e) => setSmtpHost(e.target.value)}
                  required
                  className="h-10 rounded-xl bg-muted/40"
                />
              </div>
              <div className="space-y-1.5">
                <label htmlFor="smtp-port" className="font-medium text-foreground text-sm">
                  SMTP port
                </label>
                <Input
                  id="smtp-port"
                  type="number"
                  value={smtpPort}
                  onChange={(e) => setSmtpPort(Number(e.target.value))}
                  className="h-10 rounded-xl bg-muted/40"
                />
              </div>
            </div>
          </div>

          <div className="flex items-center justify-end gap-3">
            <Button
              type="button"
              variant="ghost"
              onClick={() => setPhase('choose')}
              className="rounded-xl"
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isDeploying} className="rounded-xl">
              {isDeploying ? 'Deploying...' : 'Deploy Webmail'}
            </Button>
          </div>
        </form>
      )}
    </div>
  );
}
