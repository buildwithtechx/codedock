import { createFileRoute } from '@tanstack/react-router';
import { DomainsPage } from '#/features/dns/domains-page';

export const Route = createFileRoute('/_dashboard/domains')({
  component: DomainsPage,
});
