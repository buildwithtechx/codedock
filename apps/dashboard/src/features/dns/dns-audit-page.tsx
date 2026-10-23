import { AlertTriangle, CheckCircle2, Globe } from 'lucide-react';
import { PageHeader } from '#/components/layout/page-header';
import { QueryErrorState } from '#/components/query-error-state';
import { DomainAuditTable } from './components/domain-audit-table';
import { useListAllDomains } from '#/hooks/use-dns';

export function DnsAuditPage() {
  const { data: domainsRes, isLoading, isError, refetch } = useListAllDomains();

  const domains = domainsRes?.data || [];
  const hasDomains = domains.length > 0;
  const showStatCards = !isError && (isLoading || hasDomains);

  const totalCount = domains.length;
  const provisionedCount = domains.filter((d) => {
    const status = (d.dnsProvisionStatus || '').toLowerCase();
    return status === 'provisioned';
  }).length;
  const manualCount = domains.filter(
    (d) =>
      (d.dnsProvisionStatus || '').toLowerCase() === 'pending' ||
      (d.dnsProvisionStatus || '').toLowerCase() === 'failed' ||
      !d.dnsProvisionStatus
  ).length;

  return (
    <div className="space-y-6 pb-12">
      <PageHeader
        title="Domains"
        description="Audit configured domains, DNS provisioning, and live verification across services."
      />

      {showStatCards && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div className="flex items-center gap-4 rounded-2xl border border-border/50 bg-card p-5 shadow-sm transition-colors hover:border-border/80">
            <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
              <Globe className="h-5 w-5" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                Total Domains
              </p>
              <h3 className="mt-1 font-bold text-2xl text-foreground tracking-tight">
                {isLoading || isError ? '...' : totalCount}
              </h3>
            </div>
          </div>

          <div className="flex items-center gap-4 rounded-2xl border border-border/50 bg-card p-5 shadow-sm transition-colors hover:border-border/80">
            <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-emerald-500/10 text-emerald-500">
              <CheckCircle2 className="h-5 w-5" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                Auto-Provisioned
              </p>
              <h3 className="mt-1 font-bold text-2xl text-foreground tracking-tight">
                {isLoading || isError ? '...' : provisionedCount}
              </h3>
            </div>
          </div>

          <div className="flex items-center gap-4 rounded-2xl border border-border/50 bg-card p-5 shadow-sm transition-colors hover:border-border/80">
            <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-amber-500/10 text-amber-500">
              <AlertTriangle className="h-5 w-5" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                Pending / Action Needed
              </p>
              <h3 className="mt-1 font-bold text-2xl text-foreground tracking-tight">
                {isLoading || isError ? '...' : manualCount}
              </h3>
            </div>
          </div>
        </div>
      )}

      <section aria-label="Domains overview table">
        {isError ? (
          <QueryErrorState
            title="Domain audit is unavailable"
            description="Codedock could not load domains for the active workspace."
            onRetry={() => void refetch()}
          />
        ) : (
          <DomainAuditTable domains={domains} isLoading={isLoading} />
        )}
      </section>
    </div>
  );
}
