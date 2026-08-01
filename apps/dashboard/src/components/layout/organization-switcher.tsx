import { useNavigate } from '@tanstack/react-router';
import { ArrowUpRight, Building2, Check, ChevronsUpDown, LogOut, UserRound } from 'lucide-react';
import { useEffect } from 'react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { useLogout } from '#/features/auth';
import { useListOrganizations } from '#/features/organizations';
import { useAuthStore } from '#/stores/auth-store';
import { useOrganizationStore } from '#/stores/organization-store';

const organizationInitials = (name: string) =>
  name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0])
    .join('')
    .toUpperCase();

export function OrganizationSwitcher({ collapsed }: { collapsed: boolean }) {
  const navigate = useNavigate();
  const { data: organizations = [], isLoading } = useListOrganizations();
  const activeOrganizationId = useOrganizationStore((state) => state.activeOrganizationId);
  const setActiveOrganizationId = useOrganizationStore((state) => state.setActiveOrganizationId);
  const user = useAuthStore((state) => state.user);
  const { mutate: logout, isPending: isLoggingOut } = useLogout();
  const activeOrganization = organizations.find((org) => org.id === activeOrganizationId);

  useEffect(() => {
    if (!isLoading && organizations.length > 0 && !activeOrganization) {
      setActiveOrganizationId(organizations[0].id);
    }
  }, [activeOrganization, isLoading, organizations, setActiveOrganizationId]);

  if (isLoading || organizations.length === 0) {
    return null;
  }

  const selected = activeOrganization || organizations[0];
  const switchOrganization = (organizationId: string) => {
    setActiveOrganizationId(organizationId);
    void navigate({ to: '/' });
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button
          type="button"
          className={`group flex w-full items-center rounded-lg text-left transition-colors hover:bg-sidebar-accent ${
            collapsed ? 'justify-center p-2' : 'gap-2.5 px-2 py-2'
          }`}
          aria-label={`Switch organization. Current organization: ${selected.name}`}
        >
          <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary/15 font-semibold text-[11px] text-primary">
            {organizationInitials(selected.name) || <Building2 className="h-4 w-4" />}
          </div>
          {!collapsed && (
            <>
              <div className="min-w-0 flex-1">
                <p className="truncate font-medium text-sidebar-foreground text-sm">
                  {selected.name}
                </p>
                <p className="mt-0.5 text-[10px] text-sidebar-foreground/45 uppercase tracking-[0.14em]">
                  Organization
                </p>
              </div>
              <ChevronsUpDown className="h-4 w-4 shrink-0 text-sidebar-foreground/40 transition-colors group-hover:text-sidebar-foreground/70" />
            </>
          )}
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        className="!w-72 max-w-[calc(100vw-2rem)] rounded-xl border border-border/70 bg-popover p-1.5 shadow-xl"
        align="start"
        side="top"
        collisionPadding={16}
      >
        <DropdownMenuLabel className="px-2.5 pt-1.5 pb-2 font-medium text-[10px] uppercase tracking-[0.14em]">
          Organizations
        </DropdownMenuLabel>
        {organizations.map((organization) => (
          <DropdownMenuItem
            key={organization.id}
            onSelect={() => switchOrganization(organization.id)}
            className="min-h-11 gap-2.5 rounded-lg px-2.5 py-2"
          >
            <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-muted font-semibold text-[10px] text-muted-foreground">
              {organizationInitials(organization.name) || <Building2 className="h-3.5 w-3.5" />}
            </div>
            <span className="flex-1 truncate font-medium text-sm">{organization.name}</span>
            {organization.id === selected.id && <Check className="h-4 w-4 text-primary" />}
          </DropdownMenuItem>
        ))}
        <DropdownMenuSeparator />
        <DropdownMenuItem
          onSelect={() => void navigate({ to: '/organizations' })}
          className="min-h-10 gap-2.5 rounded-lg px-2.5 py-2 text-muted-foreground"
        >
          <Building2 className="h-4 w-4" />
          <span className="flex-1">Manage organizations</span>
          <ArrowUpRight className="h-3.5 w-3.5" />
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <div className="flex min-h-12 items-center gap-2.5 rounded-lg px-2.5 py-2">
          <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-muted font-semibold text-[10px] text-muted-foreground">
            {user?.name?.[0]?.toUpperCase() || <UserRound className="h-3.5 w-3.5" />}
          </div>
          <div className="min-w-0 flex-1">
            <p className="truncate font-medium text-sm">{user?.name || 'Account'}</p>
            <p className="truncate text-muted-foreground text-xs">
              {user?.email || 'Account settings'}
            </p>
          </div>
        </div>
        <DropdownMenuItem
          disabled={isLoggingOut}
          onSelect={() => logout()}
          variant="destructive"
          className="min-h-10 gap-2.5 rounded-lg px-2.5 py-2"
        >
          <LogOut className="h-4 w-4" />
          Sign out
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
