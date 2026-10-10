import { Link } from '@tanstack/react-router';
import { ArrowRight, Plus } from 'lucide-react';

export function HomeNextStep({ hasProjects }: { hasProjects: boolean }) {
  const title = hasProjects ? 'Keep building' : 'Your next move';
  const description = hasProjects
    ? 'Add an application or review the latest deployment activity.'
    : 'Start from a repository, folder, or ready-to-run template.';
  const href = hasProjects ? '/projects' : '/library';
  const action = hasProjects ? 'View projects' : 'Create project';

  return (
    <section className="rounded-2xl bg-card p-5">
      <div className="flex items-center gap-2">
        <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
          <Plus className="h-4 w-4" />
        </span>
        <h2 className="font-semibold text-sm">{title}</h2>
      </div>
      <p className="mt-4 text-muted-foreground/80 text-sm leading-6">{description}</p>
      <Link
        to={href}
        className="mt-4 inline-flex items-center gap-1.5 font-medium text-foreground text-sm transition-colors hover:text-primary"
      >
        {action}
        <ArrowRight className="h-3.5 w-3.5" />
      </Link>
    </section>
  );
}
