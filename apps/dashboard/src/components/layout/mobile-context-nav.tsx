import { Link } from '@tanstack/react-router';
import type { LucideIcon } from 'lucide-react';

type MobileContextNavItem = {
  title: string;
  to: string;
  icon: LucideIcon;
  active?: boolean;
};

export function MobileContextNav({
  items,
  label,
}: {
  items: MobileContextNavItem[];
  label: string;
}) {
  return (
    <nav
      aria-label={label}
      className="flex gap-1 overflow-x-auto rounded-xl border border-border/80 bg-card p-1.5 xl:hidden"
    >
      {items.map((item) => (
        <Link
          key={item.to}
          to={item.to as never}
          className={`flex shrink-0 items-center gap-2 rounded-lg px-3 py-2 font-medium text-sm transition-colors ${
            item.active
              ? 'bg-primary/12 text-foreground'
              : 'text-muted-foreground hover:bg-muted hover:text-foreground'
          }`}
        >
          <item.icon className="h-4 w-4" />
          {item.title}
        </Link>
      ))}
    </nav>
  );
}
