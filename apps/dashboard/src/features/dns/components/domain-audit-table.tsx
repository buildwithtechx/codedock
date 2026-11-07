import type { DomainConfig } from '@codedock/contracts';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';
import { CheckCircle2, Globe, Loader2, RefreshCw, Search, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { useDomainVerification } from '#/features/dns/domain-verification';
import { DnsStatusBadge } from '#/features/services/dns-status-badge';
import { useDeleteDomain } from '#/hooks/use-dns';
import { DomainEmptyState } from './domain-empty-state';

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
      <div className="flex items-center justify-center py-16 text-muted-foreground text-sm">
        <Loader2 className="mr-2 h-5 w-5 animate-spin text-primary" /> Loading configured domains...
      </div>
    );
  }

  if (domains.length === 0) {
    return <DomainEmptyState />;
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex flex-1 items-center gap-3">
          <div className="relative max-w-sm flex-1">
            <Search className="absolute top-2.5 left-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search domains..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9 text-sm"
            />
          </div>
          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-40 text-xs">
              <SelectValue placeholder="All Statuses" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Statuses</SelectItem>
              <SelectItem value="provisioned">Provisioned</SelectItem>
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
        <div className="flex flex-col items-center justify-center rounded-2xl border border-border/50 bg-card/40 py-16 text-center">
          <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-2xl bg-muted/60 text-muted-foreground">
            <Globe className="h-6 w-6" />
          </div>
          <h4 className="font-semibold text-foreground text-sm">No matching domains found</h4>
          <p className="mt-1 max-w-sm text-muted-foreground text-xs leading-relaxed">
            No domains matched your active search query or status filter.
          </p>
          <Button
            variant="outline"
            size="sm"
            onClick={() => {
              setSearch('');
              setStatusFilter('all');
            }}
            className="mt-4 h-8 font-medium text-xs"
          >
            Reset Filters
          </Button>
        </div>
      ) : (
        <div className="overflow-x-auto rounded-2xl border border-border/50 bg-card shadow-sm">
          <table className="w-full min-w-[760px] whitespace-nowrap text-left text-sm">
            <thead className="border-border/50 border-b bg-muted/30 font-semibold text-[11px] text-muted-foreground uppercase tracking-wider">
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
                  <tr key={domain.id} className="transition-colors hover:bg-muted/20">
                    <td className="px-4 py-3.5 font-medium font-mono text-sm">
                      {domain.domainName}
                    </td>
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
                          className="h-8 gap-1.5 text-xs"
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
