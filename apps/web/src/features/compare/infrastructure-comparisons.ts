import type { ComparisonContent } from './comparison-content';
export const infrastructureComparisons = {
  coolify: {
    name: 'Coolify',
    tagline: 'Two self-hosted paths. Different control planes.',
    description:
      'Both platforms deploy applications on servers you control. Compare the operational workflow and architecture that suit your team, rather than assuming one replaces every feature of the other.',
    alternativeFit:
      'Choose Coolify if its established application workflows, integrations, and existing deployment setup already fit your team.',
    codedockFit:
      'Choose Codedock if you want a Go control plane with project environments, a database browser, SQL Studio, and a visual canvas in the same workspace.',
    rows: [
      {
        criterion: 'Hosting model',
        codedock: 'Self-hosted on your Linux server',
        alternative: 'Self-hosted or Coolify Cloud',
      },
      {
        criterion: 'Control plane',
        codedock: 'A Go control plane with Docker, Traefik, and container-based build tools',
        alternative: 'Containerized Laravel dashboard and supporting services',
      },
      {
        criterion: 'Application runtime',
        codedock: 'Docker containers for applications and databases',
        alternative: 'Docker containers on connected servers',
      },
      {
        criterion: 'Operational focus',
        codedock: 'Projects, data tools, canvas, backups, and logs',
        alternative: 'Application deployment and server administration',
      },
    ],
    sources: [
      {
        title: 'Coolify architecture',
        href: 'https://coolify.io/docs/core/how-coolify-works',
      },
      {
        title: 'Coolify applications',
        href: 'https://coolify.io/docs/applications/',
      },
    ],
  },
  dokploy: {
    name: 'Dokploy',
    tagline: 'Compare the workflow, not just the feature list.',
    description:
      'Codedock and Dokploy both bring application deployment to your own infrastructure. Your existing Compose stacks and preferred operational tools should guide the choice.',
    alternativeFit:
      'Choose Dokploy if its Docker Compose deployment workflow and monitoring tools match how you already run your applications.',
    codedockFit:
      'Choose Codedock if you prefer a Go-based control plane and want database inspection, environment topology, and service operations in one dashboard.',
    rows: [
      {
        criterion: 'Hosting model',
        codedock: 'Self-hosted Linux servers',
        alternative: 'Self-hosted deployment platform',
      },
      {
        criterion: 'Deployment options',
        codedock: 'Source builds, images, and Compose import',
        alternative: 'Git, Docker images, and Docker Compose',
      },
      {
        criterion: 'Operational visibility',
        codedock: 'Container metrics, logs, and project canvas',
        alternative: 'Application and server monitoring',
      },
      {
        criterion: 'Architecture',
        codedock: 'A Go control plane with Docker, Traefik, and container-based build tools',
        alternative: 'Docker-based deployment and management workflows',
      },
    ],
    sources: [
      {
        title: 'Dokploy features',
        href: 'https://docs.dokploy.com/docs/core/features',
      },
      {
        title: 'Dokploy applications',
        href: 'https://docs.dokploy.com/docs/core/applications',
      },
    ],
  },
  portainer: {
    name: 'Portainer',
    tagline: 'Application delivery or container administration?',
    description:
      "Portainer's Docker workflows center on managing containers and stacks. Codedock organizes the developer workflow around projects, environments, deployments, and databases.",
    alternativeFit:
      'Choose Portainer if administering existing Docker environments and Compose stacks is your primary job.',
    codedockFit:
      'Choose Codedock if you want repository-to-deployment workflows, project configuration, and integrated data tools for the applications you build.',
    rows: [
      {
        criterion: 'Primary workflow',
        codedock: 'Build, deploy, and operate application projects',
        alternative: 'Manage Docker environments and stacks',
      },
      {
        criterion: 'Compose',
        codedock: 'Import definitions into project resources',
        alternative: 'Create and manage Docker stacks',
      },
      {
        criterion: 'Developer context',
        codedock: 'Services grouped by project and environment',
        alternative: 'Services grouped into container stacks',
      },
      {
        criterion: 'Control plane',
        codedock: 'A Go control plane with Docker, Traefik, and container-based build tools',
        alternative: 'Container management dashboard',
      },
    ],
    sources: [
      {
        title: 'Portainer Docker stacks',
        href: 'https://docs.portainer.io/user/docker/stacks',
      },
    ],
  },
  caprover: {
    name: 'CapRover',
    tagline: 'Choose the deployment model you want to maintain.',
    description:
      'Both tools make self-hosted application hosting accessible through a dashboard. Compare the setup, routing, and operational workflow with the needs of your applications.',
    alternativeFit:
      'Choose CapRover if its deployment definitions, dashboard, and one-click application ecosystem fit your existing workflow.',
    codedockFit:
      'Choose Codedock if you want Traefik routing alongside database browsing, backup configuration, and environment topology.',
    rows: [
      {
        criterion: 'Hosting model',
        codedock: 'Self-hosted Linux server',
        alternative: 'Self-hosted server with a dashboard',
      },
      {
        criterion: 'Setup',
        codedock: 'Installer, owner account, optional application domain',
        alternative: 'Server setup and wildcard domain configuration',
      },
      {
        criterion: 'Routing',
        codedock: 'Traefik with configured certificate email',
        alternative: 'Application domains and HTTPS setup',
      },
      {
        criterion: 'Operations',
        codedock: 'Canvas, database tools, logs, and backups',
        alternative: 'Application deployment and container configuration',
      },
    ],
    sources: [
      {
        title: 'CapRover getting started',
        href: 'https://caprover.com/docs/get-started.html',
      },
    ],
  },
  dokku: {
    name: 'Dokku',
    tagline: 'A visual workspace or a Git-first server workflow.',
    description:
      'Dokku offers a command-oriented application deployment workflow. Codedock adds a browser workspace for project resources, deployments, and data operations.',
    alternativeFit:
      'Choose Dokku if you enjoy a Git-first deployment workflow and command-line server administration.',
    codedockFit:
      'Choose Codedock if your team benefits from a visual canvas, browser data tools, and server inventory alongside deployment controls.',
    rows: [
      {
        criterion: 'Interaction',
        codedock: 'Browser dashboard, API, and CLI',
        alternative: 'Command-line and Git deployment workflow',
      },
      {
        criterion: 'Source deployment',
        codedock: 'Repository configuration and deployment controls',
        alternative: 'Git push application deployment',
      },
      {
        criterion: 'Resource organization',
        codedock: 'Projects, environments, and a visual canvas',
        alternative: 'Applications managed through server commands',
      },
      {
        criterion: 'Operations',
        codedock: 'Dashboard logs, metrics, database tools, and backups',
        alternative: 'Command-oriented application administration',
      },
    ],
    sources: [
      {
        title: 'Dokku application deployment',
        href: 'https://dokku.com/docs/deployment/application-deployment/',
      },
    ],
  },
} satisfies Record<string, ComparisonContent>;
