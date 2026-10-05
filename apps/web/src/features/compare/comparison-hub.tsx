import { ArrowUpRight } from 'lucide-react';
import { SiteLink } from '../../components/site-link';
import { comparisons } from './comparisons';

const comparisonIds = [
  'coolify',
  'dokploy',
  'portainer',
  'caprover',
  'dokku',
  'vercel',
  'render',
] as const satisfies readonly (keyof typeof comparisons)[];

export function ComparisonHub() {
  return (
    <div className="mx-auto max-w-6xl px-6 py-16 sm:py-24">
      <p className="font-semibold text-primary text-xs uppercase tracking-widest">Find your fit</p>
      <h1 className="mt-5 max-w-4xl text-balance font-extrabold text-4xl tracking-tight sm:text-6xl">
        A clearer choice for your next deployment.
      </h1>
      <p className="mt-6 max-w-3xl text-lg text-muted-foreground leading-relaxed">
        Managed hosting, self-hosted PaaS, and container administration solve different problems.
        Compare the workflow, ownership, and operational responsibility that matter to you.
      </p>
      <div className="mt-12 grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
        {comparisonIds.map((id) => (
          <SiteLink
            key={id}
            href={`/vs/${id}`}
            className="group rounded-2xl border border-border bg-card/70 p-7 transition-colors hover:border-primary/50 hover:bg-primary/5"
          >
            <div className="mb-7 flex items-center justify-between">
              <span className="rounded-lg border border-border px-3 py-1 font-mono text-muted-foreground text-xs">
                vs. {comparisons[id].name}
              </span>
              <ArrowUpRight className="size-5 text-muted-foreground group-hover:text-primary" />
            </div>
            <h2 className="font-bold text-xl">Codedock vs. {comparisons[id].name}</h2>
            <p className="mt-3 text-muted-foreground text-sm leading-relaxed">
              {comparisons[id].tagline}
            </p>
            <p className="mt-6 font-semibold text-primary text-xs">Explore the comparison →</p>
          </SiteLink>
        ))}
      </div>
    </div>
  );
}
