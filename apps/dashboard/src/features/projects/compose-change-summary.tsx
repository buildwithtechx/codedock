import { Badge } from '@codedock/ui/badge';
import { AlertTriangle } from 'lucide-react';
import type { ComposeSourceChanges } from './compose-stack-types';

interface ComposeChangeSummaryProps {
  changes?: ComposeSourceChanges;
}

export function ComposeChangeSummary({ changes }: ComposeChangeSummaryProps) {
  if (!changes) return null;

  const hasAdded = changes.added && changes.added.length > 0;
  const hasChanged = changes.changed && changes.changed.length > 0;
  const hasRemoved = changes.removed && changes.removed.length > 0;

  if (!hasAdded && !hasChanged && !hasRemoved) {
    return null;
  }

  return (
    <div className="space-y-3 rounded-lg border border-border/60 bg-muted/20 p-3.5 text-sm">
      <div className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
        Detected Changes
      </div>
      <div className="flex flex-wrap gap-2">
        {hasAdded && (
          <div className="flex items-center gap-1.5">
            <span className="text-muted-foreground text-xs">Added:</span>
            {changes.added.map((name) => (
              <Badge
                key={name}
                variant="outline"
                className="border-emerald-500/40 bg-emerald-500/10 font-mono text-emerald-600 text-xs dark:text-emerald-400"
              >
                +{name}
              </Badge>
            ))}
          </div>
        )}
        {hasChanged && (
          <div className="flex items-center gap-1.5">
            <span className="text-muted-foreground text-xs">Changed:</span>
            {changes.changed.map((name) => (
              <Badge
                key={name}
                variant="outline"
                className="border-amber-500/40 bg-amber-500/10 font-mono text-amber-600 text-xs dark:text-amber-400"
              >
                ~{name}
              </Badge>
            ))}
          </div>
        )}
        {hasRemoved && (
          <div className="flex items-center gap-1.5">
            <span className="text-muted-foreground text-xs">Removed:</span>
            {changes.removed.map((name) => (
              <Badge
                key={name}
                variant="destructive"
                className="border-destructive/30 bg-destructive/10 font-mono text-destructive text-xs"
              >
                -{name}
              </Badge>
            ))}
          </div>
        )}
      </div>

      {hasRemoved && (
        <div className="flex items-start gap-2 rounded-md border border-destructive/20 bg-destructive/5 p-2 text-destructive text-xs">
          <AlertTriangle className="mt-0.5 size-4 shrink-0" />
          <span>
            Removed services (<strong>{changes.removed.join(', ')}</strong>) will be stopped,
            unrouted from edge proxy, and orphaned containers safely pruned on deployment.
          </span>
        </div>
      )}
    </div>
  );
}
