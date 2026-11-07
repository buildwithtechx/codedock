import type * as React from 'react';
import { cn } from '#/lib/utils';
import { useAiSidebarStore } from '#/stores/ai-sidebar-store';

type ContextRailProps = {
  children: React.ReactNode;
  className?: string;
  hideWhenStacked?: boolean;
};

export function ContextRail({ children, className, hideWhenStacked = false }: ContextRailProps) {
  const aiSidebarOpen = useAiSidebarStore((s) => s.isOpen);

  return (
    <aside
      className={cn(
        aiSidebarOpen
          ? cn(
              hideWhenStacked ? 'hidden 2xl:block' : 'block',
              '2xl:sticky 2xl:top-6 2xl:max-h-[calc(100dvh-3rem)] 2xl:self-start 2xl:overflow-y-auto'
            )
          : cn(
              hideWhenStacked ? 'hidden lg:block' : 'block',
              'lg:sticky lg:top-6 lg:max-h-[calc(100dvh-3rem)] lg:self-start lg:overflow-y-auto'
            ),
        className
      )}
    >
      {children}
    </aside>
  );
}
