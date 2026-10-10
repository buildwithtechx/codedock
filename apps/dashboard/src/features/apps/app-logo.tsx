import { useState } from 'react';
import { cn } from '#/lib/utils';

export function AppLogo({
  appId,
  icon,
  name,
  className,
}: {
  appId?: string;
  icon?: string;
  name?: string;
  className?: string;
}) {
  const [failed, setFailed] = useState(false);
  const resolvedIcon = icon || appId;
  const resolvedName = name || appId || '';
  if (!resolvedIcon || failed) {
    return (
      <span
        aria-hidden="true"
        className={cn(
          'flex items-center justify-center rounded-xl bg-primary/10 font-semibold text-primary uppercase',
          className ?? 'size-10 text-sm'
        )}
      >
        {resolvedName.trim().charAt(0) || '?'}
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
        src={`/app-logos/${resolvedIcon}.svg`}
        alt=""
        aria-hidden="true"
        className="size-3/5 object-contain"
        onError={() => setFailed(true)}
      />
    </span>
  );
}
