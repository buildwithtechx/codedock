import { useEffect, useState } from 'react';
import { Input } from '@codedock/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';
import { useGetPublicSettings, useGetSettings } from '#/hooks/use-settings';
import type { OneClickEndpoint } from '#/interfaces/templates';

function slugify(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)+/g, '');
}

interface EndpointRouteState {
  mode: 'domain' | 'port' | 'internal';
  domainType: 'free' | 'custom';
  subdomain: string;
  customDomain: string;
}

export function InstallDomainsRouting({
  appId,
  endpoints = [],
  defaultPort,
  domain: _domain,
  onDomainChange,
  appName,
  hostPort,
  onHostPortChange,
  disabled = false,
}: {
  appId?: string;
  endpoints?: OneClickEndpoint[];
  defaultPort: number;
  domain: string;
  onDomainChange: (domain: string) => void;
  appName: string;
  hostPort: string;
  onHostPortChange: (port: string) => void;
  disabled?: boolean;
}) {
  const { data: publicRes } = useGetPublicSettings();
  const { data: settingsRes } = useGetSettings();
  const isCloud = Boolean(publicRes?.data?.cloudMode);
  const serverIp = settingsRes?.data?.publicIpv4 || settingsRes?.data?.traefikWildcardIp || '';
  const cleanIp = serverIp ? serverIp.replace(/\./g, '-') : '';
  const sslipSuffix = cleanIp ? `${cleanIp}.sslip.io` : 'sslip.io';
  const configuredWildcard = settingsRes?.data?.defaultWildcardDomain?.trim();
  const wildcardDomain = isCloud ? 'codedock.run' : configuredWildcard || sslipSuffix;

  const resolvedEndpoints: OneClickEndpoint[] =
    endpoints.length > 0
      ? endpoints
      : appId === 'minio'
        ? [
            { label: 'Console', port: 9001, kind: 'http' },
            { label: 'S3 API', port: 9000, kind: 'http' },
          ]
        : [{ label: appId || slugify(appName) || 'web', port: defaultPort || 3000, kind: 'http' }];

  const buildDefaultSubdomain = (ep: OneClickEndpoint) => {
    const appSlug = slugify(appName) || 'app';
    const epSlug = slugify(ep.label) || 'web';
    return `${appSlug}-${epSlug}`;
  };

  const [routes, setRoutes] = useState<Record<number, EndpointRouteState>>(() => {
    const initial: Record<number, EndpointRouteState> = {};
    for (const ep of resolvedEndpoints) {
      initial[ep.port] = {
        mode: 'domain',
        domainType: 'free',
        subdomain: buildDefaultSubdomain(ep),
        customDomain: '',
      };
    }
    return initial;
  });

  useEffect(() => {
    setRoutes((prev) => {
      const next = { ...prev };
      for (const ep of resolvedEndpoints) {
        if (!next[ep.port]) {
          next[ep.port] = {
            mode: 'domain',
            domainType: 'free',
            subdomain: buildDefaultSubdomain(ep),
            customDomain: '',
          };
        }
      }
      return next;
    });
  }, [appName, resolvedEndpoints]);

  useEffect(() => {
    const primary = resolvedEndpoints[0];
    if (!primary) return;
    const st = routes[primary.port];
    if (st?.mode !== 'domain') {
      onDomainChange('');
      return;
    }
    if (st.domainType === 'free') {
      const sub = st.subdomain || buildDefaultSubdomain(primary);
      onDomainChange(`${sub}.${wildcardDomain}`);
    } else {
      onDomainChange(st.customDomain);
    }
  }, [routes, resolvedEndpoints, appName, wildcardDomain, onDomainChange]);

  const updateRoute = (port: number, patch: Partial<EndpointRouteState>) => {
    setRoutes((prev) => ({
      ...prev,
      [port]: { ...prev[port], ...patch },
    }));
  };

  return (
    <div className="space-y-4 rounded-2xl border border-border/50 bg-card p-5 shadow-xs">
      <h2 className="font-semibold text-foreground text-sm">Domains & routing</h2>

      <div className="grid grid-cols-1 items-start gap-4 sm:grid-cols-2">
        {resolvedEndpoints.map((ep) => {
          const st = routes[ep.port] ?? {
            mode: 'domain',
            domainType: 'free',
            subdomain: buildDefaultSubdomain(ep),
            customDomain: '',
          };

          return (
            <div
              key={ep.port}
              className="space-y-3 rounded-2xl border border-border/50 bg-muted/20 p-4"
            >
              <div className="flex items-center justify-between gap-3">
                <span className="truncate font-semibold text-foreground text-sm">
                  {ep.label.toLowerCase()}
                </span>
                <span className="shrink-0 rounded-lg border border-border/50 bg-muted/60 px-2.5 py-0.5 font-mono text-muted-foreground text-xs">
                  Web {ep.port}
                </span>
              </div>

              <Select
                value={st.mode}
                onValueChange={(val) =>
                  updateRoute(ep.port, { mode: val as EndpointRouteState['mode'] })
                }
                disabled={disabled}
              >
                <SelectTrigger className="h-9 w-full rounded-xl border-border/60 bg-background/70 px-3 text-sm">
                  <SelectValue placeholder="Domain" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="domain">Domain</SelectItem>
                  <SelectItem value="port">Port only</SelectItem>
                  <SelectItem value="internal">Internal</SelectItem>
                </SelectContent>
              </Select>

              {st.mode === 'domain' && (
                <div className="space-y-3">
                  <div className="inline-flex items-center gap-1.5">
                    <button
                      type="button"
                      disabled={disabled}
                      onClick={() => updateRoute(ep.port, { domainType: 'free' })}
                      className={`rounded-xl px-3.5 py-1.5 font-medium text-xs transition-colors ${
                        st.domainType === 'free'
                          ? 'border border-border/70 bg-muted/80 text-foreground shadow-2xs'
                          : 'text-muted-foreground hover:bg-muted/40 hover:text-foreground'
                      }`}
                    >
                      Free
                    </button>
                    <button
                      type="button"
                      disabled={disabled}
                      onClick={() => updateRoute(ep.port, { domainType: 'custom' })}
                      className={`rounded-xl px-3.5 py-1.5 font-medium text-xs transition-colors ${
                        st.domainType === 'custom'
                          ? 'border border-border/70 bg-muted/80 text-foreground shadow-2xs'
                          : 'text-muted-foreground hover:bg-muted/40 hover:text-foreground'
                      }`}
                    >
                      Custom
                    </button>
                  </div>

                  {st.domainType === 'free' ? (
                    <div className="flex h-10 items-center rounded-xl border border-border/60 bg-background/80 px-3.5">
                      <Input
                        value={st.subdomain}
                        disabled={disabled}
                        onChange={(e) =>
                          updateRoute(ep.port, { subdomain: slugify(e.target.value) })
                        }
                        placeholder="slug"
                        className="h-7 border-0 bg-transparent p-0 font-mono text-xs shadow-none focus-visible:ring-0"
                      />
                      <span className="shrink-0 pl-2 font-mono text-muted-foreground text-xs">
                        .{wildcardDomain}
                      </span>
                    </div>
                  ) : (
                    <Input
                      value={st.customDomain}
                      disabled={disabled}
                      onChange={(e) => updateRoute(ep.port, { customDomain: e.target.value })}
                      placeholder="app.example.com"
                      className="h-10 rounded-xl border-border/60 bg-background/80 px-3.5 font-mono text-xs"
                    />
                  )}
                </div>
              )}

              {st.mode === 'port' && (
                <div className="space-y-1.5">
                  <Input
                    value={hostPort}
                    disabled={disabled}
                    onChange={(e) => onHostPortChange(e.target.value)}
                    placeholder={String(ep.port)}
                    inputMode="numeric"
                    className="h-10 rounded-xl border-border/60 bg-background/80 px-3.5 font-mono text-xs"
                  />
                  <p className="text-[11px] text-muted-foreground">
                    Bound directly to host port {hostPort || ep.port}.
                  </p>
                </div>
              )}

              {st.mode === 'internal' && (
                <p className="text-[11px] text-muted-foreground">
                  Internal network traffic only. No external port exposed.
                </p>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
