import { Globe, Loader2, Plus, RefreshCw, Shield, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useCreate, useDelete, useListByService } from '#/features/dns';
import { useDomainVerification } from '#/features/dns/domain-verification';
import { useGetSettings } from '#/features/settings';
import { useAuthStore } from '#/stores/auth-store';
import { DnsGuidanceBanner } from './dns-guidance-banner';
import { DnsStatusBadge } from './dns-status-badge';

export function ServiceDomains({ serviceId }: { serviceId: string }) {
  const userRole = useAuthStore((state) => state.user?.role);
  const canManageProviderDomains = userRole === 'admin' || String(userRole) === 'owner';
  const { data: domainsRes, isLoading } = useListByService(serviceId);
  const {
    data: settingsRes,
    isLoading: isLoadingSettings,
    error: settingsError,
  } = useGetSettings();
  const createDomain = useCreate();
  const deleteDomain = useDelete();
  const { verificationResults, verifyingMap, verifyOne } = useDomainVerification();

  const [domainType, setDomainType] = useState<'custom' | 'free'>('custom');
  const [newDomain, setNewDomain] = useState('');
  const [exposedPort, setExposedPort] = useState('80');
  const [includeWww, setIncludeWww] = useState(false);

  const settings = settingsRes?.data;
  const serverIp = settings?.publicIpv4 || '';
  const wildcardDomain = settings?.defaultWildcardDomain || 'codedock.app';
  const hasDnsProvider = Boolean(
    settings?.cloudflareApiToken ||
      (settings?.namecheapApiUser && settings?.namecheapApiKey) ||
      (settings?.spaceshipApiKey && settings?.spaceshipApiSecret)
  );

  const domains = domainsRes?.data || [];

  const showGuidanceBanner =
    !isLoadingSettings &&
    !settingsError &&
    (!hasDnsProvider ||
      domains.some(
        (d) =>
          !d.dnsProvisionStatus ||
          d.dnsProvisionStatus === 'pending' ||
          d.dnsProvisionStatus === 'failed'
      ));

  const handleCreate = () => {
    let hostname = newDomain.trim().toLowerCase();
    if (!hostname) return;

    if (domainType === 'free') {
      if (!hostname.includes('.')) {
        hostname = `${hostname}.${wildcardDomain}`;
      }
    }

    createDomain.mutate(
      { serviceId, payload: { domainName: hostname } },
      {
        onSuccess: () => {
          if (includeWww && domainType === 'custom' && !hostname.startsWith('www.')) {
            const wwwHost = `www.${hostname}`;
            createDomain.mutate(
              { serviceId, payload: { domainName: wwwHost } },
              {
                onSuccess: () => {
                  toast.success(`Domain ${hostname} and alias ${wwwHost} added successfully`);
                },
              }
            );
          } else {
            toast.success(`Domain ${hostname} added successfully`);
          }
          setNewDomain('');
          setIncludeWww(false);
        },
        onError: (err) => {
          toast.error(err.message || 'Failed to add domain');
        },
      }
    );
  };

  const handleVerify = async (domainId: string) => {
    await verifyOne(domainId);
  };

  return (
    <Card className="divide-y divide-border/50 rounded-2xl border-border/60 bg-card shadow-xs">
      <CardHeader className="flex flex-row items-center justify-between pb-3">
        <div className="space-y-0.5">
          <div className="flex items-center gap-2">
            <Globe className="size-4 text-primary" />
            <CardTitle className="font-semibold text-base">Service Domains & Routing</CardTitle>
          </div>
          <CardDescription className="text-xs">
            Configure custom domain aliases, SSL termination, and target container port mapping.
          </CardDescription>
        </div>
      </CardHeader>
      <CardContent className="space-y-5 pt-5">
        {showGuidanceBanner && (
          <DnsGuidanceBanner serverIp={serverIp} hasDnsProvider={hasDnsProvider} />
        )}

        {isLoading ? (
          <div className="flex items-center justify-center py-8 text-muted-foreground text-sm">
            <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Loading domains...
          </div>
        ) : domains.length === 0 ? (
          <div className="rounded-xl border border-border/70 border-dashed p-8 text-center text-muted-foreground text-sm">
            <Globe className="mx-auto mb-2 size-8 text-muted-foreground/40" />
            <p className="font-medium text-foreground text-sm">No custom domains configured</p>
            <p className="mt-1 text-xs">
              Add a custom domain or free subdomain to expose this service publicly over HTTPS.
            </p>
          </div>
        ) : (
          <div className="divide-y divide-border/40 overflow-hidden rounded-xl border border-border/50 bg-background/50">
            {domains.map((domain) => {
              const result = verificationResults[domain.id];
              const isVerifying = verifyingMap[domain.id];

              return (
                <div
                  key={domain.id}
                  className="flex flex-col gap-3 p-4 transition-colors hover:bg-muted/30 sm:flex-row sm:items-center sm:justify-between"
                >
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <span className="font-semibold text-foreground text-sm">
                        {domain.domainName}
                      </span>
                      <DnsStatusBadge status={domain.dnsProvisionStatus} />
                    </div>
                    <div className="flex flex-wrap items-center gap-3 text-xs">
                      <span className="flex items-center gap-1 text-muted-foreground">
                        <Shield className="size-3 text-emerald-500" />
                        SSL: {domain.sslCertStatus || 'active'}
                      </span>
                      {result && (
                        <span
                          className={`font-medium ${
                            result.status === 'resolves_to_server'
                              ? 'text-emerald-600 dark:text-emerald-400'
                              : result.status === 'resolves_to_different_ip'
                                ? 'text-amber-600 dark:text-amber-400'
                                : 'text-destructive'
                          }`}
                        >
                          {result.message}
                        </span>
                      )}
                    </div>
                  </div>

                  <div className="flex items-center gap-2 pt-2 sm:pt-0">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleVerify(domain.id)}
                      disabled={isVerifying}
                      className="h-8 gap-1.5 rounded-lg text-xs"
                    >
                      {isVerifying ? (
                        <Loader2 className="h-3.5 w-3.5 animate-spin" />
                      ) : (
                        <RefreshCw className="h-3.5 w-3.5" />
                      )}
                      Verify DNS
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8 rounded-lg text-destructive hover:bg-destructive/10 hover:text-destructive"
                      onClick={() => deleteDomain.mutate({ id: domain.id })}
                      disabled={deleteDomain.isPending}
                      aria-label={`Delete domain ${domain.domainName}`}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              );
            })}
          </div>
        )}

        {canManageProviderDomains && (
          <div className="space-y-4 rounded-xl border border-border/60 bg-muted/20 p-4">
            <div className="flex items-center justify-between">
              <Label className="font-semibold text-foreground text-xs">Add Domain Alias</Label>
              <div className="flex items-center gap-1 rounded-lg bg-muted p-0.5 text-xs">
                <button
                  type="button"
                  onClick={() => setDomainType('custom')}
                  className={`rounded-md px-2.5 py-1 font-medium text-xs transition-colors ${
                    domainType === 'custom'
                      ? 'bg-card text-foreground shadow-xs'
                      : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  Custom Domain
                </button>
                <button
                  type="button"
                  onClick={() => setDomainType('free')}
                  className={`rounded-md px-2.5 py-1 font-medium text-xs transition-colors ${
                    domainType === 'free'
                      ? 'bg-card text-foreground shadow-xs'
                      : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  Free Subdomain
                </button>
              </div>
            </div>

            <div className="grid gap-3 sm:grid-cols-[1fr_120px]">
              <div className="space-y-1">
                <Input
                  id="new-domain"
                  placeholder={
                    domainType === 'free'
                      ? `my-app (will use .${wildcardDomain})`
                      : 'subdomain.example.com'
                  }
                  value={newDomain}
                  onChange={(e) => setNewDomain(e.target.value)}
                  className="rounded-xl font-mono text-sm"
                />
              </div>

              <div className="space-y-1">
                <Input
                  id="exposed-port"
                  placeholder="Port (80)"
                  value={exposedPort}
                  onChange={(e) => setExposedPort(e.target.value)}
                  className="rounded-xl font-mono text-sm"
                />
              </div>
            </div>

            {domainType === 'custom' && (
              <label className="flex cursor-pointer select-none items-center gap-2 text-muted-foreground text-xs">
                <input
                  type="checkbox"
                  checked={includeWww}
                  onChange={(e) => setIncludeWww(e.target.checked)}
                  className="rounded border-border/70"
                />
                <span>Include www alias (creates both apex and www routes)</span>
              </label>
            )}

            <div className="flex justify-end pt-1">
              <Button
                onClick={handleCreate}
                disabled={createDomain.isPending || !newDomain.trim()}
                className="h-9 gap-1.5 rounded-xl text-xs"
              >
                <Plus className="size-3.5" />
                Add Domain
              </Button>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
