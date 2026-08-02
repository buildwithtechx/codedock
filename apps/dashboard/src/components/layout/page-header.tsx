import type * as React from 'react';
import { cn } from '#/lib/utils';

type PageHeaderProps = {
  title: string;
  description?: string;
  action?: React.ReactNode;
  eyebrow?: string;
  className?: string;
};

export function PageHeader({ title, description, action, eyebrow, className }: PageHeaderProps) {
  return (
    <header
      className={cn('flex flex-col justify-between gap-4 sm:flex-row sm:items-end', className)}
    >
      <div className="min-w-0">
        {eyebrow && (
          <p className="mb-1 font-medium text-muted-foreground/70 text-xs uppercase tracking-[0.14em]">
            {eyebrow}
          </p>
        )}
        <h1 className="font-medium text-2xl text-foreground/90 tracking-[-0.02em]">{title}</h1>
        {description && <p className="mt-1 text-muted-foreground/75 text-sm">{description}</p>}
      </div>
      {action && <div className="shrink-0 self-start sm:self-auto">{action}</div>}
    </header>
  );
}
