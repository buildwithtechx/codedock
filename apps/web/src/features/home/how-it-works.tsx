import { ArrowRight, Cloud, Lock, Server } from 'lucide-react';

export function HowItWorks() {
  return (
    <section className="relative border-border border-y bg-surface/30 py-24 backdrop-blur-sm">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mb-16 text-center">
          <h2 className="mb-4 font-extrabold text-3xl text-foreground md:text-4xl">How it works</h2>
          <p className="mx-auto max-w-xl text-muted-foreground text-sm">
            A secure, hybrid architecture. Codedock Cloud never touches your code or data—it only
            orchestrates the daemon.
          </p>
        </div>

        <div className="relative flex flex-col items-center justify-center gap-4 md:flex-row md:gap-8">
          <div className="relative w-full max-w-xs rounded-xl border border-border bg-background/80 p-6 text-center shadow-lg backdrop-blur-md">
            <div className="mx-auto mb-4 flex size-12 items-center justify-center rounded-lg border border-border bg-surface text-primary">
              <Server className="size-6" />
            </div>
            <h3 className="mb-2 font-bold text-foreground">1. Your Server</h3>
            <p className="text-muted-foreground text-xs leading-relaxed">
              Runs the &lt;30MB open-source daemon. Manages Docker, Traefik, and your actual
              databases &amp; apps.
            </p>
          </div>

          <div className="rotate-90 text-primary md:rotate-0">
            <ArrowRight className="size-6 md:size-8" />
          </div>

          <div className="relative w-full max-w-xs rounded-xl border border-primary/30 bg-primary/5 p-6 text-center shadow-lg shadow-primary/5 backdrop-blur-md">
            <div className="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-primary px-2 py-0.5 font-bold text-[10px] text-white uppercase tracking-wider">
              Encrypted
            </div>
            <div className="mx-auto mb-4 flex size-12 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
              <Lock className="size-6" />
            </div>
            <h3 className="mb-2 font-bold text-foreground">2. Yamux Tunnel</h3>
            <p className="text-muted-foreground text-xs leading-relaxed">
              The daemon dials out to Codedock Cloud. No port forwarding, no VPNs. Your server
              remains completely dark to the public internet.
            </p>
          </div>

          <div className="rotate-90 text-primary md:rotate-0">
            <ArrowRight className="size-6 md:size-8" />
          </div>

          <div className="relative w-full max-w-xs rounded-xl border border-border bg-background/80 p-6 text-center shadow-lg backdrop-blur-md">
            <div className="mx-auto mb-4 flex size-12 items-center justify-center rounded-lg border border-border bg-surface text-primary">
              <Cloud className="size-6" />
            </div>
            <h3 className="mb-2 font-bold text-foreground">3. Codedock Cloud</h3>
            <p className="text-muted-foreground text-xs leading-relaxed">
              The control plane. View logs, manage deployments, track metrics, and manage team
              access from a beautiful web dashboard.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
