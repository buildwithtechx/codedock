import type { PageContent } from './page-content';

export const solutionContent = {
  'self-hosted': {
    eyebrow: 'Self-hosted PaaS',
    title: 'A deployment platform that lives on your server.',
    description:
      'Install Codedock on a Linux server, create your owner account, and deploy applications with a Go control plane and a browser dashboard.',
    highlights: ['No manual .env setup', 'Apache-2.0 source', 'Your choice of server'],
    sections: [
      {
        title: 'Start with one command',
        body: 'The installer checks the host, installs Docker when needed, generates private secrets, and starts the control plane. The first account becomes the instance owner.',
        points: [
          'Automatic secrets retained across restarts',
          'Optional domain and certificate email through setup',
          'HTTP dashboard on port 8080',
        ],
      },
      {
        title: 'Understand the moving parts',
        body: 'The Go daemon serves the dashboard and coordinates Docker. Traefik handles application routing; builders and database tools run as separate containers.',
        points: [
          'SQLite stores control plane state',
          'Persistent data includes secrets and the encrypted vault',
          'Docker, Traefik, and build tools remain runtime dependencies',
        ],
      },
      {
        title: 'Own the operational decisions',
        body: 'Choose your provider, networking, data retention, and update schedule. Self-hosted billing endpoints and cloud plan limits are disabled.',
        points: [
          'No Codedock per-deployment fee in self-hosted mode',
          'Hosting, bandwidth, and storage costs come from your providers',
          'Back up control plane state and application volumes',
        ],
      },
    ],
    steps: [
      'Install on a Linux server with a reachable address.',
      'Create the owner account in the browser or terminal.',
      'Optionally configure an application domain.',
      'Deploy a service and configure backups.',
    ],
    docsPath: '/getting-started/installation/',
    docsLabel: 'Install Codedock',
  },
  enterprise: {
    eyebrow: 'Teams & fleet operations',
    title: 'Make infrastructure visible to the whole team.',
    description:
      'Organize projects, permissions, worker servers, and operational history in a shared workspace. Evaluate the capabilities you need before rolling out to a larger fleet.',
    highlights: ['Organization permissions', 'SSH worker servers', 'Operational visibility'],
    sections: [
      {
        title: 'Define who can do what',
        body: 'Use instance and organization roles to control access to projects and administrative operations. Keep production configuration in a dedicated environment.',
        points: [
          'Owner, admin, and member workflows',
          'Organization membership and project access',
          'Project API tokens with scoped permissions',
        ],
      },
      {
        title: 'Connect worker servers',
        body: 'Register reachable Linux hosts with direct SSH credentials and run the preflight check. Review health and resource metrics from the server inventory.',
        points: [
          'SSH key or password authentication',
          'Credentials encrypted in the vault',
          'Jump-host and Cloudflare transports are not currently supported',
        ],
      },
      {
        title: 'Plan an operational rollout',
        body: 'Test restore procedures, update handling, and permissions on a small deployment first. Cloud billing is a distinct mode; self-hosted instances do not require Stripe.',
        points: [
          'Back up the database, vault key, and generated secrets',
          'Monitor deployments, tasks, and backup history',
          'Yamux tunnels and usage-based fleet metering are not shipping capabilities',
        ],
      },
    ],
    steps: [
      'Create your organization and invite teammates.',
      'Assign permissions and separate environments.',
      'Register a worker server and verify SSH access.',
      'Run a deployment, backup, and restore exercise.',
    ],
    docsPath: '/fleet-management/servers/',
    docsLabel: 'Read the fleet guide',
  },
  agencies: {
    eyebrow: 'Agencies & teams',
    title: "Keep each client's work in its own workspace.",
    description:
      'Manage client projects, staging environments, configuration, and deployments with clear boundaries and a repeatable handover process.',
    highlights: ['Projects per client', 'Staging & production', 'Repeatable handovers'],
    sections: [
      {
        title: 'Give every project a clear home',
        body: 'Group applications, databases, and variables by project. Use environments to distinguish preview work from production deployments.',
        points: [
          'Project-level deployment history',
          'Environment-specific configuration',
          'Canvas overview of project resources',
        ],
      },
      {
        title: 'Share access deliberately',
        body: 'Invite the people who need access and choose organization permissions that match their responsibilities. Separate organizations or instances when clients require stronger isolation.',
        points: [
          'Team roles and scoped API access',
          'Encrypted configuration values',
          'Project grouping is not a substitute for host or network isolation',
        ],
      },
      {
        title: 'Make delivery repeatable',
        body: 'Use versioned images or source repositories, write down required configuration, and test backups before handover. Keep client domains and storage destinations documented.',
        points: [
          'Git and Docker image workflows',
          'Custom domains and certificates',
          'Documented backup and restore procedures',
        ],
      },
    ],
    steps: [
      'Choose an organization or instance boundary for the client.',
      'Create the project and staging environment.',
      'Deploy the application and its database.',
      'Review access, domains, and recovery before handover.',
    ],
    docsPath: '/operations/teams/',
    docsLabel: 'Read the collaboration guide',
  },
} satisfies Record<string, PageContent>;
