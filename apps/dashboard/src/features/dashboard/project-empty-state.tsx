import { Link } from '@tanstack/react-router';
import { ArrowRight, FolderPlus, Plus } from 'lucide-react';
import { Button } from '#/components/ui/button';

export function ProjectEmptyState() {
  return (
    <section className="flex min-h-[30rem] items-center justify-center py-10 text-center">
      <div className="max-w-lg">
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
        <Link to="/projects/new">
          <Button className="mt-6 gap-2 px-5">
            <Plus className="h-4 w-4" />
            New project
            <ArrowRight className="h-4 w-4" />
          </Button>
        </Link>
      </div>
    </section>
  );
}
