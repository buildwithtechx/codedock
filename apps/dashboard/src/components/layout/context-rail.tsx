import type * as React from 'react';
import { cn } from '#/lib/utils';

type ContextRailProps = {
  children: React.ReactNode;
  className?: string;
};

export function ContextRail({ children, className }: ContextRailProps) {
  return (
    <aside className={cn('hidden xl:sticky xl:top-6 xl:block xl:self-start', className)}>
      {children}
    </aside>
  );
}
