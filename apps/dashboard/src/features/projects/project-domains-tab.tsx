import {
  ArrowUpRight,
  CheckCircle2,
  Globe,
  Loader2,
  Plus,
  RefreshCw,
  Trash2,
  XCircle,
} from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useCreateDomain, useDeleteDomain, useVerifyDomain } from '#/features/dns/hooks';
import { useGetSettings } from '#/features/settings';
import { useListByProject as useListAppsByProject } from '#/hooks/use-apps';
import { useListProjectDomains } from './hooks';

export function ProjectDomainsTab({ projectId }: { projectId: string }) {
  const {
    data: domainsRes,
    isLoading: domainsLoading,
    refetch: refetchDomains,
  } = useListProjectDomains(projectId);
  const { data: appsRes, isLoading: appsLoading } = useListAppsByProject(
    projectId,
    undefined,
    true
  );
  const { data: settingsRes } = useGetSettings();
  const createDomain = useCreateDomain();
  const deleteDomain = useDeleteDomain();
  const verifyDomain = useVerifyDomain();

  const [domainName, setDomainName] = useState('');
  const [selectedServiceId, setSelectedServiceId] = useState('');
  const [verifyingId, setVerifyingId] = useState<string | null>(null);

  const domains = domainsRes?.data || [];
  const services = appsRes?.data || [];
  const serverIp = settingsRes?.data?.publicIpv4 || '';

  const activeServiceId = selectedServiceId || (services[0]?.id ?? '');

  const handleCreate = async () => {
    if (!domainName.trim()) {
      toast.error('Domain name is required');
      return;
    }
    if (!activeServiceId) {
      toast.error('Please select a target service');
      return;
    }
    try {
      await createDomain.mutateAsync({
        serviceId: activeServiceId,
        payload: { domainName: domainName.trim() },
      });
      setDomainName('');
      toast.success('Domain mapped successfully');
      await refetchDomains();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to add domain');
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteDomain.mutateAsync({ id });
      toast.success('Domain removed');
      await refetchDomains();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to delete domain');
    }
  };

  const handleVerify = async (id: string) => {
    setVerifyingId(id);
    try {
      const res = await verifyDomain.mutateAsync(id);
      if (res.data?.verified) {
        toast.success(res.data.message || 'Domain resolves to this server');
      } else {
        toast.warning(res.data?.message || 'Domain is not pointing to this server yet');
      }
      await refetchDomains();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to verify domain');
    } finally {
      setVerifyingId(null);
    }
  };

  const serviceNameFor = (serviceId: string) => {
    const s = services.find((srv) => srv.id === serviceId);
    return s ? s.name : 'Unknown Service';
  };

  return (
    <div className="space-y-6">
      <Card className="rounded-2xl border border-border/60 bg-card p-6 shadow-xs">
        <CardHeader className="p-0 pb-4">
          <div className="flex items-center gap-3">
            <div className="flex size-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
              <Globe className="size-4.5" />
            </div>
            <div>
              <CardTitle className="font-semibold text-base text-foreground">
                Custom Domains & Routing
              </CardTitle>
              <p className="mt-0.5 text-muted-foreground text-xs">
                Map custom hostnames to services in this project with automated TLS certificates.
              </p>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4 p-0 pt-2">
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-12 sm:items-end">
            <div className="sm:col-span-6">
              <Label className="text-xs">Domain Name</Label>
              <Input
                placeholder="app.example.com"
                value={domainName}
                onChange={(e) => setDomainName(e.target.value)}
                className="mt-1.5 font-mono text-xs"
              />
            </div>
            <div className="sm:col-span-4">
              <Label className="text-xs">Target Service</Label>
              <Select
                value={activeServiceId}
                onValueChange={setSelectedServiceId}
                disabled={appsLoading || services.length === 0}
              >
                <SelectTrigger className="mt-1.5 text-xs">
                  <SelectValue
                    placeholder={appsLoading ? 'Loading services...' : 'Select service'}
                  />
                </SelectTrigger>
                <SelectContent>
                  {services.map((srv) => (
                    <SelectItem key={srv.id} value={srv.id} className="text-xs">
                      {srv.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="sm:col-span-2">
              <Button
                type="button"
                onClick={handleCreate}
                disabled={createDomain.isPending || !domainName.trim() || !activeServiceId}
                className="w-full gap-1.5 text-xs"
              >
                {createDomain.isPending ? (
                  <Loader2 className="size-3.5 animate-spin" />
                ) : (
                  <Plus className="size-3.5" />
                )}
                Add Domain
              </Button>
            </div>
          </div>

          {serverIp && (
            <div className="rounded-xl border border-border/50 bg-muted/20 p-3 text-xs">
              <p className="font-medium text-foreground">DNS Configuration Record:</p>
              <p className="mt-1 font-mono text-muted-foreground">
                Point an <span className="text-foreground">A</span> record to{' '}
                <span className="font-bold text-foreground">{serverIp}</span> or a{' '}
                <span className="text-foreground">CNAME</span> to your server host.
              </p>
            </div>
          )}
        </CardContent>
      </Card>

      <Card className="rounded-2xl border border-border/60 bg-card p-6 shadow-xs">
        <CardHeader className="p-0 pb-4">
          <div className="flex items-center justify-between">
            <CardTitle className="font-semibold text-foreground text-sm">
              Configured Domains ({domains.length})
            </CardTitle>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => refetchDomains()}
              className="h-8 gap-1.5 text-muted-foreground text-xs hover:text-foreground"
            >
              <RefreshCw className="size-3.5" />
              Refresh
            </Button>
          </div>
        </CardHeader>
        <CardContent className="p-0 pt-2">
          {domainsLoading ? (
            <div className="flex h-32 items-center justify-center">
              <Loader2 className="size-6 animate-spin text-muted-foreground" />
            </div>
          ) : domains.length === 0 ? (
            <div className="flex min-h-40 flex-col items-center justify-center rounded-xl border border-border/70 border-dashed p-6 text-center">
              <Globe className="size-8 text-muted-foreground/50" />
              <p className="mt-2 font-medium text-foreground text-sm">No domains mapped</p>
              <p className="mt-1 max-w-sm text-muted-foreground text-xs">
                Add a custom domain above to expose your project services to the internet with TLS.
              </p>
            </div>
          ) : (
            <div className="divide-y divide-border/40 overflow-hidden rounded-xl border border-border/50">
              {domains.map((dom) => {
                const isVerifying = verifyingId === dom.id;
                const status = (dom.dnsProvisionStatus || '').toLowerCase();
                const isProvisioned = status === 'provisioned';

                return (
                  <div
                    key={dom.id}
                    className="flex flex-col gap-3 p-4 sm:flex-row sm:items-center sm:justify-between"
                  >
                    <div className="min-w-0 space-y-1">
                      <div className="flex items-center gap-2">
                        <a
                          href={`https://${dom.domainName}`}
                          target="_blank"
                          rel="noreferrer"
                          className="inline-flex items-center gap-1 font-medium font-mono text-foreground text-sm hover:text-primary hover:underline"
                        >
                          {dom.domainName}
                          <ArrowUpRight className="size-3 text-muted-foreground" />
                        </a>
                        {isProvisioned ? (
                          <span className="inline-flex items-center gap-1 rounded-md bg-emerald-500/10 px-2 py-0.5 font-medium text-[11px] text-emerald-500">
                            <CheckCircle2 className="size-3" />
                            Active
                          </span>
                        ) : (
                          <span className="inline-flex items-center gap-1 rounded-md bg-amber-500/10 px-2 py-0.5 font-medium text-[11px] text-amber-500">
                            <XCircle className="size-3" />
                            {status || 'Pending'}
                          </span>
                        )}
                      </div>
                      <p className="text-muted-foreground text-xs">
                        Routes to:{' '}
                        <span className="font-medium text-foreground">
                          {serviceNameFor(dom.serviceId)}
                        </span>
                      </p>
                    </div>

                    <div className="flex items-center gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleVerify(dom.id)}
                        disabled={isVerifying}
                        className="h-8 gap-1.5 text-xs"
                      >
                        {isVerifying ? (
                          <Loader2 className="size-3 animate-spin" />
                        ) : (
                          <RefreshCw className="size-3" />
                        )}
                        Verify DNS
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDelete(dom.id)}
                        disabled={deleteDomain.isPending}
                        className="size-8 h-8 p-0 text-muted-foreground hover:text-destructive"
                      >
                        <Trash2 className="size-3.5" />
                      </Button>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
