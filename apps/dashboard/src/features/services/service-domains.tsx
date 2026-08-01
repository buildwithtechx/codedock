import { Loader2, RefreshCw, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  type DomainVerifyResult,
  useCreate,
  useDelete,
  useListByService,
  useVerifyDomain,
} from '#/features/dns';
import { useGetSettings } from '#/features/settings';
import { DnsGuidanceBanner } from './dns-guidance-banner';
import { DnsStatusBadge } from './dns-status-badge';

export function ServiceDomains({ serviceId }: { serviceId: string }) {
  const { data: domainsRes, isLoading } = useListByService(serviceId);
  const { data: settingsRes } = useGetSettings();
  const createDomain = useCreate();
  const deleteDomain = useDelete();
  const verifyDomain = useVerifyDomain();

  const [newDomain, setNewDomain] = useState('');
  const [verifyingMap, setVerifyingMap] = useState<Record<string, boolean>>({});
  const [verificationResults, setVerificationResults] = useState<
    Record<string, DomainVerifyResult>
  >({});

  const settings = settingsRes?.data;
  const serverIp = settings?.publicIpv4 || '';
  const hasDnsProvider = Boolean(
    settings?.cloudflareApiToken || settings?.namecheapApiKey || settings?.spaceshipApiKey
  );

  const domains = domainsRes?.data || [];

  const showGuidanceBanner =
    !hasDnsProvider ||
    domains.some(
      (d) =>
        !d.dnsProvisionStatus ||
        d.dnsProvisionStatus === 'manual' ||
        d.dnsProvisionStatus === 'pending'
    );

  const handleCreate = () => {
    if (!newDomain.trim()) return;
    createDomain.mutate(
      { serviceId, payload: { domainName: newDomain } },
      {
        onSuccess: () => {
          setNewDomain('');
          toast.success('Domain added successfully');
        },
        onError: (err) => {
          toast.error(err.message || 'Failed to add domain');
        },
      }
    );
  };

  const handleVerify = async (domainId: string) => {
    setVerifyingMap((prev) => ({ ...prev, [domainId]: true }));
    try {
      const res = await verifyDomain.mutateAsync(domainId);
      const data = res.data;
      if (data) {
        setVerificationResults((prev) => ({ ...prev, [domainId]: data }));
        if (data.status === 'resolves_to_server') {
          toast.success(data.message || '✅ Resolves to server IP');
        } else if (data.status === 'resolves_to_different_ip') {
          toast.warning(data.message || '⚠️ Resolving to different IP');
        } else {
          toast.error(data.message || '❌ Unresolved');
        }
      }
    } catch (err: unknown) {
      const errorMsg = err instanceof Error ? err.message : 'Verification failed';
      toast.error(errorMsg);
    } finally {
      setVerifyingMap((prev) => ({ ...prev, [domainId]: false }));
    }
  };

  return (
    <Card className="divide-y divide-border/50">
      <CardHeader className="flex flex-row items-center justify-between pb-3">
        <CardTitle className="text-lg">Domains</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4 pt-4">
        {showGuidanceBanner && (
          <DnsGuidanceBanner serverIp={serverIp} hasDnsProvider={hasDnsProvider} />
        )}

        {isLoading ? (
          <div className="flex items-center justify-center py-6 text-muted-foreground text-sm">
            <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Loading domains...
          </div>
        ) : domains.length === 0 ? (
          <p className="py-4 text-center text-muted-foreground text-sm">
            No custom domains configured for this service.
          </p>
        ) : (
          <ul className="space-y-3">
            {domains.map((domain) => {
              const result = verificationResults[domain.id];
              const isVerifying = verifyingMap[domain.id];

              return (
                <li
                  key={domain.id}
                  className="flex flex-col gap-2 rounded-xl border border-border/60 bg-card/50 p-3 sm:flex-row sm:items-center sm:justify-between"
                >
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <span className="font-semibold text-sm">{domain.domainName}</span>
                      <DnsStatusBadge status={domain.dnsProvisionStatus} />
                    </div>
                    <div className="flex flex-wrap items-center gap-3 text-xs">
                      <span className="text-muted-foreground">
                        SSL: {domain.sslCertStatus || 'pending'}
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
                      className="h-8 gap-1.5 text-xs"
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
                      className="h-8 w-8 text-destructive hover:bg-destructive/10 hover:text-destructive"
                      onClick={() => deleteDomain.mutate({ id: domain.id })}
                      disabled={deleteDomain.isPending}
                      aria-label={`Delete domain ${domain.domainName}`}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </li>
              );
            })}
          </ul>
        )}

        <div className="mt-4 flex flex-col gap-2 border-t pt-4">
          <Label htmlFor="new-domain" className="font-medium text-xs">
            Add New Domain
          </Label>
          <div className="flex gap-2">
            <Input
              id="new-domain"
              placeholder="subdomain.example.com"
              value={newDomain}
              onChange={(e) => setNewDomain(e.target.value)}
              className="font-mono text-sm"
            />
            <Button onClick={handleCreate} disabled={createDomain.isPending || !newDomain.trim()}>
              Add Domain
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
