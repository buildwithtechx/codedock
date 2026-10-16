import { useNavigate } from '@tanstack/react-router';
import { Building2, Check, ChevronsUpDown, Loader2, LogOut, UserRound } from 'lucide-react';
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
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <button
            type="button"
            className={`group flex w-full items-center rounded-xl px-2 py-2 text-start transition-colors hover:bg-foreground/[0.06] ${
              collapsed ? 'justify-center' : 'gap-3'
            }`}
            aria-label={`Switch organization. Current organization: ${selected.name}`}
            title={collapsed ? selected.name : undefined}
          >
            <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-foreground/[0.08] text-sm font-semibold uppercase text-foreground">
              {organizationInitials(selected.name) || <Building2 className="size-4" />}
            </div>
            {!collapsed && (
              <>
                <div className="min-w-0 flex-1">
                  <p className="truncate text-[14px] font-semibold leading-tight text-foreground">
                    {selected.name}
                  </p>
                  <p className="truncate text-[12px] leading-tight text-muted-foreground">
                    {organizations.length > 1
                      ? `${organizations.length} workspaces`
                      : (user?.email ?? 'Workspace')}
                  </p>
                </div>
                <ChevronsUpDown className="size-4 shrink-0 text-muted-foreground transition-colors group-hover:text-foreground" />
              </>
            )}
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          className="!w-[240px] max-w-[calc(100vw-2rem)] overflow-hidden rounded-2xl border border-border/50 bg-popover p-0 shadow-xl"
          align="start"
          side="top"
          collisionPadding={16}
        >
          <div className="px-3 pt-3 pb-2">
            <p className="text-[11px] font-semibold uppercase tracking-widest text-muted-foreground/70">
              Switch organization
            </p>
          </div>
          <div className="max-h-60 overflow-y-auto px-1 pb-1">
            {organizations.map((organization) => (
              <DropdownMenuItem
                key={organization.id}
                onSelect={() => switchOrganization(organization.id)}
                className="flex cursor-pointer items-center gap-2.5 rounded-xl px-2.5 py-2 hover:bg-foreground/[0.05]"
              >
                <div className="flex size-7 shrink-0 items-center justify-center rounded-lg bg-foreground/[0.08] text-[12px] font-semibold uppercase text-foreground">
                  {organizationInitials(organization.name) || <Building2 className="size-3.5" />}
                </div>
                <span className="flex-1 truncate text-[13px] font-medium leading-tight text-foreground">
                  {organization.name}
                </span>
                {organization.id === selected.id && <Check className="size-4 shrink-0 text-primary" />}
              </DropdownMenuItem>
            ))}
          </div>
          <div className="border-t border-border/40 p-1.5">
            <DropdownMenuItem
              onSelect={() => void navigate({ to: '/profile' })}
              className="flex cursor-pointer items-center gap-2.5 rounded-xl px-2.5 py-2 hover:bg-foreground/[0.05]"
            >
              <div className="flex size-7 shrink-0 items-center justify-center rounded-full bg-foreground/[0.08] text-[11px] font-semibold uppercase text-foreground">
                {user?.name?.[0]?.toUpperCase() || <UserRound className="size-3.5" />}
              </div>
              <div className="min-w-0 flex-1">
                <p className="truncate text-[12px] font-medium leading-tight text-foreground">
                  {user?.name || 'Account'}
                </p>
                <p className="truncate text-[11px] leading-tight text-muted-foreground">
                  {user?.email || 'Profile & security'}
                </p>
              </div>
            </DropdownMenuItem>
            <DropdownMenuItem
              disabled={isLoggingOut}
              onSelect={() => logout()}
              variant="destructive"
              className="mt-0.5 flex cursor-pointer items-center gap-2 rounded-xl px-2.5 py-2 text-[13px] font-medium text-destructive transition-colors hover:bg-destructive/10"
            >
              <LogOut className="size-4" />
              Sign out
            </DropdownMenuItem>
          </div>
        </DropdownMenuContent>
      </DropdownMenu>
      {collapsed && (
        <button
          type="button"
          onClick={() => logout()}
          disabled={isLoggingOut}
          className="mt-2 flex w-full items-center justify-center rounded-xl py-2 text-muted-foreground transition-colors hover:bg-foreground/[0.06] hover:text-foreground disabled:opacity-50"
          aria-label="Sign out"
          title="Sign out"
        >
          {isLoggingOut ? (
            <Loader2 className="size-4 animate-spin" />
          ) : (
            <LogOut className="size-4" />
          )}
        </button>
      )}
    </>
  );
}
