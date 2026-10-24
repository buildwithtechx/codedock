import { Button } from '@codedock/ui/button';
import { Link } from '@tanstack/react-router';
import { ArrowRight } from 'lucide-react';
import { useId } from 'react';
import { useGetPublicSettings } from '#/hooks/use-settings';
import { useAuthStore } from '#/stores/auth-store';

function CloudPlanIllustration() {
  const washId = useId();
  return (
    <svg
      viewBox="0 0 224 112"
      fill="none"
      aria-hidden="true"
      focusable="false"
      className="pointer-events-none mx-auto mt-2 mb-3 block h-auto w-48 max-w-full select-none"
    >
      <defs>
        <linearGradient id={washId} x1="30" y1="40" x2="196" y2="88" gradientUnits="userSpaceOnUse">
          <stop stopColor="var(--primary)" stopOpacity="0.08" />
          <stop offset="1" stopColor="var(--primary)" stopOpacity="0.14" />
        </linearGradient>
      </defs>
      <ellipse cx="113" cy="67" rx="88" ry="35" fill={`url(#${washId})`} />
      <path
        d="M148 51h30c8 0 14-6 14-13 0-7-5-12-12-13-2-9-10-15-19-13-9 1-15 8-16 16-6 0-11 5-11 11 0 7 6 12 14 12Z"
        fill="var(--card)"
        stroke="var(--border)"
        strokeWidth="1.5"
        strokeLinejoin="round"
      />
      <path
        d="M144 76h18a10 10 0 0 0 10-10v-9"
        stroke="var(--muted-foreground)"
        strokeOpacity="0.4"
        strokeWidth="1.25"
        strokeDasharray="3 4"
        strokeLinecap="round"
      />
      <rect
        x="52"
        y="32"
        width="98"
        height="58"
        rx="11"
        fill="var(--muted)"
        fillOpacity="0.35"
        stroke="var(--border)"
      />
      <rect
        x="43"
        y="40"
        width="98"
        height="58"
        rx="11"
        fill="var(--card)"
        stroke="var(--border)"
      />
      <path d="M44 54h96" stroke="var(--border)" />
      <g fill="var(--muted-foreground)" fillOpacity="0.45">
        <circle cx="54" cy="47" r="1.5" />
        <circle cx="60" cy="47" r="1.5" />
        <circle cx="66" cy="47" r="1.5" />
      </g>
      <circle cx="65" cy="76" r="11" fill="var(--muted)" fillOpacity="0.4" />
      <circle
        cx="65"
        cy="76"
        r="7"
        stroke="var(--muted-foreground)"
        strokeOpacity="0.6"
        strokeWidth="1.6"
      />
      <rect
        x="86"
        y="68"
        width="35"
        height="3"
        rx="1.5"
        fill="var(--muted-foreground)"
        fillOpacity="0.3"
      />
      <rect x="86" y="76" width="26" height="3" rx="1.5" fill="var(--border)" />
      <rect x="86" y="84" width="19" height="3" rx="1.5" fill="var(--border)" />
      <circle cx="28" cy="48" r="2.5" fill="var(--border)" />
      <path
        d="M198 65v6m-3-3h6"
        stroke="var(--muted-foreground)"
        strokeOpacity="0.35"
        strokeWidth="1.25"
        strokeLinecap="round"
      />
      <circle cx="161" cy="76" r="2.5" fill="var(--primary)" fillOpacity="0.6" />
    </svg>
  );
}

/** First-subscription offer in Cloud mode; Self-Hosted and subscribed workspaces return null. */
export function HomeModeCard() {
  const { data: publicSettings } = useGetPublicSettings();
  const user = useAuthStore((s) => s.user);

  const isCloud = Boolean(publicSettings?.data?.cloudMode);
  const planType = user?.planType || 'free';

  if (!isCloud || planType !== 'free') {
    return null;
  }

  return (
    <section
      aria-label="Codedock Cloud"
      className="overflow-hidden rounded-2xl border border-border/50 bg-card p-5 shadow-xs"
      style={{
        backgroundImage:
          'radial-gradient(ellipse at 100% 0%, color-mix(in oklab, var(--primary) 8%, transparent), transparent 65%)',
      }}
    >
      <p className="font-medium text-muted-foreground text-sm">Codedock Cloud</p>
      <CloudPlanIllustration />
      <h2 className="font-medium text-foreground/85 text-lg tracking-tight">Scale with a plan</h2>
      <p className="mt-1.5 text-muted-foreground text-sm leading-6">
        Move production workloads to dedicated cloud resources with predictable pricing.
      </p>
      <Button
        asChild
        variant="secondary"
        className="mt-4 h-10 w-full text-foreground/85 hover:brightness-105"
      >
        <Link to="/billing">
          View plans
          <ArrowRight className="size-4 rtl:rotate-180" aria-hidden="true" />
        </Link>
      </Button>
    </section>
  );
}
