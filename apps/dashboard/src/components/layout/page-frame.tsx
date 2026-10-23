import type * as React from 'react';
import { cn } from '#/lib/utils';
import { useAiSidebarStore } from '#/stores/ai-sidebar-store';
import { ContextRail } from './context-rail';

type PageFrameProps = {
  children: React.ReactNode;
  rail?: React.ReactNode;
  className?: string;
  mainClassName?: string;
  hideRailWhenStacked?: boolean;
};

export function PageFrame({
  children,
  rail,
  className,
  mainClassName,
  hideRailWhenStacked = false,
}: PageFrameProps) {
  const aiSidebarOpen = useAiSidebarStore((s) => s.isOpen);

  if (!rail) {
    return <div className={cn('min-w-0', className)}>{children}</div>;
  }

  return (
    <div
      className={cn(
        'grid min-w-0 grid-cols-1 items-start gap-6',
        aiSidebarOpen
          ? '2xl:grid-cols-[minmax(0,1fr)_280px]'
          : 'lg:grid-cols-[minmax(0,1fr)_21.25rem]',
        className
      )}
    >
      <section className={cn('min-w-0', mainClassName)}>{children}</section>
      <ContextRail hideWhenStacked={hideRailWhenStacked}>{rail}</ContextRail>
    </div>
  );
}
