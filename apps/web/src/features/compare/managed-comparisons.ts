import type { ComparisonContent } from './comparison-content';
export const managedComparisons = {
  vercel: {
    name: 'Vercel',
    tagline: 'Managed frontend infrastructure or your own server?',
    description:
      'Vercel handles hosting infrastructure for you. Codedock gives you control of the server and container runtime, with the operational responsibility that comes with it.',
    alternativeFit:
      'Choose Vercel if you want managed frontend delivery and serverless workflows without operating the underlying server.',
    codedockFit:
      'Choose Codedock if your application needs long-running containers, databases on your own hardware, and direct control of deployment infrastructure.',
    rows: [
      {
        criterion: 'Infrastructure',
        codedock: 'You choose and operate the server',
        alternative: 'Managed Vercel infrastructure',
      },
      {
        criterion: 'Workloads',
        codedock: 'Containers, workers, databases, and static services',
        alternative: 'Frontend applications and platform function runtimes',
      },
      {
        criterion: 'Costs',
        codedock: 'Your hosting, bandwidth, storage, and operational costs',
        alternative: 'Plan and usage-based infrastructure pricing',
      },
      {
        criterion: 'Operational responsibility',
        codedock: 'You manage updates, capacity, and recovery',
        alternative: 'Platform operates the hosting infrastructure',
      },
    ],
    sources: [
      {
        title: 'Vercel pricing model',
        href: 'https://vercel.com/docs/pricing',
      },
      {
        title: 'Vercel limits',
        href: 'https://vercel.com/docs/limits',
      },
    ],
  },
  render: {
    name: 'Render',
    tagline: 'Bring your own hardware or use a managed runtime.',
    description:
      'Render supports Docker and native language runtimes on its infrastructure. Codedock runs the control plane and your services on infrastructure you choose.',
    alternativeFit:
      'Choose Render if you want managed service hosting and Docker builds with fewer server administration tasks.',
    codedockFit:
      'Choose Codedock if you need your own host, direct infrastructure access, and control over persistent storage and networking.',
    rows: [
      {
        criterion: 'Hosting model',
        codedock: 'Your own Linux server or VPS',
        alternative: 'Managed Render infrastructure',
      },
      {
        criterion: 'Docker workflow',
        codedock: 'Build source or deploy a versioned image',
        alternative: 'Dockerfile builds or prebuilt registry images',
      },
      {
        criterion: 'Hardware control',
        codedock: 'You select the host and available resources',
        alternative: 'Resources configured through platform services',
      },
      {
        criterion: 'Operations',
        codedock: 'You manage patching, backups, and server recovery',
        alternative: 'Platform manages the hosting layer',
      },
    ],
    sources: [
      {
        title: 'Docker on Render',
        href: 'https://render.com/docs/docker',
      },
    ],
  },
} satisfies Record<string, ComparisonContent>;
