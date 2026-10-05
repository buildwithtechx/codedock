import { productLinks } from './product-links';
export interface NavigationItem {
  label: string;
  href: string;
  description: string;
}
export const navigationGroups: { label: string; items: NavigationItem[] }[] = [
  {
    label: 'Features',
    items: [
      {
        label: 'Application deployment',
        href: '/features/application-deployment',
        description: 'Git, images, and build workflows',
      },
      {
        label: 'Databases',
        href: '/features/databases',
        description: 'Provision, inspect, query, and back up',
      },
      {
        label: 'Monitoring & logs',
        href: '/features/monitoring',
        description: 'Understand builds and running services',
      },
      {
        label: 'AI workloads',
        href: '/features/ai-deployment',
        description: 'Inference and vector services on your hardware',
      },
    ],
  },
  {
    label: 'Solutions',
    items: [
      {
        label: 'Self-hosted PaaS',
        href: '/solutions/self-hosted',
        description: 'Start on infrastructure you control',
      },
      {
        label: 'Teams & fleet',
        href: '/solutions/enterprise',
        description: 'Permissions and SSH worker servers',
      },
      {
        label: 'Agencies',
        href: '/solutions/agencies',
        description: 'Organize client projects and environments',
      },
    ],
  },
  {
    label: 'Resources',
    items: [
      {
        label: 'Template library',
        href: '/templates',
        description: 'Find a database workflow or Docker recipe',
      },
      {
        label: 'Compare platforms',
        href: '/vs',
        description: 'Find the right fit for your deployment',
      },
      {
        label: 'Changelog',
        href: '/changelog',
        description: 'Follow development and published releases',
      },
      {
        label: 'Philosophy',
        href: '/philosophy',
        description: 'Open source and infrastructure ownership',
      },
      {
        label: 'Documentation',
        href: productLinks.docs,
        description: 'Install, deploy, and operate Codedock',
      },
    ],
  },
];
export const footerGroups = [
  { label: 'Product', items: navigationGroups[0].items },
  {
    label: 'Enterprise',
    items: [
      { label: 'Teams & fleet', href: '/solutions/enterprise' },
      { label: 'Roles & permissions', href: `${productLinks.docs}/operations/teams/` },
      { label: 'API reference', href: `${productLinks.docs}/api/` },
    ],
  },
  { label: 'Solutions', items: navigationGroups[1].items },
  {
    label: 'Compare & learn',
    items: [
      { label: 'All comparisons', href: '/vs' },
      { label: 'vs. Coolify', href: '/vs/coolify' },
      { label: 'vs. Dokploy', href: '/vs/dokploy' },
      { label: 'vs. Vercel', href: '/vs/vercel' },
      { label: 'Template library', href: '/templates' },
      { label: 'Documentation', href: productLinks.docs },
    ],
  },
  {
    label: 'Company',
    items: [
      { label: 'Philosophy', href: '/philosophy' },
      { label: 'Changelog', href: '/changelog' },
      { label: 'GitHub', href: productLinks.github },
      { label: 'Community', href: productLinks.discord },
      { label: 'Licence', href: `${productLinks.github}/blob/main/LICENSE` },
    ],
  },
];
