import { Link } from '@tanstack/react-router';
import {
  AlertTriangle,
  ArrowRight,
  Box,
  Database,
  ExternalLink,
  Globe,
  KeyRound,
  LifeBuoy,
  Rocket,
  SlidersHorizontal,
  Sparkles,
  Terminal,
} from 'lucide-react';
import { Avatar, AvatarBadge, AvatarFallback, AvatarImage } from '@codedock/ui/avatar';
import { CodedockLogo } from '#/components/layout/codedock-logo';
import { useGetPublicSettings } from '#/hooks/use-settings';
import type { AiProviderStatusSummary } from '#/interfaces/ai';

interface StarterCard {
  id: string;
  title: string;
  subtitle: string;
  prompt: string;
  icon: typeof Box;
}

const STARTER_CARDS: StarterCard[] = [
  {
    id: 'deploy',
    title: 'Deploy a Service or Template',
    subtitle: 'Inspect projects & deploy from Git or Catalog',
    prompt:
      'Show me the current projects and services in this instance and walk me through deploying a new service from Git or the One-Click App Catalog.',
    icon: Rocket,
  },
  {
    id: 'diagnose',
    title: 'Diagnose Containers & Logs',
    subtitle: 'Check unhealthy services & recent builds',
    prompt:
      'Inspect all services, deployments, and active issues on this instance and highlight anything unhealthy or failing.',
    icon: Terminal,
  },
  {
    id: 'env',
    title: 'Environment Variables & Secrets',
    subtitle: 'Configure runtime variables across services',
    prompt:
      'How do I add, bulk-import, or update environment variables and secrets for my services in Codedock?',
    icon: SlidersHorizontal,
  },
  {
    id: 'domain',
    title: 'Custom Domains & Wildcard TLS',
    subtitle: 'Connect DNS records, Traefik & SSL',
    prompt:
      'Walk me through connecting a custom domain, configuring DNS records, and enabling automated wildcard TLS or Cloudflare Tunnels.',
    icon: Globe,
  },
  {
    id: 'storage',
    title: 'Databases & S3 Backups',
    subtitle: 'Provision PostgreSQL/Redis & backup policies',
    prompt:
      'List the databases and backup destinations on this instance and show me how to connect PostgreSQL/Redis or schedule S3 backups.',
    icon: Database,
  },
];

function getTimeGreeting(): string {
  const hour = new Date().getHours();
  if (hour < 12) return 'Good morning';
  if (hour < 18) return 'Good afternoon';
  return 'Good evening';
}

function CodedockAiEmblem() {
  return (
    <div className="relative mx-auto flex h-20 w-20 items-center justify-center select-none">
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 rounded-2xl bg-primary/15 blur-xl"
      />
      <Avatar
        size="lg"
        className="size-14 overflow-visible rounded-2xl shadow-lg ring-1 ring-primary/30"
      >
        <AvatarImage src="/favicon.svg" alt="Codedock Copilot" className="rounded-2xl" />
        <AvatarFallback className="rounded-2xl bg-transparent p-0">
          <CodedockLogo className="size-14 rounded-2xl" />
        </AvatarFallback>
        <AvatarBadge className="-top-1.5 -right-1.5 bottom-auto size-5 border border-border bg-card text-emerald-500 shadow-xs ring-0">
          <Sparkles className="size-3 text-emerald-500" />
        </AvatarBadge>
      </Avatar>
    </div>
  );
}

