import { PanelLeft, PanelLeftClose, Search, X } from 'lucide-react';
import { useEffect, useState } from 'react';
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@codedock/ui';
import { cn } from 'cn';
import { CodedockLogo } from './codedock-logo';
import { NavItem } from './nav-item';
import { getSidebarSections } from './sidebar-items';
import { SidebarUserPopover } from './sidebar-user-popover';
import { useSidebarCounts } from './use-sidebar-counts';

interface AppSidebarProps {
  collapsed: boolean;
  onToggle: () => void;
  mobileOpen: boolean;
  onMobileClose: () => void;
  onOpenCommand: () => void;
}

export function AppSidebar({
  collapsed,
  onToggle,
  mobileOpen,
  onMobileClose,
  onOpenCommand,
}: AppSidebarProps) {
  const [autoCollapse, setAutoCollapse] = useState(() => {
    if (typeof window === 'undefined') return false;
    return localStorage.getItem('codedock_auto_collapse_sidebar') === 'true';
  });

  const counts = useSidebarCounts();
  const navCollapsed = collapsed && !mobileOpen;

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key.toLowerCase() === 'b' && (event.metaKey || event.ctrlKey)) {
        event.preventDefault();
        onToggle();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [onToggle]);

  const handleToggleAutoCollapse = () => {
    const next = !autoCollapse;
    setAutoCollapse(next);
    localStorage.setItem('codedock_auto_collapse_sidebar', next ? 'true' : 'false');
    if (next !== collapsed) {
      onToggle();
    }
  };

  const { workspaceItems, infrastructureItems, manageItems } = getSidebarSections(counts);

  return (
    <TooltipProvider delayDuration={150}>
      <aside
        className={cn(
          'fixed inset-y-0 left-0 z-40 flex h-dvh max-h-dvh flex-col border-sidebar-border border-r bg-sidebar text-sidebar-foreground transition-all duration-200 ease-in-out md:static',
          navCollapsed ? 'w-16' : 'w-60',
          mobileOpen ? 'translate-x-0 shadow-2xl' : '-translate-x-full md:translate-x-0'
        )}
      >
        {/* Sidebar Header */}
        <div className="flex h-13 shrink-0 items-center justify-between border-sidebar-border/70 border-b px-3">
          {!navCollapsed ? (
            <div className="flex min-w-0 items-center gap-2">
              <Avatar size="sm" className="size-6 rounded-lg shadow-xs">
                <AvatarImage src="/favicon.svg" alt="Codedock" className="rounded-lg" />
                <AvatarFallback className="rounded-lg bg-transparent p-0">
                  <CodedockLogo className="size-6 rounded-lg" />
                </AvatarFallback>
              </Avatar>
              <span className="truncate font-semibold text-sidebar-foreground text-sm tracking-tight">
                Codedock
              </span>
              <span className="rounded bg-muted/70 px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
                v0.1.0
              </span>
            </div>
          ) : (
            <Avatar className="mx-auto size-7 rounded-lg shadow-xs">
              <AvatarImage src="/favicon.svg" alt="Codedock" className="rounded-lg" />
              <AvatarFallback className="rounded-lg bg-transparent p-0">
                <CodedockLogo className="size-7 rounded-lg" />
              </AvatarFallback>
            </Avatar>
          )}

          {!navCollapsed && (
            <button
              type="button"
              onClick={onToggle}
              aria-label="Collapse sidebar (⌘B)"
              className="hidden size-7 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground md:flex"
            >
              <PanelLeftClose className="size-4" />
            </button>
          )}

          {mobileOpen && (
            <button
              type="button"
              onClick={onMobileClose}
              className="flex size-7 items-center justify-center rounded-lg text-muted-foreground hover:bg-sidebar-accent md:hidden"
            >
              <X className="size-4" />
            </button>
          )}
        </div>

        {/* Quick Search & Collapsed Expand Toggle */}
        <div className="p-2.5 pb-1 space-y-1.5">
          {!navCollapsed ? (
            <button
              type="button"
              onClick={onOpenCommand}
              className="flex w-full items-center justify-between rounded-lg border border-sidebar-border/80 bg-sidebar-accent/40 px-2.5 py-1.5 text-muted-foreground text-xs transition-colors hover:border-sidebar-border hover:bg-sidebar-accent hover:text-sidebar-foreground"
            >
              <span className="flex items-center gap-2">
                <Search className="size-3.5 text-muted-foreground" />
                <span>Search</span>
              </span>
              <kbd className="rounded border border-sidebar-border/80 bg-sidebar px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
                Ctrl+K
              </kbd>
            </button>
          ) : (
            <>
              <Tooltip>
                <TooltipTrigger asChild>
                  <button
                    type="button"
                    onClick={onOpenCommand}
                    className="mx-auto flex size-9 items-center justify-center rounded-lg border border-sidebar-border/80 bg-sidebar-accent/40 text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
                  >
                    <Search className="size-3.5" />
                  </button>
                </TooltipTrigger>
                <TooltipContent side="right">Search (Ctrl+K)</TooltipContent>
              </Tooltip>

              {/* Show Expand button right below search in collapsed state */}
              <Tooltip delayDuration={150}>
                <TooltipTrigger asChild>
                  <button
                    type="button"
                    onClick={onToggle}
                    aria-label="Expand sidebar (⌘B)"
                    className="mx-auto flex size-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
                  >
                    <PanelLeft className="size-4" />
                  </button>
                </TooltipTrigger>
                <TooltipContent side="right">Expand sidebar (⌘B)</TooltipContent>
              </Tooltip>
            </>
          )}
        </div>

        {/* Scrollable Navigation List */}
        <div className="min-h-0 flex-1 overflow-y-auto overscroll-contain px-2.5 py-2">
          <div className="space-y-3">
            {/* Workspace Group */}
            <div className="space-y-0.5">
              {!navCollapsed && (
                <p className="px-2.5 py-1 font-semibold text-[10px] text-muted-foreground/75 uppercase tracking-wider">
                  Workspace
                </p>
              )}
              {workspaceItems.map((item) => (
                <NavItem key={item.url + item.title} item={item} collapsed={navCollapsed} />
              ))}
            </div>

            {/* Fleet & Infrastructure Group */}
            <div className="space-y-0.5">
              {!navCollapsed && (
                <p className="px-2.5 py-1 font-semibold text-[10px] text-muted-foreground/75 uppercase tracking-wider">
                  Infrastructure
                </p>
              )}
              {infrastructureItems.map((item) => (
                <NavItem key={item.url + item.title} item={item} collapsed={navCollapsed} />
              ))}
            </div>

            {/* Manage Group */}
            <div className="space-y-0.5">
              {!navCollapsed && (
                <p className="px-2.5 py-1 font-semibold text-[10px] text-muted-foreground/75 uppercase tracking-wider">
                  Manage
                </p>
              )}
              {manageItems.map((item) => (
                <NavItem key={item.url + item.title} item={item} collapsed={navCollapsed} />
              ))}
            </div>
          </div>
        </div>

        {/* Sidebar Footer with Active User Popover */}
        <div className="shrink-0 border-sidebar-border/70 border-t bg-sidebar p-2 pb-[max(0.75rem,env(safe-area-inset-bottom))]">
          <SidebarUserPopover
            collapsed={navCollapsed}
            onToggleCollapse={onToggle}
            autoCollapse={autoCollapse}
            onToggleAutoCollapse={handleToggleAutoCollapse}
          />
        </div>
      </aside>
    </TooltipProvider>
  );
}
