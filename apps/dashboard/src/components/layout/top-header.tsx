import { Link, useRouterState } from '@tanstack/react-router';
import {
  Check,
  ChevronsUpDown,
  Menu,
  Moon,
  Plus,
  Search,
  Sun,
  Terminal as TerminalIcon,
  Users,
  X,
} from 'lucide-react';
import {
  Avatar,
  AvatarFallback,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
  Button,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  useTheme,
} from '@codedock/ui';
import { EnvironmentSwitcher } from '#/features/projects/environment-switcher';
import { useListOrganizations } from '#/hooks/use-organizations';
import { useGetProject } from '#/hooks/use-projects';
import { useServer } from '#/hooks/use-servers';
import { useAuthStore } from '#/stores/auth-store';
import { useDockStore } from '#/stores/dock-store';
import { useOrganizationStore } from '#/stores/organization-store';

import { HEADER_SECTION_TITLES } from './sidebar-items';

interface TopHeaderProps {
  mobileOpen: boolean;
  onToggleMobile: () => void;
  onOpenCommand: () => void;
}

export function TopHeader({ mobileOpen, onToggleMobile, onOpenCommand }: TopHeaderProps) {
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });
  const { theme, resolvedTheme, setTheme } = useTheme();
  const user = useAuthStore((s) => s.user);
  const openDock = useDockStore((s) => s.openDock);

  const activeOrgId = useOrganizationStore((s) => s.activeOrganizationId);
  const setActiveOrgId = useOrganizationStore((s) => s.setActiveOrganizationId);
  const { data: orgs } = useListOrganizations();

  const activeOrg = orgs?.find((o) => o.id === activeOrgId);
  const teamName = activeOrg?.name || user?.name || 'Solomon Olatunji';

  const segments = pathname.split('/').filter(Boolean);
  const rootSegment = segments[0] || '';
  const secondSegment = segments[1] || '';

  const isProjectPage =
    rootSegment === 'projects' && Boolean(secondSegment) && secondSegment !== 'new';
  const isServerPage =
    rootSegment === 'servers' && Boolean(secondSegment) && secondSegment !== 'new';

  const { data: projectRes } = useGetProject(isProjectPage ? secondSegment : '');
  const { data: serverRes } = useServer(isServerPage ? secondSegment : undefined);
  const projectName = projectRes?.data?.name || 'Project';
  const serverName = serverRes?.name || 'Server';
  const userInitial =
    user?.name?.charAt(0)?.toUpperCase() || user?.email?.charAt(0)?.toUpperCase() || 'S';

  const toggleTheme = () => {
    setTheme(theme === 'light' ? 'dim' : theme === 'dim' ? 'dark' : 'light');
  };

  const currentSectionTitle =
    HEADER_SECTION_TITLES[rootSegment] ||
    (rootSegment ? rootSegment.charAt(0).toUpperCase() + rootSegment.slice(1) : 'Dashboard');

  return (
    <header className="sticky top-0 z-30 flex h-14 shrink-0 items-center justify-between border-sidebar-border/80 border-b bg-sidebar px-4 text-sidebar-foreground sm:px-6">
      {/* Left: Mobile Toggle + Breadcrumb Hierarchy */}
      <div className="flex min-w-0 items-center gap-3">
        <button
          type="button"
          onClick={onToggleMobile}
          aria-label={mobileOpen ? 'Close menu' : 'Open menu'}
          className="flex size-8 items-center justify-center rounded-lg border border-sidebar-border text-muted-foreground transition-colors hover:bg-sidebar-accent md:hidden"
        >
          {mobileOpen ? <X className="size-4" /> : <Menu className="size-4" />}
        </button>

        <Breadcrumb className="hidden sm:block">
          <BreadcrumbList className="font-medium text-muted-foreground text-xs">
            {/* Team / Workspace Dropdown Segment */}
            <BreadcrumbItem>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <button
                    type="button"
                    className="inline-flex cursor-pointer items-center gap-1.5 rounded-md px-1.5 py-1 text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
                  >
                    <span className="font-medium text-sidebar-foreground">{teamName}</span>
                    <ChevronsUpDown className="size-3 text-muted-foreground/70" />
                  </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="start" className="w-56 p-1.5">
                  <DropdownMenuLabel className="font-normal text-[11px] text-muted-foreground">
                    Workspaces & Teams
                  </DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  {orgs && orgs.length > 0 ? (
                    orgs.map((org) => (
                      <DropdownMenuItem
                        key={org.id}
                        onClick={() => setActiveOrgId(org.id)}
                        className="flex cursor-pointer items-center justify-between py-1.5 text-xs"
                      >
                        <span className="truncate">{org.name}</span>
                        {org.id === activeOrgId && <Check className="size-3.5 text-primary" />}
                      </DropdownMenuItem>
                    ))
                  ) : (
                    <DropdownMenuItem
                      onClick={() => setActiveOrgId(null)}
                      className="flex justify-between py-1.5 text-xs"
                    >
                      <span className="truncate">{teamName}</span>
                      <Check className="size-3.5 text-primary" />
                    </DropdownMenuItem>
                  )}
                  <DropdownMenuSeparator />
                  <DropdownMenuItem asChild>
                    <Link
                      to="/team"
                      className="flex cursor-pointer items-center gap-2 py-1.5 text-xs"
                    >
                      <Users className="size-3.5 text-muted-foreground" />
                      <span>Manage Team</span>
                    </Link>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </BreadcrumbItem>

            <BreadcrumbSeparator className="text-muted-foreground/40">/</BreadcrumbSeparator>

            {/* Current Section or Project */}
            {!rootSegment ? (
              <BreadcrumbItem>
                <BreadcrumbPage className="inline-flex items-center gap-1 font-semibold text-sidebar-foreground">
                  <span>Dashboard</span>
                </BreadcrumbPage>
              </BreadcrumbItem>
            ) : isProjectPage ? (
              <>
                <BreadcrumbItem>
                  <BreadcrumbLink asChild>
                    <Link
                      to="/projects"
                      className="text-muted-foreground hover:text-sidebar-foreground"
                    >
                      Projects
                    </Link>
                  </BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator className="text-muted-foreground/40">/</BreadcrumbSeparator>
                <BreadcrumbItem>
                  <BreadcrumbPage className="inline-flex items-center gap-1.5 font-semibold text-sidebar-foreground">
                    <span>{projectName}</span>
                    <EnvironmentSwitcher projectId={secondSegment} onSelectEnvironment={() => {}} />
                  </BreadcrumbPage>
                </BreadcrumbItem>
              </>
            ) : isServerPage ? (
              <>
                <BreadcrumbItem>
                  <BreadcrumbLink asChild>
                    <Link
                      to="/servers"
                      className="text-muted-foreground hover:text-sidebar-foreground"
                    >
                      Servers
                    </Link>
                  </BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator className="text-muted-foreground/40">/</BreadcrumbSeparator>
                <BreadcrumbItem>
                  <BreadcrumbPage className="font-semibold text-sidebar-foreground">
                    {serverName}
                  </BreadcrumbPage>
                </BreadcrumbItem>
              </>
            ) : (
              <BreadcrumbItem>
                <BreadcrumbPage className="inline-flex items-center gap-1 font-semibold text-sidebar-foreground">
                  <span>{currentSectionTitle}</span>
                </BreadcrumbPage>
              </BreadcrumbItem>
            )}
          </BreadcrumbList>
        </Breadcrumb>
      </div>

      {/* Right: Actions, Search, Terminal Dock Trigger & Theme */}
      <div className="flex items-center gap-2">
        <Button asChild size="sm" className="hidden h-8 gap-1.5 text-xs sm:inline-flex">
          <Link to="/projects" search={{ new: 'true' } as never}>
            <Plus className="size-3.5" />
            <span>New Project</span>
          </Link>
        </Button>

        <Button
          variant="outline"
          size="sm"
          onClick={onOpenCommand}
          className="h-8 gap-2 border-sidebar-border bg-sidebar-accent/50 px-2.5 text-muted-foreground text-xs hover:bg-sidebar-accent hover:text-sidebar-foreground"
        >
          <Search className="size-3.5 text-muted-foreground" />
          <span className="hidden md:inline">Search...</span>
          <kbd className="rounded border border-sidebar-border bg-sidebar px-1 py-0.2 font-mono text-[9px] text-muted-foreground">
            ⌘K
          </kbd>
        </Button>

        <Button
          variant="ghost"
          size="sm"
          onClick={() => openDock({ tab: 'terminal' })}
          className="hidden h-8 gap-1 px-2 text-muted-foreground text-xs hover:text-sidebar-foreground sm:inline-flex"
          title="Toggle Terminal Panel (⌘J)"
        >
          <TerminalIcon className="size-3.5" />
          <kbd className="rounded border border-sidebar-border px-1 py-0.2 font-mono text-[9px]">
            ⌘J
          </kbd>
        </Button>

        <Button
          variant="ghost"
          size="icon"
          onClick={toggleTheme}
          aria-label="Toggle theme"
          className="size-8 text-muted-foreground hover:text-sidebar-foreground"
        >
          {resolvedTheme === 'dark' ? <Sun className="size-4" /> : <Moon className="size-4" />}
        </Button>

        <Link to="/settings" search={{ tab: 'profile' } as never} className="inline-flex">
          <Avatar className="size-7 ring-1 ring-sidebar-border">
            <AvatarFallback className="bg-primary/10 font-semibold text-[11px] text-primary">
              {userInitial}
            </AvatarFallback>
          </Avatar>
        </Link>
      </div>
    </header>
  );
}
