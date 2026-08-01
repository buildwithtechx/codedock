import { Link } from '@tanstack/react-router';
import { UserRound } from 'lucide-react';
import { useAuthStore } from '#/stores/auth-store';

export function UserMenu({ collapsed }: { collapsed: boolean }) {
  const user = useAuthStore((state) => state.user);
  const initials = user?.name
    ? user.name
        .split(' ')
        .map((part) => part[0])
        .join('')
        .toUpperCase()
    : 'U';

  return (
    <Link
      to="/profile"
      className={`mt-1 flex items-center rounded-xl transition-colors hover:bg-sidebar-accent ${
        collapsed ? 'justify-center px-1 py-2' : 'gap-3 px-2 py-2'
      }`}
      title={collapsed ? user?.name || 'Account' : undefined}
    >
      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-sidebar-accent font-semibold text-[10px] text-sidebar-foreground">
        {initials || <UserRound className="h-4 w-4" />}
      </div>
      {!collapsed && (
        <div className="min-w-0 flex-1">
          <p className="truncate font-medium text-sidebar-foreground text-sm">
            {user?.name || 'Account'}
          </p>
          <p className="truncate text-sidebar-foreground/45 text-xs">
            {user?.email || 'Account settings'}
          </p>
        </div>
      )}
    </Link>
  );
}
