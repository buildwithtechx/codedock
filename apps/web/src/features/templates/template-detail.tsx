import { ArrowLeft, ArrowUpRight, BadgeCheck, Boxes } from 'lucide-react';
import { SiteLink } from '../../components/site-link';
import { productLinks } from '../../lib/product-links';
import { appCatalog } from './app-catalog.generated';

export function TemplateDetail({ templateId }: { templateId: string }) {
  const entry = appCatalog.find((item) => item.id === templateId);

  if (!entry) {
    return (
      <div className="mx-auto max-w-6xl px-6 py-16 sm:py-24">
        <h1 className="font-extrabold text-4xl tracking-tight">Template not found</h1>
        <p className="mt-4 max-w-3xl text-lg text-muted-foreground leading-relaxed">
          No catalog entry matches “{templateId}”.
        </p>
        <SiteLink
          href="/templates"
          className="mt-8 inline-flex items-center gap-2 font-semibold text-primary text-sm"
        >
          <ArrowLeft className="size-3.5" />
          Back to templates
        </SiteLink>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-6xl px-6 py-16 sm:py-24">
      <SiteLink
        href="/templates"
        className="inline-flex items-center gap-2 font-semibold text-primary text-xs"
      >
        <ArrowLeft className="size-3.5" />
        All templates
      </SiteLink>
      <div className="mt-6 flex items-start gap-5">
        <span className="flex size-14 items-center justify-center rounded-2xl border border-primary/20 bg-primary/10">
          <Boxes className="size-6 text-primary" />
        </span>
        <div>
          <div className="flex flex-wrap items-center gap-2">
            <p className="font-semibold text-primary text-xs uppercase tracking-widest">
              {entry.category}
            </p>
            {entry.verified && (
              <span className="inline-flex items-center gap-1 rounded-full border border-border px-2 py-0.5 text-[10px] text-muted-foreground">
                <BadgeCheck className="size-3 text-primary" />
                Verified
              </span>
            )}
          </div>
          <h1 className="mt-2 font-extrabold text-4xl tracking-tight sm:text-5xl">{entry.name}</h1>
          <p className="mt-4 max-w-3xl text-lg text-muted-foreground leading-relaxed">
            {entry.description}
          </p>
        </div>
      </div>

      <div className="mt-10 grid gap-5 lg:grid-cols-2">
        <section className="rounded-2xl border border-border bg-card/70 p-6">
          <h2 className="font-bold text-xl">Services in this stack</h2>
          <ul className="mt-4 flex flex-wrap gap-2">
            {entry.services.map((service) => (
              <li key={service} className="rounded-lg bg-muted px-3 py-1.5 font-mono text-xs">
                {service}
              </li>
            ))}
          </ul>
        </section>
        <section className="rounded-2xl border border-border bg-card/70 p-6">
          <h2 className="font-bold text-xl">Secrets you provide</h2>
          {entry.secrets.length === 0 ? (
            <p className="mt-4 text-muted-foreground text-sm">
              No secrets required — deploys with defaults.
            </p>
          ) : (
            <ul className="mt-4 divide-y divide-border/40">
              {entry.secrets.map((secret) => (
                <li key={secret.var} className="flex items-center justify-between gap-3 py-2.5">
                  <div>
                    <p className="font-mono text-sm">{secret.var}</p>
                    <p className="text-muted-foreground text-xs">{secret.label}</p>
                  </div>
                  {secret.length > 0 && (
                    <span className="shrink-0 text-muted-foreground text-xs">
                      {secret.length} chars
                    </span>
                  )}
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>

      <div className="mt-10 flex flex-wrap gap-3">
        <SiteLink
          href={`${productLinks.docs}/deployments/templates/#${entry.id}`}
          className="inline-flex items-center gap-2 rounded-xl bg-primary px-5 py-3 font-semibold text-sm text-white"
        >
          Open setup guide
          <ArrowUpRight className="size-4" />
        </SiteLink>
        <SiteLink
          href="/templates"
          className="inline-flex items-center gap-2 rounded-xl border border-border px-5 py-3 font-semibold text-sm"
        >
          Browse all templates
        </SiteLink>
      </div>
    </div>
  );
}
