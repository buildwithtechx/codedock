import { Link } from '@tanstack/react-router';
import {
  ArrowRight,
  Boxes,
  FolderPlus,
  GitBranch,
  Globe2,
  LayoutTemplate,
  Plus,
  Rocket,
} from 'lucide-react';
import { Button } from '#/components/ui/button';

const projectCapabilities = [
  { title: 'Git-driven', detail: 'Connect source control', icon: GitBranch },
  { title: 'Applications', detail: 'Group runtime services', icon: Boxes },
  { title: 'Release history', detail: 'Track every deployment', icon: Rocket },
  { title: 'Domains', detail: 'Manage DNS and routing', icon: Globe2 },
];

export function ProjectEmptyState() {
  return (
    <section className="flex min-h-[34rem] items-center justify-center py-12 text-center">
      <div className="w-full max-w-3xl">
        <div className="relative mx-auto h-32 w-56" aria-hidden="true">
          <span className="absolute top-8 left-7 h-16 w-24 rounded-2xl bg-card" />
          <span className="absolute top-3 left-16 h-18 w-28 rounded-2xl border border-border/70 bg-card" />
          <span className="absolute top-10 right-6 flex h-14 w-14 items-center justify-center rounded-2xl border border-primary/30 bg-primary/10 text-primary">
            <FolderPlus className="h-6 w-6" />
          </span>
          <span className="absolute top-1/2 right-1 h-7 w-7 rounded-full border border-primary/45 border-dashed" />
          <Plus className="absolute top-[57px] right-[17px] h-3 w-3 text-primary" />
        </div>
        <h2 className="mt-5 font-medium text-foreground/90 text-xl tracking-[-0.02em]">
          Create your first project
        </h2>
        <p className="mx-auto mt-2 max-w-md text-muted-foreground/75 text-sm leading-6">
          Bring together the services, environments, and deployment activity that belong to one
          product.
        </p>
        <div className="mt-6 flex flex-col items-center justify-center gap-2 sm:flex-row">
          <Link to="/projects/new">
            <Button className="gap-2 px-5">
              <Plus className="h-4 w-4" />
              New project
              <ArrowRight className="h-4 w-4" />
            </Button>
          </Link>
          <Link to="/projects/new" search={{ template: 'one-click' }}>
            <Button variant="secondary" className="gap-2 px-5">
              <LayoutTemplate className="h-4 w-4" />
              Browse templates
            </Button>
          </Link>
        </div>
        <div className="mt-11">
          <p className="mb-4 font-medium text-[10px] text-muted-foreground uppercase tracking-[0.16em]">
            A project keeps the work together
          </p>
          <div className="grid gap-3 text-left sm:grid-cols-4">
            {projectCapabilities.map(({ title, detail, icon: Icon }) => (
              <div key={title} className="rounded-2xl bg-card p-4">
                <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-muted text-muted-foreground">
                  <Icon className="h-4 w-4" />
                </span>
                <p className="mt-4 font-semibold text-sm">{title}</p>
                <p className="mt-1 text-muted-foreground text-xs leading-5">{detail}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
