import { cva, type VariantProps } from 'class-variance-authority';
import type * as React from 'react';
import { cn } from 'cn';

export const statusDotVariants = cva('rounded-full shrink-0', {
  variants: {
    status: {
      ready: 'bg-emerald-500',
      running: 'bg-blue-500',
      warning: 'bg-amber-500',
      error: 'bg-red-500',
      pending: 'bg-yellow-500',
      muted: 'bg-zinc-400 dark:bg-zinc-600',
    },
    size: {
      sm: 'size-1.5',
      md: 'size-2',
      lg: 'size-2.5',
    },
  },
  defaultVariants: {
    status: 'ready',
    size: 'md',
  },
});

export interface StatusDotProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof statusDotVariants> {
  pulse?: boolean;
}

export function StatusDot({ status, size, pulse = false, className, ...props }: StatusDotProps) {
  return (
    <span className="relative flex items-center justify-center">
      {pulse && (
        <span
          className={cn(
            'absolute inline-flex animate-ping rounded-full opacity-75',
            statusDotVariants({ status, size })
          )}
        />
      )}
      <span
        className={cn('relative inline-flex', statusDotVariants({ status, size }), className)}
        {...props}
      />
    </span>
  );
}
