import type { ReactNode } from 'react';

export function Section({
  title,
  icon,
  children,
}: {
  title: string;
  icon?: ReactNode;
  children: ReactNode;
}) {
  return (
    <section className="rounded-2xl border border-border/50 bg-card">
      <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
        {icon && (
          <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
            {icon}
          </div>
        )}
        <h2 className="font-semibold text-[15px] text-foreground">{title}</h2>
      </div>
      <div className="space-y-4 p-5">{children}</div>
    </section>
  );
}

export function SummaryRow({
  label,
  value,
  mono,
}: {
  label: string;
  value: string;
  mono?: boolean;
}) {
  return (
    <div className="flex items-center justify-between gap-3">
      <span className="shrink-0 text-muted-foreground text-sm">{label}</span>
      <span
        className={`min-w-0 truncate text-right text-foreground text-sm ${mono ? 'font-mono' : ''}`}
      >
        {value}
      </span>
    </div>
  );
}
