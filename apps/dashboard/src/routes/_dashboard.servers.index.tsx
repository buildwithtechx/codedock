import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { Cpu, Plus, RefreshCw, Search, Server as ServerIcon } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { Input } from '#/components/ui/input';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { createServer } from '#/features/servers/api';
import { ServerDeleteDialog } from '#/features/servers/server-delete-dialog';
import { ServerEmptyState } from '#/features/servers/server-empty-state';
import { ServerFleetPanel } from '#/features/servers/server-fleet-panel';
import { ServerListRow } from '#/features/servers/server-list-row';
import { ServerListSkeleton } from '#/features/servers/server-list-skeleton';
import { useListServers } from '#/hooks/use-servers';
import type { Server } from '#/interfaces/server';

export const Route = createFileRoute('/_dashboard/servers/')({
  component: ServersPage,
});

function ServersPage() {
  const navigate = useNavigate();
  const { data: servers, isLoading, isError, refetch, isRefetching } = useListServers();
  const [search, setSearch] = useState('');
  const [removeServer, setRemoveServer] = useState<Server | null>(null);
  const [addingLocal, setAddingLocal] = useState(false);

  const rows = servers ?? [];
  const hasLocalServer = rows.some((server) => server.isLocal);
  const showFilters = rows.length > 6;
  const query = search.trim().toLowerCase();
  const visible = rows.filter(
    (server) =>
      !query ||
      server.name.toLowerCase().includes(query) ||
      server.ipAddress.toLowerCase().includes(query) ||
      (server.sshHost ?? '').toLowerCase().includes(query)
  );

  const addThisMachine = async () => {
    setAddingLocal(true);
    try {
      const server = await createServer({
        name: 'Local Daemon Node',
        isLocal: true,
        ipAddress: '127.0.0.1',
      });
      toast.success('Local Docker node registered');
      await refetch();
      navigate({ to: '/servers/$serverId', params: { serverId: server.id } });
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to register local node');
    } finally {
      setAddingLocal(false);
    }
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Servers"
        description={isLoading ? 'Loading infrastructure...' : `${rows.length} connected servers`}
        action={
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="icon"
              disabled={isRefetching}
              aria-label="Refresh servers"
              onClick={() => void refetch()}
            >
              <RefreshCw className={`size-4 ${isRefetching ? 'animate-spin' : ''}`} />
            </Button>
            {hasLocalServer ? (
              <Button className="gap-2" onClick={() => navigate({ to: '/servers/new' })}>
                <Plus className="h-4 w-4" />
                Add server
              </Button>
            ) : (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button className="gap-2">
                    <Plus className="h-4 w-4" />
                    Add server
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={() => navigate({ to: '/servers/new' })}>
                    <ServerIcon className="size-4" />
                    Add remote server
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => void addThisMachine()} disabled={addingLocal}>
                    <Cpu className="size-4" />
                    {addingLocal ? 'Connecting…' : 'Use this machine'}
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        }
      />

      {isError ? (
        <QueryErrorState
          title="Servers are unavailable"
          description="Codedock could not load the current runtime fleet."
          onRetry={() => void refetch()}
        />
      ) : !isLoading && rows.length === 0 ? (
        <ServerEmptyState onAddLocal={() => void addThisMachine()} addingLocal={addingLocal} />
      ) : (
        <div className="grid grid-cols-1 items-start gap-6 xl:grid-cols-[minmax(0,1fr)_340px]">
          <div className="min-w-0">
            {!isLoading && showFilters && (
              <div className="relative mb-4">
                <Search className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  value={search}
                  onChange={(event) => setSearch(event.target.value)}
                  placeholder="Search servers by name or host…"
                  className="pl-9"
                  aria-label="Search servers"
                />
              </div>
            )}
            <div className="divide-y divide-border/50 rounded-2xl bg-card">
              {isLoading ? (
                <ServerListSkeleton />
              ) : visible.length === 0 ? (
                <p className="px-5 py-8 text-center text-muted-foreground text-sm">
                  No servers match this search.
                </p>
              ) : (
                visible.map((server) => (
                  <ServerListRow key={server.id} server={server} onRemove={setRemoveServer} />
                ))
              )}
            </div>
          </div>

          <div className="xl:sticky xl:top-6 xl:self-start">
            <ServerFleetPanel servers={rows} loading={isLoading} />
          </div>
        </div>
      )}

      <ServerDeleteDialog server={removeServer} onClose={() => setRemoveServer(null)} />
    </div>
  );
}
