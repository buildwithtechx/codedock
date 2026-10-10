import { Globe, ShieldCheck } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Badge } from '#/components/ui/badge';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useGetSettings } from '#/features/settings';

function slugify(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)+/g, '');
}

export function InstallDomainsRouting({
  hostPort,
  onHostPortChange,
  defaultPort,
  domain,
  onDomainChange,
  appName,
  disabled = false,
}: {
  hostPort: string;
  onHostPortChange: (port: string) => void;
  defaultPort: number;
  domain: string;
  onDomainChange: (domain: string) => void;
  appName: string;
  disabled?: boolean;
}) {
  const { data: settingsRes } = useGetSettings();
  const wildcardDomain = settingsRes?.data?.defaultWildcardDomain?.trim() || 'codedock.run';

  const [mode, setMode] = useState<'free' | 'custom' | 'none'>('free');
  const [subdomain, setSubdomain] = useState(() => slugify(appName) || 'app');

  useEffect(() => {
    if (mode === 'free') {
      const slug = slugify(appName) || 'app';
      setSubdomain(slug);
      onDomainChange(`${slug}.${wildcardDomain}`);
    } else if (mode === 'none') {
      onDomainChange('');
    }
  }, [appName, mode, wildcardDomain]);

  const handleSubdomainChange = (val: string) => {
    const cleaned = slugify(val);
    setSubdomain(cleaned);
    onDomainChange(cleaned ? `${cleaned}.${wildcardDomain}` : '');
  };

  const handleModeChange = (nextMode: 'free' | 'custom' | 'none') => {
    setMode(nextMode);
    if (nextMode === 'free') {
      const slug = subdomain || slugify(appName) || 'app';
      onDomainChange(`${slug}.${wildcardDomain}`);
    } else if (nextMode === 'none') {
      onDomainChange('');
    } else {
      onDomainChange('');
    }
  };

  return (
    <div className="space-y-4 rounded-2xl bg-card p-5">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="font-semibold text-foreground text-sm">Domains & Routing</h2>
          <p className="text-muted-foreground text-xs">
            Configure how traffic routes to this application.
          </p>
        </div>
        <Badge variant="outline" className="gap-1 font-mono text-[10px]">
          <Globe className="size-3" />
          HTTP / HTTPS
        </Badge>
      </div>

      <div className="flex rounded-lg bg-muted/40 p-1">
        {(['free', 'custom', 'none'] as const).map((m) => (
          <button
            key={m}
            type="button"
            disabled={disabled}
            onClick={() => handleModeChange(m)}
            className={`flex-1 rounded-md py-1.5 font-medium text-xs capitalize transition-colors ${
              mode === m
                ? 'bg-card text-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            {m === 'free' ? 'Free subdomain' : m === 'custom' ? 'Custom domain' : 'No domain'}
          </button>
        ))}
      </div>

      {mode === 'free' && (
        <div className="space-y-2">
          <Label htmlFor="wizard-subdomain">Subdomain</Label>
          <div className="flex items-center rounded-lg border border-border/60 bg-background px-3 py-1.5">
            <Input
              id="wizard-subdomain"
              value={subdomain}
              disabled={disabled}
              onChange={(e) => handleSubdomainChange(e.target.value)}
              placeholder="my-app"
              className="h-7 border-0 p-0 text-xs shadow-none focus-visible:ring-0"
            />
            <span className="shrink-0 font-mono text-muted-foreground text-xs">
              .{wildcardDomain}
            </span>
          </div>
          <div className="flex items-center gap-1.5 text-emerald-600 text-xs dark:text-emerald-400">
            <ShieldCheck className="size-3.5 shrink-0" />
            <span className="truncate">
              Live at https://{subdomain || 'my-app'}.{wildcardDomain} with automated SSL
            </span>
          </div>
        </div>
      )}

      {mode === 'custom' && (
        <div className="space-y-2">
          <Label htmlFor="wizard-domain">Custom Domain</Label>
          <Input
            id="wizard-domain"
            placeholder="app.mycompany.com"
            value={domain}
            disabled={disabled}
            onChange={(e) => onDomainChange(e.target.value)}
            className="text-xs"
          />
          <p className="text-muted-foreground text-xs">
            Point a CNAME or A record from your DNS provider to this server.
          </p>
        </div>
      )}

      {mode === 'none' && (
        <p className="text-muted-foreground text-xs">
          This service will not be accessible via a public domain name. It will only receive traffic
          internally or via its host port.
        </p>
      )}

      <div className="border-border/40 border-t pt-4">
        <div className="space-y-2">
          <Label htmlFor="wizard-port">Host port</Label>
          <Input
            id="wizard-port"
            inputMode="numeric"
            placeholder={String(defaultPort)}
            value={hostPort}
            disabled={disabled}
            onChange={(event) => onHostPortChange(event.target.value)}
            className="text-xs"
          />
          <p className="text-[11px] text-muted-foreground">
            Leave empty to let Docker assign an available port, or override with a specific host
            port.
          </p>
        </div>
      </div>
    </div>
  );
}
