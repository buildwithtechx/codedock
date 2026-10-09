import { createFileRoute, Link, useNavigate } from '@tanstack/react-router';
import {
  ArrowLeft,
  Container,
  EllipsisVertical,
  LayoutGrid,
  Loader2,
  RefreshCw,
  Server as ServerIcon,
  Settings,
  Shield,
  TerminalSquare,
  Trash2,
} from 'lucide-react';
import { useState } from 'react';
import { z } from 'zod';
import { Button } from '#/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '#/components/ui/empty';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import { ServerComponentsTab } from '#/features/servers/server-components-tab';
import { ServerConnectionBanner } from '#/features/servers/server-connection-banner';
import { ServerConnectionCard } from '#/features/servers/server-connection-card';
import { ServerDeleteDialog } from '#/features/servers/server-delete-dialog';
import { reachabilityOf } from '#/features/servers/server-list-row';
import { ServerOverviewTab } from '#/features/servers/server-overview-tab';
import { ServerSecurityTab } from '#/features/servers/server-security-tab';
import { ServerSettingsTab } from '#/features/servers/server-settings-tab';
import { ServerTerminalTab } from '#/features/servers/server-terminal-tab';
import { useServer } from '#/hooks/use-servers';
import type { Server } from '#/interfaces/server';
import { cn } from '#/lib/utils';

const tabs = [
  { key: 'overview', label: 'Overview', icon: LayoutGrid },
  { key: 'components', label: 'Components', icon: Container },
  { key: 'security', label: 'Security', icon: Shield },
  { key: 'terminal', label: 'Terminal', icon: TerminalSquare },
  { key: 'settings', label: 'Settings', icon: Settings },
] as const;

type ServerTab = (typeof tabs)[number]['key'];

const tabKeys = tabs.map((tab) => tab.key);

export const Route = createFileRoute('/_dashboard/servers/$serverId')({
  validateSearch: z.object({
    tab: z.string().optional(),
  }),
  component: ServerDetailsPage,
});

function ServerDetailsPage() {
  const { serverId } = Route.useParams();
  const { tab: tabParam } = Route.useSearch();
  const navigate = useNavigate();
  const { data: server, isLoading, isRefetching, refetch } = useServer(serverId);
  const [removeTarget, setRemoveTarget] = useState<Server | null>(null);

  const activeTab: ServerTab = tabKeys.includes(tabParam as ServerTab)
    ? (tabParam as ServerTab)
    : 'overview';
  const changeTab = (value: string) =>
    navigate({
      to: '/servers/$serverId',
      params: { serverId },
      search: { tab: value },
      replace: true,
    });

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!server) {
    return (
      <Empty>
        <EmptyHeader>
          <EmptyMedia variant="icon">
            <ServerIcon />
          </EmptyMedia>
          <EmptyTitle>Server not found</EmptyTitle>
          <EmptyDescription>
            This server no longer exists or you don't have access to it.
          </EmptyDescription>
        </EmptyHeader>
        <EmptyContent>
          <Button asChild variant="outline">
            <Link to="/servers">Back to Servers</Link>
          </Button>
        </EmptyContent>
      </Empty>
    );
  }

  const reach = reachabilityOf(server);
  const statusTone =
    reach === 'online'
      ? 'bg-success/10 text-success'
      : reach === 'offline'
        ? 'bg-destructive/10 text-destructive'
        : 'bg-warning/10 text-warning';
  const statusLabel = reach === 'online' ? 'Online' : reach === 'offline' ? 'Offline' : 'Unknown';

  return (
    <div>
      <div className="mb-6 grid grid-cols-[auto_minmax(0,1fr)_auto] items-center gap-x-3 gap-y-1">
        <Link
          to="/servers"
          className="flex h-8 w-8 items-center justify-center rounded-lg transition-colors hover:bg-muted"
          aria-label="Back to servers"
        >
          <ArrowLeft className="h-4 w-4 text-muted-foreground" />
        </Link>
        <div className="min-w-0">
          <h1 className="truncate font-medium text-2xl text-foreground/80 tracking-[-0.2px]">
            {server.name}
          </h1>
        </div>
        <div className="flex items-center gap-1.5">
          <Button variant="secondary" size="sm" onClick={() => changeTab('settings')}>
            <Settings className="size-4" />
            <span className="hidden sm:inline">Edit</span>
          </Button>
          <Button
            variant="ghost"
            size="icon"
            disabled={isRefetching}
            aria-label="Refresh server"
            onClick={() => void refetch()}
          >
            <RefreshCw className={`size-4 ${isRefetching ? 'animate-spin' : ''}`} />
          </Button>
          {!server.isLocal && (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" aria-label="Server actions">
                  <EllipsisVertical className="size-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem variant="destructive" onClick={() => setRemoveTarget(server)}>
                  <Trash2 className="size-4" />
                  Remove server
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          )}
        </div>
        <div className="col-span-2 col-start-2 flex min-w-0 flex-wrap items-center gap-2">
          <p className="min-w-0 break-all font-mono text-muted-foreground text-sm">
            {server.isLocal
              ? 'local deployment runtime'
              : `${server.sshUser ?? 'root'}@${server.sshHost ?? server.ipAddress}`}
          </p>
          <span className={cn('rounded-full px-2 py-0.5 font-medium text-xs', statusTone)}>
            {statusLabel}
          </span>
        </div>
      </div>

      <ServerConnectionBanner
        server={server}
        retrying={isRefetching}
        onRetry={() => void refetch()}
      />

      <Tabs value={activeTab} onValueChange={changeTab} className="gap-6">
        <TabsList variant="line">
          {tabs.map((tab) => (
            <TabsTrigger key={tab.key} value={tab.key} className="gap-1.5">
              <tab.icon className="size-4" />
              {tab.label}
            </TabsTrigger>
          ))}
        </TabsList>

        <div className="grid grid-cols-1 items-start gap-6 xl:grid-cols-[minmax(0,1fr)_340px]">
          <div className="min-w-0 xl:col-start-1 xl:row-start-1">
            <TabsContent value="overview">
              <ServerOverviewTab server={server} />
            </TabsContent>
            <TabsContent value="components">
              <ServerComponentsTab />
            </TabsContent>
            <TabsContent value="security">
              <ServerSecurityTab />
            </TabsContent>
            <TabsContent value="terminal">
              <ServerTerminalTab server={server} />
            </TabsContent>
            <TabsContent value="settings">
              <ServerSettingsTab server={server} />
            </TabsContent>
          </div>

          <div className="space-y-4 xl:sticky xl:top-6 xl:col-start-2 xl:row-start-1 xl:self-start">
            <ServerConnectionCard server={server} />
          </div>
        </div>
      </Tabs>

      <ServerDeleteDialog
        server={removeTarget}
        onClose={() => setRemoveTarget(null)}
        onRemoved={() => navigate({ to: '/servers' })}
      />
    </div>
  );
}
