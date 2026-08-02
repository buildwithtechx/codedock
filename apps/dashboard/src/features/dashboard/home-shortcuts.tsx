import { Link } from '@tanstack/react-router';
import { ArrowUpRight, CloudCog, GitBranch, Globe2, Server } from 'lucide-react';

const shortcuts = [
  {
    title: 'Connect source',
    description: 'Link Git providers and registries.',
    to: '/settings',
    search: { tab: 'sources' },
    icon: GitBranch,
  },
  {
    title: 'Add server',
    description: 'Bring another runtime online.',
    to: '/servers',
    icon: Server,
  },
  {
    title: 'Configure domains',
    description: 'Review DNS and certificates.',
    to: '/dns',
    icon: Globe2,
  },
  {
    title: 'Platform settings',
    description: 'Tune your Codedock instance.',
    to: '/settings',
    icon: CloudCog,
  },
];

export function HomeShortcuts() {
  return (
    <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
      {shortcuts.map((shortcut) => (
        <Link
          key={shortcut.to}
          to={shortcut.to}
          search={shortcut.search as never}
          className="group rounded-xl bg-card p-4 transition-colors hover:bg-muted/60"
        >
          <div className="flex items-start justify-between">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-muted text-muted-foreground transition-colors group-hover:bg-primary/12 group-hover:text-primary">
              <shortcut.icon className="h-4 w-4" />
            </div>
            <ArrowUpRight className="h-3.5 w-3.5 text-muted-foreground transition-colors group-hover:text-primary" />
          </div>
          <h3 className="mt-5 font-semibold text-sm">{shortcut.title}</h3>
          <p className="mt-1 text-muted-foreground text-xs leading-5">{shortcut.description}</p>
        </Link>
      ))}
    </section>
  );
}
