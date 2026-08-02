import type * as React from 'react';
import { cn } from '#/lib/utils';
import { ContextRail } from './context-rail';

type PageFrameProps = {
  children: React.ReactNode;
  rail?: React.ReactNode;
  className?: string;
  mainClassName?: string;
};

export function PageFrame({ children, rail, className, mainClassName }: PageFrameProps) {
  if (!rail) {
    return <div className={cn('min-w-0', className)}>{children}</div>;
  }

  return (
    <div className={cn('grid min-w-0 gap-6 lg:grid-cols-[minmax(0,1fr)_20rem]', className)}>
      <section className={cn('min-w-0', mainClassName)}>{children}</section>
      <ContextRail>{rail}</ContextRail>
    </div>
  );
}
