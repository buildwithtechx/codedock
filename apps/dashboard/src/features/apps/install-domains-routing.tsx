import { Globe, ShieldCheck } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Badge } from '#/components/ui/badge';
import { Input } from '#/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useGetPublicSettings, useGetSettings } from '#/features/settings';
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
  disabled = false,
}: {
  appId?: string;
  endpoints?: OneClickEndpoint[];
  defaultPort: number;
  domain: string;
  onDomainChange: (domain: string) => void;
  appName: string;
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
        : [{ label: appName.toLowerCase() || 'web', port: defaultPort || 3000, kind: 'http' }];

  const [routes, setRoutes] = useState<Record<number, EndpointRouteState>>(() => {
    const initial: Record<number, EndpointRouteState> = {};
    for (const ep of resolvedEndpoints) {
      const epSlug =
        resolvedEndpoints.length > 1
          ? `${slugify(appName)}-${slugify(ep.label)}`
          : slugify(appName) || 'app';
      initial[ep.port] = {
        mode: 'domain',
        domainType: 'free',
        subdomain: epSlug,
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
          const epSlug =
            resolvedEndpoints.length > 1
              ? `${slugify(appName)}-${slugify(ep.label)}`
              : slugify(appName) || 'app';
          next[ep.port] = {
            mode: 'domain',
            domainType: 'free',
            subdomain: epSlug,
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
      const sub = st.subdomain || slugify(appName) || 'app';
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
    <div className="space-y-4 rounded-2xl bg-card p-5">
      <div className="flex items-center justify-between">
        <h2 className="font-semibold text-foreground text-sm">Domains & routing</h2>
        <Badge variant="outline" className="gap-1 font-mono text-[10px]">
          <Globe className="size-3" />
          HTTP / HTTPS
        </Badge>
      </div>

      <div
        className={
          resolvedEndpoints.length > 1
            ? 'grid grid-cols-1 items-start gap-4 sm:grid-cols-2'
            : 'space-y-4'
        }
      >
        {resolvedEndpoints.map((ep) => {
          const st = routes[ep.port] ?? {
            mode: 'domain',
            domainType: 'free',
            subdomain:
              resolvedEndpoints.length > 1
                ? `${slugify(appName)}-${slugify(ep.label)}`
                : slugify(appName) || 'app',
            customDomain: '',
          };

          return (
            <div
              key={ep.port}
              className="space-y-3 rounded-xl border border-border/40 bg-muted/20 p-4"
            >
              <div className="flex items-center justify-between gap-3">
                <span className="truncate font-medium text-foreground text-sm">{ep.label}</span>
                <span className="shrink-0 rounded-md bg-muted/60 px-2 py-0.5 font-mono text-muted-foreground text-xs">
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
                <SelectTrigger className="w-full bg-background text-xs">
                  <SelectValue placeholder="HTTP (domain)" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="domain">HTTP (domain)</SelectItem>
                  <SelectItem value="port">Port only</SelectItem>
                  <SelectItem value="internal">Internal</SelectItem>
                </SelectContent>
              </Select>

              {st.mode === 'domain' && (
                <div className="space-y-2">
                  <div className="flex rounded-lg bg-muted/50 p-0.5">
                    <button
                      type="button"
                      disabled={disabled}
                      onClick={() => updateRoute(ep.port, { domainType: 'free' })}
                      className={`flex-1 rounded-md py-1 font-medium text-xs transition-colors ${
                        st.domainType === 'free'
                          ? 'bg-card text-foreground shadow-sm'
                          : 'text-muted-foreground hover:text-foreground'
                      }`}
                    >
                      Free
                    </button>
                    <button
                      type="button"
                      disabled={disabled}
                      onClick={() => updateRoute(ep.port, { domainType: 'custom' })}
                      className={`flex-1 rounded-md py-1 font-medium text-xs transition-colors ${
                        st.domainType === 'custom'
                          ? 'bg-card text-foreground shadow-sm'
                          : 'text-muted-foreground hover:text-foreground'
                      }`}
                    >
                      Custom
                    </button>
                  </div>

                  {st.domainType === 'free' ? (
                    <div className="space-y-1.5">
                      <div className="flex items-center rounded-lg border border-border/60 bg-background px-3 py-1.5">
                        <Input
                          value={st.subdomain}
                          disabled={disabled}
                          onChange={(e) =>
                            updateRoute(ep.port, { subdomain: slugify(e.target.value) })
                          }
                          placeholder="slug"
                          className="h-6 border-0 p-0 font-mono text-xs shadow-none focus-visible:ring-0"
                        />
                        <span className="shrink-0 font-mono text-muted-foreground text-xs">
                          .{wildcardDomain}
                        </span>
                      </div>
                      <div className="flex items-center gap-1.5 text-[11px] text-emerald-600 dark:text-emerald-400">
                        <ShieldCheck className="size-3 shrink-0" />
                        <span className="truncate">
                          https://{st.subdomain || 'app'}.{wildcardDomain}
                        </span>
                      </div>
                    </div>
                  ) : (
                    <div className="space-y-1">
                      <Input
                        value={st.customDomain}
                        disabled={disabled}
                        onChange={(e) => updateRoute(ep.port, { customDomain: e.target.value })}
                        placeholder="app.example.com"
                        className="h-8 font-mono text-xs"
                      />
                      <p className="text-[11px] text-muted-foreground">
                        Point a CNAME record to your host IP.
                      </p>
                    </div>
                  )}
                </div>
              )}

              {st.mode === 'port' && (
                <p className="text-[11px] text-muted-foreground">
                  Bound directly to host port {ep.port} without domain routing.
                </p>
              )}

              {st.mode === 'internal' && (
                <p className="text-[11px] text-muted-foreground">
                  Internal east-west traffic only. No external port exposed.
                </p>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
