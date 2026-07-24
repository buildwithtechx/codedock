import { createFileRoute } from '@tanstack/react-router';
import { Loader2, Plus, ServerIcon, Terminal } from 'lucide-react';
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
import { useCreateServer, useListServers } from '#/hooks/useServers';
import type { Server } from '#/interfaces/server';

export const Route = createFileRoute('/_dashboard/servers')({
  component: ServersPage,
});

function ServersPage() {
  const { data: servers, isLoading } = useListServers();
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
      ) : servers?.length === 0 ? (
        <div className="flex h-64 flex-col items-center justify-center rounded-xl border border-border border-dashed bg-card/40">
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

function ServerCard({ server }: { server: Server }) {
  const statusColor = server.status === 'online' ? 'bg-green-500' : 'bg-yellow-500';

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="font-medium text-sm">{server.name}</CardTitle>
        <ServerIcon className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="mt-2 flex items-center gap-2">
          <div className={`h-2 w-2 rounded-full ${statusColor}`} />
          <span className="font-medium text-sm capitalize">{server.status}</span>
        </div>
        <CardDescription className="mt-2 text-xs">IP: {server.ipAddress}</CardDescription>
        <CardDescription className="text-xs">
          Last seen: {server.lastSeenAt ? new Date(server.lastSeenAt).toLocaleString() : 'Never'}
        </CardDescription>
      </CardContent>
    </Card>
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
