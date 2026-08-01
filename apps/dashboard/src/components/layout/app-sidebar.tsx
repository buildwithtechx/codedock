import {
  Bot,
  Building,
  Cloud,
  Code,
  Download,
  FolderKanban,
  HardDrive,
  Key,
  LayoutDashboard,
  Network,
  PanelLeft,
  RefreshCw,
  Rocket,
  ScrollText,
  Server,
  Settings,
  Sparkles,
  Users,
  Wrench,
  X,
} from 'lucide-react';
import { useState } from 'react';
import { CreateProjectModal } from '#/features/projects/create-project-modal';
import { ServerConnectionSwitcher } from '../server-connection-switcher';
import { NavItem, type NavItemProps } from './nav-item';
import { OrganizationSwitcher } from './organization-switcher';
import { UserMenu } from './user-menu';

type NavGroup = {
  title?: string;
  items: (NavItemProps & { exact?: boolean })[];
};

const navGroups: NavGroup[] = [
  {
    title: 'Workspace',
    items: [
      { title: 'Home', url: '/', icon: LayoutDashboard, exact: true },
      { title: 'Projects', url: '/projects', icon: FolderKanban },
      { title: 'Apps', url: '/apps', icon: Sparkles },
      { title: 'Deployments', url: '/deployments', icon: Rocket },
    ],
  },
  {
    title: 'Infrastructure',
    items: [
      { title: 'Servers', url: '/servers', icon: Server },
      { title: 'Storage', url: '/s3-destinations', icon: HardDrive },
      { title: 'Domains & DNS', url: '/dns', icon: Network },
    ],
  },
  {
    title: 'Connect',
    items: [
      { title: 'Sources', url: '/sources', icon: Code },
      { title: 'Automation', url: '/ai', icon: Bot },
    ],
  },
  {
    title: 'Administration',
    items: [
      { title: 'Organizations', url: '/organizations', icon: Building },
      { title: 'Members', url: '/users', icon: Users },
      { title: 'API Access', url: '/api-access', icon: Key },
      { title: 'Settings', url: '/settings', icon: Settings, exact: true },
      { title: 'Maintenance', url: '/maintenance', icon: Wrench },
      { title: 'Updates', url: '/updates', icon: RefreshCw },
      { title: 'Migration', url: '/migrations', icon: Download },
    ],
  },
];

const bottomNav = [
  {
    title: 'Docs',
    url: 'https://docs.codedock.com',
    icon: ScrollText,
    external: true,
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
          collapsed ? 'md:w-16' : 'md:w-68'
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
          <div className={collapsed ? 'px-2 py-3' : 'px-3 py-3'}>
            <div className={`flex items-center ${collapsed ? 'justify-center' : 'gap-2.5 py-2'}`}>
              <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-primary/12 text-primary">
                <Cloud className="h-4 w-4" />
              </div>
              {!collapsed && (
                <>
                  <span className="flex-1 truncate font-medium text-sidebar-foreground text-sm">
                    Codedock
                  </span>
                  <span className="rounded-md bg-sidebar-accent px-1.5 py-0.5 font-medium text-[10px] text-muted-foreground">
                    v0.1
                  </span>
                  <button
                    type="button"
                    onClick={onToggle}
                    className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg text-muted-foreground transition-all duration-150 hover:bg-sidebar-accent hover:text-sidebar-foreground"
                  >
                    <PanelLeft className="h-4 w-4" />
                  </button>
                </>
              )}
            </div>
          </div>

          {collapsed && (
            <button
              type="button"
              onClick={onToggle}
              className="absolute top-3 right-0 z-30 hidden h-7 w-7 translate-x-1/2 items-center justify-center rounded-lg border border-sidebar-border bg-card text-muted-foreground shadow-md transition-all duration-300 hover:bg-sidebar-accent hover:text-sidebar-foreground active:scale-[0.95] md:flex"
            >
              <PanelLeft className="h-4 w-4" />
            </button>
          )}
        </div>

        <nav className="flex flex-1 flex-col gap-5 overflow-y-auto px-3 pt-2 pb-3">
          <OrganizationSwitcher collapsed={navCollapsed} />
          {navGroups.map((group, i) => (
            <div key={i} className="flex flex-col gap-0.5">
              {!navCollapsed && group.title && (
                <h4 className="px-2 pb-1.5 font-medium text-[10px] text-sidebar-foreground/35 uppercase tracking-[0.14em]">
                  {group.title}
                </h4>
              )}
              {group.items.map((item) => (
                <NavItem key={item.url} item={item} exact={item.exact} collapsed={navCollapsed} />
              ))}
            </div>
          ))}
        </nav>

        {!navCollapsed && (
          <div className="px-3 pb-3">
            <button
              type="button"
              onClick={() => setCreateProjectOpen(true)}
              className="flex h-10 w-full items-center justify-center gap-2 rounded-lg bg-primary font-semibold text-primary-foreground text-sm transition-colors hover:bg-primary-hover"
            >
              <FolderKanban className="h-4 w-4" />
              New project
            </button>
          </div>
        )}

        <div
          className={`mt-auto flex flex-col gap-0.5 ${navCollapsed ? 'px-1 py-1' : 'px-3 py-2'}`}
        >
          {bottomNav.map((item) => (
            <NavItem key={item.url} item={item} collapsed={navCollapsed} />
          ))}
        </div>

        {!navCollapsed && (
          <div className="border-sidebar-border/30 border-t px-2 py-1.5">
            <ServerConnectionSwitcher />
          </div>
        )}
        <UserMenu collapsed={navCollapsed} />
        <CreateProjectModal open={createProjectOpen} onOpenChange={setCreateProjectOpen} />
      </aside>
    </>
  );
}
