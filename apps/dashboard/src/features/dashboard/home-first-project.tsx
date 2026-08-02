import { Link } from '@tanstack/react-router';
import { Database, GitBranch, Globe2, Plus, Rocket } from 'lucide-react';
import { Button } from '#/components/ui/button';

export function HomeFirstProject({ onCreateProject }: { onCreateProject: () => void }) {
  return (
    <section className="flex min-h-[25rem] items-center justify-center px-6 py-10 text-center">
      <div className="max-w-lg">
        <div className="relative mx-auto h-42 w-72" aria-hidden="true">
          <span className="absolute top-5 left-5 flex h-11 w-11 items-center justify-center rounded-xl bg-background text-muted-foreground">
            <GitBranch className="h-5 w-5" />
          </span>
          <span className="absolute top-5 right-5 flex h-11 w-11 items-center justify-center rounded-xl bg-background text-muted-foreground">
            <Globe2 className="h-5 w-5" />
          </span>
          <span className="absolute bottom-4 left-10 flex h-11 w-11 items-center justify-center rounded-xl bg-background text-muted-foreground">
            <Rocket className="h-5 w-5" />
          </span>
          <span className="absolute right-10 bottom-4 flex h-11 w-11 items-center justify-center rounded-xl bg-background text-muted-foreground">
            <Database className="h-5 w-5" />
          </span>
          <span className="absolute top-1/2 left-1/2 h-20 w-20 -translate-x-1/2 -translate-y-1/2 rounded-full border border-primary/25 bg-primary/8" />
          <span className="absolute top-1/2 left-1/2 flex h-14 w-14 -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-2xl bg-primary text-primary-foreground shadow-lg shadow-primary/20">
            <Plus className="h-6 w-6" />
          </span>
          <span className="absolute top-1/2 left-13 h-px w-16 -rotate-24 bg-border" />
          <span className="absolute top-1/2 right-13 h-px w-16 rotate-24 bg-border" />
          <span className="absolute bottom-11 left-18 h-px w-16 rotate-27 bg-border" />
          <span className="absolute right-18 bottom-11 h-px w-16 -rotate-27 bg-border" />
        </div>
        <h2 className="mt-6 font-medium text-foreground/90 text-xl tracking-[-0.02em]">
          Build the first project
        </h2>
        <p className="mx-auto mt-2 max-w-md text-muted-foreground/75 text-sm leading-6">
          Projects are the operational boundary for your applications, environments, releases, and
          domains.
        </p>
        <div className="mt-6 flex flex-col items-center justify-center gap-2 sm:flex-row">
          <Button className="gap-2 px-5" onClick={onCreateProject}>
            <Plus className="h-4 w-4" />
            Create project
          </Button>
          <Link
            to="/settings"
            search={{ tab: 'sources' }}
            className="inline-flex h-9 items-center justify-center rounded-xl bg-muted px-4 font-medium text-sm transition-colors hover:bg-secondary"
          >
            Connect a source
          </Link>
        </div>
      </div>
    </section>
  );
}
