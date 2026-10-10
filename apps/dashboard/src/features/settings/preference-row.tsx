import type { ReactNode } from 'react';
import { Switch } from '#/components/ui/switch';

export function PreferenceRow({
  title,
  hint,
  checked,
  onChange,
}: {
  title: string;
  hint: ReactNode;
  checked: boolean;
  onChange: (v: boolean) => void;
}) {
  return (
    <div className="flex items-center justify-between gap-4 rounded-xl border border-border/50 bg-muted/10 p-4">
      <div className="min-w-0">
        <p className="font-medium text-foreground text-sm">{title}</p>
        <p className="text-muted-foreground text-xs">{hint}</p>
      </div>
      <Switch checked={checked} onCheckedChange={onChange} aria-label={title} />
    </div>
  );
}
