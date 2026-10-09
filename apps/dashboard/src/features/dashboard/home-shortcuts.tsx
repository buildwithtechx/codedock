import { Link } from '@tanstack/react-router';
import { ArrowUpRight, BookOpen, Bot, GitBranch, Settings } from 'lucide-react';

const shortcuts = [
  {
    title: 'Import from Git',
    description: 'Deploy a repository in minutes.',
    to: '/library',
    icon: GitBranch,
    external: false,
  },
  {
    title: 'MCP deploy',
    description: 'Deploy via AI assistant.',
    to: '/settings?tab=mcp',
    icon: Bot,
    external: false,
  },
  {
    title: 'Settings',
    description: 'Configure your instance.',
    to: '/settings',
    icon: Settings,
    external: false,
  },
  {
    title: 'Documentation',
    description: 'Guides and references.',
    to: 'https://docs.codedock.run',
    icon: BookOpen,
    external: true,
  },
];

export function HomeShortcuts() {
  return (
    <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
      {shortcuts.map((shortcut) => {
        const content = (
          <>
            <div className="flex items-start justify-between">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-muted text-muted-foreground transition-colors group-hover:bg-primary/12 group-hover:text-primary">
                <shortcut.icon className="h-4 w-4" />
              </div>
              <ArrowUpRight className="h-3.5 w-3.5 text-muted-foreground transition-colors group-hover:text-primary" />
            </div>
            <h3 className="mt-5 font-semibold text-sm">{shortcut.title}</h3>
            <p className="mt-1 text-muted-foreground text-xs leading-5">{shortcut.description}</p>
          </>
        );
        const className = 'group rounded-xl bg-card p-4 transition-colors hover:bg-muted/60';
        return shortcut.external ? (
          <a
            key={shortcut.title}
            href={shortcut.to}
            target="_blank"
            rel="noopener noreferrer"
            className={className}
          >
            {content}
          </a>
        ) : (
          <Link key={shortcut.title} to={shortcut.to} className={className}>
            {content}
          </Link>
        );
      })}
    </section>
  );
}
