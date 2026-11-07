import { ChevronDown } from 'lucide-react';
import { type ReactNode, useState } from 'react';

export function SettingsSection({
  icon,
  title,
  description,
  children,
  action,
  iconBg = 'bg-primary/10',
  iconColor = 'text-primary',
  collapsible = false,
  defaultOpen = false,
}: {
  icon: ReactNode;
  title: string;
  description?: string;
  children: ReactNode;
  action?: ReactNode;
  iconBg?: string;
  iconColor?: string;
  collapsible?: boolean;
  defaultOpen?: boolean;
}) {
  const [open, setOpen] = useState(defaultOpen);
  const expanded = collapsible ? open : true;

  const header = (
    <>
      <div
        className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-xl ${iconBg} ${iconColor}`}
      >
        {icon}
      </div>
      <div className="min-w-0 flex-1">
        <h2 className="font-semibold text-[15px] text-foreground">{title}</h2>
        {description && <p className="text-muted-foreground text-xs">{description}</p>}
      </div>
      {collapsible && (
        <ChevronDown
          className={`size-4 shrink-0 text-muted-foreground transition-transform ${
            open ? 'rotate-180' : ''
          }`}
        />
      )}
    </>
  );

  return (
    <div className="rounded-2xl border border-border/50 bg-card">
      <div className={`flex items-center ${expanded ? 'border-border/50 border-b' : ''}`}>
        {collapsible ? (
          <button
            type="button"
            onClick={() => setOpen((v) => !v)}
            aria-expanded={open}
            className="flex min-w-0 flex-1 items-center gap-3 px-5 py-4 text-start"
          >
            {header}
          </button>
        ) : (
          <div className="flex min-w-0 flex-1 items-center gap-3 px-5 py-4">{header}</div>
        )}
        {action && <div className="shrink-0 ps-3 pe-5">{action}</div>}
      </div>
      {expanded && <div className="p-5">{children}</div>}
    </div>
  );
}
