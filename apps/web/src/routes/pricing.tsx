import { createFileRoute } from '@tanstack/react-router';
import { InstallCommand } from '../components/install-command';
import { productLinks } from '../lib/product-links';
import { pageHead } from '../lib/seo';
export const Route = createFileRoute('/pricing')({
  component: PricingPage,
  head: () =>
    pageHead(
      '/pricing',
      'Pricing',
      'Run the Apache-2.0 licensed Codedock software on your own servers, or visit Codedock Cloud for current plan availability and pricing.'
    ),
});
function PricingPage() {
  return (
    <section className="mx-auto max-w-5xl px-6 py-24">
      <h1 className="font-bold text-4xl md:text-6xl">Your infrastructure. Your choice.</h1>
      <p className="mt-5 max-w-2xl text-lg text-muted-foreground">
        Choose where you run Codedock. Your hosting and external service costs are separate from the
        software.
      </p>
      <div className="mt-12 grid gap-6 md:grid-cols-2">
        <article className="rounded-2xl border border-border bg-card p-8">
          <h2 className="font-semibold text-2xl">Self-hosted</h2>
          <p className="mt-5 font-bold text-4xl">
            $0{' '}
            <span className="font-normal text-muted-foreground text-sm">software licence cost</span>
          </p>
          <p className="mt-5 text-muted-foreground">
            Run the Apache-2.0 licensed code on your Linux servers. You maintain the host, Docker,
            networking, updates and backups.
          </p>
          <ul className="my-6 list-inside list-disc space-y-2 text-sm">
            <li>Applications and database workflows</li>
            <li>SSH worker servers and team access</li>
            <li>Logs, metrics and backup configuration</li>
            <li>No Stripe or email provider required for initial setup</li>
          </ul>
          <InstallCommand />
          <a href={productLinks.installation} className="mt-6 inline-block text-primary">
            Installation guide →
          </a>
        </article>
        <article className="rounded-2xl border border-primary/40 bg-primary/5 p-8">
          <h2 className="font-semibold text-2xl">Codedock Cloud</h2>
          <p className="mt-5 font-bold text-3xl">Check current plans</p>
          <p className="mt-5 text-muted-foreground">
            Cloud account plans and availability are shown in the hosted application. Review the
            current offer, limits and terms before subscribing.
          </p>
          <p className="mt-5 text-muted-foreground text-sm">
            Cloud billing and account configuration are separate from self-hosted installation.
            Connecting a worker server still requires SSH access and compatible infrastructure.
          </p>
          <a
            href={productLinks.cloud}
            className="mt-8 inline-block rounded-lg bg-primary px-6 py-3 font-semibold text-white"
          >
            Open Codedock Cloud →
          </a>
        </article>
      </div>
      <p className="mt-8 text-muted-foreground text-sm">
        Budget separately for servers, domains, object storage and any AI or email providers you
        enable. Self-hosting gives you control and responsibility for operating the stack.
      </p>
    </section>
  );
}
