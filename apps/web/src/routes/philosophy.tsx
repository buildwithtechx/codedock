import { createFileRoute } from '@tanstack/react-router';
import { CheckCircle, Cloud } from 'lucide-react';
import { GithubIcon } from '../components/icons/github-icon';
import { createMeta } from '../lib/seo';

export const Route = createFileRoute('/philosophy')({
  component: PhilosophyComponent,
  head: () => ({
    meta: createMeta({
      title: 'Our Philosophy — Codedock',
      description:
        'Why we built Codedock as MIT-licensed infrastructure software, and how we plan to keep it that way forever.',
    }),
  }),
});

function PhilosophyComponent() {
  const commitments = [
    {
      title: 'Every feature ships in the open-source daemon.',
      body: "Traefik routing, S3 backups, canvas view, AI settings, cron jobs — all of it is in the public repo. We don't maintain a secret pro tier that does more. If we build it, you get it.",
    },
    {
      title: 'MIT license, no exceptions.',
      body: "We chose MIT deliberately. Apache 2.0 has patent clauses that make commercial redistribution complicated. MIT is the plainest, broadest grant we can give. If you want to wrap Codedock in a product and sell it, you can. We'd appreciate a mention, not a cheque.",
    },
    {
      title: 'No paywalled stability.',
      body: "We don't release a buggy community edition while shipping a polished enterprise build. The daemon you download is the only daemon we maintain. Bug fixes go in; everyone gets them.",
    },
    {
      title: 'These rules are public and permanent.',
      body: "This page is version-controlled. If our philosophy ever changes, the git history will show it. We don't expect it to — but the audit trail is there if you want it.",
    },
  ];

  const stats = [
    { label: 'Cloud MRR', value: '$4,200+', href: 'https://app.codedock.run' },
    { label: 'GitHub Stars', value: '1,800+', href: 'https://github.com/techxteam/codedock' },
    { label: 'Active Servers', value: '1,200+', href: 'https://app.codedock.run' },
  ];

  return (
    <div className="mx-auto max-w-3xl px-6 text-center">
      <div className="pt-16 pb-6 md:pt-24">
        <h1 className="text-balance font-extrabold text-3xl text-foreground leading-tight tracking-tight md:text-4xl">
          Infrastructure that belongs to you<span className="text-primary">.</span>
        </h1>
        <p className="mx-auto mt-5 max-w-2xl text-base text-muted-foreground leading-relaxed">
          We built Codedock because we kept seeing the same pattern: a tool starts free, gains
          traction, then slowly walls off its best features behind enterprise plans. The moment a
          self-hosting tool needs a licence key to use SSL, something has gone wrong.
        </p>
        <p className="mx-auto mt-4 max-w-2xl text-base text-muted-foreground leading-relaxed">
          Codedock is MIT-licensed. Not "source-available". Not "fair-code". MIT. You can read it,
          fork it, patch it, and run it commercially without asking us. That constraint is
          intentional — it means we can't quietly change our minds later.
        </p>
      </div>

      <div className="pt-14">
        <h2 className="mb-2 font-bold text-2xl text-foreground">
          What we commit to<span className="text-primary">.</span>
        </h2>
        <p className="mb-8 text-muted-foreground text-sm">
          Not just promises — these are structural constraints.
        </p>

        <ul className="space-y-7 text-left">
          {commitments.map((p) => (
            <li key={p.title} className="flex gap-4">
              <CheckCircle className="mt-1 size-5 shrink-0 text-primary" />
              <div>
                <strong className="font-semibold text-foreground">{p.title}</strong>
                <p className="mt-1 text-muted-foreground text-sm leading-relaxed">{p.body}</p>
              </div>
            </li>
          ))}
        </ul>
      </div>

      <div className="pt-16">
        <h2 className="mb-2 font-bold text-2xl text-foreground">
          How we keep the lights on<span className="text-primary">.</span>
        </h2>
        <p className="mb-6 text-muted-foreground text-sm">
          Free software still has servers to pay for.
        </p>

        <div className="space-y-4 text-left text-muted-foreground text-sm leading-relaxed">
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

      <div className="flex flex-row flex-wrap items-center justify-center gap-3 pt-10 pb-4 sm:gap-4">
        <a
          href="https://github.com/sponsors/techxteam"
          target="_blank"
          rel="noreferrer"
          className="flex items-center justify-center gap-2.5 whitespace-nowrap rounded-lg border border-border bg-background/80 px-6 py-3.5 font-semibold text-foreground text-sm backdrop-blur-md transition-colors hover:bg-border sm:px-10"
        >
          <GithubIcon className="size-5 shrink-0" />
          <span>GitHub Sponsors</span>
        </a>

        <a
          href="https://app.codedock.run"
          className="flex items-center justify-center gap-2.5 whitespace-nowrap rounded-lg border border-primary/30 bg-primary/10 px-6 py-3.5 font-semibold text-foreground text-sm transition-colors hover:bg-primary/20 sm:px-10"
        >
          <Cloud className="size-5 shrink-0 text-primary" />
          <span>Try Codedock Cloud</span>
        </a>
      </div>

      <div className="pt-16 pb-4">
        <h2 className="mb-2 font-bold text-2xl text-foreground">
          Where we stand today<span className="text-primary">.</span>
        </h2>
        <p className="mb-8 text-muted-foreground text-sm">
          We publish these numbers because transparency is part of the deal.
        </p>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
          {stats.map((stat) => (
            <a
              key={stat.label}
              href={stat.href}
              target="_blank"
              rel="noreferrer"
              className="group rounded-xl border border-border bg-background/80 p-6 text-center backdrop-blur-md transition-colors hover:border-primary/40"
            >
              <div className="font-extrabold text-3xl text-foreground transition-colors group-hover:text-primary">
                {stat.value}
              </div>
              <div className="mt-1 text-muted-foreground text-sm">{stat.label}</div>
            </a>
          ))}
        </div>
      </div>

      <div className="space-y-4 pt-10 pb-20 text-muted-foreground text-sm leading-relaxed">
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
