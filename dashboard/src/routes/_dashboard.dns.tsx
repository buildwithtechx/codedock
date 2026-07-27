import { createFileRoute } from '@tanstack/react-router';
import { DnsSettings } from '#/features/dns/dns-settings';

export const Route = createFileRoute('/_dashboard/dns')({
  component: DnsSettings,
});
