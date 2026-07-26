// @ts-check

import starlight from '@astrojs/starlight';
import { defineConfig } from 'astro/config';

export default defineConfig({
  site: 'https://docs.codedock.run',
  integrations: [
    starlight({
      title: 'Codedock Docs',
      favicon: '/favicon.svg',
      customCss: ['./src/styles/theme.css'],
      sidebar: [
        {
          label: 'Getting Started',
          collapsed: false,
          items: [
            { label: 'Documentation Home', slug: 'index' },
            { label: 'Installation', slug: 'getting-started/installation' },
            { label: 'Browser Onboarding', slug: 'getting-started/onboarding' },
            { label: 'Quick Start Guide', slug: 'getting-started/quick-start' },
            { label: 'Deploy Your First App', slug: 'tutorial' },
          ],
        },
        {
          label: 'Core Concepts',
          collapsed: false,
          items: [
            { label: 'Architecture', slug: 'core-concepts/architecture' },
            { label: 'Projects & Environments', slug: 'projects/overview' },
            { label: 'Project Settings', slug: 'projects/settings' },
          ],
        },
        {
          label: 'Deployments',
          collapsed: false,
          items: [
            { label: 'Build Strategies', slug: 'deployments/build-strategies' },
            { label: 'Templates & Examples', slug: 'deployments/templates' },
            { label: 'CI/CD & Git Webhooks', slug: 'deployments/ci-cd' },
            { label: 'Service Types', slug: 'deployments/service-types' },
            { label: 'Scheduled Cron Jobs', slug: 'deployments/jobs' },
            {
              label: 'Environment Variables',
              slug: 'deployments/environment-variables',
            },
            {
              label: 'Deployment Lifecycle',
              slug: 'deployments/deployment-lifecycle',
            },
          ],
        },
        {
          label: 'Databases',
          collapsed: false,
          items: [
            { label: 'Database Provisioning', slug: 'databases/provisioning' },
            { label: 'Native Data Browser', slug: 'databases/data-browser' },
            { label: 'SQL Studio & Queries', slug: 'databases/sql-studio' },
            { label: 'Data Imports & Migration', slug: 'databases/data-imports' },
            {
              label: 'Public Access & TLS',
              slug: 'databases/public-access-and-tls',
            },
          ],
        },
        {
          label: 'Storage & Backups',
          collapsed: false,
          items: [
            {
              label: 'Database Backups',
              slug: 'storage-and-backups/database-backups',
            },
            {
              label: 'MinIO Object Storage',
              slug: 'storage-and-backups/minio-storage',
            },
            {
              label: 'Cloudflare R2 Storage',
              slug: 'storage-and-backups/r2-storage',
            },
            {
              label: 'Restore & Downloads',
              slug: 'storage-and-backups/restore-and-download',
            },
          ],
        },
        {
          label: 'Networking & SSL',
          collapsed: false,
          items: [
            { label: 'Domains & DNS Setup', slug: 'operations/domains-and-dns' },
            { label: 'Migration & Bundles', slug: 'migration/codedock-bundles' },
          ],
        },
        {
          label: 'Security & Operations',
          collapsed: false,
          items: [
            { label: 'Teams & Collaboration', slug: 'operations/teams' },
            { label: 'Admin Dashboard', slug: 'admin' },
            { label: 'Third-party Integrations', slug: 'integrations' },
            { label: 'Observability & Logs', slug: 'operations/observability' },
            { label: 'Account Security & 2FA', slug: 'operations/account-security' },
            {
              label: 'Maintenance & Updates',
              slug: 'operations/maintenance-and-updates',
            },
            { label: 'Troubleshooting Guide', slug: 'operations/troubleshooting' },
          ],
        },
        {
          label: 'Reference',
          collapsed: false,
          items: [
            { label: 'System Configuration', slug: 'configuration' },
            { label: 'CLI Reference', slug: 'cli' },
            { label: 'Full REST API Reference', slug: 'api' },
            { label: 'System Settings API', slug: 'reference/system-settings' },
            { label: 'API Authentication Access', slug: 'reference/api-access' },
            { label: 'API Projects Reference', slug: 'reference/api-projects' },
            { label: 'API Services Reference', slug: 'reference/api-services' },
            { label: 'API Deployments Reference', slug: 'reference/api-deployments' },
            { label: 'API Databases Reference', slug: 'reference/api-databases' },
            {
              label: 'API Env & Domains',
              slug: 'reference/api-environment-and-domains',
            },
            { label: 'Zero Lock-in Commitment', slug: 'adopt' },
          ],
        },
      ],
      components: {
        SiteTitle: './src/components/docs-site-title.astro',
        ThemeSelect: './src/components/docs-theme-select.astro',
      },
      social: [
        {
          icon: 'github',
          label: 'GitHub',
          href: 'https://github.com/buildwithtechx/codedock',
        },
      ],
    }),
  ],
});
