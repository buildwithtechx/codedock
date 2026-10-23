// @ts-check

import starlight from "@astrojs/starlight";
import { defineConfig } from "astro/config";

export default defineConfig({
  site: "https://docs.codedock.run",
  redirects: {
    "/tutorial": "/getting-started/tutorial",
    "/admin": "/operations/admin",
    "/integrations": "/operations/integrations",
    "/cli": "/reference/cli",
    "/configuration": "/reference/configuration",
    "/api": "/reference/api",
    "/adopt": "/reference/adopt",
  },
  integrations: [
    starlight({
      title: "Codedock Docs",
      favicon: "/favicon.svg",
      customCss: ["./src/styles/theme.css"],
      sidebar: [
        {
          label: "Getting Started",
          collapsed: false,
          items: [
            { label: "Documentation Home", slug: "index" },
            { label: "Quick Start Guide", slug: "getting-started/quick-start" },
            { label: "Installation", slug: "getting-started/installation" },
            { label: "Browser Onboarding", slug: "getting-started/onboarding" },
            { label: "Deploy Your First App", slug: "getting-started/tutorial" },
            { label: "Architecture", slug: "core-concepts/architecture" },
            { label: "Projects & Environments", slug: "projects/overview" },
            { label: "Project Settings", slug: "projects/settings" },
          ],
        },
        {
          label: "Deployments & Workflows",
          collapsed: false,
          items: [
            { label: "Build Strategies", slug: "deployments/build-strategies" },
            { label: "Docker Compose Import", slug: "deployments/compose" },
            { label: "CI/CD & Git Webhooks", slug: "deployments/ci-cd" },
            { label: "Actions & CI/CD Runners", slug: "deployments/actions-and-runners" },
            { label: "Service Types", slug: "deployments/service-types" },
            { label: "Scheduled Cron Jobs", slug: "deployments/jobs" },
            { label: "Environment Variables", slug: "deployments/environment-variables" },
            { label: "Rollbacks & Attention Health", slug: "deployments/rollbacks-and-health" },
            { label: "Deployment Lifecycle", slug: "deployments/deployment-lifecycle" },
            { label: "AI Workloads", slug: "deployments/ai-workloads" },
            { label: "Templates & Examples", slug: "deployments/templates" },
          ],
        },
        {
          label: "Databases & Storage",
          collapsed: false,
          items: [
            { label: "Database Provisioning", slug: "databases/provisioning" },
            { label: "Native Data Browser", slug: "databases/data-browser" },
            { label: "SQL Studio & Queries", slug: "databases/sql-studio" },
            { label: "Database Backups", slug: "storage-and-backups/database-backups" },
            { label: "Restore & Downloads", slug: "storage-and-backups/restore-and-download" },
            { label: "MinIO Object Storage", slug: "storage-and-backups/minio-storage" },
            { label: "Cloudflare R2 Storage", slug: "storage-and-backups/r2-storage" },
            { label: "Public Access & TLS", slug: "databases/public-access-and-tls" },
            { label: "Data Imports & Migration", slug: "databases/data-imports" },
          ],
        },
        {
          label: "Networking & Routing",
          collapsed: false,
          items: [
            { label: "Domains & DNS Setup", slug: "operations/domains-and-dns" },
            { label: "Ingress & L7 Route Rules", slug: "networking/ingress-and-routing" },
            { label: "Software-Defined Networks", slug: "networking/software-defined-networks" },
            { label: "Server Takeover & Adoption", slug: "fleet-management/server-takeover" },
            { label: "Migration & Bundles", slug: "migration/codedock-bundles" },
          ],
        },
        {
          label: "Fleet Management",
          collapsed: false,
          items: [
            { label: "Worker Servers", slug: "fleet-management/servers" },
            { label: "Managed Cloud (Hetzner)", slug: "fleet-management/managed-servers" },
            { label: "High Availability & Clustering", slug: "fleet-management/clustering-ha" },
            { label: "Canvas Topology", slug: "fleet-management/canvas" },
          ],
        },
        {
          label: "Security & Operations",
          collapsed: false,
          items: [
            { label: "Teams & Collaboration", slug: "operations/teams" },
            { label: "Account Security & 2FA", slug: "operations/account-security" },
            { label: "Codedock Mail Server", slug: "operations/mail-server" },
            { label: "Push Notifications & Mobile", slug: "operations/push-notifications" },
            { label: "Observability & Logs", slug: "operations/observability" },
            { label: "Admin Dashboard", slug: "operations/admin" },
            { label: "Third-party Integrations", slug: "operations/integrations" },
            { label: "Maintenance & Updates", slug: "operations/maintenance-and-updates" },
            { label: "Troubleshooting Guide", slug: "operations/troubleshooting" },
          ],
        },
        {
          label: "Developer Tools & AI",
          collapsed: false,
          items: [
            { label: "MCP Server & Ask AI Copilot", slug: "reference/mcp-server" },
            { label: "Developer Tools & Mobile Apps", slug: "reference/developer-tools" },
          ],
        },
        {
          label: "Reference",
          collapsed: false,
          items: [
            { label: "CLI Reference", slug: "reference/cli" },
            { label: "REST API Overview", slug: "reference/api" },
            { label: "API Endpoint Catalogue", slug: "reference/api-endpoints" },
            { label: "System Settings API", slug: "reference/system-settings" },
            { label: "API Authentication Access", slug: "reference/api-access" },
            { label: "API Projects Reference", slug: "reference/api-projects" },
            { label: "API Services Reference", slug: "reference/api-services" },
            { label: "API Deployments Reference", slug: "reference/api-deployments" },
            { label: "API Databases Reference", slug: "reference/api-databases" },
            { label: "API Env & Domains", slug: "reference/api-environment-and-domains" },
            { label: "System Configuration", slug: "reference/configuration" },
            { label: "Zero Lock-in Commitment", slug: "reference/adopt" },
          ],
        },
      ],
      components: {
        SiteTitle: "./src/components/docs-site-title.astro",
      },
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/buildwithtechx/codedock",
        },
      ],
    }),
  ],
});
