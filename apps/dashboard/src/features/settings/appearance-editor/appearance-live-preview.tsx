import { Button } from '@codedock/ui/button';
import { LayoutGrid } from 'lucide-react';

interface AppearanceLivePreviewProps {
  tokens: Record<string, string>;
  activeThemeName?: string;
}

export function AppearanceLivePreview({
  tokens,
  activeThemeName = 'codedock',
}: AppearanceLivePreviewProps) {
  const cardBg = tokens['--th-bg-card'] || 'var(--card)';
  const accentColor = tokens['--th-accent'] || 'var(--primary)';
  const borderCol = tokens['--th-bd-default'] || 'var(--border)';

  return (
    <div className="sticky top-6 space-y-2">
      <h3 className="font-semibold text-foreground text-sm">Live preview</h3>

      <div
        className="rounded-2xl border p-5 shadow-lg transition-colors"
        style={{
          backgroundColor: cardBg,
          borderColor: borderCol,
        }}
      >
        {/* Workspace Title */}
        <div className="flex items-center gap-2">
          <LayoutGrid className="size-4 text-foreground/80" />
          <span className="font-semibold text-foreground text-sm">Your workspace</span>
        </div>

        {/* Workspace status */}
        <div className="mt-4 rounded-xl border border-border/50 bg-background/40 p-3.5">
          <div className="flex items-center gap-2">
            <span className="size-2 rounded-full bg-emerald-500 shadow-emerald-500/50 shadow-xs" />
            <span className="font-medium text-foreground text-xs">{activeThemeName}</span>
          </div>
          <p className="mt-1 text-[11px] text-muted-foreground">
            Changes apply to this account only.
          </p>

          {/* Badges Row */}
          <div className="mt-3 flex items-center gap-2">
            <span className="rounded-full bg-emerald-500/15 px-2.5 py-0.5 font-medium text-[10px] text-emerald-500">
              Active
            </span>
            <span className="rounded-full bg-amber-500/15 px-2.5 py-0.5 font-medium text-[10px] text-amber-500">
              Queued
            </span>
            <span className="rounded-full bg-rose-500/15 px-2.5 py-0.5 font-medium text-[10px] text-rose-500">
              Failed
            </span>
          </div>
        </div>

        {/* Domain and CTA */}
        <div className="mt-4 space-y-2.5">
          <p className="font-mono text-muted-foreground text-xs">app.example.com</p>

          <Button
            type="button"
            className="w-full justify-center font-medium text-xs shadow-md"
            style={{
              backgroundColor: accentColor,
            }}
          >
            Create
          </Button>
        </div>

        {/* Deploy snippet */}
        <div className="mt-4 rounded-xl border border-border/40 bg-background/60 p-3 font-mono text-[11px] text-muted-foreground leading-relaxed">
          <p className="flex items-center gap-1.5 text-emerald-500">
            <span>✓</span>
            <span>build completed</span>
          </p>
          <p className="flex items-center gap-1.5 text-foreground/75">
            <span>→</span>
            <span>ready to deploy</span>
          </p>
        </div>

        {/* 7 Vertical Metric Bars Chart */}
        <div className="mt-5 flex h-18 items-end justify-between gap-2 px-1">
          <div className="h-6 w-full rounded-sm bg-sky-500/80 transition-all" />
          <div className="h-10 w-full rounded-sm bg-sky-500/80 transition-all" />
          <div className="h-8 w-full rounded-sm bg-sky-500/80 transition-all" />
          <div className="h-13 w-full rounded-sm bg-sky-500/80 transition-all" />
          <div className="h-9 w-full rounded-sm bg-sky-500/80 transition-all" />
          <div className="h-16 w-full rounded-sm bg-sky-500/80 transition-all" />
          <div className="h-11 w-full rounded-sm bg-sky-500/80 transition-all" />
        </div>
      </div>

      <p className="pt-1 text-[11px] text-muted-foreground/70">
        Changes apply to this account only.
      </p>
    </div>
  );
}
