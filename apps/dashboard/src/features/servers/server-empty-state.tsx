import { Link } from '@tanstack/react-router';
import { Activity, ArrowRight, Container, Globe2, Plus, Server } from 'lucide-react';
import { Button } from '#/components/ui/button';

const serverCapabilities = [
  { label: 'Runtime', detail: 'Container workloads', icon: Container },
  { label: 'Networking', detail: 'Public application routes', icon: Globe2 },
  { label: 'Health', detail: 'Live capacity and status', icon: Activity },
];

export function ServerEmptyState() {
  return (
    <section className="flex min-h-[34rem] items-center justify-center py-12 text-center">
      <div className="w-full max-w-3xl">
        <div className="relative mx-auto h-32 w-56" aria-hidden="true">
          <span className="absolute top-8 left-8 h-17 w-26 rounded-2xl bg-card" />
          <span className="absolute top-3 left-15 flex h-17 w-26 items-center justify-center rounded-2xl bg-muted text-muted-foreground">
            <Server className="h-6 w-6" />
          </span>
          <span className="absolute top-10 right-6 flex h-11 w-11 items-center justify-center rounded-xl border border-primary/45 border-dashed bg-primary/8 text-primary">
            <Plus className="h-4 w-4" />
          </span>
          <span className="absolute top-1/2 right-15 w-8 border-border border-t border-dashed" />
        </div>
        <h2 className="mt-5 font-medium text-foreground/90 text-xl tracking-[-0.02em]">
          Add a server
        </h2>
        <p className="mx-auto mt-2 max-w-md text-muted-foreground/75 text-sm leading-6">
          Connect a runtime so Codedock can build, deploy, and observe applications on your
          infrastructure.
        </p>
        <Link to="/servers/new">
          <Button className="mt-6 gap-2 px-5">
            Add server
            <ArrowRight className="h-4 w-4" />
          </Button>
        </Link>
        <div className="mt-11">
          <p className="mb-4 font-medium text-[10px] text-muted-foreground uppercase tracking-[0.16em]">
            Connected servers provide
          </p>
          <div className="grid gap-3 text-left sm:grid-cols-3">
            {serverCapabilities.map(({ label, detail, icon: Icon }) => (
              <div key={label} className="rounded-2xl bg-card p-4">
                <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                  <Icon className="h-4 w-4" />
                </span>
                <p className="mt-4 font-semibold text-sm">{label}</p>
                <p className="mt-1 text-muted-foreground text-xs leading-5">{detail}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
