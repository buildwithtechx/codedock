import { Link } from '@tanstack/react-router';
import { GitBranch, Plus } from 'lucide-react';
import { Button } from '#/components/ui/button';

export function HomeFirstProject({ onCreateProject }: { onCreateProject: () => void }) {
  return (
    <section className="px-6 pt-5 pb-7 text-center sm:pt-7 sm:pb-8">
      <div className="relative mx-auto mb-3 h-36 w-72 max-w-full">
        <svg
          className="absolute inset-0 h-full w-full"
          viewBox="0 0 288 132"
          fill="none"
          aria-hidden="true"
        >
          <g
            stroke="currentColor"
            strokeWidth="1.5"
            strokeDasharray="4 4"
            className="text-border"
            fill="none"
          >
            <path d="M71 32 Q 100 42 123 58" />
            <path d="M71 96 Q 100 90 123 74" />
            <path d="M165 58 Q 188 42 217 32" />
            <path d="M165 74 Q 188 90 217 96" />
          </g>

          <circle cx="98.5" cy="43.5" r="2.4" className="fill-muted-foreground/40" />
          <circle cx="98.5" cy="87.5" r="2" className="fill-muted-foreground/30" />
          <circle cx="189.5" cy="43.5" r="2.4" className="fill-muted-foreground/40" />
          <circle cx="189.5" cy="87.5" r="2" className="fill-muted-foreground/30" />

          <g transform="translate(52 32)">
            <rect
              x="-19"
              y="-17"
              width="38"
              height="34"
              rx="10"
              className="fill-card stroke-border"
              strokeWidth="1.5"
            />
            <g
              fill="none"
              className="stroke-muted-foreground"
              strokeWidth="1.6"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <circle cx="-6" cy="5" r="3.2" />
              <circle cx="6" cy="-6" r="3.2" />
              <path d="M-6 1.8v-3.8a4 4 0 0 1 4-4h4.8" />
            </g>
          </g>
          <text
            x="52"
            y="60"
            textAnchor="middle"
            className="fill-muted-foreground font-medium text-[7px]"
          >
            Repo
          </text>

          <g transform="translate(52 96)">
            <rect
              x="-19"
              y="-17"
              width="38"
              height="34"
              rx="10"
              className="fill-card stroke-border"
              strokeWidth="1.5"
            />
            <g
              fill="none"
              className="stroke-muted-foreground"
              strokeWidth="1.6"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M0 7V-6" />
              <path d="M-5 -1l5 -5.5 5 5.5" />
              <path d="M-7 9.5h14" />
            </g>
          </g>
          <text
            x="52"
            y="124"
            textAnchor="middle"
            className="fill-muted-foreground font-medium text-[7px]"
          >
            Deploy
          </text>

          <g transform="translate(236 32)">
            <rect
              x="-19"
              y="-17"
              width="38"
              height="34"
              rx="10"
              className="fill-card stroke-border"
              strokeWidth="1.5"
            />
            <g
              fill="none"
              className="stroke-muted-foreground"
              strokeWidth="1.6"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <circle cx="0" cy="0" r="8" />
              <path d="M-8 0h16" />
              <path d="M0 -8c4 3.6 4 12.4 0 16" />
              <path d="M0 -8c-4 3.6-4 12.4 0 16" />
            </g>
          </g>
          <text
            x="236"
            y="60"
            textAnchor="middle"
            className="fill-muted-foreground font-medium text-[7px]"
          >
            Domain
          </text>

          <g transform="translate(236 96)">
            <rect
              x="-19"
              y="-17"
              width="38"
              height="34"
              rx="10"
              className="fill-card stroke-border"
              strokeWidth="1.5"
            />
            <g
              fill="none"
              className="stroke-muted-foreground"
              strokeWidth="1.6"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <ellipse cx="0" cy="-5" rx="8" ry="3" />
              <path d="M-8 -5v9.5c0 1.7 3.6 3 8 3s8-1.3 8-3V-5" />
              <path d="M-8 -0.5c0 1.7 3.6 3 8 3s8-1.3 8-3" />
            </g>
          </g>
          <text
            x="236"
            y="124"
            textAnchor="middle"
            className="fill-muted-foreground font-medium text-[7px]"
          >
            Data
          </text>

          <circle cx="144" cy="66" r="27" className="fill-primary/5" />
          <circle
            cx="144"
            cy="66"
            r="18"
            fill="none"
            className="stroke-foreground/80"
            strokeWidth="3.5"
          />
          <circle cx="158" cy="52" r="4.2" className="fill-card" />
          <circle cx="158" cy="52" r="2.8" className="fill-emerald-500" />
        </svg>
      </div>

      <h2
        className="mt-4 font-medium text-foreground/85 text-xl"
        style={{ letterSpacing: '-0.2px' }}
      >
        Ship your first project
      </h2>
      <p className="mx-auto mt-1.5 mb-5 max-w-sm text-muted-foreground/80 text-sm leading-relaxed">
        Connect a repository and Codedock builds, deploys, and secures it. No pipelines to wire by
        hand.
      </p>
      <div className="flex flex-col items-center justify-center gap-3 sm:flex-row">
        <Button className="gap-2 px-6" onClick={onCreateProject}>
          <Plus className="h-4 w-4" />
          Create project
        </Button>
        <Button asChild variant="secondary" className="gap-2 px-6">
          <Link to="/library">
            <GitBranch className="h-4 w-4" />
            Import from Git
          </Link>
        </Button>
      </div>
      <p className="mt-5 text-muted-foreground/60 text-xs">
        Tip: press{' '}
        <kbd className="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px]">&#8984; K</kbd> to
        jump anywhere
      </p>
    </section>
  );
}
