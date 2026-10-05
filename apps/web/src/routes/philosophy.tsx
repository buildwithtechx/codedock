import { createFileRoute } from '@tanstack/react-router';
import { productLinks } from '../lib/product-links';
import { pageHead } from '../lib/seo';

const principles = [
  {
    title: 'Open code, a clear licence',
    body: 'Codedock is published under Apache-2.0. Read the source and licence, modify the software and run it on your own infrastructure under those terms.',
  },
  {
    title: 'Familiar infrastructure',
    body: 'The Go daemon manages Docker workloads and Traefik routing. Databases, volumes and backups use tools you can inspect and operate.',
  },
  {
    title: 'Self-hosting should start simply',
    body: 'The installer configures initial settings and generates local secrets. Create the first owner account in the browser, then add only the integrations you need.',
  },
  {
    title: 'Decisions you can inspect',
    body: 'The repository, documentation and development history are public. Capabilities and limits should be explained with the code that implements them.',
  },
];
export const Route = createFileRoute('/philosophy')({
  component: PhilosophyPage,
  head: () =>
    pageHead(
      '/philosophy',
      'Our philosophy',
      'Open-source infrastructure under Apache-2.0, built around tools you can inspect and servers you control.'
    ),
});
function PhilosophyPage() {
  return (
    <section className="mx-auto max-w-3xl px-6 py-24">
      <p className="font-semibold text-primary text-sm">WHY CODEDOCK</p>
      <h1 className="mt-4 font-bold text-4xl md:text-5xl">
        Infrastructure you can understand and own.
      </h1>
      <p className="mt-6 text-lg text-muted-foreground">
        Deploying should be approachable without hiding the infrastructure that keeps your services
        running.
      </p>
      <div className="mt-12 space-y-8">
        {principles.map((item) => (
          <article key={item.title}>
            <h2 className="font-semibold text-xl">{item.title}</h2>
            <p className="mt-3 text-muted-foreground leading-relaxed">{item.body}</p>
          </article>
        ))}
      </div>
      <div className="mt-10 flex flex-wrap gap-5 text-primary">
        <a href={productLinks.github}>Explore the source ?</a>
        <a href={`${productLinks.github}/blob/main/LICENSE`}>Read the Apache-2.0 licence ?</a>
      </div>
    </section>
  );
}
