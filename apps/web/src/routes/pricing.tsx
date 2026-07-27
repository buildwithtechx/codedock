import { createFileRoute } from "@tanstack/react-router";
import { ArrowRight, Book, Check, ChevronDown, Cloud, Info, Server } from "lucide-react";
import { createMeta } from "../lib/seo";

export const Route = createFileRoute("/pricing")({
  component: PricingComponent,
  head: () => ({
    meta: createMeta({
      title: "Pricing — Codedock",
      description:
        "Self-hosted is MIT-licensed and free forever. Codedock Cloud adds fleet management and metering from $5/month.",
    }),
  }),
});

function PricingComponent() {
  const tableRows = [
    { feature: "App Deployments", sh: true, cloud: true },
    { feature: "Database Deployments", sh: true, cloud: true },
    { feature: "Cron Jobs", sh: true, cloud: true },
    { feature: "Canvas View", sh: true, cloud: true },
    { feature: "Backup to S3", sh: true, cloud: true },
    { feature: "Traefik Routing + SSL", sh: true, cloud: true },
    { feature: "Scoped Environment Secrets", sh: true, cloud: true },
    { feature: "AI Settings (BYO key)", sh: true, cloud: true },
    { feature: "Git-native Deploys", sh: true, cloud: true },
    { feature: "Log Tailing", sh: true, cloud: true },
    { feature: "Notification Channels", sh: true, cloud: true },
    { feature: "Codebase", sh: "MIT open source", cloud: "MIT open source" },
    { feature: "Daemon Hosting", sh: "Your server", cloud: "Your server" },
    { feature: "Control Plane", sh: "Local dashboard", cloud: "Managed (app.codedock.run)" },
    { feature: "Fleet Management", sh: "—", cloud: true },
    { feature: "Usage Metering + Alerts", sh: "—", cloud: true },
    { feature: "Yamux Tunnel Handshake", sh: "Manual token", cloud: "Managed" },
    { feature: "Team Seats", sh: "Unlimited", cloud: "Unlimited" },
    { feature: "Upcoming Features", sh: "Free forever", cloud: "Included" },
  ];

  const faqs = [
    {
      q: "Do Cloud features replace the self-hosted daemon?",
      a: "No. Codedock Cloud is a control plane — your apps still run on your server via the Codedock daemon. Cloud adds fleet management, usage metering, and team billing on top. The daemon itself is unchanged.",
    },
    {
      q: "What exactly runs on my server vs. Codedock's infrastructure?",
      a: "Everything that matters — Docker containers, Traefik proxy, databases, backups — runs on your machine. Codedock Cloud only stores metadata: deployment history, team config, and usage metrics. We never touch your data.",
    },
    {
      q: "Is the self-hosted edition actually free forever, or will you change your mind?",
      a: "We built self-hosted on MIT. Revoking that would break the license itself. It's free, period — read our philosophy page for the full commitment.",
    },
    {
      q: "Can I move from Cloud back to pure self-hosted?",
      a: "Yes. It's a config change — point your daemon at your own dashboard instead of app.codedock.run. No vendor lock-in, no export-import dance.",
    },
    {
      q: "How does billing work when I add more servers?",
      a: "The base plan covers 2 servers. Each additional server is $3/month. There are no per-deployment or per-seat charges on top of that.",
    },
    {
      q: "Is there a trial?",
      a: "Yes — 14 days, no credit card required. Connect your first server and explore the full Cloud dashboard before committing.",
    },
  ];

  return (
    <div className="max-w-5xl mx-auto px-6">
      <div className="text-center pt-16 md:pt-24 pb-12">
        <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight text-foreground text-balance leading-tight">
          Simple pricing<span className="text-primary">.</span>
        </h1>
        <p className="mt-3 text-muted-foreground text-base max-w-lg mx-auto">
          The daemon is MIT-licensed and free forever. Pay only when you want the fleet control
          plane and metering on top.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-background/80 backdrop-blur-md p-7 rounded-xl border border-border flex flex-col">
          <div className="mb-6">
            <div className="flex items-center gap-2 mb-4">
              <Server className="size-5 text-primary" />
              <h2 className="text-xl font-bold text-foreground">Self-hosted</h2>
            </div>
            <p className="text-3xl font-extrabold text-foreground">
              Free <span className="text-sm font-normal text-muted-foreground ml-1">forever</span>
            </p>
            <p className="text-sm text-muted-foreground mt-3 leading-relaxed">
              One curl command. The daemon bootstraps Traefik, sets up your network, and is ready to
              deploy in seconds. All features included, no expiry.
            </p>
          </div>

          <ul className="space-y-3 text-sm flex-1">
            {[
              "Full daemon — all features, forever",
              "Traefik proxy + auto SSL",
              "Postgres, MySQL, Redis, MongoDB",
              "S3 backup scheduling",
              "Canvas topology view",
              "Cron jobs + env secrets",
              "AI settings per team",
            ].map((f) => (
              <li key={f} className="flex items-start gap-3 text-foreground">
                <Check className="size-4 text-primary shrink-0 mt-0.5" />
                {f}
              </li>
            ))}
          </ul>

          <a
            href="https://docs.codedock.run/self-host"
            target="_blank"
            rel="noreferrer"
            className="mt-8 flex items-center justify-center gap-2 py-3.5 rounded-lg border border-border bg-background hover:bg-border text-foreground font-semibold text-sm transition-colors"
          >
            <Book className="size-4 text-primary" />
            Read the Docs
          </a>
        </div>

        <div className="bg-background/80 backdrop-blur-md p-7 rounded-xl border border-primary/30 flex flex-col relative overflow-hidden">
          <div
            className="pointer-events-none absolute inset-0 [background:radial-gradient(ellipse_80%_60%_at_50%_-10%,rgba(124,58,237,0.07),transparent)]"
            aria-hidden="true"
          />

          <div className="relative mb-6">
            <div className="flex items-center justify-between gap-3 flex-wrap mb-4">
              <div className="flex items-center gap-2">
                <Cloud className="size-5 text-primary" />
                <h2 className="text-xl font-bold text-foreground">Cloud</h2>
              </div>
              <div className="inline-flex gap-1 p-1 bg-background rounded text-xs">
                <span className="px-3 py-1 bg-primary text-white rounded font-semibold">
                  Monthly
                </span>
                <span className="px-3 py-1 text-muted-foreground font-semibold">
                  Annually <span className="text-primary">−20%</span>
                </span>
              </div>
            </div>
            <p className="text-3xl font-extrabold text-foreground">
              $5{" "}
              <span className="text-sm font-normal text-muted-foreground ml-1">
                /month · 2 servers included
              </span>
            </p>
            <p className="text-sm text-primary font-semibold mt-1">+ $3 /month per extra server</p>
            <p className="text-sm text-muted-foreground mt-3 leading-relaxed">
              All the self-hosted goodness, plus a managed fleet dashboard. Connect servers via
              Yamux tunnel — no ports to open, no VPN.
            </p>
          </div>

          <ul className="space-y-3 text-sm relative flex-1">
            {[
              "Everything in Self-hosted",
              "Fleet dashboard (Yamux tunnel)",
              "Usage metering — container-hours, deploys, bandwidth",
              "Billing alerts before you overspend",
              "Managed email notifications",
              "Team billing + per-server limits",
              "14-day free trial · no card required",
            ].map((f) => (
              <li key={f} className="flex items-start gap-3 text-foreground">
                <Check className="size-4 text-primary shrink-0 mt-0.5" />
                {f}
              </li>
            ))}
          </ul>

          <div className="flex items-start gap-2 mt-4 p-3 rounded-lg bg-background border border-border text-xs text-muted-foreground relative">
            <Info className="size-3.5 text-primary shrink-0 mt-0.5" />
            Your apps run on <em>your</em> servers. Codedock Cloud only manages the control plane.
          </div>

          <a
            href="https://app.codedock.run/register"
            className="mt-6 flex items-center justify-center gap-2 py-3.5 rounded-lg bg-primary hover:bg-primary-hover text-white font-semibold text-sm transition-colors relative"
          >
            Start 14-day trial
            <ArrowRight className="size-4" />
          </a>
        </div>
      </div>

      <div className="mt-16 overflow-x-auto rounded-xl border border-border bg-background/80 backdrop-blur-md">
        <table className="w-full text-left text-sm">
          <thead>
            <tr className="bg-surface border-b border-border">
              <th className="px-5 py-4 font-semibold text-foreground">Feature</th>
              <th className="px-5 py-4 font-semibold text-foreground min-w-37.5">Self-hosted</th>
              <th className="px-5 py-4 font-semibold text-foreground min-w-25">Cloud</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {tableRows.map((row) => (
              <tr key={row.feature} className="hover:bg-surface/60 transition-colors">
                <td className="px-5 py-3.5 text-foreground">{row.feature}</td>
                <td className="px-5 py-3.5">
                  {row.sh === true ? (
                    <Check className="size-4 text-primary" />
                  ) : (
                    <span className="text-muted-foreground text-xs">{row.sh}</span>
                  )}
                </td>
                <td className="px-5 py-3.5">
                  {row.cloud === true ? (
                    <Check className="size-4 text-primary" />
                  ) : (
                    <span className="text-muted-foreground text-xs">{row.cloud}</span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="mt-16 max-w-3xl mx-auto pb-16">
        <div className="bg-background/80 backdrop-blur-md border border-border rounded-xl p-6 md:p-8">
          <h2 className="text-2xl font-bold text-foreground mb-6">Common questions</h2>
          <div className="divide-y divide-border border-t border-border">
            {faqs.map((faq) => (
              <details key={faq.q} className="py-4 group cursor-pointer">
                <summary className="w-full text-left text-sm font-semibold text-foreground flex justify-between items-center list-none gap-4">
                  {faq.q}
                  <ChevronDown className="size-4 text-muted-foreground group-open:rotate-180 transition-transform duration-200 shrink-0" />
                </summary>
                <p className="mt-3 text-sm text-muted-foreground leading-relaxed pr-8">{faq.a}</p>
              </details>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
