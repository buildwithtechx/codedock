import { Check, Cloud, Copy, Server } from "lucide-react";
import { useState } from "react";

export function Hero() {
  const [copied, setCopied] = useState(false);

  const serverRows = [
    { k: "HOST", v: "prod-eu-1.acme.com" },
    { k: "AGENT", v: "v0.5.1" },
    { k: "PROXY", v: "traefik v3.1" },
    { k: "NETWORK", v: "codedock-network" },
    { k: "TUNNEL", v: "yamux / active", accent: true },
  ];

  const capabilities = [
    "apps",
    "databases",
    "cron-jobs",
    "backups → S3",
    "domains + SSL",
    "env secrets",
    "canvas view",
    "log tailing",
  ];

  const handleCopy = () => {
    navigator.clipboard.writeText("curl -fsSL https://get.codedock.run | sh");
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <section className="relative isolate overflow-hidden py-24 lg:py-32">
      <div className="relative mx-auto max-w-4xl px-6 text-center">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-border bg-background/80 backdrop-blur-md text-xs text-muted-foreground mb-8 font-medium">
          <span className="size-1.5 rounded-full bg-primary animate-pulse" />
          Open source &amp; free forever
        </div>

        <h1 className="text-4xl sm:text-5xl lg:text-6xl font-extrabold tracking-tight leading-tight sm:leading-none text-foreground text-balance">
          The open-source Heroku & Vercel alternative<span className="text-primary">.</span>
        </h1>

        <p className="mt-6 text-muted-foreground text-base md:text-lg leading-relaxed mx-auto max-w-2xl">
          Imagine having the ease of a cloud but with your own servers. Codedock is a sub-30MB
          daemon that turns any Linux server into a full deployment platform—with zero vendor
          lock-in.
        </p>

        <div className="mt-10 mx-auto max-w-lg flex items-center gap-2 rounded-lg border border-border bg-background/80 backdrop-blur-md px-4 py-3 font-mono text-sm">
          <span className="text-primary font-bold shrink-0">$</span>
          <code className="flex-1 text-left text-foreground truncate text-[0.8rem]">
            curl -fsSL https://get.codedock.run | sh
          </code>
          <button
            type="button"
            onClick={handleCopy}
            aria-label="Copy install command"
            className="shrink-0 p-1.5 rounded text-muted-foreground hover:text-foreground hover:bg-border transition-colors"
          >
            {copied ? <Check className="size-3.5 text-green-400" /> : <Copy className="size-3.5" />}
          </button>
        </div>

        <div className="mt-6 flex flex-row items-start justify-center gap-3 sm:gap-4">
          <div className="flex flex-col items-center">
            <a
              href="https://app.codedock.run"
              className="flex items-center justify-center gap-2 px-6 sm:px-10 py-3.5 rounded-lg font-semibold text-sm bg-background/80 backdrop-blur-md hover:bg-border border border-border text-foreground transition-colors whitespace-nowrap"
            >
              <Cloud className="size-4 text-primary" />
              Cloud
            </a>
            <p className="text-primary text-[10px] sm:text-xs font-bold mt-2 text-center">
              1,200+{" "}
              <span className="text-muted-foreground font-normal block sm:inline">
                servers online
              </span>
            </p>
          </div>
          <div className="flex flex-col items-center">
            <a
              href="https://docs.codedock.run/self-host"
              target="_blank"
              rel="noreferrer"
              className="flex items-center justify-center gap-2 px-6 sm:px-10 py-3.5 rounded-lg font-semibold text-sm bg-background/80 backdrop-blur-md hover:bg-border border border-border text-foreground transition-colors whitespace-nowrap"
            >
              <Server className="size-4 text-primary" />
              Self-hosted
            </a>
            <p className="text-primary text-[10px] sm:text-xs font-bold mt-2 text-center">
              &lt;30MB{" "}
              <span className="text-muted-foreground font-normal block sm:inline">
                single binary
              </span>
            </p>
          </div>
        </div>

        <div className="mt-12 rounded-xl border border-border bg-background/80 backdrop-blur-md overflow-hidden text-left font-mono text-xs shadow-2xl relative z-10">
          <div className="flex items-center justify-between px-4 py-3 border-b border-border bg-background/40">
            <span className="text-muted-foreground text-[10px] font-bold uppercase tracking-widest">
              codedock daemon
            </span>
            <span className="flex items-center gap-1.5 text-green-400 text-[10px] font-semibold">
              <span className="size-1.5 rounded-full bg-green-400 animate-pulse" />
              CONNECTED
            </span>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 divide-y sm:divide-y-0 sm:divide-x divide-border">
            <div className="px-5 py-4 space-y-2.5">
              {serverRows.map((r) => (
                <div key={r.k} className="flex items-center gap-3 text-[0.71rem]">
                  <span className="text-muted-foreground w-20 shrink-0">{r.k}</span>
                  <span className={r.accent ? "text-primary font-semibold" : "text-foreground"}>
                    {r.v}
                  </span>
                </div>
              ))}
            </div>
            <div className="px-5 py-4">
              <div className="text-[9px] font-bold tracking-widest uppercase text-muted-foreground mb-3">
                Capabilities on this node
              </div>
              <div className="grid grid-cols-2 gap-y-2 gap-x-4">
                {capabilities.map((cap) => (
                  <div
                    key={cap}
                    className="flex items-center gap-2 text-[0.68rem] text-muted-foreground"
                  >
                    <span className="size-1.5 rounded-full bg-primary/70 shrink-0" />
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
