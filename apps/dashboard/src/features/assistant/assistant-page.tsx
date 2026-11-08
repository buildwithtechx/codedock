import { Button } from '@codedock/ui/button';
import { Link } from '@tanstack/react-router';
import { ArrowRight, Bot, Cpu, Key, Server, Sparkles } from 'lucide-react';
import { PageHeader } from '#/components/layout/page-header';
import { useAiSidebarStore } from '#/stores/ai-sidebar-store';
import { AssistantIntroPreview } from './assistant-intro-preview';

export function AssistantPage() {
  const openSidebar = useAiSidebarStore((s) => s.openSidebar);

  return (
    <div className="space-y-6">
      <PageHeader
        title="AI Agents"
        description="Autonomous assistants and specialist agents to deploy, investigate, and maintain your services."
      />

      <div className="@container min-h-0 flex-1 overflow-y-auto">
        <section
          aria-labelledby="assistant-intro-title"
          className="mx-auto grid w-full max-w-[1040px] @[880px]:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)] items-center @[880px]:gap-12 gap-8 py-6 sm:py-10"
        >
          <div className="mx-auto w-full max-w-md">
            <h1
              id="assistant-intro-title"
              className="font-medium text-3xl text-foreground leading-tight tracking-tight sm:text-4xl"
            >
              Put your agents
              <br />
              to work.
            </h1>
            <p className="mt-4 text-muted-foreground text-sm leading-6">
              Deploy apps, investigate failures, and automate routine work. Connect the agents you
              already use, right here in Codedock.
            </p>

            <div className="mt-5 flex flex-wrap items-center gap-x-5 gap-y-3 text-foreground text-sm">
              <span className="inline-flex items-center gap-2">
                <Cpu className="h-4 w-4 text-primary" />
                Codex
              </span>
              <span className="inline-flex items-center gap-2">
                <Sparkles className="h-4 w-4 text-amber-500" />
                Claude Code
              </span>
              <span className="text-muted-foreground text-xs">or your own API</span>
            </div>

            <div className="mt-7 rounded-2xl border border-border/50 bg-card p-5 shadow-xs sm:p-6">
              <div className="space-y-4">
                <div>
                  <p className="font-medium text-foreground text-sm">
                    Run agents with zero lock-in
                  </p>
                  <p className="mt-1 text-muted-foreground text-xs leading-5">
                    Connect your own server, use this host, or configure external AI providers with
                    your own keys.
                  </p>
                </div>

                <div className="flex flex-col gap-2.5">
                  <Button
                    size="lg"
                    onClick={() => openSidebar()}
                    className="h-11 w-full gap-2.5 rounded-xl font-medium shadow-sm"
                  >
                    <Bot className="h-4 w-4" />
                    Open Assistant Co-pilot
                    <ArrowRight className="h-4 w-4" />
                  </Button>

                  <div className="flex items-center justify-between pt-2 text-xs">
                    <Link
                      to="/servers"
                      className="inline-flex items-center gap-1.5 text-muted-foreground transition-colors hover:text-foreground"
                    >
                      <Server className="h-3.5 w-3.5" />
                      Manage servers
                    </Link>

                    <Link
                      to="/settings"
                      search={{ tab: 'assistant' } as never}
                      className="inline-flex items-center gap-1.5 text-muted-foreground transition-colors hover:text-foreground"
                    >
                      <Key className="h-3.5 w-3.5" />
                      Configure AI models ↗
                    </Link>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <AssistantIntroPreview />
        </section>
      </div>
    </div>
  );
}
