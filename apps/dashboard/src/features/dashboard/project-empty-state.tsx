import { Link } from '@tanstack/react-router';
import { GitBranch, Globe, MousePointerClick, Plus, RotateCcw, Zap } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { EmptyIllustration } from './empty-illustration';

const highlights = [
  { icon: Zap, title: 'Instant deploys', description: 'Push and go live' },
  { icon: Globe, title: 'Domains & TLS', description: 'Auto certificates' },
  { icon: MousePointerClick, title: 'Previews', description: 'Per-branch URLs' },
  { icon: RotateCcw, title: 'Rollbacks', description: 'One-click restore' },
];

export function ProjectEmptyState() {
  return (
    <div className="py-16 text-center">
      <EmptyIllustration className="relative mx-auto mb-8 h-44 w-64" />
      <h3
        className="mb-2 font-medium text-2xl text-foreground/80"
        style={{ letterSpacing: '-0.2px' }}
      >
        Create your first project
      </h3>
      <p className="mx-auto mb-8 max-w-sm text-muted-foreground/70 text-sm leading-relaxed">
        Bring together the services, environments, and deployment activity that belong to one
        product.
      </p>
      <div className="mb-10 flex flex-col items-center justify-center gap-3 sm:flex-row">
        <Button asChild size="lg" className="gap-2 px-6">
          <Link to="/projects/new">
            <Plus className="h-4 w-4" />
            New project
          </Link>
        </Button>
        <Button asChild size="lg" variant="secondary" className="gap-2 px-6">
          <Link to="/library">
            <GitBranch className="h-4 w-4" />
            Import repository
          </Link>
        </Button>
      </div>
      <div className="mx-auto max-w-2xl">
        <p className="mb-4 text-muted-foreground/60 text-xs uppercase tracking-wider">
          Zero-config deployments
        </p>
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
          {highlights.map((highlight) => (
            <div
              key={highlight.title}
              className="rounded-xl border border-border/50 bg-card p-4 text-start"
            >
              <div className="mb-3 flex h-8 w-8 items-center justify-center rounded-lg bg-muted">
                <highlight.icon className="h-4 w-4 text-muted-foreground" />
              </div>
              <p className="font-medium text-foreground text-sm">{highlight.title}</p>
              <p className="mt-0.5 text-muted-foreground text-xs">{highlight.description}</p>
            </div>
          ))}
        </div>
      </div>
      <p className="mt-8 text-muted-foreground/60 text-xs">
        Press <kbd className="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px]">&#8984; K</kbd>{' '}
        to jump anywhere
      </p>
    </div>
  );
}
