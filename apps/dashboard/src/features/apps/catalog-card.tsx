import { Link } from '@tanstack/react-router';
import { ArrowRight, BadgeCheck } from 'lucide-react';
import { Badge } from '#/components/ui/badge';
import type { OneClickAppDetails } from '#/interfaces/templates';
import { AppLogo } from './app-logo';

export function CatalogCard({ app }: { app: OneClickAppDetails }) {
  return (
    <Link
      to="/apps/new/$appId"
      params={{ appId: app.id }}
      className="group flex w-full items-start gap-3 rounded-2xl border border-border/50 bg-card p-5 text-left transition-all hover:border-primary/40 hover:shadow-md"
    >
      <AppLogo icon={app.icon} name={app.name} className="size-10 shrink-0" />
      <div className="min-w-0 flex-1">
        <div className="flex items-center justify-between gap-2">
          <span className="inline-flex min-w-0 items-center gap-1.5 font-medium text-foreground">
            <span className="truncate">{app.name}</span>
            {app.verified && (
              <Badge variant="secondary" className="shrink-0 gap-1">
                <BadgeCheck className="size-3" />
                Verified
              </Badge>
            )}
          </span>
          <ArrowRight className="size-4 shrink-0 text-muted-foreground/40 transition-colors group-hover:text-foreground" />
        </div>
        <p className="mt-1 line-clamp-2 text-muted-foreground text-sm">{app.description}</p>
        {app.category && (
          <p className="mt-2 font-mono text-[10px] text-muted-foreground uppercase tracking-widest">
            {app.category}
          </p>
        )}
      </div>
    </Link>
  );
}

export function CatalogShortcut({ app }: { app: OneClickAppDetails }) {
  return (
    <Link
      to="/apps/new/$appId"
      params={{ appId: app.id }}
      className="group flex min-w-0 items-center gap-3 rounded-xl bg-card p-4 text-start transition-colors hover:bg-muted/50"
    >
      <AppLogo icon={app.icon} name={app.name} className="size-10 shrink-0" />
      <div className="min-w-0 flex-1">
        <p className="truncate font-medium text-foreground text-sm">{app.name}</p>
        <p className="mt-0.5 line-clamp-2 text-muted-foreground text-xs leading-relaxed">
          {app.description}
        </p>
      </div>
      <ArrowRight className="size-4 shrink-0 text-muted-foreground transition-colors group-hover:text-foreground" />
    </Link>
  );
}
