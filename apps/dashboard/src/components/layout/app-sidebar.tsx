import {
  Cloud,
  FolderKanban,
  HardDrive,
  Key,
  LayoutDashboard,
  Moon,
  PanelLeftClose,
  PanelLeftOpen,
  Rocket,
  Settings,
  Sparkles,
  Sun,
  X,
} from 'lucide-react';
import { useTheme } from 'next-themes';
import { useState } from 'react';
import { CreateProjectModal } from '#/features/projects/create-project-modal';
import { NavItem, type NavItemProps } from './nav-item';
import { OrganizationSwitcher } from './organization-switcher';

type NavGroup = {
  title?: string;
  items: (NavItemProps & { exact?: boolean })[];
};

const navGroups: NavGroup[] = [
  {
    title: 'Main',
    items: [
      { title: 'Home', url: '/', icon: LayoutDashboard, exact: true },
      { title: 'Projects', url: '/projects', icon: FolderKanban },
      { title: 'Apps', url: '/apps', icon: Sparkles },
      { title: 'Deployments', url: '/deployments', icon: Rocket },
    ],
  },
  {
    title: 'Settings',
    items: [
      {
        title: 'Backups',
        url: '/backups',
        icon: HardDrive,
        exact: true,
      },
      { title: 'API Access', url: '/api-access', icon: Key },
      { title: 'Settings', url: '/settings', icon: Settings, exact: true },
    ],
  },
];

interface AppSidebarProps {
  collapsed: boolean;
  onToggle: () => void;
  mobileOpen: boolean;
  onMobileClose: () => void;
}

export function AppSidebar({ collapsed, onToggle, mobileOpen, onMobileClose }: AppSidebarProps) {
  const navCollapsed = collapsed && !mobileOpen;
  const [createProjectOpen, setCreateProjectOpen] = useState(false);
  const { resolvedTheme, setTheme } = useTheme();

  const toggleTheme = () => setTheme(resolvedTheme === 'dark' ? 'light' : 'dark');

  return (
    <>
      {mobileOpen && (
        <button
          type="button"
          className="fixed inset-0 z-30 cursor-default bg-black/50 md:hidden"
          onClick={onMobileClose}
          aria-label="Close menu"
        />
      )}

      <aside
        className={`fixed inset-y-3 left-3 z-40 flex flex-col overflow-hidden rounded-2xl border border-sidebar-border/70 bg-sidebar shadow-2xl shadow-black/15 transition-all duration-300 md:z-20 ${
          collapsed ? 'md:w-[72px]' : 'md:w-[260px]'
        } ${mobileOpen ? 'w-[min(19rem,calc(100vw-1.5rem))] translate-x-0' : 'w-[min(19rem,calc(100vw-1.5rem))] -translate-x-[calc(100%+1rem)] md:translate-x-0'}`}
      >
        <div className="flex items-center justify-between px-3 pt-3 pb-2 md:hidden">
          <div className="flex items-center gap-2.5 py-2">
            <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-primary/12 text-primary">
              <Cloud className="h-4 w-4" />
            </div>
            <span className="truncate font-medium text-sidebar-foreground text-sm">Codedock</span>
          </div>
          <button
            type="button"
            onClick={onMobileClose}
            className="flex h-7 w-7 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        <div className="hidden md:block">
          <div className={collapsed ? 'px-2 py-3' : 'px-5 py-5'}>
            <div
              className={`flex ${collapsed ? 'flex-col items-center gap-2' : 'items-center justify-between gap-2.5 py-2'}`}
            >
              <div className="flex min-w-0 items-center gap-2.5">
                <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-primary/12 text-primary">
                  <Cloud className="h-4 w-4" />
                </div>
                {!collapsed && (
                  <span className="flex-1 truncate font-medium text-sidebar-foreground text-sm">
                    Codedock
                  </span>
                )}
              </div>
              <div className={`flex items-center ${collapsed ? 'flex-col gap-1' : 'gap-1'}`}>
                <button
                  type="button"
                  onClick={toggleTheme}
                  className="flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
                  aria-label="Toggle theme"
                >
                  {resolvedTheme === 'dark' ? (
                    <Sun className="h-4 w-4" />
                  ) : (
                    <Moon className="h-4 w-4" />
                  )}
                </button>
                <button
                  type="button"
                  onClick={onToggle}
                  className="flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
                  aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
                >
                  {collapsed ? (
                    <PanelLeftOpen className="h-4 w-4" />
                  ) : (
                    <PanelLeftClose className="h-4 w-4" />
                  )}
                </button>
              </div>
            </div>
          </div>
          <div className="mx-3 h-px bg-sidebar-border" />
        </div>

        <nav className="flex flex-1 flex-col gap-5 overflow-y-auto px-3 pt-3 pb-12">
          {navGroups.map((group, i) => (
            <div key={i} className="flex flex-col gap-0.5">
              {!navCollapsed && group.title && (
                <h4 className="px-2 pb-1.5 font-medium text-[10px] text-sidebar-foreground/50 uppercase tracking-[0.14em]">
                  {group.title}
                </h4>
              )}
              {group.items.map((item) => (
                <NavItem key={item.url} item={item} exact={item.exact} collapsed={navCollapsed} />
              ))}
            </div>
          ))}
        </nav>

        <div className="px-3 pb-2">
          <button
            type="button"
            onClick={() => setCreateProjectOpen(true)}
            className={`flex w-full items-center justify-center gap-2 rounded-xl bg-primary font-semibold text-primary-foreground text-sm transition-colors hover:bg-primary-hover ${
              navCollapsed ? 'h-10 px-0' : 'h-10 px-3'
            }`}
            aria-label="New project"
            title={navCollapsed ? 'New project' : undefined}
          >
            <FolderKanban className="h-4 w-4" />
            {!navCollapsed && 'New project'}
          </button>
        </div>

        <div className={`mt-auto px-3 pt-1 pb-3 ${navCollapsed ? 'px-2' : ''}`}>
          <div className="mx-2 mb-3 h-px bg-sidebar-border/60" />
          {!navCollapsed && (
            <p className="mb-2 px-2 font-medium text-[10px] text-sidebar-foreground/50 uppercase tracking-[0.14em]">
              Workspace
            </p>
          )}
          <OrganizationSwitcher collapsed={navCollapsed} />
        </div>
        <CreateProjectModal open={createProjectOpen} onOpenChange={setCreateProjectOpen} />
      </aside>
    </>
  );
}
