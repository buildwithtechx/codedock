import { Link } from '@tanstack/react-router';
import { ArrowLeft, BadgeCheck } from 'lucide-react';
import { useState } from 'react';
import type { OneClickAppDetails } from '#/interfaces/templates';
import { AppLogo } from './app-logo';

export function WizardHeader({ app, projectId }: { app: OneClickAppDetails; projectId?: string }) {
  const [expanded, setExpanded] = useState(false);
  const long = app.description.length > 160;
  return (
    <div>
      {projectId ? (
        <Link
          to="/projects/$projectId/new"
          params={{ projectId }}
          search={{ tab: 'one-click' }}
          className="mb-4 inline-flex items-center gap-1.5 text-muted-foreground text-sm transition-colors hover:text-foreground"
        >
          <ArrowLeft className="size-4" />
          Back
        </Link>
      ) : (
        <Link
          to="/apps/new"
          className="mb-4 inline-flex items-center gap-1.5 text-muted-foreground text-sm transition-colors hover:text-foreground"
        >
          <ArrowLeft className="size-4" />
          Back
        </Link>
      )}
      <div className="flex items-center gap-4">
        <div className="flex size-12 shrink-0 items-center justify-center rounded-2xl border border-border/50 bg-card shadow-xs">
          <AppLogo icon={app.icon} name={app.name} className="size-7" />
        </div>
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-2">
            <h1 className="font-semibold text-foreground text-xl tracking-tight">{app.name}</h1>
            {app.verified && (
              <BadgeCheck className="size-4 shrink-0 text-primary" aria-label="Verified app" />
            )}
          </div>
          <p className={`mt-1 text-muted-foreground text-sm ${expanded ? '' : 'line-clamp-2'}`}>
            {app.description}
          </p>
          {long && (
            <button
              type="button"
              onClick={() => setExpanded((value) => !value)}
              className="mt-0.5 font-medium text-muted-foreground/80 text-xs transition-colors hover:text-foreground"
            >
              {expanded ? 'Less' : 'More'}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
