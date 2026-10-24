export function EmptyIllustration({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      viewBox="0 0 248 168"
      fill="none"
      aria-hidden="true"
      role="presentation"
    >
      <rect x="72" y="52" width="120" height="86" rx="13" className="fill-muted/50" />
      <rect
        x="63"
        y="43"
        width="120"
        height="86"
        rx="13"
        className="fill-muted/70 stroke-border/60"
        strokeWidth="1"
      />
      <rect
        x="54"
        y="34"
        width="120"
        height="86"
        rx="13"
        className="fill-card stroke-border"
        strokeWidth="1"
      />
      <rect x="54" y="34" width="120" height="26" rx="13" className="fill-muted/60" />
      <circle cx="70" cy="47" r="3.6" className="fill-rose-500/60" />
      <circle cx="81" cy="47" r="3.6" className="fill-amber-400/70" />
      <circle cx="92" cy="47" r="3.6" className="fill-emerald-500/70" />
      <rect x="68" y="72" width="46" height="5" rx="2.5" className="fill-foreground/15" />
      <rect x="68" y="82" width="80" height="4" rx="2" className="fill-foreground/10" />
      <rect x="68" y="90" width="60" height="4" rx="2" className="fill-foreground/10" />
      <circle
        cx="82"
        cy="108"
        r="9"
        className="fill-primary/10 stroke-primary/40"
        strokeWidth="1"
      />
      <path
        d="M82 103.5v9M77.5 108h9"
        className="stroke-primary"
        strokeWidth="1.6"
        strokeLinecap="round"
      />
      <circle cx="204" cy="86" r="21" className="fill-primary/10" />
      <circle
        cx="204"
        cy="86"
        r="15"
        className="fill-card stroke-primary/50"
        strokeWidth="1.8"
        strokeDasharray="4 3"
      />
      <path
        d="M204 79v14M197 86h14"
        className="stroke-primary"
        strokeWidth="1.8"
        strokeLinecap="round"
      />
      <path
        d="M178 92 Q 186 90 189 88"
        className="stroke-border"
        strokeWidth="1.4"
        strokeDasharray="3 3"
      />
      <circle cx="30" cy="58" r="3.5" className="fill-foreground/10" />
      <circle cx="38" cy="134" r="5.5" className="fill-foreground/8" />
      <circle cx="222" cy="38" r="3" className="fill-foreground/12" />
      <circle cx="236" cy="124" r="4.5" className="fill-foreground/8" />
    </svg>
  );
}
