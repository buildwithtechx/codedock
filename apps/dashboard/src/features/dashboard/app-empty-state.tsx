import { Link } from '@tanstack/react-router';
import { Box, Database, Plus, Workflow } from 'lucide-react';
import { Button } from '#/components/ui/button';

const appCapabilities = [
  { label: 'Web service', icon: Box },
  { label: 'Database', icon: Database },
  { label: 'Worker', icon: Workflow },
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
        {appCapabilities.map(({ label, icon: Icon }) => (
          <div
            key={label}
            className="flex items-center gap-3 rounded-xl bg-card px-4 py-3.5 text-left"
          >
            <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-muted/70 text-muted-foreground">
              <Icon className="h-4 w-4" />
            </span>
            <p className="font-medium text-sm">{label}</p>
          </div>
        ))}
      </div>
    </section>
  );
}
