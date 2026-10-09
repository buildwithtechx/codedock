import { Link } from '@tanstack/react-router';
import { ArrowLeft, BadgeCheck } from 'lucide-react';
import { useState } from 'react';
import { Badge } from '#/components/ui/badge';
import type { OneClickAppDetails } from '#/interfaces/templates';
import { AppLogo } from './app-logo';

export function WizardHeader({ app }: { app: OneClickAppDetails }) {
  const [expanded, setExpanded] = useState(false);
  const long = app.description.length > 160;
  return (
    <div>
      <Link
        to="/apps/new"
        className="mb-4 inline-flex items-center gap-1.5 text-muted-foreground text-sm transition-colors hover:text-foreground"
      >
        <ArrowLeft className="size-4" />
        Back to catalog
      </Link>
      <div className="flex items-center gap-4">
        <AppLogo icon={app.icon} name={app.name} className="size-12 shrink-0" />
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-2">
            <h1 className="font-semibold text-foreground text-xl">{app.name}</h1>
            {app.verified && (
              <Badge variant="secondary" className="gap-1">
                <BadgeCheck className="size-3" />
                Verified
              </Badge>
            )}
            {app.category && <Badge variant="outline">{app.category}</Badge>}
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
              {expanded ? 'Show less' : 'Show more'}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
