import { Button } from '@codedock/ui/button';
import { Link } from '@tanstack/react-router';
import { ArrowRight, Plus } from 'lucide-react';
import { useEffect, useId, useRef, useState } from 'react';
import { AppLogo } from '#/features/apps/app-logo';
import { useInstalledApps } from '#/hooks/use-apps';

const FEATURED_APPS = [
  'supabase',
  'pocketbase',
  'convex',
  'minio',
  'grafana',
  'neon',
  'qdrant',
  'gitea',
  'vaultwarden',
  'code-server',
  'posthog',
  'mail',
];

const TILES = [
  { position: 'start-0 top-7 size-7 -rotate-12', logo: 'size-4' },
  { position: 'start-10 top-1 size-9 rotate-6', logo: 'size-5' },
  { position: 'start-24 top-3 size-11 -rotate-3', logo: 'size-7' },
  { position: 'end-10 top-0 size-8 rotate-12', logo: 'size-5' },
  { position: 'end-0 top-8 size-7 -rotate-6', logo: 'size-4' },
];

function HomeAppsIllustration({ loading = false }: { loading?: boolean }) {
  const rootRef = useRef<HTMLDivElement>(null);
  const [apps, setApps] = useState(() => FEATURED_APPS.slice(0, TILES.length));

  useEffect(() => {
    const element = rootRef.current;
    if (!element || typeof IntersectionObserver === 'undefined') return;
    const motion = window.matchMedia('(prefers-reduced-motion: reduce)');
    let visible = false;
    let previousSlot = -1;
    let timer: ReturnType<typeof setInterval> | undefined;

    const stop = () => {
      clearInterval(timer);
      timer = undefined;
    };
    const sync = () => {
      stop();
      if (!visible || document.hidden || motion.matches) return;
      timer = setInterval(() => {
        const slot =
          (previousSlot + 1 + Math.floor(Math.random() * (TILES.length - 1))) % TILES.length;
        const choice = Math.random();
        previousSlot = slot;
        setApps((current) => {
          const available = FEATURED_APPS.filter((app) => !current.includes(app));
          return current.map((app, index) =>
            index === slot ? available[Math.floor(choice * available.length)] : app
          );
        });
      }, 5000);
    };

    const observer = new IntersectionObserver(([entry]) => {
      visible = Boolean(entry?.isIntersecting);
      sync();
    });
    observer.observe(element);
    document.addEventListener('visibilitychange', sync);
    motion.addEventListener('change', sync);
    return () => {
      stop();
      observer.disconnect();
      document.removeEventListener('visibilitychange', sync);
      motion.removeEventListener('change', sync);
    };
  }, []);

  return (
    <div
      ref={rootRef}
      aria-hidden="true"
      aria-busy={loading || undefined}
      className="relative mx-auto mb-2 h-16 w-full max-w-56"
    >
      <svg
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 size-full rtl:-scale-x-100"
        viewBox="0 0 224 64"
        fill="none"
      >
        <path
          d="M14 42L58 22L118 34L168 16L210 46"
          stroke="var(--border)"
          strokeWidth="1.5"
          strokeDasharray="3 4"
          strokeLinecap="round"
        />
        <circle cx="18" cy="12" r="2" fill="var(--border)" />
        <circle cx="84" cy="54" r="2" fill="var(--border)" />
        <circle cx="194" cy="8" r="2.5" fill="var(--border)" />
      </svg>
      {TILES.map((tile, index) => {
        const appId = apps[index];
        return (
          <div
            key={index}
            aria-hidden="true"
            className={`absolute flex items-center justify-center rounded-xl bg-card shadow-sm ring-1 ring-border/60 ${tile.position}`}
          >
            {loading ? (
              <span className={`rounded-md bg-muted motion-safe:animate-pulse ${tile.logo}`} />
            ) : (
              <span
                key={appId}
                className={`flex items-center justify-center rounded-md ${tile.logo}`}
              >
                <AppLogo appId={appId} className="size-full" />
              </span>
            )}
          </div>
        );
      })}
    </div>
  );
}

export function HomeAppInventory({ isLoading: parentLoading = false }: { isLoading?: boolean }) {
  const headingId = useId();
  const installedQuery = useInstalledApps();
  const loading = parentLoading || installedQuery.isLoading;
  const apps = installedQuery.data ?? [];
  const hasApps = !loading && apps.length > 0;

  return (
    <section
      aria-labelledby={headingId}
      aria-busy={loading}
      className="rounded-2xl border border-border/50 bg-card p-5 shadow-xs"
    >
      <div className="mb-3 flex items-center justify-between gap-3">
        <div className="flex min-w-0 items-center gap-2">
          <h3 id={headingId} className="font-semibold text-foreground text-sm">
            <Link to="/apps" className="hover:text-muted-foreground">
              Apps
            </Link>
          </h3>
          {hasApps && (
            <span className="text-muted-foreground text-xs tabular-nums">{apps.length}</span>
          )}
        </div>
        {hasApps && (
          <Button asChild variant="ghost" size="icon" className="size-7 rounded-lg">
            <Link to="/apps/new" aria-label="New app" title="New app">
              <Plus className="size-4" aria-hidden="true" />
            </Link>
          </Button>
        )}
      </div>

      {loading ? (
        <p className="text-muted-foreground text-sm">Loading…</p>
      ) : apps.length === 0 ? (
        <div className="text-center">
          <HomeAppsIllustration />
          <p className="font-medium text-foreground text-sm">No apps installed</p>
          <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
            Deploy ready-to-use apps like Supabase, Grafana, and PocketBase in one click.
          </p>
        </div>
      ) : (
        <ul className="max-h-60 divide-y divide-border/40 overflow-y-auto rounded-xl bg-muted/30">
          {apps.map((app) => {
            const online = (app.status ?? '').toLowerCase() === 'running';
            return (
              <li
                key={app.id}
                className="relative flex min-w-0 items-center gap-3 px-3 py-2.5 transition-colors hover:bg-muted/50"
              >
                <Link
                  to="/projects/$projectId"
                  params={{ projectId: app.id }}
                  aria-label={app.name}
                  className="absolute inset-0 rounded-xl focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-inset"
                />
                <div className="flex size-9 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-background">
                  <AppLogo appId={app.appId} className="size-5" />
                </div>
                <p
                  className="min-w-0 flex-1 truncate font-medium text-foreground text-sm"
                  title={app.name}
                >
                  {app.name}
                </p>
                <span
                  className={`shrink-0 rounded-md px-1.5 py-0.5 font-medium text-xs capitalize ${
                    online
                      ? 'bg-emerald-500/10 text-emerald-500'
                      : 'bg-muted text-muted-foreground'
                  }`}
                >
                  {app.status || 'deployed'}
                </span>
              </li>
            );
          })}
        </ul>
      )}

      {!hasApps && (
        <Button asChild variant="secondary" size="sm" className="mt-4 w-full text-sm">
          <Link to="/apps/new">
            Browse all apps
            <ArrowRight className="size-4 rtl:rotate-180" aria-hidden="true" />
          </Link>
        </Button>
      )}
    </section>
  );
}
