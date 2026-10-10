import type { LucideIcon } from 'lucide-react';
import { RefreshCw } from 'lucide-react';
import type { ReactNode } from 'react';
import { Button } from '#/components/ui/button';

export function ManagedControlPanel({
  title,
  description,
  icon: IconComponent,
  busy,
  error,
  refresh,
  children,
}: {
  title: string;
  description: string;
  icon: LucideIcon;
  busy?: boolean;
  error?: string | null;
  refresh?: () => void;
  children: ReactNode;
}) {
  return (
    <section className="min-w-0 space-y-5 rounded-2xl bg-card p-5" aria-busy={busy}>
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <h2 className="flex items-center gap-2 font-medium text-base">
            <IconComponent className="size-4 shrink-0 text-muted-foreground" />
            {title}
          </h2>
          <p className="mt-1 text-muted-foreground text-sm">{description}</p>
        </div>
        {refresh && (
          <Button
            variant="ghost"
            size="icon"
            aria-label="Refresh"
            disabled={busy}
            onClick={refresh}
          >
            <RefreshCw className={`size-4 ${busy ? 'animate-spin' : ''}`} />
          </Button>
        )}
      </div>
      {error && (
        <div
          role="alert"
          className="space-y-2 rounded-xl bg-destructive/10 p-4 text-destructive text-sm"
        >
          <p className="break-words">{error}</p>
          {refresh && (
            <Button variant="secondary" size="sm" disabled={busy} onClick={refresh}>
              Retry
            </Button>
          )}
        </div>
      )}
      {children}
    </section>
  );
}