export function AskAiWelcome({
  providerStatus,
  onSelectPrompt,
}: {
  providerStatus: AiProviderStatusSummary;
  onSelectPrompt: (prompt: string) => void;
}) {
  const greeting = getTimeGreeting();
  const { data: publicRes } = useGetPublicSettings();
  const isCloud = Boolean(publicRes?.data?.cloudMode);

  return (
    <div className="flex flex-1 flex-col justify-between gap-3 overflow-y-auto px-4 py-4">
      <div className="space-y-2.5">
        {/* Need more help? [Support] Banner -> routes to /support on Cloud */}
        <div className="flex items-center justify-between rounded-xl border border-border/60 bg-card/90 px-3.5 py-2.5">
          <span className="inline-flex items-center gap-2 font-medium text-foreground text-xs">
            <LifeBuoy className="size-3.5 text-muted-foreground" />
            Need more help?
          </span>
          {isCloud ? (
            <Link
              to="/support"
              className="inline-flex h-7 items-center rounded-lg border border-border/80 bg-muted/40 px-2.5 font-medium text-foreground text-xs transition-colors hover:bg-muted"
            >
              Support
            </Link>
          ) : (
            <a
              href="https://docs.codedock.run"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex h-7 items-center gap-1 rounded-lg border border-border/80 bg-muted/40 px-2.5 font-medium text-foreground text-xs transition-colors hover:bg-muted"
            >
              <span>Docs & Support</span>
              <ExternalLink className="size-3 text-muted-foreground" />
            </a>
          )}
        </div>

        {/* Provider Configuration Status Banner */}
        {providerStatus.state === 'unconfigured' ? (
          <div className="rounded-xl border border-amber-500/35 bg-amber-500/10 p-3.5 text-xs">
            <div className="flex items-start gap-2.5">
              <KeyRound className="mt-0.5 size-4 shrink-0 text-amber-500" />
              <div className="min-w-0 flex-1">
                <p className="font-semibold text-foreground">No AI Provider Configured</p>
                <p className="mt-1 text-[11px] text-muted-foreground leading-relaxed">
                  Add an API key (OpenAI, Anthropic, Gemini, Groq, DeepSeek, Mistral, xAI, or
                  Moonshot) and choose a default provider to unlock full LLM reasoning. Local
                  Instance Copilot is active in the meantime.
                </p>
                <Link
                  to="/settings"
                  search={{ tab: 'ai' } as never}
                  className="mt-2.5 inline-flex h-7 items-center gap-1.5 rounded-lg bg-primary px-2.5 font-medium text-[11px] text-primary-foreground transition-opacity hover:opacity-90"
                >
                  <span>Configure AI Provider</span>
                  <ArrowRight className="size-3" />
                </Link>
              </div>
            </div>
          </div>
        ) : providerStatus.state === 'no_default' ? (
          <div className="rounded-xl border border-primary/30 bg-primary/5 p-3 text-xs">
            <div className="flex items-start justify-between gap-2.5">
              <div className="flex items-start gap-2">
                <AlertTriangle className="mt-0.5 size-3.5 shrink-0 text-amber-500" />
                <div className="min-w-0">
                  <p className="font-medium text-foreground">No Default Provider Selected</p>
                  <p className="mt-0.5 text-[11px] text-muted-foreground">
                    Using{' '}
                    <span className="font-medium text-foreground">
                      {providerStatus.activeProviderName}
                    </span>{' '}
                    ({providerStatus.activeModel}) as fallback.
                  </p>
                </div>
              </div>
              <Link
                to="/settings"
                search={{ tab: 'ai' } as never}
                className="shrink-0 rounded-lg border border-border/80 bg-card px-2.5 py-1 font-medium text-[11px] text-foreground transition-colors hover:bg-muted"
              >
                Set Default
              </Link>
            </div>
          </div>
        ) : (
          <div className="flex items-center justify-between rounded-xl border border-border/60 bg-card/90 px-3.5 py-2">
            <span className="inline-flex items-center gap-1.5 text-[11px] text-muted-foreground">
              <span className="size-1.5 rounded-full bg-emerald-500" />
              Active Model
            </span>
            <Link
              to="/settings"
              search={{ tab: 'ai' } as never}
              className="font-mono font-medium text-[11px] text-foreground hover:text-primary"
            >
              {providerStatus.activeProviderName} · {providerStatus.activeModel}
            </Link>
          </div>
        )}
      </div>

      {/* Center Codedock Hero + Starter Prompts */}
      <div className="my-auto py-4">
        <CodedockAiEmblem />

        <div className="mt-2 text-center">
          <h2 className="font-semibold text-foreground text-base tracking-tight">
            {greeting} — Codedock Copilot
          </h2>
          <p className="mt-0.5 text-muted-foreground text-xs">
            Inspect live workloads, diagnose logs, or run MCP actions on this instance.
          </p>
        </div>

        <div className="mt-5 space-y-2">
          {STARTER_CARDS.map((card) => {
            const Icon = card.icon;
            return (
              <button
                key={card.id}
                type="button"
                onClick={() => onSelectPrompt(card.prompt)}
                className="group flex w-full items-center gap-3 rounded-xl border border-border/60 bg-card/90 px-3.5 py-2.5 text-left transition-all hover:border-primary/40 hover:bg-muted/40"
              >
                <div className="flex size-8 shrink-0 items-center justify-center rounded-lg border border-border/60 bg-muted/30 text-muted-foreground transition-colors group-hover:border-primary/30 group-hover:text-primary">
                  <Icon className="size-4" />
                </div>
                <div className="min-w-0 flex-1">
                  <p className="truncate font-medium text-foreground text-xs">{card.title}</p>
                  <p className="mt-0.5 truncate text-[11px] text-muted-foreground">
                    {card.subtitle}
                  </p>
                </div>
              </button>
            );
          })}
        </div>
      </div>
    </div>
  );
}
