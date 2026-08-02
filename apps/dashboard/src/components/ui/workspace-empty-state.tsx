import type { LucideIcon } from 'lucide-react';
import type * as React from 'react';
import { cn } from '#/lib/utils';

type WorkspaceEmptyStateProps = {
  icon: LucideIcon;
  title: string;
  description: string;
  action?: React.ReactNode;
  secondaryAction?: React.ReactNode;
  className?: string;
};

export function WorkspaceEmptyState({
  icon: Icon,
  title,
  description,
  action,
  secondaryAction,
  className,
}: WorkspaceEmptyStateProps) {
  return (
    <section
      className={cn(
        'flex min-h-[25rem] items-center justify-center px-6 py-12 text-center',
        className
      )}
    >
      <div className="max-w-md">
        <div
          className="relative mx-auto mb-7 flex h-28 w-44 items-center justify-center"
          aria-hidden="true"
        >
          <span className="absolute top-3 left-3 h-2 w-2 rounded-full bg-primary/35" />
          <span className="absolute right-5 bottom-5 h-2.5 w-2.5 rounded-full bg-primary/20" />
          <span className="absolute top-6 right-7 h-1.5 w-1.5 rounded-full bg-primary/50" />
          <span className="absolute top-1/2 right-12 left-12 border-primary/30 border-t border-dashed" />
          <span className="absolute inset-x-9 top-5 bottom-5 rounded-[2rem] border border-primary/10" />
          <span className="relative flex h-14 w-14 items-center justify-center rounded-2xl border border-primary/25 bg-primary/10 text-primary shadow-[0_0_0_9px_color-mix(in_oklch,var(--primary)_8%,transparent)]">
            <Icon className="h-6 w-6" />
          </span>
        </div>
        <h2 className="font-semibold text-xl tracking-tight">{title}</h2>
        <p className="mx-auto mt-2 max-w-sm text-muted-foreground text-sm leading-6">
          {description}
        </p>
        {(action || secondaryAction) && (
          <div className="mt-6 flex flex-col items-center justify-center gap-2 sm:flex-row">
            {action}
            {secondaryAction}
          </div>
        )}
      </div>
    </section>
  );
}
