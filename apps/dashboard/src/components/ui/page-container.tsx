import type * as React from 'react';
import { cn } from '#/lib/utils';

interface PageContainerProps {
  children: React.ReactNode;
  className?: string;
  outerClassName?: string;
  fullScreen?: boolean;
}

export function PageContainer({
  children,
  className,
  outerClassName,
  fullScreen = false,
}: PageContainerProps) {
  const inner = (
    <div className={cn('mx-auto max-w-[1600px] px-4 py-6 sm:px-6 lg:px-8', className)}>
      {children}
    </div>
  );

  if (!fullScreen) return inner;

  return <div className={cn('min-h-screen bg-background', outerClassName)}>{inner}</div>;
}
