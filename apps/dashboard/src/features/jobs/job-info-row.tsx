import type { ReactNode } from 'react';

interface JobInfoRowProps {
  icon: ReactNode;
  label: string;
  value: ReactNode;
}

export function JobInfoRow({ icon, label, value }: JobInfoRowProps) {
  return (
    <div className="flex items-start gap-3 rounded-xl border border-border/50 bg-card px-4 py-3">
      <span className="mt-0.5 text-muted-foreground/60">{icon}</span>
      <span className="w-28 shrink-0 text-[12px] text-muted-foreground/70">{label}</span>
      <span className="min-w-0 flex-1 break-words text-[13px] text-foreground">{value}</span>
    </div>
  );
}
