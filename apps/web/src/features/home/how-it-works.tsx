import { GitBranch, Server, UserRound } from 'lucide-react';

const steps = [
  {
    icon: Server,
    title: 'Install on Linux',
    body: 'Run the installer on your server. It configures the daemon and container infrastructure needed to start.',
  },
  {
    icon: UserRound,
    title: 'Create your workspace',
    body: 'Open the dashboard and create the first owner account. Add a project and an environment.',
  },
  {
    icon: GitBranch,
    title: 'Deploy a service',
    body: 'Connect a repository or choose an image. Set variables and the listening port, then deploy and check the logs.',
  },
];
export function HowItWorks() {
  return (
    <section className="border-border border-y bg-card/40 py-20">
      <div className="mx-auto max-w-7xl px-6">
        <h2 className="font-bold text-3xl">From a server to your first deployment.</h2>
        <div className="mt-10 grid gap-8 md:grid-cols-3">
          {steps.map(({ icon: Icon, title, body }, index) => (
            <article key={title}>
              <Icon className="mb-5 size-6 text-primary" aria-hidden="true" />
              <h3 className="font-semibold text-lg">
                {index + 1}. {title}
              </h3>
              <p className="mt-3 text-muted-foreground text-sm leading-relaxed">{body}</p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
