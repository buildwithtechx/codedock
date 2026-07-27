import { createFileRoute } from "@tanstack/react-router";
import { CheckCircle, Cloud } from "lucide-react";
import { GithubIcon } from "../components/github-icon";
import { createMeta } from "../lib/seo";

export const Route = createFileRoute("/philosophy")({
  component: PhilosophyComponent,
  head: () => ({
    meta: createMeta({
      title: "Our Philosophy — Codedock",
      description:
        "Why we built Codedock as MIT-licensed infrastructure software, and how we plan to keep it that way forever.",
    }),
  }),
});

function PhilosophyComponent() {
  const commitments = [
    {
      title: "Every feature ships in the open-source daemon.",
      body: "Traefik routing, S3 backups, canvas view, AI settings, cron jobs — all of it is in the public repo. We don't maintain a secret pro tier that does more. If we build it, you get it.",
    },
    {
      title: "MIT license, no exceptions.",
      body: "We chose MIT deliberately. Apache 2.0 has patent clauses that make commercial redistribution complicated. MIT is the plainest, broadest grant we can give. If you want to wrap Codedock in a product and sell it, you can. We'd appreciate a mention, not a cheque.",
    },
    {
      title: "No paywalled stability.",
      body: "We don't release a buggy community edition while shipping a polished enterprise build. The daemon you download is the only daemon we maintain. Bug fixes go in; everyone gets them.",
    },
    {
      title: "These rules are public and permanent.",
      body: "This page is version-controlled. If our philosophy ever changes, the git history will show it. We don't expect it to — but the audit trail is there if you want it.",
    },
  ];

  const stats = [
    { label: "Cloud MRR", value: "$4,200+", href: "https://app.codedock.run" },
    { label: "GitHub Stars", value: "1,800+", href: "https://github.com/techxteam/codedock" },
    { label: "Active Servers", value: "1,200+", href: "https://app.codedock.run" },
  ];

  return (
    <div className="max-w-3xl mx-auto px-6 text-center">
      <div className="pt-16 md:pt-24 pb-6">
        <h1 className="text-3xl md:text-4xl font-extrabold tracking-tight text-foreground text-balance leading-tight">
          Infrastructure that belongs to you<span className="text-primary">.</span>
        </h1>
        <p className="mt-5 text-muted-foreground text-base leading-relaxed max-w-2xl mx-auto">
          We built Codedock because we kept seeing the same pattern: a tool starts free, gains
          traction, then slowly walls off its best features behind enterprise plans. The moment a
          self-hosting tool needs a licence key to use SSL, something has gone wrong.
        </p>
        <p className="mt-4 text-muted-foreground text-base leading-relaxed max-w-2xl mx-auto">
          Codedock is MIT-licensed. Not "source-available". Not "fair-code". MIT. You can read it,
          fork it, patch it, and run it commercially without asking us. That constraint is
          intentional — it means we can't quietly change our minds later.
        </p>
      </div>

      <div className="pt-14">
        <h2 className="text-2xl font-bold text-foreground mb-2">
          What we commit to<span className="text-primary">.</span>
        </h2>
        <p className="text-sm text-muted-foreground mb-8">
          Not just promises — these are structural constraints.
        </p>

        <ul className="space-y-7 text-left">
          {commitments.map((p) => (
            <li key={p.title} className="flex gap-4">
              <CheckCircle className="size-5 text-primary shrink-0 mt-1" />
              <div>
                <strong className="text-foreground font-semibold">{p.title}</strong>
                <p className="mt-1 text-sm text-muted-foreground leading-relaxed">{p.body}</p>
              </div>
            </li>
          ))}
        </ul>
      </div>

      <div className="pt-16">
        <h2 className="text-2xl font-bold text-foreground mb-2">
          How we keep the lights on<span className="text-primary">.</span>
        </h2>
        <p className="text-sm text-muted-foreground mb-6">
          Free software still has servers to pay for.
        </p>

        <div className="text-muted-foreground text-sm leading-relaxed space-y-4 text-left">
          <p>
            <strong className="text-foreground">Codedock Cloud.</strong> The daemon is free. The
            fleet control plane — Yamux tunnel handshake, usage metering, billing alerts, team
            management — is what we charge for. You get everything the daemon offers; the Cloud plan
            is for teams who want centralized visibility across many servers.
          </p>
          <p>
            <strong className="text-foreground">GitHub Sponsors.</strong> If Codedock saves you
            money on Heroku or Railway, consider sponsoring. Individual contributions fund the
            engineering time that keeps the daemon healthy.
          </p>
          <p>
            We don't take VC money. External funding always comes with growth pressure, and growth
            pressure eventually cracks open-source commitments. We'd rather grow slower and stay
            independent.
          </p>
        </div>
      </div>

      <div className="pt-10 pb-4 flex flex-row flex-wrap items-center justify-center gap-3 sm:gap-4">
        <a
          href="https://github.com/sponsors/techxteam"
          target="_blank"
          rel="noreferrer"
          className="flex items-center justify-center gap-2.5 px-6 sm:px-10 py-3.5 rounded-lg border border-border bg-background/80 backdrop-blur-md hover:bg-border text-foreground font-semibold text-sm transition-colors whitespace-nowrap"
        >
          <GithubIcon className="size-5 shrink-0" />
          <span>GitHub Sponsors</span>
        </a>

        <a
          href="https://app.codedock.run"
          className="flex items-center justify-center gap-2.5 px-6 sm:px-10 py-3.5 rounded-lg border border-primary/30 bg-primary/10 hover:bg-primary/20 text-foreground font-semibold text-sm transition-colors whitespace-nowrap"
        >
          <Cloud className="size-5 text-primary shrink-0" />
          <span>Try Codedock Cloud</span>
        </a>
      </div>

      <div className="pt-16 pb-4">
        <h2 className="text-2xl font-bold text-foreground mb-2">
          Where we stand today<span className="text-primary">.</span>
        </h2>
        <p className="text-sm text-muted-foreground mb-8">
          We publish these numbers because transparency is part of the deal.
        </p>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          {stats.map((stat) => (
            <a
              key={stat.label}
              href={stat.href}
              target="_blank"
              rel="noreferrer"
              className="bg-background/80 backdrop-blur-md rounded-xl p-6 border border-border hover:border-primary/40 transition-colors group text-center"
            >
              <div className="text-3xl font-extrabold text-foreground group-hover:text-primary transition-colors">
                {stat.value}
              </div>
              <div className="text-sm text-muted-foreground mt-1">{stat.label}</div>
            </a>
          ))}
        </div>
      </div>

      <div className="pt-10 pb-20 text-sm text-muted-foreground leading-relaxed space-y-4">
        <p>
          This philosophy is not separate from the engineering. The reason the daemon is small
          (&lt;30MB) is that we refuse to bundle bloat to justify a premium tier. The reason every
          feature is in the open-source build is that we don't want two codebases to maintain. The
          choices that make Codedock good to use are the same choices that make it genuinely open.
        </p>
        <p>Thank you for using it — and for holding us to this.</p>
      </div>
    </div>
  );
}
