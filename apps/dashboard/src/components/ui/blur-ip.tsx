import type * as React from 'react';
import { useDemoMode } from '#/lib/demo-mode';
import { cn } from '#/lib/utils';

type BlurIpProps = React.HTMLAttributes<HTMLSpanElement> & {
  children: React.ReactNode;
};

export function BlurIp({ children, className, ...props }: BlurIpProps) {
  const [demoMode] = useDemoMode();

  return (
    <span
      data-demo-blur={demoMode ? 'true' : undefined}
      className={cn(
        'inline-block transition-[filter] duration-200',
        demoMode && 'select-none blur-[5px] hover:blur-none',
        className
      )}
      {...props}
    >
      {children}
    </span>
  );
}
