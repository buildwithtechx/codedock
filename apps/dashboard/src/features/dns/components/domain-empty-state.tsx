import { Link } from '@tanstack/react-router';
import {
  BookOpen,
  CheckCircle2,
  ExternalLink,
  FolderGit2,
  Globe,
  Network,
  ShieldCheck,
} from 'lucide-react';
import { Button } from '@codedock/ui/button';

const DOCS_URL = 'https://docs.codedock.run/networking/domains';

const capabilities = [
  {
    label: 'Automated TLS',
    desc: "Let's Encrypt certificates provisioned automatically with hands-free renewal.",
    icon: ShieldCheck,
  },
  {
    label: 'Edge Ingress',
    desc: 'Traefik proxy directs apex and subdomains straight to your container ports.',
    icon: Network,
  },
  {
    label: 'Live Probing',
    desc: 'Instant propagation checks verify A and CNAME records with live reporting.',
    icon: CheckCircle2,
  },
  {
    label: 'Wildcard SSL',
    desc: 'Configure catch-all wildcard rules (*.domain.com) for preview environments.',
    icon: Globe,
  },
];

export function DomainEmptyState() {
  return (
    <div className="py-10 text-center">
      <div className="relative mx-auto mb-8 h-44 w-64" aria-hidden="true">
        <svg className="absolute inset-0 h-full w-full" viewBox="0 0 260 180" fill="none">
          {/* Background card */}
          <rect x="58" y="52" width="124" height="84" rx="14" className="fill-muted/70" />
          {/* Middle card */}
          <rect
            x="48"
            y="40"
            width="124"
            height="84"
            rx="14"
            className="fill-card stroke-border/70"
            strokeWidth="1"
          />
          {/* Front card */}
          <rect
            x="38"
            y="28"
            width="124"
            height="84"
            rx="14"
            className="fill-card stroke-border"
            strokeWidth="1"
          />

          {/* Domain bar row 1 */}
          <rect x="52" y="42" width="80" height="13" rx="4" className="fill-muted" />
          <rect
            x="58"
            y="47"
            width="44"
            height="3.5"
            rx="1.5"
            className="fill-muted-foreground/40"
          />
          <circle cx="146" cy="48.5" r="3.5" className="fill-emerald-500" fillOpacity="0.9" />

          {/* Domain bar row 2 */}
          <rect x="52" y="61" width="80" height="13" rx="4" className="fill-muted/60" />
          <rect
            x="58"
            y="66"
            width="36"
            height="3.5"
            rx="1.5"
            className="fill-muted-foreground/30"
          />
          <circle cx="146" cy="67.5" r="3.5" className="fill-primary" fillOpacity="0.8" />

          {/* Domain bar row 3 */}
          <rect x="52" y="80" width="80" height="13" rx="4" className="fill-muted/60" />
          <rect
            x="58"
            y="85"
            width="48"
            height="3.5"
            rx="1.5"
            className="fill-muted-foreground/30"
          />
          <circle cx="146" cy="86.5" r="3.5" className="fill-muted-foreground/40" />

          {/* Connection line */}
          <path
            d="M162 70 Q 185 68 198 70"
            className="stroke-muted-foreground/40"
            strokeWidth="1.5"
            strokeDasharray="3 3"
          />

          {/* Globe node */}
          <circle cx="218" cy="70" r="22" className="fill-primary/10" />
          <circle
            cx="218"
            cy="70"
            r="16"
            className="fill-card stroke-primary/40"
            strokeWidth="1.5"
            strokeDasharray="3 3"
          />
          <circle cx="218" cy="70" r="7" className="fill-primary/30" />

          {/* Decorative floating dots */}
          <circle cx="26" cy="50" r="3.5" className="fill-muted" />
          <circle cx="34" cy="136" r="5" className="fill-muted/70" />
          <circle cx="236" cy="36" r="3" className="fill-muted" />
          <circle cx="242" cy="120" r="4.5" className="fill-muted/70" />
          <rect x="18" y="96" width="5" height="5" rx="1" className="fill-primary/40" />
          <rect x="222" y="136" width="5" height="5" rx="1" className="fill-emerald-500/30" />
        </svg>
      </div>

      <h3 className="mb-2 font-medium text-2xl text-foreground/90 tracking-[-0.2px]">
        No domains configured
      </h3>
      <p className="mx-auto mb-8 max-w-lg text-muted-foreground/75 text-sm leading-relaxed">
        Route custom domains and subdomains to your application services with automatic SSL/TLS
        certificates, zero downtime, and real-time DNS verification.
      </p>

      <div className="mb-12 flex flex-col items-center justify-center gap-3 sm:flex-row">
        <Button asChild size="lg" className="gap-2 px-6">
          <Link to="/projects">
            <FolderGit2 className="h-4 w-4" />
            View Projects & Services
          </Link>
        </Button>
        <Button asChild size="lg" variant="secondary" className="gap-2 px-6">
          <a href={DOCS_URL} target="_blank" rel="noopener noreferrer">
            <BookOpen className="h-4 w-4" />
            Domain Docs
            <ExternalLink className="h-3.5 w-3.5 opacity-60" />
          </a>
        </Button>
      </div>

      <div className="mx-auto max-w-5xl">
        <p className="mb-4 text-muted-foreground/60 text-xs uppercase tracking-wider">
          What gets configured
        </p>
        <div className="grid grid-cols-1 gap-3.5 text-start sm:grid-cols-2 lg:grid-cols-4">
          {capabilities.map((item) => (
            <div
              key={item.label}
              className="flex flex-col justify-between rounded-2xl border border-border/50 bg-card p-4 shadow-xs transition-colors hover:border-border/80"
            >
              <div>
                <div className="mb-3 flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
                  <item.icon className="h-4 w-4" />
                </div>
                <h4 className="font-semibold text-[13px] text-foreground tracking-tight leading-snug">
                  {item.label}
                </h4>
                <p className="mt-1 text-muted-foreground text-xs leading-relaxed">{item.desc}</p>
              </div>
            </div>
          ))}
        </div>
        <p className="mt-8 inline-flex items-center gap-2 text-muted-foreground/65 text-xs">
          <Globe className="h-4 w-4 text-primary/70" />
          Domains are assigned to individual services under each project's networking settings.
        </p>
      </div>
    </div>
  );
}
