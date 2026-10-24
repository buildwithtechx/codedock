import type { PageContent } from './page-content';

export const featureContent = {
  'application-deployment': {
    eyebrow: 'Application deployment',
    title: 'From a repository to a running service.',
    description:
      'Deploy source code, container images, and multi-service projects on infrastructure you control. Keep configuration, deployment history, and logs in the same workspace.',
    highlights: ['Git & Docker images', 'Web apps & workers', 'Deployment history'],
    sections: [
      {
        title: 'Choose how you build',
        body: 'Use a Dockerfile for an explicit build, or select an automatic builder for your source project. Set the root directory for a monorepo and tune install, build, and start commands.',
        points: [
          'Dockerfile and prebuilt image workflows',
          'Automatic builds with Nixpacks and Buildpacks',
          'Railpack selection currently uses the Nixpacks build path',
        ],
      },
      {
        title: 'Keep your environments organized',
        body: 'Group services into projects and environments. Define variables at the right scope and review a deployment before changing live traffic.',
        points: [
          'Separate staging and production configuration',
          'Git branches and deployment logs',
          'Import Compose definitions into project services',
        ],
      },
      {
        title: 'Operate after the first deploy',
        body: 'Use health checks, service metrics, logs, and deployment history to understand each rollout. Scheduled tasks and background workers live alongside your applications.',
        points: [
          'Start, stop, restart, and redeploy services',
          'Review failed build output',
          'Custom domains through Traefik',
        ],
      },
    ],
    steps: [
      'Create a project and an environment.',
      'Add a repository or a versioned container image.',
      'Set the port, commands, and environment variables.',
      'Deploy, read the logs, and attach your domain.',
    ],
    docsPath: '/deployments/build-strategies/',
    docsLabel: 'Read the deployment guide',
  },
  databases: {
    eyebrow: 'Database management',
    title: 'Your data. Your server. One workspace.',
    description:
      'Provision databases, inspect supported relational tables, run queries, and configure backups without leaving your deployment dashboard.',
    highlights: ['SQL & NoSQL workflows', 'Data browser & SQL Studio', 'S3-compatible backups'],
    sections: [
      {
        title: 'Give applications a persistent database',
        body: 'Create PostgreSQL, MySQL, MariaDB, Redis, MongoDB, or ClickHouse services with persistent storage. Keep credentials and connection details with the project.',
        points: [
          'Choose an engine and version',
          'Connect through your private runtime network',
          'Review resource limits and public-access settings',
        ],
      },
      {
        title: 'Inspect the data behind a deployment',
        body: 'Use the database browser to inspect schemas and records. SQL Studio supports query execution for compatible engines; the available tools depend on the database.',
        points: [
          'Browse PostgreSQL, MySQL and MariaDB tables',
          'Run SQL and inspect returned rows',
          'Import data with engine-specific workflows',
        ],
      },
      {
        title: 'Make recovery part of the setup',
        body: 'Create backup configurations with manual or cron schedules. Store backups locally or send them to an S3-compatible destination such as R2 or MinIO.',
        points: [
          'Configure retention and review backup records',
          'Download a backup for an independent copy',
          'Test a restore before relying on it',
        ],
      },
    ],
    steps: [
      'Create a database in a project.',
      'Copy the connection values into your application variables.',
      'Inspect a table or run a query.',
      'Configure a backup destination and test recovery.',
    ],
    docsPath: '/databases/provisioning/',
    docsLabel: 'Read the database guide',
  },
  monitoring: {
    eyebrow: 'Monitoring & logs',
    title: 'Know what happened. See what is running.',
    description:
      'Move from deployment history to build output, live service logs, and container metrics without switching tools.',
    highlights: ['Build & runtime logs', 'Container metrics', 'Notification settings'],
    sections: [
      {
        title: 'Follow a deployment from start to finish',
        body: 'Review build output and deployment status to find which step failed. Keep the context of the service, project, and environment.',
        points: [
          'Deployment history and failure output',
          'Streaming logs for active services',
          'Browser terminal for supported containers',
        ],
      },
      {
        title: 'Watch the runtime',
        body: 'Inspect CPU and memory usage, service status, and server health. Use the environment canvas to put individual services in context.',
        points: [
          'Container resource metrics',
          'Server inventory and health status',
          'Canvas view of applications and databases',
        ],
      },
      {
        title: 'Connect your operational workflow',
        body: 'Configure notification channels and review backup and scheduled-task history. AI diagnosis can help interpret deployment failures when a provider is configured.',
        points: [
          'Email and webhook integrations',
          'Backup execution records',
          'Bring your own AI provider credentials',
        ],
      },
    ],
    steps: [
      'Open an application and its deployment history.',
      'Inspect the build output or service logs.',
      'Check container metrics and server health.',
      'Configure notification channels for your team.',
    ],
    docsPath: '/operations/observability/',
    docsLabel: 'Read the monitoring guide',
  },
  'ai-deployment': {
    eyebrow: 'AI workloads',
    title: 'Bring AI workloads closer to your data.',
    description:
      "Run containerized AI services on your own hardware and connect your applications to them. Use Codedock's AI settings separately for deployment diagnosis.",
    highlights: [
      'Container image deployment',
      'Private service connections',
      'Bring your own hardware',
    ],
    sections: [
      {
        title: 'Deploy an inference service',
        body: 'Use a Docker image service for an inference runtime such as Ollama. Choose hardware and model sizes that fit your workload, and configure persistent model storage.',
        points: [
          'Set the image and exposed container port',
          'Attach a persistent volume for model data',
          'GPU configuration requires your own compatible Docker host setup',
        ],
      },
      {
        title: 'Add a retrieval layer',
        body: 'A vector database such as Qdrant can run as a container service. PostgreSQL extensions such as pgvector require a compatible image and extension configuration.',
        points: [
          'Connect services using private network addresses',
          'Manage credentials as project or service variables',
          'Back up persistent data with an appropriate workflow',
        ],
      },
      {
        title: 'Use AI to diagnose deployments',
        body: 'Configure a supported AI provider in the dashboard to analyze deployment output. This is a separate integration from running your own model container.',
        points: [
          'Bring your own provider key or endpoint',
          'Review suggestions before changing a deployment',
          "Keep the application's inference configuration separate",
        ],
      },
    ],
    steps: [
      'Choose a versioned inference or vector database image.',
      'Check CPU, memory, storage, and GPU requirements.',
      'Deploy the service and configure persistence.',
      'Connect your application and verify a real request.',
    ],
    docsPath: '/deployments/ai-workloads/',
    docsLabel: 'Read the AI workload guide',
  },
} satisfies Record<string, PageContent>;
