import { CheckCircle2, Globe, Loader2, RefreshCw, Search, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useDomainVerification } from '#/features/dns/domain-verification';
import { useDeleteDomain } from '#/features/dns/hooks';
import type { DomainConfig } from '#/features/projects';
import { DnsStatusBadge } from '#/features/services/dns-status-badge';

interface DomainAuditTableProps {
  domains: DomainConfig[];
  isLoading: boolean;
}

export function DomainAuditTable({ domains, isLoading }: DomainAuditTableProps) {
  const deleteDomain = useDeleteDomain();
  const { verificationResults, verifyingMap, isVerifyingAll, verifyAll, verifyOne } =
    useDomainVerification();

  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const hasActiveVerification = Object.values(verifyingMap).some(Boolean);

  const handleVerifyOne = async (domainId: string) => {
    await verifyOne(domainId);
  };

  const handleVerifyAll = async () => {
    if (domains.length === 0) return;
    const { verifiedCount, failedCount, totalCount } = await verifyAll(domains);
    if (failedCount > 0) {
      toast.warning(
        `Completed verification: ${verifiedCount}/${totalCount} domains resolving, ${failedCount} failed`
      );
    } else {
      toast.success(`Completed verification: ${verifiedCount}/${totalCount} domains resolving`);
    }
  };

  const filteredDomains = domains.filter((d) => {
    const matchesSearch = d.domainName.toLowerCase().includes(search.toLowerCase());
    const matchesStatus =
      statusFilter === 'all' ||
      (d.dnsProvisionStatus || 'pending').toLowerCase() === statusFilter.toLowerCase() ||
      (statusFilter === 'provisioned' &&
        (d.dnsProvisionStatus || 'pending').toLowerCase() === 'provisioned');
    return matchesSearch && matchesStatus;
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        <Loader2 className="mr-2 h-5 w-5 animate-spin" /> Loading configured domains...
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex flex-1 items-center gap-3">
          <div className="relative max-w-sm flex-1">
            <Search className="absolute top-2.5 left-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search domain..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9 text-sm"
            />
          </div>
          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-36 text-xs">
              <SelectValue placeholder="All Statuses" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Statuses</SelectItem>
              <SelectItem value="provisioned">Provisioned</SelectItem>
              <SelectItem value="manual">Manual</SelectItem>
              <SelectItem value="pending">Pending</SelectItem>
              <SelectItem value="failed">Failed</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <Button
          variant="outline"
          size="sm"
          onClick={handleVerifyAll}
          disabled={isVerifyingAll || hasActiveVerification || domains.length === 0}
          className="h-9 gap-2 font-medium text-xs"
        >
          {isVerifyingAll ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <RefreshCw className="h-4 w-4" />
          )}
          Verify All Domains
        </Button>
      </div>

      {filteredDomains.length === 0 ? (
        <div className="rounded-2xl border border-dashed p-8 text-center text-muted-foreground text-sm">
          <Globe className="mx-auto mb-2 h-8 w-8 opacity-40" />
          No domains found matching criteria.
        </div>
      ) : (
        <div className="overflow-x-auto rounded-2xl border border-border/60 bg-card/40">
          <table className="w-full min-w-[760px] whitespace-nowrap text-left text-sm">
            <thead className="border-b bg-muted/40 font-semibold text-muted-foreground text-xs uppercase tracking-wider">
              <tr>
                <th className="px-4 py-3">Domain</th>
                <th className="px-4 py-3">Service</th>
                <th className="px-4 py-3">Provision Status</th>
                <th className="px-4 py-3">DNS Verification</th>
                <th className="px-4 py-3 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border/40">
              {filteredDomains.map((domain) => {
                const result = verificationResults[domain.id];
                const isVerifying = verifyingMap[domain.id];

                return (
                  <tr key={domain.id} className="transition-colors hover:bg-card/60">
                    <td className="px-4 py-3.5 font-medium font-mono">{domain.domainName}</td>
                    <td className="px-4 py-3.5 text-muted-foreground text-xs">
                      {domain.serviceId ? (
                        <span className="font-mono text-xs">{domain.serviceId.slice(0, 8)}...</span>
                      ) : (
                        '—'
                      )}
                    </td>
                    <td className="px-4 py-3.5">
                      <DnsStatusBadge status={domain.dnsProvisionStatus} />
                    </td>
                    <td className="px-4 py-3.5 text-xs">
                      {result ? (
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
                      ) : (
                        <span className="text-muted-foreground">Not verified</span>
                      )}
                    </td>
                    <td className="px-4 py-3.5 text-right">
                      <div className="flex items-center justify-end gap-2">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleVerifyOne(domain.id)}
                          disabled={isVerifying || isVerifyingAll}
                          className="h-8 gap-1 text-xs"
                        >
                          {isVerifying ? (
                            <Loader2 className="h-3.5 w-3.5 animate-spin" />
                          ) : (
                            <CheckCircle2 className="h-3.5 w-3.5" />
                          )}
                          Verify
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => {
                            if (!window.confirm(`Delete domain ${domain.domainName}?`)) return;
                            deleteDomain.mutate(
                              { id: domain.id },
                              {
                                onSuccess: () => toast.success('Domain deleted successfully'),
                                onError: () => toast.error('Failed to delete domain'),
                              }
                            );
                          }}
                          disabled={deleteDomain.isPending}
                          className="h-8 w-8 text-destructive hover:bg-destructive/10 hover:text-destructive"
                          aria-label={`Delete domain ${domain.domainName}`}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
