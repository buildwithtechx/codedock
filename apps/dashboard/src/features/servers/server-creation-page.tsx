import { Link, useNavigate } from '@tanstack/react-router';
import { ArrowLeft, Check, Copy, ServerIcon } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { useCreateServer } from '#/hooks/use-servers';
import type { Server } from '#/interfaces/server';

export function ServerCreationPage() {
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [ipAddress, setIPAddress] = useState('');
  const [server, setServer] = useState<Server | null>(null);
  const { mutateAsync: createServer, isPending } = useCreateServer();

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    try {
      const created = await createServer({ name, ipAddress });
      setServer(created);
      toast.success('Server created');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to create server');
    }
  };

  const command = server
    ? `curl -sL https://get.codedock.dev | bash -s -- --key ${server.workerToken}`
    : '';

  const copyCommand = async () => {
    await navigator.clipboard.writeText(command);
    toast.success('Install command copied');
  };

  return (
    <div className="grid min-w-0 gap-6 xl:grid-cols-[minmax(0,1fr)_21.25rem]">
      <main className="min-w-0">
        <Link
          to="/servers"
          className="inline-flex items-center gap-2 text-muted-foreground text-sm transition-colors hover:text-foreground"
        >
          <ArrowLeft className="h-4 w-4" />
          Servers
        </Link>
        <header className="mt-6">
          <p className="font-medium text-muted-foreground text-sm">Deployment target</p>
          <h1 className="mt-1 font-semibold text-2xl tracking-tight">Add a server</h1>
          <p className="mt-1 max-w-2xl text-muted-foreground text-sm">
            Register a worker, then install the Codedock agent using a one-time connection command.
          </p>
        </header>

        {server ? (
          <section className="mt-8 max-w-3xl rounded-2xl border border-border bg-card p-6 shadow-sm">
            <div className="flex items-center gap-3">
              <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-emerald-500/12 text-emerald-500">
                <Check className="h-4 w-4" />
              </div>
              <div>
                <h2 className="font-semibold text-sm">Connect {server.name}</h2>
                <p className="text-muted-foreground text-xs">
                  Run this command on the server to complete registration.
                </p>
              </div>
            </div>
            <div className="mt-6 break-all rounded-xl border border-border bg-background p-4 font-mono text-sm">
              {command}
            </div>
            <div className="mt-5 flex flex-wrap gap-3">
              <Button onClick={() => void copyCommand()} className="gap-2">
                <Copy className="h-4 w-4" />
                Copy install command
              </Button>
              <Button
                variant="outline"
                onClick={() =>
                  void navigate({ to: '/servers/$serverId', params: { serverId: server.id } })
                }
              >
                View server
              </Button>
            </div>
          </section>
        ) : (
          <form
            onSubmit={handleSubmit}
            className="mt-8 max-w-3xl rounded-2xl border border-border bg-card p-6 shadow-sm"
          >
            <div className="flex items-center gap-3">
              <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/12 text-primary">
                <ServerIcon className="h-4 w-4" />
              </div>
              <div>
                <h2 className="font-semibold text-sm">Server details</h2>
                <p className="text-muted-foreground text-xs">Use the server's reachable address.</p>
              </div>
            </div>
            <div className="mt-6 grid gap-5 sm:grid-cols-2">
              <div className="space-y-2">
                <label htmlFor="server-name" className="font-medium text-sm">
                  Name
                </label>
                <Input
                  id="server-name"
                  value={name}
                  onChange={(event) => setName(event.target.value)}
                  placeholder="EU production node"
                  required
                />
              </div>
              <div className="space-y-2">
                <label htmlFor="server-address" className="font-medium text-sm">
                  IP address
                </label>
                <Input
                  id="server-address"
                  value={ipAddress}
                  onChange={(event) => setIPAddress(event.target.value)}
                  placeholder="198.51.100.1"
                  required
                />
              </div>
            </div>
            <div className="mt-6 flex justify-end gap-3">
              <Link to="/servers">
                <Button type="button" variant="ghost">
                  Cancel
                </Button>
              </Link>
              <Button type="submit" disabled={isPending}>
                {isPending ? 'Creating server...' : 'Create server'}
              </Button>
            </div>
          </form>
        )}
      </main>

      <aside className="hidden xl:sticky xl:top-6 xl:block xl:self-start">
        <section className="rounded-2xl border border-border bg-card p-5 shadow-sm">
          <div className="flex items-center gap-2">
            <ServerIcon className="h-4 w-4 text-primary" />
            <h2 className="font-semibold text-sm">Server onboarding</h2>
          </div>
          <p className="mt-3 text-muted-foreground text-sm leading-6">
            Keep the install command private. It authorizes the worker to report health and receive
            deployment work.
          </p>
        </section>
      </aside>
    </div>
  );
}
