import type { SupportCategory } from '#/interfaces/support';

export interface SupportTopicDefinition {
  category: SupportCategory;
  label: string;
  description: string;
}

export const SUPPORT_TOPICS: SupportTopicDefinition[] = [
  {
    category: 'deployment',
    label: 'Deployments',
    description: 'Build failures, runtime logs, cache mounts, and Docker containers.',
  },
  {
    category: 'billing',
    label: 'Billing',
    description: 'Invoices, plans, compute quota, and payment methods.',
  },
  {
    category: 'account',
    label: 'Account & Security',
    description: 'Organization membership, access tokens, 2FA, and permissions.',
  },
  {
    category: 'general',
    label: 'General Guidance',
    description: 'Platform architecture, migrations, and best practices.',
  },
];
