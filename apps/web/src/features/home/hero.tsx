import { ArrowRight, Cloud } from 'lucide-react';
import { DiscordIcon } from '../../components/icons/discord-icon';
import { GithubIcon } from '../../components/icons/github-icon';
import { InstallCommand } from '../../components/install-command';
import { productLinks } from '../../lib/product-links';
import { ControlPlanePreview } from './control-plane-preview';

export function Hero() {
  return (
    <section className="relative isolate overflow-hidden px-6 pt-18 pb-20 lg:pt-24 lg:pb-28">
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 -z-10 bg-[radial-gradient(ellipse_at_70%_20%,rgba(124,58,237,0.18),transparent_60%)]"
      />
      <div className="mx-auto grid max-w-7xl items-center gap-14 lg:grid-cols-[1.05fr_1fr] lg:gap-16">
        <div>
          <p className="mb-6 inline-flex items-center gap-2 rounded-full border border-primary/25 bg-primary/5 px-3 py-1.5 font-medium text-primary text-xs">
            <span className="size-1.5 rounded-full bg-primary" />
            Open source. Built for your servers.
          </p>
          <h1 className="max-w-2xl text-balance font-extrabold text-5xl leading-[1.08] tracking-[-0.045em] sm:text-6xl xl:text-7xl">
            Ship fast.
            <br />
            Own your <span className="text-primary">infrastructure.</span>
          </h1>
          <p className="mt-6 max-w-xl text-pretty text-base text-muted-foreground leading-relaxed sm:text-lg">
            Deploy applications, databases, and background workers on your own Linux server. One
            workspace for your builds, data, logs, and backups.
          </p>
          <div className="mt-8 flex flex-wrap items-center gap-3">
            <a
              href={productLinks.installation}
              className="inline-flex items-center gap-2 rounded-xl bg-primary px-5 py-3 font-semibold text-sm text-white shadow-lg shadow-primary/20 hover:bg-primary/90"
            >
              Start self-hosting
              <ArrowRight className="size-4" />
            </a>
            <a
              href={productLinks.cloud}
              className="inline-flex items-center gap-2 rounded-xl border border-border px-5 py-3 font-semibold text-sm hover:bg-muted"
            >
              <Cloud className="size-4" />
              Explore Cloud
            </a>
          </div>
          <div className="mt-7">
            <InstallCommand />
          </div>
          <div className="mt-4 flex flex-wrap items-center gap-5 text-muted-foreground text-xs">
            <a href={productLinks.github} className="flex items-center gap-2 hover:text-foreground">
              <GithubIcon className="size-4" />
              View the source
            </a>
            <a
              href={productLinks.discord}
              className="flex items-center gap-2 hover:text-foreground"
            >
              <DiscordIcon className="size-4" />
              Join the community
            </a>
            <span>Apache-2.0 licence</span>
          </div>
        </div>
        <ControlPlanePreview />
      </div>
    </section>
  );
}
