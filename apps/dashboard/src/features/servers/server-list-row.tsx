import { Link } from '@tanstack/react-router';
import {
  ChevronRight,
  Copy,
  Cpu,
  EllipsisVertical,
  KeyRound,
  LayoutGrid,
  Lock,
  Server as ServerIcon,
  Settings,
  TerminalSquare,
  Trash2,
} from 'lucide-react';
import { toast } from 'sonner';
import { BlurIp } from '#/components/blur-ip';
import { Button } from '@codedock/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@codedock/ui/dropdown-menu';
import type { Server } from '#/interfaces/server';
import { cn } from '#/lib/utils';

export type ServerReachability = 'online' | 'offline' | 'unknown';

export function reachabilityOf(server: Server): ServerReachability {
  const status = server.status.toLowerCase();
  if (status === 'online') return 'online';
  if (status === 'offline') return 'offline';
  return 'unknown';
}

const reachabilityTone: Record<ServerReachability, { ring: string; text: string; label: string }> =
  {
    online: { ring: 'border-success', text: 'text-success', label: 'Online' },
    offline: { ring: 'border-destructive', text: 'text-destructive', label: 'Offline' },
    unknown: { ring: 'border-warning animate-pulse', text: 'text-warning', label: 'Unknown' },
  };

interface ServerListRowProps {
  server: Server;
  onRemove: (server: Server) => void;
}

export function ServerListRow({ server, onRemove }: ServerListRowProps) {
  const reach = reachabilityOf(server);
  const tone = reachabilityTone[reach];
  const host = server.sshHost || server.ipAddress;
  const AuthIcon = server.sshAuthMethod === 'password' ? Lock : KeyRound;

  const handleCopyHost = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (!host) return;
    void navigator.clipboard.writeText(host);
    toast.success('Host IP copied to clipboard');
  };

  return (
    <div className="flex items-center pe-3">
      <Link
        to="/servers/$serverId"
        params={{ serverId: server.id }}
        className="group grid min-w-0 flex-1 grid-cols-[auto_minmax(0,1fr)] items-center gap-x-3 gap-y-2 px-4 py-3 text-start transition-colors hover:bg-muted/40 sm:flex sm:gap-3.5 sm:px-5"
      >
        <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-muted/60 transition-colors group-hover:bg-muted">
          {server.isLocal ? (
            <Cpu className="size-[18px] text-foreground/70" />
          ) : (
            <ServerIcon className="size-[18px] text-foreground/70" />
          )}
        </div>

        <div className="min-w-0 text-start sm:min-w-44 sm:max-w-xs lg:max-w-sm">
          <div className="flex items-center gap-2">
            <span className="truncate font-medium text-foreground text-sm">{server.name}</span>
            {server.isLocal && (
              <span className="shrink-0 whitespace-nowrap rounded bg-primary/10 px-1.5 py-0.5 font-medium text-[10px] text-primary">
                This server
              </span>
            )}
          </div>
          <p className="mt-0.5 truncate font-mono text-muted-foreground text-xs">
            {server.isLocal ? 'Local deployment runtime' : <BlurIp>{host}</BlurIp>}
          </p>
        </div>

        <div className="col-span-2 flex min-w-0 flex-1 flex-wrap items-center gap-2 overflow-hidden sm:flex-nowrap sm:gap-3">
          {!server.isLocal && server.sshAuthMethod && (
            <span className="inline-flex shrink-0 items-center gap-1.5 rounded-md bg-muted/60 px-2 py-0.5 text-muted-foreground text-xs">
              <AuthIcon className="size-3.5" />
              {server.sshAuthMethod === 'password' ? 'Password' : 'SSH key'}
            </span>
          )}
          {!server.isLocal && server.sshTransport && server.sshTransport !== 'direct' && (
            <span className="inline-flex shrink-0 items-center rounded-md bg-muted/60 px-2 py-0.5 text-muted-foreground text-xs capitalize">
              {server.sshTransport}
            </span>
          )}
          {server.region && (
            <span className="hidden shrink-0 items-center rounded-md bg-muted/60 px-2 py-0.5 text-muted-foreground text-xs md:inline-flex">
              {server.region}
            </span>
          )}
        </div>

        <div className="col-span-2 flex shrink-0 items-center gap-4">
          <span className={cn('inline-flex items-center gap-1.5 font-medium text-xs', tone.text)}>
            <span className={cn('size-2.5 rounded-full border-2', tone.ring)} />
            {tone.label}
          </span>
          <ChevronRight className="size-4 text-muted-foreground/40 transition-colors group-hover:text-muted-foreground" />
        </div>
      </Link>

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon" aria-label={`${server.name} actions`}>
            <EllipsisVertical className="size-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          align="end"
          className="w-48 rounded-xl border border-border/60 bg-popover/95 p-1 text-popover-foreground shadow-xl backdrop-blur-md"
        >
          <DropdownMenuItem asChild>
            <Link
              to="/servers/$serverId"
              params={{ serverId: server.id }}
              className="flex cursor-pointer items-center gap-2 rounded-lg px-2.5 py-1.5 text-xs font-medium"
            >
              <LayoutGrid className="size-3.5 text-muted-foreground" />
              <span className="whitespace-nowrap">View server</span>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem asChild>
            <Link
              to="/servers/$serverId"
              params={{ serverId: server.id }}
              search={{ tab: 'terminal' }}
              className="flex cursor-pointer items-center gap-2 rounded-lg px-2.5 py-1.5 text-xs font-medium"
            >
              <TerminalSquare className="size-3.5 text-muted-foreground" />
              <span className="whitespace-nowrap">Open terminal</span>
            </Link>
          </DropdownMenuItem>
          {host && (
            <DropdownMenuItem
              onClick={handleCopyHost}
              className="flex cursor-pointer items-center gap-2 rounded-lg px-2.5 py-1.5 text-xs font-medium"
            >
              <Copy className="size-3.5 text-muted-foreground" />
              <span className="whitespace-nowrap">Copy host address</span>
            </DropdownMenuItem>
          )}
          <DropdownMenuItem asChild>
            <Link
              to="/servers/$serverId"
              params={{ serverId: server.id }}
              search={{ tab: 'settings' }}
              className="flex cursor-pointer items-center gap-2 rounded-lg px-2.5 py-1.5 text-xs font-medium"
            >
              <Settings className="size-3.5 text-muted-foreground" />
              <span className="whitespace-nowrap">Server settings</span>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuSeparator className="my-1 bg-border/40" />
          <DropdownMenuItem
            variant="destructive"
            onClick={() => onRemove(server)}
            className="flex cursor-pointer items-center gap-2 rounded-lg px-2.5 py-1.5 text-xs font-medium text-destructive hover:bg-destructive/10"
          >
            <Trash2 className="size-3.5" />
            <span className="whitespace-nowrap">Delete server</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
