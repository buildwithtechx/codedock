import { createFileRoute, Link } from '@tanstack/react-router';
import { AlertCircle, Cpu, HardDrive, Loader2, MemoryStick, Plus, ServerIcon } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Progress } from '#/components/ui/progress';
import { useCreateServer, useListServers } from '#/hooks/use-servers';
import { parseServerMetrics, type Server } from '#/interfaces/server';

export const Route = createFileRoute('/_dashboard/servers')({
  component: ServersPage,
});

function ServersPage() {
  const { data: servers, isLoading, isError, refetch } = useListServers();
  const [createOpen, setCreateOpen] = useState(false);
  const [newServer, setNewServer] = useState<Server | null>(null);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <ServerIcon className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Servers</h1>
            <p className="text-muted-foreground text-sm">
              Manage your distributed fleet of worker servers.
            </p>
          </div>
        </div>

        <Button onClick={() => setCreateOpen(true)} className="gap-2">
          <Plus className="h-4 w-4" />
          NEW SERVER
        </Button>
      </div>

      {isLoading ? (
        <div className="flex justify-center p-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : isError ? (
        <div className="flex h-64 flex-col items-center justify-center rounded-xl border border-destructive/30 bg-card p-6 text-center">
          <AlertCircle className="mb-4 h-8 w-8 text-destructive" />
          <h3 className="font-bold text-foreground text-lg tracking-tight">
            Could not load servers
          </h3>
          <p className="mt-1 text-muted-foreground text-sm">Check your connection and try again.</p>
          <Button className="mt-6" variant="outline" onClick={() => void refetch()}>
            Try again
          </Button>
        </div>
      ) : servers?.length === 0 ? (
        <div className="flex h-64 flex-col items-center justify-center rounded-xl border border border-dashed bg-card/40">
          <ServerIcon className="mb-4 h-8 w-8 text-muted-foreground" />
          <h3 className="font-bold text-foreground text-lg tracking-tight">No servers yet</h3>
          <p className="mt-1 text-center text-muted-foreground text-sm">
            Add a server to start deploying your applications globally.
          </p>
          <Button className="mt-6 gap-2" onClick={() => setCreateOpen(true)}>
            <Plus className="h-4 w-4" />
            ADD SERVER
          </Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          {servers?.map((server) => (
            <ServerCard key={server.id} server={server} />
          ))}
        </div>
      )}

      <CreateServerModal
        open={createOpen}
        onOpenChange={(open) => {
          setCreateOpen(open);
          if (!open) setNewServer(null);
        }}
        onCreated={(server) => setNewServer(server)}
        newServer={newServer}
      />
    </div>
  );
}

