import { useId } from 'react';

export function ActionsIllustration({
  kind = 'workflow',
  className,
}: {
  kind?: 'workflow' | 'runner' | 'history';
  className?: string;
}) {
  const id = useId();
  const runner = kind === 'runner';

  return (
    <svg
      viewBox="0 0 440 256"
      fill="none"
      className={className}
      aria-hidden="true"
      focusable="false"
    >
      <defs>
        <pattern id={`${id}-grid`} width="18" height="18" patternUnits="userSpaceOnUse">
          <circle cx="1" cy="1" r="1" fill="currentColor" fillOpacity="0.12" />
        </pattern>
        <radialGradient id={`${id}-fade`}>
          <stop offset="0.35" stopColor="white" />
          <stop offset="1" stopColor="black" />
        </radialGradient>
        <mask id={`${id}-mask`}>
          <rect width="440" height="256" fill={`url(#${id}-fade)`} />
        </mask>
      </defs>

      <rect width="440" height="256" fill={`url(#${id}-grid)`} mask={`url(#${id}-mask)`} />
      <ellipse cx="211" cy="211" rx="153" ry="17" fill="currentColor" fillOpacity="0.04" />

      {/* Layered source window */}
      <rect x="91" y="61" width="230" height="150" rx="17" fill="currentColor" fillOpacity="0.06" />
      <rect x="80" y="48" width="230" height="150" rx="17" className="fill-card stroke-border/40" />
      <rect x="68" y="35" width="230" height="150" rx="17" className="fill-card stroke-border" />
      <path
        d="M85 35h196a17 17 0 0 1 17 17v17H68V52a17 17 0 0 1 17-17Z"
        fill="currentColor"
        fillOpacity="0.04"
      />
      <circle cx="85" cy="52" r="3" fill="currentColor" fillOpacity="0.25" />
      <circle cx="97" cy="52" r="3" fill="currentColor" fillOpacity="0.15" />
      <circle cx="109" cy="52" r="3" fill="currentColor" fillOpacity="0.1" />
      <rect x="216" y="48" width="62" height="7" rx="3.5" fill="currentColor" fillOpacity="0.1" />

      {runner ? (
        <>
          {[87, 119, 151].map((y, index) => (
            <g key={y}>
              <rect
                x="87"
                y={y}
                width="174"
                height="23"
                rx="6"
                fill="currentColor"
                fillOpacity="0.04"
                className="stroke-border/40"
              />
              <circle
                cx="100"
                cy={y + 11.5}
                r="3"
                className={index === 0 ? 'fill-emerald-500' : 'fill-muted-foreground/30'}
              />
              <path
                d={`M113 ${y + 11.5}h44`}
                stroke="currentColor"
                strokeOpacity="0.2"
                strokeWidth="4"
                strokeLinecap="round"
              />
              <path
                d={`M228 ${y + 8}v7m8-7v7m8-7v7`}
                stroke="currentColor"
                strokeOpacity="0.25"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </g>
          ))}
        </>
      ) : (
        <>
          <path
            d="m88 88 5 4-5 4m13-4h15"
            className="stroke-primary"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
          <rect
            x="126"
            y="89"
            width="73"
            height="5"
            rx="2.5"
            fill="currentColor"
            fillOpacity="0.3"
          />
          <path d="M95 111v42" stroke="currentColor" strokeOpacity="0.12" strokeWidth="1.5" />
          <rect
            x="105"
            y="109"
            width="53"
            height="5"
            rx="2.5"
            fill="currentColor"
            fillOpacity="0.15"
          />
          <rect
            x="105"
            y="124"
            width="87"
            height="5"
            rx="2.5"
            fill="currentColor"
            fillOpacity="0.1"
          />
          <rect
            x="116"
            y="139"
            width="66"
            height="5"
            rx="2.5"
            fill="currentColor"
            fillOpacity="0.15"
          />
          <rect
            x="116"
            y="154"
            width="92"
            height="5"
            rx="2.5"
            fill="currentColor"
            fillOpacity="0.1"
          />
        </>
      )}

      {/* Connected jobs */}
      <path
        d="M298 114h19a14 14 0 0 0 14-14V84m-33 30h19a14 14 0 0 1 14 14v23"
        stroke="currentColor"
        strokeOpacity="0.2"
        strokeWidth="1.5"
      />
      <circle
        cx="298"
        cy="114"
        r="4"
        className="fill-card stroke-muted-foreground/30"
        strokeWidth="1.5"
      />
      <rect x="307" y="43" width="87" height="42" rx="12" className="fill-card stroke-border" />
      <rect x="319" y="55" width="18" height="18" rx="5" className="fill-emerald-500/15" />
      <path
        d="m324 64 3 3 5-6"
        className="stroke-emerald-500"
        strokeWidth="1.6"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <rect x="345" y="59" width="33" height="4" rx="2" fill="currentColor" fillOpacity="0.3" />
      <rect x="345" y="67" width="23" height="3" rx="1.5" fill="currentColor" fillOpacity="0.1" />

      <rect x="307" y="151" width="87" height="42" rx="12" className="fill-card stroke-border" />
      <rect x="319" y="163" width="18" height="18" rx="5" className="fill-blue-500/15" />
      {kind === 'history' ? (
        <>
          <circle cx="328" cy="172" r="5" className="stroke-blue-500" strokeWidth="1.5" />
          <path
            d="M328 169v3l2 1"
            className="stroke-blue-500"
            strokeWidth="1.5"
            strokeLinecap="round"
          />
        </>
      ) : (
        <path d="m326 168 5 4-5 4v-8Z" className="fill-blue-500" />
      )}
      <rect x="345" y="167" width="29" height="4" rx="2" fill="currentColor" fillOpacity="0.3" />
      <rect x="345" y="175" width="19" height="3" rx="1.5" fill="currentColor" fillOpacity="0.1" />

      {/* Source badge */}
      <rect x="36" y="97" width="52" height="52" rx="14" className="fill-card stroke-border" />
      <path
        d="M55 116v16m0-8h7a8 8 0 0 0 8-8"
        stroke="currentColor"
        strokeOpacity="0.4"
        strokeWidth="2"
        strokeLinecap="round"
      />
      <circle cx="55" cy="113" r="3" className="fill-card stroke-foreground/40" strokeWidth="2" />
      <circle cx="55" cy="135" r="3" className="fill-card stroke-foreground/40" strokeWidth="2" />
      <circle cx="70" cy="113" r="3" className="fill-card stroke-foreground/40" strokeWidth="2" />
      <rect x="160" y="178" width="123" height="39" rx="12" className="fill-card stroke-border" />
      <circle cx="181" cy="197.5" r="10" className="fill-emerald-500/15" />
      <path
        d="m177 198 3 3 5-6"
        className="stroke-emerald-500"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <rect x="200" y="191" width="64" height="5" rx="2.5" fill="currentColor" fillOpacity="0.3" />
      <rect x="200" y="201" width="42" height="4" rx="2" fill="currentColor" fillOpacity="0.1" />
    </svg>
  );
}
