import { CheckCircle2, Globe, HelpCircle } from 'lucide-react';
import { PageHeader } from '#/components/layout/page-header';
import { Card, CardContent } from '#/components/ui/card';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { DomainAuditTable } from './components/domain-audit-table';
import { useListAllDomains } from './hooks';

export function DnsAuditPage() {
  const { data: domainsRes, isLoading, isError, refetch } = useListAllDomains();

  const domains = domainsRes?.data || [];

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

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <Card className="border-border/60 bg-card">
          <CardContent className="flex items-center gap-4 p-4">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <Globe className="h-5 w-5" />
            </div>
            <div>
              <p className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                Total Domains
              </p>
              <h3 className="mt-0.5 font-bold text-2xl">
                {isLoading || isError ? '...' : totalCount}
              </h3>
            </div>
          </CardContent>
        </Card>

        <Card className="border-border/60 bg-card">
          <CardContent className="flex items-center gap-4 p-4">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-500">
              <CheckCircle2 className="h-5 w-5" />
            </div>
            <div>
              <p className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                Auto-Provisioned
              </p>
              <h3 className="mt-0.5 font-bold text-2xl">
                {isLoading || isError ? '...' : provisionedCount}
              </h3>
            </div>
          </CardContent>
        </Card>

        <Card className="border-border/60 bg-card">
          <CardContent className="flex items-center gap-4 p-4">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-amber-500/10 text-amber-500">
              <HelpCircle className="h-5 w-5" />
            </div>
            <div>
              <p className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                Manual / Pending / Failed Setup
              </p>
              <h3 className="mt-0.5 font-bold text-2xl">
                {isLoading || isError ? '...' : manualCount}
              </h3>
            </div>
          </CardContent>
        </Card>
      </div>

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
