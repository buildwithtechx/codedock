import {
  ArrowUpRight,
  CheckCircle2,
  ChevronDown,
  ChevronUp,
  Globe,
  Loader2,
  RefreshCw,
  Trash2,
  XCircle,
} from 'lucide-react';
import { useMemo, useState } from 'react';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import type { DomainBase } from '@codedock/contracts';
import { RouteProtectionButton } from '#/features/dns';

interface ProjectDomainsListProps {
  domains: DomainBase[];
  projectId: string;
  verifyingId: string | null;
  onVerify: (id: string) => void;
  onDelete: (id: string) => void;
  deletePending?: boolean;
  serviceNameFor: (id?: string) => string;
}

export function ProjectDomainsList({
  domains,
  projectId,
  verifyingId,
  onVerify,
  onDelete,
  deletePending,
  serviceNameFor,
}: ProjectDomainsListProps) {
  const [expandedGroups, setExpandedGroups] = useState<Record<string, boolean>>({});

  const grouped = useMemo(() => {
    const map = new Map<string, DomainBase[]>();
    for (const dom of domains) {
      const key = dom.serviceId || 'unassigned';
      if (!map.has(key)) map.set(key, []);
      map.get(key)!.push(dom);
    }
    return map;
  }, [domains]);

  const toggleGroup = (key: string) => {
    setExpandedGroups((prev) => ({ ...prev, [key]: !prev[key] }));
  };

  const renderDomainRow = (dom: DomainBase) => {
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
            <span className="font-medium text-foreground">{serviceNameFor(dom.serviceId)}</span>
          </p>
        </div>

        <div className="flex items-center gap-2">
          <RouteProtectionButton hostname={dom.domainName} projectId={projectId} />
          <Button
            variant="outline"
            size="sm"
            onClick={() => onVerify(dom.id)}
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
            onClick={() => onDelete(dom.id)}
            disabled={deletePending}
            className="size-8 h-8 p-0 text-muted-foreground hover:text-destructive"
          >
            <Trash2 className="size-3.5" />
          </Button>
        </div>
      </div>
    );
  };

  if (domains.length === 0) {
    return (
      <div className="flex min-h-40 flex-col items-center justify-center rounded-xl border border-border/70 border-dashed p-6 text-center">
        <Globe className="size-8 text-muted-foreground/50" />
        <p className="mt-2 font-medium text-foreground text-sm">No domains mapped</p>
        <p className="mt-1 max-w-sm text-muted-foreground text-xs">
          Add a custom domain above to expose your project services to the internet with TLS.
        </p>
      </div>
    );
  }

  if (grouped.size <= 1) {
    return (
      <div className="divide-y divide-border/40 overflow-hidden rounded-xl border border-border/50">
        {domains.map(renderDomainRow)}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {Array.from(grouped.entries()).map(([serviceId, groupDomains]) => {
        const isExpanded = expandedGroups[serviceId];
        const displayDomains = isExpanded ? groupDomains : groupDomains.slice(0, 3);
        const hasMore = groupDomains.length > 3;

        return (
          <div
            key={serviceId}
            className="overflow-hidden rounded-xl border border-border/50 bg-background/50"
          >
            <div className="flex items-center justify-between border-b border-border/40 bg-muted/20 px-4 py-2.5">
              <div className="flex items-center gap-2">
                <span className="font-semibold text-foreground text-xs">
                  {serviceNameFor(serviceId === 'unassigned' ? undefined : serviceId)}
                </span>
                <Badge variant="outline" className="text-[10px]">
                  {groupDomains.length} {groupDomains.length === 1 ? 'domain' : 'domains'}
                </Badge>
              </div>

              {hasMore && (
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => toggleGroup(serviceId)}
                  className="h-6 gap-1 px-2 text-[11px] text-muted-foreground hover:text-foreground"
                >
                  {isExpanded ? (
                    <>
                      <span>Show less</span>
                      <ChevronUp className="size-3" />
                    </>
                  ) : (
                    <>
                      <span>View all {groupDomains.length}</span>
                      <ChevronDown className="size-3" />
                    </>
                  )}
                </Button>
              )}
            </div>

            <div className="divide-y divide-border/40">{displayDomains.map(renderDomainRow)}</div>
          </div>
        );
      })}
    </div>
  );
}
