import { Activity, Box, Database, GitBranch, Layers } from 'lucide-react';

export function ControlPlanePreview() {
  const services = [
    { name: 'web', type: 'Application', icon: Box },
    { name: 'api', type: 'Application', icon: Box },
    { name: 'postgres', type: 'Database', icon: Database },
  ];
  return (
    <div className="relative mx-auto w-full max-w-xl lg:rotate-1">
      <div
        aria-hidden="true"
        className="absolute -inset-8 -z-10 rounded-full bg-primary/15 blur-3xl"
      />
      <div className="overflow-hidden rounded-2xl border border-primary/25 bg-card/85 shadow-2xl shadow-primary/10 backdrop-blur-xl">
        <div className="flex items-center justify-between border-border border-b px-5 py-4">
          <span className="flex items-center gap-2 font-semibold text-sm">
            <Layers className="size-4 text-primary" />
            Example workspace
          </span>
          <span className="rounded-full border border-border px-2.5 py-1 font-mono text-[10px] text-muted-foreground">
            PRODUCTION
          </span>
        </div>
        <div className="p-5">
          <div className="mb-5 flex items-start justify-between">
            <div>
              <p className="text-muted-foreground text-xs">Project</p>
              <p className="mt-1 font-bold text-lg">Your next application</p>
            </div>
            <span className="flex items-center gap-1 rounded-lg bg-emerald-500/10 px-2 py-1 text-emerald-400 text-xs">
              <Activity className="size-3" />
              Running
            </span>
          </div>
          <div className="grid grid-cols-3 gap-2">
            {services.map((service) => (
              <div
                key={service.name}
                className="rounded-xl border border-border bg-background/70 p-3"
              >
                <service.icon className="mb-3 size-5 text-primary" />
                <p className="font-semibold text-xs">{service.name}</p>
                <p className="mt-1 text-[10px] text-muted-foreground">{service.type}</p>
                <div className="mt-3 h-1 rounded-full bg-primary/20">
                  <div className="h-full w-2/3 rounded-full bg-primary/70" />
                </div>
              </div>
            ))}
          </div>
          <div className="mt-5 rounded-xl border border-border bg-background/60 p-4">
            <p className="mb-3 flex items-center gap-2 text-muted-foreground text-xs">
              <GitBranch className="size-3.5" />
              Deployment activity
            </p>
            <div className="space-y-2 font-mono text-[11px]">
              <p>
                <span className="text-muted-foreground">01</span>{' '}
                <span className="text-primary">build</span> container image ready
              </p>
              <p>
                <span className="text-muted-foreground">02</span>{' '}
                <span className="text-primary">deploy</span> health check passed
              </p>
              <p>
                <span className="text-muted-foreground">03</span>{' '}
                <span className="text-emerald-400">route</span> service connected
              </p>
            </div>
          </div>
        </div>
        <div className="flex justify-between border-border border-t px-5 py-3 text-[10px] text-muted-foreground">
          <span>Apps / databases / logs / backups</span>
          <span>Illustrative preview</span>
        </div>
      </div>
    </div>
  );
}
