import { Link } from '@tanstack/react-router';
import { ChevronRight } from 'lucide-react';
import type { AppService } from '#/features/services';
import { cn } from '#/lib/utils';
import { AppLogo } from './app-logo';

const statusTone: Record<string, string> = {
  running: 'bg-emerald-500',
  building: 'bg-amber-500',
  created: 'bg-amber-500',
  error: 'bg-destructive',
  stopped: 'bg-zinc-500',
};

export function InstalledAppRow({ app }: { app: AppService }) {
  const meta = [app.branch, app.domain].filter(Boolean).join(' · ');
  return (
    <Link
      to="/projects/$projectId"
      params={{ projectId: app.projectId }}
      className="group flex items-center gap-4 px-5 py-4 transition-colors hover:bg-muted/40"
    >
      <AppLogo icon={app.icon} name={app.name} className="size-10 shrink-0" />
      <div className="min-w-0 flex-1">
        <p className="truncate font-medium text-foreground text-sm">{app.name}</p>
        <p className="mt-0.5 flex items-center gap-2 text-muted-foreground text-xs">
          <span className={cn('size-1.5 rounded-full', statusTone[app.status] ?? 'bg-zinc-500')} />
          <span className="capitalize">{app.status}</span>
          {meta && (
            <>
              <span aria-hidden="true">·</span>
              <span className="truncate">{meta}</span>
            </>
          )}
        </p>
      </div>
      <ChevronRight className="size-4 shrink-0 text-muted-foreground/40 transition-colors group-hover:text-muted-foreground" />
    </Link>
  );
}
