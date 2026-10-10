import { useState } from 'react';
import { cn } from '#/lib/utils';

const DARK_INVERT_LOGOS = new Set([
  'buzz',
  'calcom',
  'directus',
  'ghost',
  'kafka',
  'posthog',
  'umami',
  'vaultwarden',
]);

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

  const cleanIcon = resolvedIcon.toLowerCase().replace(/[^a-z0-9]/g, '');
  const cleanAppId = (appId ?? '').toLowerCase().replace(/[^a-z0-9]/g, '');
  const shouldInvert = DARK_INVERT_LOGOS.has(cleanIcon) || DARK_INVERT_LOGOS.has(cleanAppId);

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
        className={cn(
          'size-3/5 object-contain',
          shouldInvert && 'dim:brightness-0 dim:invert dark:brightness-0 dark:invert'
        )}
        onError={() => setFailed(true)}
      />
    </span>
  );
}
