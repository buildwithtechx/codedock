import { Link } from '@tanstack/react-router';
import { Eye, GitBranch, Globe, Plus, RotateCcw, Zap } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { EmptyIllustration } from './empty-illustration';

const highlights = [
  { icon: Zap, title: 'Instant setup', description: 'Auto-detected stack and configuration' },
  { icon: Globe, title: 'Custom domains', description: 'Free automatic SSL certificates' },
  {
    icon: Eye,
    title: 'Live previews',
    description: 'Git branch previews with isolated environments',
  },
  {
    icon: RotateCcw,
    title: 'Zero-downtime rollbacks',
    description: 'Instant restore to any previous version',
  },
];

export function ProjectEmptyState() {
  return (
    <div className="py-16 text-center">
      <EmptyIllustration className="relative mx-auto mb-8 h-44 w-64" />
      <h3
        className="mb-2 font-medium text-2xl text-foreground/80"
        style={{ letterSpacing: '-0.2px' }}
      >
        No projects yet
      </h3>
      <p className="mx-auto mb-8 max-w-sm text-muted-foreground/70 text-sm leading-relaxed">
        Deploy a Git repository, Docker container, or start with one of our ready-to-run templates.
      </p>
      <div className="mb-10 flex flex-col items-center justify-center gap-3 sm:flex-row">
        <Button asChild size="lg" className="gap-2 px-6">
          <Link to="/library">
            <Plus className="h-4 w-4" />
            Create project
          </Link>
        </Button>
        <Button asChild size="lg" variant="secondary" className="gap-2 px-6">
          <Link to="/library" search={{ tab: 'examples' }}>
            <GitBranch className="h-4 w-4" />
            Browse templates
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
        Press <kbd className="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px]">⌘ K</kbd> to
        open the command palette
      </p>
    </div>
  );
}
