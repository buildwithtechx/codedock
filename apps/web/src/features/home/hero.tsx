import { Check, Cloud, Copy, Server } from 'lucide-react';
import { useState } from 'react';

export function Hero() {
  const [copied, setCopied] = useState(false);

  const serverRows = [
    { k: 'HOST', v: 'prod-eu-1.acme.com' },
    { k: 'AGENT', v: 'v0.5.1' },
    { k: 'PROXY', v: 'traefik v3.1' },
    { k: 'NETWORK', v: 'codedock-network' },
    { k: 'TUNNEL', v: 'yamux / active', accent: true },
  ];

  const capabilities = [
    'apps',
    'databases',
    'cron-jobs',
    'backups → S3',
    'domains + SSL',
    'env secrets',
    'canvas view',
    'log tailing',
  ];

  const handleCopy = () => {
    navigator.clipboard.writeText('curl -fsSL https://get.codedock.run | bash');
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
    if (
      typeof window !== 'undefined' &&
      (
        window as {
          posthog?: { capture: (event: string, props?: Record<string, unknown>) => void };
        }
      ).posthog
    ) {
      (
        window as {
          posthog?: { capture: (event: string, props?: Record<string, unknown>) => void };
        }
      ).posthog?.capture('install_command_copied', {
        source: 'hero_section',
        command: 'curl -fsSL https://get.codedock.run | bash',
      });
    }
  };

  return (
    <section className="relative isolate overflow-hidden py-24 lg:py-32">
      <div className="relative mx-auto max-w-4xl px-6 text-center">
        <div className="mb-8 inline-flex items-center gap-2 rounded-full border border-border bg-background/80 px-3 py-1 font-medium text-muted-foreground text-xs backdrop-blur-md">
          <span className="size-1.5 animate-pulse rounded-full bg-primary" />
          Open source &amp; free forever
        </div>

        <h1 className="text-balance font-extrabold text-4xl text-foreground leading-tight tracking-tight sm:text-5xl sm:leading-none lg:text-6xl">
          The open-source Heroku & Vercel alternative<span className="text-primary">.</span>
        </h1>

        <p className="mx-auto mt-6 max-w-2xl text-base text-muted-foreground leading-relaxed md:text-lg">
          Imagine having the ease of a cloud but with your own servers. Codedock is a sub-30MB
          daemon that turns any Linux server into a full deployment platform—with zero vendor
          lock-in.
        </p>

        <div className="mx-auto mt-10 flex max-w-lg items-center gap-2 rounded-lg border border-border bg-background/80 px-4 py-3 font-mono text-sm backdrop-blur-md">
          <span className="shrink-0 font-bold text-primary">$</span>
          <code className="flex-1 truncate text-left text-[0.8rem] text-foreground">
            curl -fsSL https://get.codedock.run | bash
          </code>
          <button
            type="button"
            onClick={handleCopy}
            aria-label="Copy install command"
            className="shrink-0 rounded p-1.5 text-muted-foreground transition-colors hover:bg-border hover:text-foreground"
          >
            {copied ? <Check className="size-3.5 text-green-400" /> : <Copy className="size-3.5" />}
          </button>
        </div>

        <div className="mt-6 flex flex-row items-start justify-center gap-3 sm:gap-4">
          <div className="flex flex-col items-center">
            <a
              href="https://app.codedock.run"
              className="flex items-center justify-center gap-2 whitespace-nowrap rounded-lg border border-border bg-background/80 px-6 py-3.5 font-semibold text-foreground text-sm backdrop-blur-md transition-colors hover:bg-border sm:px-10"
            >
              <Cloud className="size-4 text-primary" />
              Cloud
            </a>
            <p className="mt-2 text-center font-bold text-[10px] text-primary sm:text-xs">
              1,200+{' '}
              <span className="block font-normal text-muted-foreground sm:inline">
                servers online
              </span>
            </p>
          </div>
          <div className="flex flex-col items-center">
            <a
              href="https://docs.codedock.run/self-host"
              target="_blank"
              rel="noreferrer"
              className="flex items-center justify-center gap-2 whitespace-nowrap rounded-lg border border-border bg-background/80 px-6 py-3.5 font-semibold text-foreground text-sm backdrop-blur-md transition-colors hover:bg-border sm:px-10"
            >
              <Server className="size-4 text-primary" />
              Self-hosted
            </a>
            <p className="mt-2 text-center font-bold text-[10px] text-primary sm:text-xs">
              &lt;30MB{' '}
              <span className="block font-normal text-muted-foreground sm:inline">
                single binary
              </span>
            </p>
          </div>
        </div>

        <div className="relative z-10 mt-12 overflow-hidden rounded-xl border border-border bg-background/80 text-left font-mono text-xs shadow-2xl backdrop-blur-md">
          <div className="flex items-center justify-between border-border border-b bg-background/40 px-4 py-3">
            <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
              codedock daemon
            </span>
            <span className="flex items-center gap-1.5 font-semibold text-[10px] text-green-400">
              <span className="size-1.5 animate-pulse rounded-full bg-green-400" />
              CONNECTED
            </span>
          </div>
          <div className="grid grid-cols-1 divide-y divide-border sm:grid-cols-2 sm:divide-x sm:divide-y-0">
            <div className="space-y-2.5 px-5 py-4">
              {serverRows.map((r) => (
                <div key={r.k} className="flex items-center gap-3 text-[0.71rem]">
                  <span className="w-20 shrink-0 text-muted-foreground">{r.k}</span>
                  <span className={r.accent ? 'font-semibold text-primary' : 'text-foreground'}>
                    {r.v}
                  </span>
                </div>
              ))}
            </div>
            <div className="px-5 py-4">
              <div className="mb-3 font-bold text-[9px] text-muted-foreground uppercase tracking-widest">
                Capabilities on this node
              </div>
              <div className="grid grid-cols-2 gap-x-4 gap-y-2">
                {capabilities.map((cap) => (
                  <div
                    key={cap}
                    className="flex items-center gap-2 text-[0.68rem] text-muted-foreground"
                  >
                    <span className="size-1.5 shrink-0 rounded-full bg-primary/70" />
                    {cap}
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
