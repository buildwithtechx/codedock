import { createFileRoute } from '@tanstack/react-router';
import { productLinks } from '../lib/product-links';
import { pageHead } from '../lib/seo';
export const Route = createFileRoute('/changelog')({
  component: ChangelogPage,
  head: () =>
    pageHead(
      '/changelog',
      'Development updates',
      'Follow Codedock development, release notes and source changes in the public repository.'
    ),
});
function ChangelogPage() {
  return (
    <section className="mx-auto max-w-3xl px-6 py-24">
      <p className="font-semibold text-primary text-sm">DEVELOPMENT UPDATES</p>
      <h1 className="mt-4 font-bold text-4xl">Follow what is changing.</h1>
      <p className="mt-6 text-lg text-muted-foreground">
        Published releases and their notes live in the GitHub repository. The highlights below
        describe development changes, rather than a versioned release.
      </p>
      <article className="mt-10 rounded-2xl border border-border bg-card p-8">
        <h2 className="font-semibold text-2xl">Self-hosted setup and configuration</h2>
        <ul className="mt-5 list-inside list-disc space-y-3 text-muted-foreground">
          <li>Generate and persist the initial self-hosted secrets automatically.</li>
          <li>Create the first owner account through browser onboarding.</li>
          <li>Separate optional cloud billing, email and OAuth integrations from local setup.</li>
          <li>Improve deployment persistence and configuration validation.</li>
        </ul>
        <a href={`${productLinks.github}/pull/25`} className="mt-6 inline-block text-primary">
          Read the merged changes →
        </a>
      </article>
      <div className="mt-8 flex flex-wrap gap-6 text-primary">
        <a href={`${productLinks.github}/releases`}>Published releases ?</a>
        <a href={`${productLinks.github}/commits/main/`}>Development history ?</a>
      </div>
    </section>
  );
}
