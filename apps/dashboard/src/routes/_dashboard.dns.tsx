import { createFileRoute } from '@tanstack/react-router';
import { DnsAuditPage } from '#/features/dns/dns-audit-page';

export const Route = createFileRoute('/_dashboard/dns')({
  component: DnsAuditPage,
});
