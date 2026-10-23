import { Link } from '@tanstack/react-router';
import {
  BookOpen,
  Check,
  ChevronsUpDown,
  Heart,
  LogOut,
  MessageSquare,
  Moon,
  Palette,
  PanelLeftClose,
  Sparkles,
  Sun,
  SunMedium,
  User as UserIcon,
} from 'lucide-react';
import { useState } from 'react';
import {
  Avatar,
  AvatarFallback,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
  Tooltip,
  TooltipContent,
  TooltipTrigger,
  useTheme,
} from '@codedock/ui';
import { useAuthStore } from '#/stores/auth-store';

interface SidebarUserPopoverProps {
  collapsed: boolean;
  onToggleCollapse: () => void;
  autoCollapse?: boolean;
  onToggleAutoCollapse?: () => void;
}

const THEME_OPTIONS = [
  { id: 'light', label: 'Light', icon: Sun },
  { id: 'dim', label: 'Dim', icon: SunMedium },
  { id: 'dark', label: 'Dark', icon: Moon },
  { id: 'system', label: 'System', icon: Palette },
] as const;

export function SidebarUserPopover({
  collapsed,
  onToggleCollapse,
  autoCollapse = false,
  onToggleAutoCollapse,
}: SidebarUserPopoverProps) {
  const { theme, setTheme } = useTheme();
  const user = useAuthStore((s) => s.user);
  const logout = useAuthStore((s) => s.logout);
  const [open, setOpen] = useState(false);

  const userName = user?.name || user?.email?.split('@')[0] || 'Solomon Olatunji';
  const userEmail = user?.email || 'iamsolomonolatunji@gmail.com';
  const userInitial = userName.charAt(0).toUpperCase() || 'S';

  return (
    <div className="flex w-full items-center gap-1">
      <DropdownMenu open={open} onOpenChange={setOpen}>
        <DropdownMenuTrigger asChild>
          {collapsed ? (
            <Tooltip delayDuration={150}>
              <TooltipTrigger asChild>
                <button
                  type="button"
                  className="mx-auto flex size-9 items-center justify-center rounded-xl transition-colors hover:bg-sidebar-accent"
                  aria-label="User account menu"
                >
                  <Avatar className="size-8 shrink-0 ring-1 ring-sidebar-border">
                    <AvatarFallback className="bg-primary/10 font-semibold text-primary text-xs">
                      {userInitial}
                    </AvatarFallback>
                  </Avatar>
                </button>
              </TooltipTrigger>
              <TooltipContent side="right" className="text-xs">
                {userName}
              </TooltipContent>
            </Tooltip>
          ) : (
            <button
              type="button"
              className="group flex min-w-0 flex-1 items-center justify-between gap-2 rounded-xl p-1.5 text-left transition-colors hover:bg-sidebar-accent/80"
              aria-label="User account menu"
            >
              <div className="flex min-w-0 items-center gap-2.5">
                <Avatar className="size-8 shrink-0 ring-1 ring-sidebar-border">
                  <AvatarFallback className="bg-primary/10 font-semibold text-primary text-xs">
                    {userInitial}
                  </AvatarFallback>
                </Avatar>
                <div className="min-w-0 flex-1 text-left">
                  <p className="truncate font-medium text-sidebar-foreground text-xs">{userName}</p>
                  <p className="truncate text-[10px] text-muted-foreground">{userEmail}</p>
                </div>
              </div>
              <ChevronsUpDown className="size-3.5 shrink-0 text-muted-foreground transition-transform group-hover:text-sidebar-foreground" />
            </button>
          )}
        </DropdownMenuTrigger>

        <DropdownMenuContent
          side="top"
          align={collapsed ? 'start' : 'center'}
          sideOffset={8}
          className="w-64 rounded-xl border border-border/80 bg-popover p-1.5 text-popover-foreground shadow-2xl"
        >
          {/* User Card Header */}
          <DropdownMenuLabel className="p-2 font-normal">
            <div className="flex items-center gap-2.5">
              <Avatar className="size-9 shrink-0 ring-1 ring-border/50">
                <AvatarFallback className="bg-primary/10 font-semibold text-primary text-sm">
                  {userInitial}
                </AvatarFallback>
              </Avatar>
              <div className="min-w-0 flex-1">
                <p className="truncate font-medium text-foreground text-xs leading-none">
                  {userName}
                </p>
                <p className="mt-1 truncate text-[11px] text-muted-foreground leading-none">
                  {userEmail}
                </p>
              </div>
            </div>
          </DropdownMenuLabel>

          <DropdownMenuSeparator className="bg-border/60" />

          {/* Navigation links */}
          <DropdownMenuItem asChild>
            <Link
              to="/settings"
              search={{ tab: 'profile' } as never}
              className="flex cursor-pointer items-center gap-2.5 rounded-lg px-2.5 py-1.5 text-xs hover:bg-muted"
            >
              <UserIcon className="size-4 text-muted-foreground" />
              <span className="flex-1">Profile</span>
            </Link>
          </DropdownMenuItem>

          {/* Appearance Switcher Submenu */}
          <DropdownMenuSub>
            <DropdownMenuSubTrigger className="flex cursor-pointer items-center gap-2.5 rounded-lg px-2.5 py-1.5 text-xs hover:bg-muted">
              <Palette className="size-4 text-muted-foreground" />
              <span className="flex-1">Appearance</span>
              <span className="capitalize text-[11px] text-muted-foreground">{theme}</span>
            </DropdownMenuSubTrigger>
            <DropdownMenuSubContent className="w-36 rounded-xl border border-border/80 bg-popover p-1 shadow-xl">
              {THEME_OPTIONS.map(({ id, label, icon: Icon }) => (
                <DropdownMenuItem
                  key={id}
                  onClick={() => setTheme(id)}
                  className="flex cursor-pointer items-center justify-between rounded-lg px-2.5 py-1.5 text-xs hover:bg-muted"
                >
                  <span className="flex items-center gap-2">
                    <Icon className="size-3.5 text-muted-foreground" />
                    <span>{label}</span>
                  </span>
                  {theme === id && <Check className="size-3.5 text-primary" />}
                </DropdownMenuItem>
              ))}
            </DropdownMenuSubContent>
          </DropdownMenuSub>

          {/* Auto-collapse sidebar toggle */}
          {onToggleAutoCollapse && (
            <DropdownMenuItem
              onClick={onToggleAutoCollapse}
              className="flex cursor-pointer items-center justify-between rounded-lg px-2.5 py-1.5 text-xs hover:bg-muted"
            >
              <span className="flex items-center gap-2.5">
                <PanelLeftClose className="size-4 text-muted-foreground" />
                <span>Auto-collapse sidebar</span>
              </span>
              {autoCollapse && <Check className="size-3.5 text-primary" />}
            </DropdownMenuItem>
          )}

          <DropdownMenuSeparator className="bg-border/60" />

          {/* What's New */}
          <DropdownMenuItem asChild>
            <Link
              to="/settings"
              search={{ tab: 'instance' } as never}
              className="flex cursor-pointer items-center justify-between rounded-lg px-2.5 py-1.5 text-xs hover:bg-muted"
            >
              <span className="flex items-center gap-2.5">
                <Sparkles className="size-4 text-amber-500" />
                <span>What's New</span>
              </span>
              <span className="rounded-full bg-rose-500/20 px-1.5 py-0.2 font-semibold text-[10px] text-rose-500">
                9+
              </span>
            </Link>
          </DropdownMenuItem>

          {/* Documentation */}
          <DropdownMenuItem asChild>
            <a
              href="https://codedock.io/docs"
              target="_blank"
              rel="noopener noreferrer"
              className="flex cursor-pointer items-center gap-2.5 rounded-lg px-2.5 py-1.5 text-xs hover:bg-muted"
            >
              <BookOpen className="size-4 text-muted-foreground" />
              <span className="flex-1">Documentation</span>
            </a>
          </DropdownMenuItem>

          {/* Feedback */}
          <DropdownMenuItem asChild>
            <Link
              to="/settings"
              search={{ tab: 'general' } as never}
              className="flex cursor-pointer items-center gap-2.5 rounded-lg px-2.5 py-1.5 text-xs hover:bg-muted"
            >
              <MessageSquare className="size-4 text-muted-foreground" />
              <span className="flex-1">Feedback</span>
            </Link>
          </DropdownMenuItem>

          {/* Sponsor us */}
          <DropdownMenuItem asChild>
            <a
              href="https://github.com/sponsors"
              target="_blank"
              rel="noopener noreferrer"
              className="flex cursor-pointer items-center gap-2.5 rounded-lg px-2.5 py-1.5 text-xs hover:bg-muted"
            >
              <Heart className="size-4 text-rose-500" />
              <span className="flex-1">Sponsor us</span>
            </a>
          </DropdownMenuItem>

          <DropdownMenuSeparator className="bg-border/60" />

          {/* Log out */}
          <DropdownMenuItem
            onClick={() => logout()}
            className="flex cursor-pointer items-center gap-2.5 rounded-lg px-2.5 py-1.5 text-destructive text-xs hover:bg-destructive/10"
          >
            <LogOut className="size-4" />
            <span>Log out</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {!collapsed && (
        <button
          type="button"
          onClick={onToggleCollapse}
          aria-label="Collapse sidebar (⌘B)"
          className="flex size-7 shrink-0 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-foreground"
        >
          <PanelLeftClose className="size-3.5" />
        </button>
      )}
    </div>
  );
}
