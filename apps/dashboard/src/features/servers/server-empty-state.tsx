import { Link } from '@tanstack/react-router';
import {
  Activity,
  BookOpen,
  Container,
  Cpu,
  ExternalLink,
  GitBranch,
  Globe,
  Plus,
  Server as ServerIcon,
} from 'lucide-react';
import { Button } from '#/components/ui/button';

const DOCS_URL = 'https://docs.codedock.run';

const capabilities = [
  { label: 'Containers', desc: 'Docker workloads on your host', icon: Container },
  { label: 'Networking', desc: 'Public routes for your apps', icon: Globe },
  { label: 'Monitoring', desc: 'Live CPU, memory and disk', icon: Activity },
  { label: 'Git', desc: 'Deploy straight from a repo', icon: GitBranch },
];

interface ServerEmptyStateProps {
  onAddLocal?: () => void;
  addingLocal?: boolean;
}

export function ServerEmptyState({ onAddLocal, addingLocal }: ServerEmptyStateProps) {
  return (
    <div className="py-16 text-center">
      <div className="relative mx-auto mb-8 h-44 w-64" aria-hidden="true">
        <svg className="absolute inset-0 h-full w-full" viewBox="0 0 260 180" fill="none">
          <rect x="62" y="56" width="112" height="84" rx="14" className="fill-muted" />
          <rect
            x="50"
            y="44"
            width="112"
            height="84"
            rx="14"
            className="fill-card stroke-border"
            strokeWidth="1"
          />
          <rect
            x="38"
            y="32"
            width="112"
            height="84"
            rx="14"
            className="fill-card stroke-border"
            strokeWidth="1"
          />
          <rect x="52" y="46" width="70" height="14" rx="4" className="fill-muted" />
          <rect x="58" y="51" width="34" height="4" rx="2" className="fill-muted-foreground/40" />
          <circle cx="134" cy="53" r="3" className="fill-success" fillOpacity="0.8" />
          <rect x="52" y="66" width="70" height="14" rx="4" className="fill-muted/60" />
          <rect x="58" y="71" width="34" height="4" rx="2" className="fill-muted-foreground/25" />
          <circle cx="134" cy="73" r="3" className="fill-muted-foreground/30" />
          <rect x="52" y="86" width="70" height="14" rx="4" className="fill-muted/60" />
          <rect x="58" y="91" width="24" height="4" rx="2" className="fill-muted-foreground/25" />
          <circle cx="134" cy="93" r="3" className="fill-muted-foreground/30" />
          <path
            d="M150 84 Q 172 82 186 84"
            className="stroke-muted-foreground/40"
            strokeWidth="1.5"
            strokeDasharray="3 3"
          />
          <circle cx="208" cy="84" r="22" className="fill-muted/60" />
          <circle
            cx="208"
            cy="84"
            r="15"
            className="fill-card stroke-muted-foreground/40"
            strokeWidth="2"
            strokeDasharray="4 3"
          />
          <path
            d="M208 77v14M201 84h14"
            className="stroke-muted-foreground"
            strokeWidth="2"
            strokeLinecap="round"
          />
          <circle cx="24" cy="56" r="4" className="fill-muted" />
          <circle cx="32" cy="144" r="6" className="fill-muted/70" />
          <circle cx="238" cy="38" r="3.5" className="fill-muted" />
          <circle cx="246" cy="130" r="4.5" className="fill-muted/70" />
          <rect x="20" y="102" width="5" height="5" rx="1" className="fill-primary/40" />
          <rect x="226" y="148" width="5" height="5" rx="1" className="fill-primary/30" />
        </svg>
      </div>

      <h3 className="mb-2 font-medium text-2xl text-foreground/80 tracking-[-0.2px]">
        No servers connected
      </h3>
      <p className="mx-auto mb-8 max-w-sm text-muted-foreground/70 text-sm leading-relaxed">
        Connect a server so Codedock can deploy, run, and observe your applications on your own
        infrastructure.
      </p>

      <div className="mb-10 flex flex-col items-center justify-center gap-3 sm:flex-row">
        <Link to="/servers/new">
          <Button size="lg" className="gap-2 px-6">
            <Plus className="h-4 w-4" />
            Add your first server
          </Button>
        </Link>
        {onAddLocal && (
          <Button
            size="lg"
            variant="secondary"
            className="gap-2 px-6"
            onClick={onAddLocal}
            disabled={addingLocal}
          >
            <Cpu className="h-4 w-4" />
            {addingLocal ? 'Connecting…' : 'Use this machine'}
          </Button>
        )}
        <Button asChild size="lg" variant="secondary" className="gap-2 px-6">
          <a href={DOCS_URL} target="_blank" rel="noopener noreferrer">
            <BookOpen className="h-4 w-4" />
            Docs
            <ExternalLink className="h-3.5 w-3.5 opacity-60" />
          </a>
        </Button>
      </div>

      <div className="mx-auto max-w-2xl">
        <p className="mb-4 text-muted-foreground/60 text-xs uppercase tracking-wider">
          What gets configured
        </p>
        <div className="grid grid-cols-2 gap-3 text-start sm:grid-cols-4">
          {capabilities.map((feature) => (
            <div key={feature.label} className="rounded-xl border border-border/50 bg-card p-4">
              <div className="mb-3 flex h-8 w-8 items-center justify-center rounded-lg bg-muted">
                <feature.icon className="h-4 w-4 text-muted-foreground" />
              </div>
              <p className="font-medium text-foreground text-sm">{feature.label}</p>
              <p className="mt-0.5 text-muted-foreground text-xs">{feature.desc}</p>
            </div>
          ))}
        </div>
        <p className="mt-6 inline-flex items-center gap-1.5 text-muted-foreground/60 text-xs">
          <ServerIcon className="h-3.5 w-3.5" />
          Bring any host with SSH access, or start with the local Docker node.
        </p>
      </div>
    </div>
  );
}
