import { SiteLink } from '../../components/site-link';

const workflows = [
  {
    title: 'Your own projects',
    body: 'Keep an API, database and background services together on a Linux server you control.',
    href: '/solutions/self-hosted',
  },
  {
    title: 'Client work',
    body: 'Organize client stacks into projects and environments, with roles and deployment histories for your team.',
    href: '/solutions/agencies',
  },
  {
    title: 'A growing fleet',
    body: 'Connect additional servers through SSH and choose where your services run.',
    href: '/solutions/enterprise',
  },
] as const;
export function SocialProof() {
  return (
    <section className="border-border border-y py-20">
      <div className="mx-auto max-w-7xl px-6">
        <h2 className="font-bold text-3xl">Start with the way you work.</h2>
        <div className="mt-8 grid gap-6 md:grid-cols-3">
          {workflows.map((item) => (
            <SiteLink
              key={item.title}
              href={item.href}
              className="rounded-xl border border-border bg-card/60 p-6 hover:border-primary/50"
            >
              <h3 className="font-semibold text-lg">{item.title}</h3>
              <p className="mt-3 text-muted-foreground text-sm leading-relaxed">{item.body}</p>
              <span className="mt-5 block text-primary text-sm">Explore the workflow →</span>
            </SiteLink>
          ))}
        </div>
      </div>
    </section>
  );
}
