import { ArrowRight, Check, Layers } from 'lucide-react';
import { InstallCommand } from '../../components/install-command';
import { SiteLink } from '../../components/site-link';
import { productLinks } from '../../lib/product-links';
import type { PageContent } from './page-content';

export function DetailPage({ content }: { content: PageContent }) {
  return (
    <div className="mx-auto max-w-6xl px-6 py-16 sm:py-24">
      <p className="mb-5 font-semibold text-primary text-xs uppercase tracking-[0.2em]">
        {content.eyebrow}
      </p>
      <h1 className="max-w-4xl text-balance font-extrabold text-4xl leading-tight tracking-tight sm:text-6xl">
        {content.title}
      </h1>
      <p className="mt-6 max-w-3xl text-lg text-muted-foreground leading-relaxed">
        {content.description}
      </p>
      <div className="mt-8 flex flex-wrap gap-3">
        {content.highlights.map((item) => (
          <span
            key={item}
            className="flex items-center gap-2 rounded-full border border-border bg-card/60 px-4 py-2 text-xs"
          >
            <Check className="size-3 text-primary" />
            {item}
          </span>
        ))}
      </div>
      <div className="mt-9 flex flex-wrap items-center gap-4">
        <SiteLink
          href={productLinks.installation}
          className="rounded-xl bg-primary px-5 py-3 font-semibold text-sm text-white hover:bg-primary/90"
        >
          Get started
        </SiteLink>
        <SiteLink
          href={`${productLinks.docs}${content.docsPath}`}
          className="flex items-center gap-2 text-sm hover:text-primary"
        >
          {content.docsLabel}
          <ArrowRight className="size-4" />
        </SiteLink>
      </div>
      <div className="mt-16 grid gap-5 lg:grid-cols-3">
        {content.sections.map((section, index) => (
          <section key={section.title} className="rounded-2xl border border-border bg-card/70 p-6">
            <p className="mb-6 font-mono text-primary text-xs">
              0{index + 1} / <Layers className="ml-1 inline size-3.5" />
            </p>
            <h2 className="font-bold text-xl">{section.title}</h2>
            <p className="mt-4 text-muted-foreground text-sm leading-relaxed">{section.body}</p>
            <ul className="mt-5 space-y-3">
              {section.points.map((point) => (
                <li key={point} className="flex gap-2 text-sm leading-relaxed">
                  <Check className="mt-1 size-3.5 shrink-0 text-primary" />
                  {point}
                </li>
              ))}
            </ul>
          </section>
        ))}
      </div>
      <section className="mt-14 grid gap-8 rounded-2xl border border-border bg-primary/5 p-6 sm:p-9 lg:grid-cols-2">
        <div>
          <p className="font-semibold text-primary text-xs uppercase tracking-widest">
            Make it concrete
          </p>
          <h2 className="mt-3 font-bold text-2xl">Your first workflow</h2>
          <ol className="mt-6 space-y-4">
            {content.steps.map((step, index) => (
              <li key={step} className="flex gap-3 text-muted-foreground text-sm">
                <span className="flex size-6 shrink-0 items-center justify-center rounded-full border border-primary/25 text-primary text-xs">
                  {index + 1}
                </span>
                <span className="pt-0.5">{step}</span>
              </li>
            ))}
          </ol>
        </div>
        <div className="flex flex-col justify-center">
          <h3 className="mb-4 font-semibold">Start on your server.</h3>
          <InstallCommand source="recipe_detail" />
          <SiteLink
            href={`${productLinks.docs}${content.docsPath}`}
            className="mt-4 flex items-center gap-2 font-semibold text-primary text-sm"
          >
            Open the complete guide
            <ArrowRight className="size-4" />
          </SiteLink>
        </div>
      </section>
    </div>
  );
}
