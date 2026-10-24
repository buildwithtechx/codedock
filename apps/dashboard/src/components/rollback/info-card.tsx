import type { LucideIcon } from 'lucide-react';
import type React from 'react';
import { cn } from '#/lib/utils';

export function InfoCard({
  icon: Icon,
  title,
  value,
  description,
  action,
  footer,
  tone = 'neutral',
}: {
  icon: LucideIcon;
  title: string;
  value: string;
  description: string;
  action?: React.ReactNode;
  footer?: React.ReactNode;
  tone?: 'neutral' | 'success';
}) {
  return (
    <div className="rounded-xl border border-border/50 bg-muted/20 px-4 py-3">
      <div className="flex items-start justify-between gap-3">
        <div className="flex min-w-0 items-start gap-3">
          <div
            className={cn(
              'flex h-8 w-8 shrink-0 items-center justify-center rounded-lg',
              tone === 'success' ? 'bg-success/15 text-success' : 'bg-primary/10 text-primary'
            )}
          >
            <Icon className="size-4" />
          </div>
          <div className="min-w-0">
            <p className="font-medium text-[13px] text-foreground">{title}</p>
            <p className="mt-1 font-semibold text-[13px] text-foreground">{value}</p>
            <p className="mt-1 text-[12px] text-muted-foreground">{description}</p>
          </div>
        </div>
        {action ? <div className="shrink-0">{action}</div> : null}
      </div>
      {footer ? <div className="mt-3 border-border/40 border-t pt-3">{footer}</div> : null}
    </div>
  );
}

export default InfoCard;
