import { AlertCircle, CheckCircle2, Clock, HelpCircle } from 'lucide-react';
import type { DNSProvisionStatus } from '#/features/dns';

interface DnsStatusBadgeProps {
  status?: DNSProvisionStatus | string;
}

export function DnsStatusBadge({ status = 'pending' }: DnsStatusBadgeProps) {
  const normalizedStatus = (status || 'pending').toLowerCase();

  switch (normalizedStatus) {
    case 'provisioned':
    case 'success':
      return (
        <span className="inline-flex items-center gap-1 rounded-full border border-emerald-500/20 bg-emerald-500/10 px-2 py-0.5 font-medium text-emerald-600 text-xs dark:text-emerald-400">
          <CheckCircle2 className="h-3 w-3" />
          Provisioned
        </span>
      );
    case 'failed':
      return (
        <span className="inline-flex items-center gap-1 rounded-full border border-destructive/20 bg-destructive/10 px-2 py-0.5 font-medium text-destructive text-xs">
          <AlertCircle className="h-3 w-3" />
          Failed
        </span>
      );
    case 'manual':
      return (
        <span className="inline-flex items-center gap-1 rounded-full border border-blue-500/20 bg-blue-500/10 px-2 py-0.5 font-medium text-blue-600 text-xs dark:text-blue-400">
          <HelpCircle className="h-3 w-3" />
          Manual
        </span>
      );
    case 'pending':
      return (
        <span className="inline-flex items-center gap-1 rounded-full border border-amber-500/20 bg-amber-500/10 px-2 py-0.5 font-medium text-amber-600 text-xs dark:text-amber-400">
          <Clock className="h-3 w-3" />
          Pending
        </span>
      );
    default:
      return (
        <span className="inline-flex items-center gap-1 rounded-full border border-muted bg-muted/40 px-2 py-0.5 font-medium text-muted-foreground text-xs">
          <HelpCircle className="h-3 w-3" />
          {status}
        </span>
      );
  }
}
