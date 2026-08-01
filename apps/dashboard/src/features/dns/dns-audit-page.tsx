import { CheckCircle2, Globe, HelpCircle, Network, ShieldCheck } from 'lucide-react';
import { useState } from 'react';
import { Card, CardContent } from '#/components/ui/card';
import { DomainAuditTable } from './components/domain-audit-table';
import { DnsSettings } from './dns-settings';
import { DomainsPage } from './domains-page';
import { useListAllDomains } from './hooks';

export function DnsAuditPage() {
  const { data: domainsRes, isLoading } = useListAllDomains();
  const [activeTab, setActiveTab] = useState<'audit' | 'providers' | 'global'>('audit');

  const domains = domainsRes?.data || [];

  const totalCount = domains.length;
  const provisionedCount = domains.filter(
    (d) => (d.dnsProvisionStatus || '').toLowerCase() === 'provisioned'
  ).length;
  const manualCount = domains.filter(
    (d) =>
      (d.dnsProvisionStatus || '').toLowerCase() === 'manual' ||
      (d.dnsProvisionStatus || '').toLowerCase() === 'pending' ||
      !d.dnsProvisionStatus
  ).length;

  return (
    <div className="space-y-6 pb-12">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl border border-primary/20 bg-primary/10 text-primary">
            <Network className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl tracking-tight">DNS & Domain Audit</h1>
            <p className="text-muted-foreground text-sm">
              Overview and audit of all configured domains across your services, DNS provision
              status, and live verification.
            </p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <Card className="border/60 bg-card/40">
          <CardContent className="flex items-center gap-4 p-4">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <Globe className="h-5 w-5" />
            </div>
            <div>
              <p className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                Total Domains
              </p>
              <h3 className="mt-0.5 font-bold text-2xl">{isLoading ? '...' : totalCount}</h3>
            </div>
          </CardContent>
        </Card>

        <Card className="border/60 bg-card/40">
          <CardContent className="flex items-center gap-4 p-4">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-500">
              <CheckCircle2 className="h-5 w-5" />
            </div>
            <div>
              <p className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                Auto-Provisioned
              </p>
              <h3 className="mt-0.5 font-bold text-2xl">{isLoading ? '...' : provisionedCount}</h3>
            </div>
          </CardContent>
        </Card>

        <Card className="border/60 bg-card/40">
          <CardContent className="flex items-center gap-4 p-4">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-amber-500/10 text-amber-500">
              <HelpCircle className="h-5 w-5" />
            </div>
            <div>
              <p className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                Manual / Pending Setup
              </p>
              <h3 className="mt-0.5 font-bold text-2xl">{isLoading ? '...' : manualCount}</h3>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="flex border-border/60 border-b">
        <button
          type="button"
          onClick={() => setActiveTab('audit')}
          className={`flex items-center gap-2 border-b-2 px-4 py-2.5 font-medium text-sm transition-colors ${
            activeTab === 'audit'
              ? 'border-primary text-primary'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <ShieldCheck className="h-4 w-4" />
          Domain Overview & Audit
        </button>
        <button
          type="button"
          onClick={() => setActiveTab('providers')}
          className={`flex items-center gap-2 border-b-2 px-4 py-2.5 font-medium text-sm transition-colors ${
            activeTab === 'providers'
              ? 'border-primary text-primary'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <Network className="h-4 w-4" />
          DNS Provider API Credentials
        </button>
        <button
          type="button"
          onClick={() => setActiveTab('global')}
          className={`flex items-center gap-2 border-b-2 px-4 py-2.5 font-medium text-sm transition-colors ${
            activeTab === 'global'
              ? 'border-primary text-primary'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <Globe className="h-4 w-4" />
          Global Domain Settings
        </button>
      </div>

      {activeTab === 'audit' && <DomainAuditTable domains={domains} isLoading={isLoading} />}

      {activeTab === 'providers' && <DnsSettings />}

      {activeTab === 'global' && <DomainsPage />}
    </div>
  );
}