function formatBytes(bytes: number) {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${Number.parseFloat((bytes / k ** i).toFixed(1))} ${sizes[i]}`;
}

function ServerCard({ server }: { server: Server }) {
  const statusColor = server.status === 'online' ? 'bg-green-500' : 'bg-yellow-500';
  const m = parseServerMetrics(server.metrics);

  const memPercent =
    m && m.memory_limit_bytes > 0 ? (m.memory_usage_bytes / m.memory_limit_bytes) * 100 : 0;
  const diskPercent =
    m && m.disk_total_bytes > 0 ? (m.disk_usage_bytes / m.disk_total_bytes) * 100 : 0;
  const cpuPercent = m ? m.cpu_usage_percentage : 0;

  return (
    <Link to="/servers/$serverId" params={{ serverId: server.id }} className="block">
      <Card className="flex h-full flex-col overflow-hidden transition-all hover:shadow-md">
        <CardHeader className="flex flex-row items-center justify-between border-b bg-muted/20 px-6 py-4 pb-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
              <ServerIcon className="h-5 w-5 text-primary" />
            </div>
            <div>
              <CardTitle className="font-semibold text-base">{server.name}</CardTitle>
              <CardDescription className="text-xs">
                {server.isControlPlane ? 'Local deployment runtime' : server.ipAddress}
              </CardDescription>
            </div>
          </div>
          <div className="flex items-center gap-2 rounded-full border bg-background px-3 py-1 shadow-sm">
            <div className={`h-2.5 w-2.5 rounded-full ${statusColor} animate-pulse shadow-sm`} />
            <span className="font-semibold text-xs uppercase tracking-wider">{server.status}</span>
          </div>
        </CardHeader>

        <CardContent className="flex-1 space-y-5 p-6">
          <div className="grid grid-cols-3 gap-4">
            {/* CPU */}
            <div className="space-y-2">
              <div className="flex items-center justify-between text-xs">
                <span className="flex items-center gap-1.5 font-medium text-muted-foreground">
                  <Cpu className="h-3.5 w-3.5" /> CPU
                </span>
                <span className="font-semibold">{cpuPercent.toFixed(1)}%</span>
              </div>
              <Progress value={cpuPercent} className="h-1.5" />
            </div>

            {/* RAM */}
            <div className="space-y-2">
              <div className="flex items-center justify-between text-xs">
                <span className="flex items-center gap-1.5 font-medium text-muted-foreground">
                  <MemoryStick className="h-3.5 w-3.5" /> RAM
                </span>
                <span className="font-semibold">{memPercent.toFixed(1)}%</span>
              </div>
              <Progress value={memPercent} className="h-1.5" />
              <div className="text-right text-[10px] text-muted-foreground">
                {m
                  ? `${formatBytes(m.memory_usage_bytes)} / ${formatBytes(m.memory_limit_bytes)}`
                  : 'N/A'}
              </div>
            </div>

            {/* Disk */}
            <div className="space-y-2">
              <div className="flex items-center justify-between text-xs">
                <span className="flex items-center gap-1.5 font-medium text-muted-foreground">
                  <HardDrive className="h-3.5 w-3.5" /> Disk
                </span>
                <span className="font-semibold">{diskPercent.toFixed(1)}%</span>
              </div>
              <Progress value={diskPercent} className="h-1.5" />
              <div className="text-right text-[10px] text-muted-foreground">
                {m
                  ? `${formatBytes(m.disk_usage_bytes)} / ${formatBytes(m.disk_total_bytes)}`
                  : 'N/A'}
              </div>
            </div>
          </div>

          <div className="pt-2">
            <CardDescription className="text-center text-[11px]">
              Last heartbeat:{' '}
              {server.lastSeenAt ? new Date(server.lastSeenAt).toLocaleString() : 'Never'}
            </CardDescription>
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}

function CreateServerModal({
  open,
  onOpenChange,
  onCreated,
  newServer,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onCreated: (server: Server) => void;
  newServer: Server | null;
}) {
  const [name, setName] = useState('');
  const [ipAddress, setIpAddress] = useState('');
  const { mutateAsync: createServer, isPending } = useCreateServer();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !ipAddress) return;
    try {
      const server = await createServer({ name, ipAddress });
      toast.success('Server created successfully');
      onCreated(server);
    } catch (err: any) {
      toast.error(err.message || 'Failed to create server');
    }
  };

  const curlCommand = `curl -sL https://get.codedock.dev | bash -s -- --key ${newServer?.workerToken}`;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        {newServer ? (
          <div>
            <DialogHeader>
              <DialogTitle>Connect your Server</DialogTitle>
              <DialogDescription>
                Run the following command on your server to install the Codedock Worker and connect
                it to your dashboard.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-6">
              <div className="relative break-all rounded-md bg-muted p-4 font-mono text-sm">
                {curlCommand}
              </div>
              <p className="text-muted-foreground text-xs">
                Note: Keep this token secret. It gives full access to your server.
              </p>
            </div>
            <DialogFooter>
              <Button onClick={() => onOpenChange(false)}>Done</Button>
            </DialogFooter>
          </div>
        ) : (
          <form onSubmit={handleSubmit}>
            <DialogHeader>
              <DialogTitle>Add Server</DialogTitle>
              <DialogDescription>Add a new worker node to your Codedock cluster.</DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-6">
              <div className="space-y-2">
                <label htmlFor="name" className="font-medium text-sm">
                  Name
                </label>
                <Input
                  id="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="EU Production Node 1"
                  disabled={isPending}
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="ipAddress" className="font-medium text-sm">
                  IP Address
                </label>
                <Input
                  id="ipAddress"
                  value={ipAddress}
                  onChange={(e) => setIpAddress(e.target.value)}
                  placeholder="198.51.100.1"
                  disabled={isPending}
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isPending}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={!name || !ipAddress || isPending}>
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Create Server
              </Button>
            </DialogFooter>
          </form>
        )}
      </DialogContent>
    </Dialog>
  );
}
