import { CalendarClock, Plus } from 'lucide-react';
import { Button } from '#/components/ui/button';

export function JobsEmptyState({ onCreate }: { onCreate: () => void }) {
  return (
    <div className="py-8 text-center">
      <div className="relative mx-auto mb-6 h-44 w-64 max-w-full">
        <svg className="absolute inset-0 h-full w-full" viewBox="0 0 260 180" fill="none">
          <rect x="60" y="46" width="140" height="98" rx="14" className="fill-muted" />
          <rect
            x="48"
            y="34"
            width="140"
            height="98"
            rx="14"
            className="fill-card stroke-border"
            strokeWidth="1"
          />
          <rect x="48" y="34" width="140" height="22" rx="11" className="fill-muted/60" />
          <rect x="48" y="45" width="140" height="11" className="fill-muted/60" />
          <circle cx="62" cy="45" r="3.5" className="fill-muted-foreground/30" />
          <circle cx="73" cy="45" r="3.5" className="fill-muted-foreground/20" />
          <circle cx="84" cy="45" r="3.5" className="fill-muted-foreground/10" />
          <rect x="150" y="42" width="26" height="6" rx="3" className="fill-muted-foreground/15" />
          <rect x="66" y="68" width="64" height="7" rx="3.5" className="fill-foreground/15" />
          <rect x="66" y="80" width="40" height="5" rx="2.5" className="fill-foreground/10" />
          <rect x="66" y="93" width="56" height="7" rx="3.5" className="fill-foreground/10" />
          <rect x="66" y="105" width="46" height="5" rx="2.5" className="fill-foreground/10" />
          <rect x="66" y="117" width="60" height="7" rx="3.5" className="fill-foreground/10" />
          <line
            x1="150"
            y1="68"
            x2="150"
            y2="120"
            className="stroke-border"
            strokeWidth="1.5"
            strokeLinecap="round"
          />
          <circle cx="150" cy="118" r="6" className="fill-card stroke-border" strokeWidth="1.5" />
          <path
            d="M147.5 118l1.8 1.8 3.2-3.6"
            className="stroke-muted-foreground"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
          <circle cx="150" cy="94" r="6" className="fill-card stroke-border" strokeWidth="1.5" />
          <path
            d="M147.5 94l1.8 1.8 3.2-3.6"
            className="stroke-muted-foreground"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
          <circle cx="150" cy="70" r="9" className="fill-success/15" />
          <circle cx="150" cy="70" r="6" className="fill-card stroke-success" strokeWidth="1.5" />
          <path
            d="M147.5 70l1.8 1.8 3.2-3.6"
            className="stroke-success"
            strokeWidth="1.6"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
          <circle cx="212" cy="100" r="24" className="fill-primary/10" />
          <circle
            cx="212"
            cy="100"
            r="17"
            className="fill-card stroke-primary"
            strokeWidth="2"
            strokeDasharray="5 4"
          />
          <path d="M208.5 92.5v15l11-7.5z" className="fill-primary" />
          <path
            d="M188 106 Q196 104 198 102"
            className="stroke-border"
            strokeWidth="1.5"
            strokeDasharray="3 3"
          />
          <circle cx="212" cy="100" r="1.6" className="fill-primary" />
          <circle cx="26" cy="70" r="4.5" className="fill-muted-foreground/15" />
          <circle cx="40" cy="152" r="5" className="fill-muted-foreground/10" />
          <circle cx="238" cy="52" r="3.5" className="fill-muted-foreground/15" />
          <circle cx="230" cy="150" r="4" className="fill-muted-foreground/10" />
        </svg>
      </div>

      <h3 className="mb-2 font-medium text-foreground/80 text-xl">No jobs yet</h3>
      <p className="mx-auto mb-6 max-w-md text-muted-foreground text-sm leading-relaxed">
        Schedule recurring commands on your services — cleanups, reports, syncs — and track every
        run from one place.
      </p>

      <div className="flex items-center justify-center">
        <Button onClick={onCreate} size="lg" className="rounded-xl px-6">
          <Plus className="size-4" />
          Create your first job
        </Button>
      </div>

      <p className="mt-5 flex items-center justify-center gap-1.5 text-muted-foreground/70 text-xs">
        <CalendarClock className="size-3.5" />
        Cron schedules, manual runs, and per-job output included
      </p>
    </div>
  );
}
