import type * as React from 'react';
import { cn } from '#/lib/utils';

type ContextRailProps = {
  children: React.ReactNode;
  className?: string;
};

export function ContextRail({ children, className }: ContextRailProps) {
  return (
    <aside className={cn('hidden lg:sticky lg:top-6 lg:block lg:self-start', className)}>
      {children}
    </aside>
  );
}
