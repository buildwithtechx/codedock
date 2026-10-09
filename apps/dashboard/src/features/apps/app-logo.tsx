import { useState } from 'react';
import { cn } from '#/lib/utils';

export function AppLogo({
  icon,
  name,
  className,
}: {
  icon?: string;
  name: string;
  className?: string;
}) {
  const [failed, setFailed] = useState(false);
  if (!icon || failed) {
    return (
      <span
        aria-hidden="true"
        className={cn(
          'flex items-center justify-center rounded-xl bg-primary/10 font-semibold text-primary uppercase',
          className ?? 'size-10 text-sm'
        )}
      >
        {name.trim().charAt(0) || '?'}
      </span>
    );
  }
  return (
    <span
      className={cn(
        'flex items-center justify-center overflow-hidden rounded-xl bg-muted/60',
        className ?? 'size-10'
      )}
    >
      <img
        src={`/app-logos/${icon}.svg`}
        alt=""
        aria-hidden="true"
        className="size-3/5 object-contain"
        onError={() => setFailed(true)}
      />
    </span>
  );
}
