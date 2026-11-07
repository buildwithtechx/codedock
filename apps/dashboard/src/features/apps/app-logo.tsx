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
  'linux',
]);

const PNG_LOGOS = new Set([
  'homepage',
  'it-tools',
  'wiki-js',
  'ollama',
  'openwebui',
  'open-webui',
  'vllm',
]);

function resolveLogoSrc(icon: string): string {
  if (icon.startsWith('http://') || icon.startsWith('https://') || icon.startsWith('/')) {
    return icon;
  }
  if (icon.endsWith('.png') || icon.endsWith('.svg')) {
    return `/app-logos/${icon}`;
  }
  const clean = icon.toLowerCase().replace(/[^a-z0-9_-]/g, '');
  if (PNG_LOGOS.has(clean) || PNG_LOGOS.has(icon.toLowerCase())) {
    return `/app-logos/${icon}.png`;
  }
  return `/app-logos/${icon}.svg`;
}

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
        src={resolveLogoSrc(resolvedIcon)}
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
