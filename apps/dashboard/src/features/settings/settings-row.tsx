import type { ReactNode } from 'react';
import { Label } from '#/components/ui/label';

export function SettingsRow({
  label,
  description,
  children,
}: {
  label: string;
  description?: string;
  children: ReactNode;
}) {
  return (
    <div className="flex items-start justify-between gap-8 border-border/70 border-b py-4 last:border-0">
      <div className="min-w-0 flex-1">
        <Label className="font-medium text-sm">{label}</Label>
        {description && <p className="mt-0.5 text-muted-foreground text-xs">{description}</p>}
      </div>
      <div className="w-80 shrink-0">{children}</div>
    </div>
  );
}
