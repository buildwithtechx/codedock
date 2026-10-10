import type { ReactNode } from 'react';
import { Label } from '#/components/ui/label';
import { cn } from '#/lib/utils';

export function SettingsRow({
  label,
  description,
  children,
  bordered = true,
}: {
  label: string;
  description?: string;
  children: ReactNode;
  bordered?: boolean;
}) {
  return (
    <div
      className={cn(
        'flex items-start justify-between gap-8 py-4 last:border-0',
        bordered && 'border-border/70 border-b'
      )}
    >
      <div className="min-w-0 flex-1">
        <Label className="font-medium text-sm">{label}</Label>
        {description && <p className="mt-0.5 text-muted-foreground text-xs">{description}</p>}
      </div>
      <div className="flex w-80 shrink-0 items-center justify-end">{children}</div>
    </div>
  );
}
