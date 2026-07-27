import { ArrowRight, Cloud, Lock, Server } from "lucide-react";

export function HowItWorks() {
  return (
    <section className="py-24 border-y border-border bg-surface/30 backdrop-blur-sm relative">
      <div className="max-w-7xl mx-auto px-6">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-extrabold text-foreground mb-4">How it works</h2>
          <p className="text-muted-foreground text-sm max-w-xl mx-auto">
            A secure, hybrid architecture. Codedock Cloud never touches your code or data—it only
            orchestrates the daemon.
          </p>
        </div>

        <div className="flex flex-col md:flex-row items-center justify-center gap-4 md:gap-8 relative">
          <div className="bg-background/80 backdrop-blur-md border border-border rounded-xl p-6 w-full max-w-xs text-center shadow-lg relative">
            <div className="size-12 mx-auto bg-surface border border-border rounded-lg flex items-center justify-center mb-4 text-primary">
              <Server className="size-6" />
            </div>
            <h3 className="font-bold text-foreground mb-2">1. Your Server</h3>
            <p className="text-xs text-muted-foreground leading-relaxed">
              Runs the &lt;30MB open-source daemon. Manages Docker, Traefik, and your actual
              databases &amp; apps.
            </p>
          </div>

          <div className="md:rotate-0 rotate-90 text-primary">
            <ArrowRight className="size-6 md:size-8" />
          </div>

          <div className="bg-primary/5 backdrop-blur-md border border-primary/30 rounded-xl p-6 w-full max-w-xs text-center shadow-lg shadow-primary/5 relative">
            <div className="absolute -top-3 left-1/2 -translate-x-1/2 px-2 py-0.5 bg-primary text-white text-[10px] font-bold rounded-full uppercase tracking-wider">
              Encrypted
            </div>
            <div className="size-12 mx-auto bg-primary/10 border border-primary/20 rounded-lg flex items-center justify-center mb-4 text-primary">
              <Lock className="size-6" />
            </div>
            <h3 className="font-bold text-foreground mb-2">2. Yamux Tunnel</h3>
            <p className="text-xs text-muted-foreground leading-relaxed">
              The daemon dials out to Codedock Cloud. No port forwarding, no VPNs. Your server
              remains completely dark to the public internet.
            </p>
          </div>

          <div className="md:rotate-0 rotate-90 text-primary">
            <ArrowRight className="size-6 md:size-8" />
          </div>

          <div className="bg-background/80 backdrop-blur-md border border-border rounded-xl p-6 w-full max-w-xs text-center shadow-lg relative">
            <div className="size-12 mx-auto bg-surface border border-border rounded-lg flex items-center justify-center mb-4 text-primary">
              <Cloud className="size-6" />
            </div>
            <h3 className="font-bold text-foreground mb-2">3. Codedock Cloud</h3>
            <p className="text-xs text-muted-foreground leading-relaxed">
              The control plane. View logs, manage deployments, track metrics, and manage team
              access from a beautiful web dashboard.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
