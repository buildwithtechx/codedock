import { Link } from '@tanstack/react-router';
import { ArrowRight, Box, Database, GitBranch, Plus, Server, Workflow } from 'lucide-react';
import { Button } from '#/components/ui/button';

const appCapabilities = [
  { label: 'Web service', detail: 'Build and deploy an application', icon: Box },
  { label: 'Background worker', detail: 'Run asynchronous workloads', icon: Workflow },
  { label: 'Database', detail: 'Provision application data', icon: Database },
  { label: 'Docker image', detail: 'Deploy a ready-made container', icon: Server },
  { label: 'Git repository', detail: 'Build from a source branch', icon: GitBranch },
];

export function AppEmptyState() {
  return (
    <section className="py-12 sm:py-16">
      <div className="text-center">
        <div className="mx-auto flex items-center justify-center gap-2" aria-hidden="true">
          {appCapabilities.map(({ icon: Icon }, index) => (
            <div key={index} className="flex items-center gap-2">
              {index > 0 && <span className="w-8 border-border border-t border-dashed" />}
              <span className="flex h-13 w-13 items-center justify-center rounded-2xl bg-card text-muted-foreground">
                <Icon className="h-5 w-5" />
              </span>
            </div>
          ))}
          <span className="mx-1 w-8 border-border border-t border-dashed" />
          <span className="flex h-13 w-13 items-center justify-center rounded-2xl border border-primary/45 border-dashed bg-primary/8 text-primary">
            <Plus className="h-5 w-5" />
          </span>
        </div>
        <h2 className="mt-7 font-medium text-foreground/90 text-xl tracking-[-0.02em]">
          Deploy your first app
        </h2>
        <p className="mx-auto mt-2 max-w-md text-muted-foreground/75 text-sm leading-6">
          Choose a runtime and connect it to a project. Codedock will keep its releases and health
          in one place.
        </p>
        <Link to="/apps/new">
          <Button className="mt-6 gap-2 px-5">
            <Plus className="h-4 w-4" />
            Deploy app
          </Button>
        </Link>
      </div>
      <div className="mx-auto mt-10 grid max-w-2xl gap-3 sm:grid-cols-3">
        {appCapabilities.map(({ label, detail, icon: Icon }) => (
          <Link
            key={label}
            to="/apps/new"
            className="group flex items-center gap-3 rounded-xl bg-card px-4 py-3.5 text-left transition-colors hover:bg-muted/60"
          >
            <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-muted/70 text-muted-foreground">
              <Icon className="h-4 w-4" />
            </span>
            <span className="min-w-0 flex-1">
              <span className="block font-medium text-sm">{label}</span>
              <span className="mt-0.5 block truncate text-muted-foreground text-xs">{detail}</span>
            </span>
            <ArrowRight className="h-4 w-4 text-muted-foreground transition-colors group-hover:text-primary" />
          </Link>
        ))}
        <Link
          to="/apps/new"
          className="group flex items-center gap-3 rounded-xl border border-border/80 border-dashed bg-card/60 px-4 py-3.5 text-left transition-colors hover:bg-card"
        >
          <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-muted/70 text-muted-foreground">
            <Plus className="h-4 w-4" />
          </span>
          <span className="min-w-0 flex-1">
            <span className="block font-medium text-sm">Explore deployment options</span>
            <span className="mt-0.5 block truncate text-muted-foreground text-xs">
              Choose a project to continue
            </span>
          </span>
          <ArrowRight className="h-4 w-4 text-muted-foreground transition-colors group-hover:text-primary" />
        </Link>
      </div>
    </section>
  );
}
