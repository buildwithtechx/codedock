import { Cpu } from 'lucide-react';
import { cn } from '#/lib/utils';

export function AppGpuBadge({
  gpu,
  className,
}: {
  gpu?: 'required' | 'recommended' | string | null;
  className?: string;
}) {
  if (!gpu) return null;

  const label = gpu === 'required' ? 'GPU Required' : 'GPU Recommended';

  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 rounded-md px-1.5 py-0.5 font-medium text-[11px] leading-none',
        gpu === 'required'
          ? 'bg-amber-500/15 text-amber-600 dark:text-amber-400'
          : 'bg-muted/70 text-muted-foreground',
        className
      )}
    >
      <Cpu className="size-3 shrink-0" aria-hidden="true" />
      <span>{label}</span>
    </span>
  );
}
