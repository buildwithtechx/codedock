import { ArrowRight, ExternalLink } from 'lucide-react';
import { productLinks } from '../../lib/product-links';
import type { ComparisonContent } from './comparison-content';

export function ComparisonPage({ content }: { content: ComparisonContent }) {
  return (
    <div className="mx-auto max-w-6xl px-6 py-16 sm:py-24">
      <a href="/vs" className="text-muted-foreground text-sm hover:text-primary">
        Back to all comparisons
      </a>
      <p className="mt-9 font-semibold text-primary text-xs uppercase tracking-widest">
        Codedock vs. {content.name}
      </p>
      <h1 className="mt-5 max-w-4xl text-balance font-extrabold text-4xl leading-tight tracking-tight sm:text-6xl">
        {content.tagline}
      </h1>
      <p className="mt-6 max-w-3xl text-lg text-muted-foreground leading-relaxed">
        {content.description}
      </p>
      <div className="mt-12 overflow-x-auto rounded-2xl border border-border bg-card/70">
        <table className="w-full min-w-[640px] border-collapse text-left text-sm">
          <caption className="sr-only">
            Codedock and {content.name}: deployment and operations comparison
          </caption>
          <thead>
            <tr className="border-border border-b">
              <th scope="col" className="w-1/4 p-5 font-medium text-muted-foreground">
                What matters
              </th>
              <th scope="col" className="w-[37.5%] bg-primary/5 p-5 font-bold text-primary">
                Codedock
              </th>
              <th scope="col" className="p-5 font-bold">
                {content.name}
              </th>
            </tr>
          </thead>
          <tbody>
            {content.rows.map((row) => (
              <tr key={row.criterion} className="border-border border-b last:border-0">
                <th scope="row" className="p-5 font-medium">
                  {row.criterion}
                </th>
                <td className="bg-primary/5 p-5 text-muted-foreground leading-relaxed">
                  {row.codedock}
                </td>
                <td className="p-5 text-muted-foreground leading-relaxed">{row.alternative}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="mt-8 grid gap-5 md:grid-cols-2">
        <section className="rounded-2xl border border-primary/25 bg-primary/5 p-7">
          <h2 className="font-bold text-xl">When Codedock fits</h2>
          <p className="mt-4 text-muted-foreground text-sm leading-relaxed">
            {content.codedockFit}
          </p>
          <a
            href={productLinks.installation}
            className="mt-6 inline-flex items-center gap-2 font-semibold text-primary text-sm"
          >
            Start self-hosting
            <ArrowRight className="size-4" />
          </a>
        </section>
        <section className="rounded-2xl border border-border bg-card/70 p-7">
          <h2 className="font-bold text-xl">When {content.name} fits</h2>
          <p className="mt-4 text-muted-foreground text-sm leading-relaxed">
            {content.alternativeFit}
          </p>
        </section>
      </div>
      <section className="mt-10 border-border border-t pt-7">
        <h2 className="font-semibold text-sm">Sources & scope</h2>
        <p className="mt-2 max-w-3xl text-muted-foreground text-xs leading-relaxed">
          Reviewed against official documentation on October 5, 2026. This compares the workflows
          described here, not every feature or plan. Check the vendor's current documentation before
          choosing a platform. Self-hosting still has hosting, bandwidth, storage, and maintenance
          costs.
        </p>
        <ul className="mt-4 flex flex-wrap gap-4">
          {content.sources.map((source) => (
            <li key={source.href}>
              <a
                href={source.href}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center gap-1 text-primary text-xs hover:underline"
              >
                {source.title}
                <ExternalLink className="size-3" />
              </a>
            </li>
          ))}
          <li>
            <a href={productLinks.docs} className="text-primary text-xs hover:underline">
              Codedock documentation
            </a>
          </li>
        </ul>
      </section>
    </div>
  );
}
