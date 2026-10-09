import { Link } from '@tanstack/react-router';
import { Database, GitBranch, Globe2, Plus, Rocket } from 'lucide-react';
import { Button } from '#/components/ui/button';

export function HomeFirstProject({ onCreateProject }: { onCreateProject: () => void }) {
  return (
    <section className="px-6 pt-5 pb-7 text-center sm:pt-7 sm:pb-8">
      <div className="relative mx-auto h-36 w-72 max-w-full" aria-hidden="true">
        <span className="absolute top-4 left-4 flex h-11 w-11 items-center justify-center rounded-xl border border-border bg-card text-muted-foreground">
          <GitBranch className="h-5 w-5" />
        </span>
        <span className="absolute top-4 right-4 flex h-11 w-11 items-center justify-center rounded-xl border border-border bg-card text-muted-foreground">
          <Globe2 className="h-5 w-5" />
        </span>
        <span className="absolute bottom-3 left-9 flex h-11 w-11 items-center justify-center rounded-xl border border-border bg-card text-muted-foreground">
          <Rocket className="h-5 w-5" />
        </span>
        <span className="absolute right-9 bottom-3 flex h-11 w-11 items-center justify-center rounded-xl border border-border bg-card text-muted-foreground">
          <Database className="h-5 w-5" />
        </span>
        <span className="absolute top-1/2 left-1/2 h-20 w-20 -translate-x-1/2 -translate-y-1/2 rounded-full border border-primary/25 bg-primary/8" />
        <span className="absolute top-1/2 left-1/2 flex h-14 w-14 -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-2xl bg-primary text-primary-foreground shadow-lg shadow-primary/20">
          <Plus className="h-6 w-6" />
        </span>
      </div>
      <h2
        className="mt-4 font-medium text-foreground/85 text-xl"
        style={{ letterSpacing: '-0.2px' }}
      >
        Ship your first project
      </h2>
      <p className="mx-auto mt-1.5 mb-5 max-w-sm text-muted-foreground/80 text-sm leading-relaxed">
        Connect a repository and Codedock builds, deploys, and secures it. No pipelines to wire by
        hand.
      </p>
      <div className="flex flex-col items-center justify-center gap-3 sm:flex-row">
        <Button className="gap-2 px-6" onClick={onCreateProject}>
          <Plus className="h-4 w-4" />
          Create project
        </Button>
        <Button asChild variant="secondary" className="gap-2 px-6">
          <Link to="/library">
            <GitBranch className="h-4 w-4" />
            Import from Git
          </Link>
        </Button>
      </div>
      <p className="mt-5 text-muted-foreground/60 text-xs">
        Tip: press{' '}
        <kbd className="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px]">&#8984; K</kbd> to
        jump anywhere
      </p>
    </section>
  );
}
