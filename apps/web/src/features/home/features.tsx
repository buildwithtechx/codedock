import {
  ArrowUpRight,
  Database,
  GitBranch,
  LayoutDashboard,
  Server,
  Shield,
  Timer,
} from 'lucide-react';

const features = [
  {
    icon: GitBranch,
    title: 'Application deployments',
    body: 'Deploy from Git, a Dockerfile or an existing image. Track builds and manage service lifecycles from the dashboard.',
    href: '/features/application-deployment',
  },
  {
    icon: Database,
    title: 'Databases and backups',
    body: 'Provision supported database engines, inspect data and configure scheduled backups to S3-compatible storage.',
    href: '/features/databases',
  },
  {
    icon: LayoutDashboard,
    title: 'A clear view of your stack',
    body: 'Use Canvas, deployment history, logs and metrics to understand what is running and investigate failures.',
    href: '/features/monitoring',
  },
  {
    icon: Server,
    title: 'Your servers, connected',
    body: 'Connect worker servers over SSH. Manage workloads across your infrastructure from one dashboard.',
    href: '/solutions/enterprise',
  },
  {
    icon: Shield,
    title: 'Projects and access',
    body: 'Organize services by project and environment, manage team roles and store sensitive configuration in the encrypted vault.',
    href: '/solutions/agencies',
  },
  {
    icon: Timer,
    title: 'Operations in one place',
    body: 'Manage recurring jobs, domains, notification settings and service operations alongside your applications.',
    href: '/features/monitoring',
  },
];
export function Features() {
  return (
    <section id="features" className="mx-auto max-w-7xl px-6 py-24">
      <p className="mb-3 text-center font-semibold text-primary text-sm">FROM BUILD TO BACKUP</p>
      <h2 className="text-center font-bold text-3xl md:text-4xl">
        Everything your stack needs to run.
      </h2>
      <p className="mx-auto mt-4 max-w-xl text-center text-muted-foreground">
        A Go control plane, Docker containers and Traefik routing. Familiar tools, brought together
        on your infrastructure.
      </p>
      <div className="mt-12 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {features.map(({ icon: Icon, title, body, href }) => (
          <a
            key={title}
            href={href}
            className="group rounded-2xl border border-border bg-card/70 p-7 transition-colors hover:border-primary/50"
          >
            <Icon className="mb-6 size-6 text-primary" aria-hidden="true" />
            <h3 className="flex items-center justify-between gap-3 font-semibold text-lg">
              {title}
              <ArrowUpRight className="size-4 text-muted-foreground" aria-hidden="true" />
            </h3>
            <p className="mt-3 text-muted-foreground text-sm leading-relaxed">{body}</p>
          </a>
        ))}
      </div>
    </section>
  );
}
