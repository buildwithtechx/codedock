import { useNavigate } from '@tanstack/react-router';
import { Building2, Check, ChevronsUpDown, Plus } from 'lucide-react';
import { useEffect } from 'react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { useListOrganizations } from '#/features/organizations';
import { useOrganizationStore } from '#/stores/organization-store';

export function OrganizationSwitcher({ collapsed }: { collapsed: boolean }) {
  const navigate = useNavigate();
  const { data: organizations = [], isLoading } = useListOrganizations();
  const activeOrganizationId = useOrganizationStore((state) => state.activeOrganizationId);
  const setActiveOrganizationId = useOrganizationStore((state) => state.setActiveOrganizationId);
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
          className={`flex w-full items-center rounded-xl border border-sidebar-border/50 bg-sidebar-accent/35 text-left transition-colors hover:bg-sidebar-accent ${
            collapsed ? 'justify-center p-2' : 'gap-2.5 px-2.5 py-2'
          }`}
          aria-label={`Switch organization. Current organization: ${selected.name}`}
        >
          <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-primary/15 text-primary">
            <Building2 className="h-3.5 w-3.5" />
          </div>
          {!collapsed && (
            <>
              <div className="min-w-0 flex-1">
                <p className="truncate font-medium text-sidebar-foreground text-xs">
                  {selected.name}
                </p>
                <p className="text-[10px] text-sidebar-foreground/50 uppercase tracking-wider">
                  Workspace
                </p>
              </div>
              <ChevronsUpDown className="h-3.5 w-3.5 shrink-0 text-sidebar-foreground/50" />
            </>
          )}
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-60" align="start">
        <DropdownMenuLabel>Switch workspace</DropdownMenuLabel>
        {organizations.map((organization) => (
          <DropdownMenuItem
            key={organization.id}
            onSelect={() => switchOrganization(organization.id)}
          >
            <Building2 className="h-4 w-4 text-muted-foreground" />
            <span className="flex-1 truncate">{organization.name}</span>
            {organization.id === selected.id && <Check className="h-4 w-4 text-primary" />}
          </DropdownMenuItem>
        ))}
        <DropdownMenuSeparator />
        <DropdownMenuItem onSelect={() => void navigate({ to: '/organizations' })}>
          <Plus className="h-4 w-4" />
          Manage organizations
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
