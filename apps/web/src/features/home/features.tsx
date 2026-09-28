import {
  BarChart2,
  Bell,
  Clock,
  Code,
  Database,
  FileText,
  Layers,
  Rocket,
  Server,
  Shield,
  Cloud as TunnelCloud,
  Users,
} from 'lucide-react';

export function Features() {
  const featuresList = [
    {
      icon: FileText,
      title: 'Canvas View',
      desc: 'See your entire environment at a glance — apps, databases, and their connections laid out as a live topology map, not a flat list.',
    },
    {
      icon: Server,
      title: 'Traefik-backed Routing',
      desc: 'Codedock spins up and manages Traefik automatically — zero-config load balancing, path routing, and automatic SSL for every service you deploy.',
    },
    {
      icon: Rocket,
      title: 'Git-native Deployments',
      desc: 'Connect GitHub, GitLab, or Gitea. Push a branch, get a deployment. No YAML pipelines. No external CI required.',
    },
    {
      icon: Layers,
      title: 'Environment Scoping',
      desc: 'Group apps and services into named environments — staging, production, preview — each with isolated variables and their own deployment flow.',
    },
    {
      icon: Database,
      title: 'Database Snapshots',
      desc: 'Schedule automatic database backups to any S3-compatible storage. Configure retention windows and restore from any checkpoint.',
    },
    {
      icon: Code,
      title: 'AI Settings Layer',
      desc: 'Each team configures its own LLM provider — OpenAI, Anthropic, local Ollama. Codedock routes completions, never stores keys plaintext.',
    },
    {
      icon: BarChart2,
      title: 'Usage Metering',
      desc: 'The cloud tracks container-hours, deployments, and bandwidth per team. Set thresholds, get notified, never hit an unexpected bill.',
    },
    {
      icon: TunnelCloud,
      title: 'Yamux Tunnel Fleet',
      desc: 'Codedock Cloud connects to your daemons via an encrypted Yamux tunnel — no port-forwarding, no VPN. Your servers are never directly exposed.',
    },
    {
      icon: Shield,
      title: 'Scoped Secret Vars',
      desc: 'Inject environment variables per-service or per-environment. Secrets encrypted at rest, never visible in logs or build output.',
    },
    {
      icon: Clock,
      title: 'Cron Jobs',
      desc: 'Define recurring jobs alongside your services — same project, same env vars, same log stream. No separate scheduler needed.',
    },
    {
      icon: Bell,
      title: 'Notification Channels',
      desc: 'Route events — build failure, cert expiry, usage spikes — to email, webhook, or Discord. Configurable per project.',
    },
    {
      icon: Users,
      title: 'Isolated Projects',
      desc: 'Invite members, assign roles, and isolate applications per project. The Cloud dashboard enforces plan limits so billing stays predictable.',
    },
  ];

  return (
    <section id="features" className="scroll-mt-20 py-24">
      <div className="mx-auto max-w-7xl px-6">
        <h2 className="mb-4 text-center font-extrabold text-[clamp(1.8rem,4vw,3.5rem)] text-foreground">
          Everything on your server<span className="text-primary">.</span>
        </h2>
        <p className="mx-auto mb-16 max-w-xl text-center text-muted-foreground text-sm">
          Codedock isn't a wrapper — it's a daemon with real opinions: Traefik for routing, Docker
          for isolation, S3 for backups, Yamux for the fleet tunnel.
        </p>
        <dl className="grid grid-cols-1 gap-x-12 gap-y-10 md:grid-cols-2 lg:grid-cols-3">
          {featuresList.map((f) => {
            const IconComponent = f.icon;
            return (
              <div key={f.title} className="flex flex-col">
                <dt className="mb-2 flex items-center gap-3 font-semibold text-foreground text-sm">
                  <span className="flex size-8 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10">
                    <IconComponent className="size-4 text-primary" />
                  </span>
                  {f.title}
                </dt>
                <dd className="pl-11 text-muted-foreground text-sm leading-relaxed">{f.desc}</dd>
              </div>
            );
          })}
        </dl>
      </div>
    </section>
  );
}
