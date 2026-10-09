import { BookOpen, ExternalLink, Loader2, RefreshCw } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { cn } from '#/lib/utils';
import { useEvaluateIssues } from './hooks';

const DOCS_URL = 'https://docs.codedock.run';

type EmptyVariant = 'open' | 'resolved' | 'filtered';

const COPY: Record<EmptyVariant, { title: string; body: string }> = {
  open: {
    title: 'All clear',
    body: 'No open issues. Deployments, services, backups, migrations, and capacity all report healthy.',
  },
  resolved: {
    title: 'Nothing resolved yet',
    body: 'Issues you resolve, acknowledge-clear, or that disappear on rescan will be listed here.',
  },
  filtered: {
    title: 'No matching issues',
    body: 'Nothing matches this search or severity filter. Clear the filters to see the full feed.',
  },
};

export function EmptyIssues({
  filtered,
  resolved,
  onClearFilters,
}: {
  filtered: boolean;
  resolved: boolean;
  onClearFilters?: () => void;
}) {
  const variant: EmptyVariant = filtered ? 'filtered' : resolved ? 'resolved' : 'open';
  const copy = COPY[variant];
  const narrow = variant === 'filtered';
  const evaluate = useEvaluateIssues();
  const showClear = variant === 'filtered' && onClearFilters !== undefined;

  return (
    <div
      className={cn(
        'rounded-2xl border border-border/50 bg-card px-6 text-center',
        narrow ? 'py-10' : 'py-14'
      )}
    >
      <IssuesIllustration
        variant={variant}
        className={cn(
          'relative mx-auto max-w-full text-muted-foreground/50',
          narrow ? 'mb-4 h-32 w-56' : 'mb-6 h-36 w-64'
        )}
      />
      <h3 className="font-medium text-[15px] tracking-[-0.2px]">{copy.title}</h3>
      <p className="mx-auto mt-1.5 max-w-md text-[13px] text-muted-foreground/80 leading-relaxed">
        {copy.body}
      </p>
      <div className="mt-5 flex flex-col items-center justify-center gap-3 sm:flex-row">
        {showClear ? (
          <Button onClick={onClearFilters}>Clear filters</Button>
        ) : (
          <Button onClick={() => evaluate.mutate()} disabled={evaluate.isPending}>
            {evaluate.isPending ? (
              <Loader2 className="size-4 animate-spin" />
            ) : (
              <RefreshCw className="size-4" />
            )}
            {evaluate.isPending ? 'Rescanning' : 'Run rescan'}
          </Button>
        )}
        <Button asChild variant="secondary">
          <a href={DOCS_URL} target="_blank" rel="noopener noreferrer">
            <BookOpen className="size-4" />
            Docs
            <ExternalLink className="size-3.5 opacity-60" />
          </a>
        </Button>
      </div>
    </div>
  );
}

export function IssuesIllustration({
  variant,
  className,
}: {
  variant: EmptyVariant;
  className?: string;
}) {
  return (
    <svg
      viewBox="0 0 256 144"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      role="img"
      aria-label={COPY[variant].title}
    >
      <rect x="48" y="18" width="160" height="108" rx="14" strokeOpacity="0.35" />
      <rect x="48" y="18" width="160" height="30" rx="14" strokeOpacity="0.25" />
      <line x1="48" y1="48" x2="208" y2="48" strokeOpacity="0.25" />
      <circle cx="62" cy="33" r="3" fill="currentColor" stroke="none" opacity="0.5" />
      <circle cx="74" cy="33" r="3" fill="currentColor" stroke="none" opacity="0.3" />
      <rect
        x="88"
        y="29"
        width="52"
        height="8"
        rx="4"
        fill="currentColor"
        stroke="none"
        opacity="0.2"
      />
      {variant === 'open' && (
        <g>
          <circle cx="128" cy="88" r="20" strokeOpacity="0.6" />
          <path d="M119 88.5l6.5 6.5L138 82" strokeOpacity="0.9" />
        </g>
      )}
      {variant === 'resolved' && (
        <g strokeOpacity="0.6">
          <rect x="76" y="66" width="104" height="12" rx="6" />
          <rect x="76" y="86" width="72" height="12" rx="6" strokeOpacity="0.35" />
          <circle cx="172" cy="92" r="14" />
          <path d="M172 85.5v6.5l4.5 3" />
        </g>
      )}
      {variant === 'filtered' && (
        <g>
          <line x1="76" y1="72" x2="150" y2="72" strokeOpacity="0.35" />
          <line x1="76" y1="86" x2="132" y2="86" strokeOpacity="0.25" />
          <line x1="76" y1="100" x2="142" y2="100" strokeOpacity="0.25" />
          <circle cx="164" cy="86" r="16" strokeOpacity="0.7" />
          <line x1="176" y1="98" x2="188" y2="110" strokeOpacity="0.7" />
        </g>
      )}
    </svg>
  );
}
